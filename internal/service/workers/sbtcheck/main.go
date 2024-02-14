package sbtcheck

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	iden3 "github.com/iden3/go-iden3-core/v2"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/referralid"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/sbtcheck/verifiers"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

type runner struct {
	network
	db    *pgdb.DB
	types evtypes.Types
	log   *logan.Entry
}

type network struct {
	name         string
	filterer     filterer
	blockHandler blockHandler
	timeout      time.Duration
	fromBlock    uint64
	blockWindow  uint64
	maxBlocks    uint64
	disabled     bool
}

type blockHandler interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type filterer interface {
	FilterSBTIdentityProved(*bind.FilterOpts, []*big.Int) (*verifiers.SBTIdentityVerifierSBTIdentityProvedIterator, error)
}

type extConfig interface {
	comfig.Logger
	pgdb.Databaser
	evtypes.EventTypeser
	SbtChecker
}

func Run(ctx context.Context, cfg extConfig) {
	log := cfg.Log().WithField("who", "sbt-checker")
	if cfg.EventTypes().IsExpired(evtypes.TypeGetPoH) {
		log.Warn("PoH event is expired, SBT check will not run")
		return
	}

	var wg sync.WaitGroup
	for name, net := range cfg.SbtCheck().networks {
		if net.disabled {
			log.Infof("SBT check: network %s disabled", name)
			continue
		}

		log.Infof("SBT check: running for network %s", name)
		wg.Add(1)

		r := &runner{
			network: net,
			db:      cfg.DB(),
			types:   cfg.EventTypes(),
			log:     log.WithField("network", name),
		}

		runnerName := fmt.Sprintf("sbt-checker[%s]", net.name)
		go func() {
			running.WithBackOff(ctx, r.log, runnerName, r.subscription,
				30*time.Second, 5*time.Second, 30*time.Second)
			wg.Done()
		}()
	}

	wg.Wait()
	log.Infof("SBT check: all network checkers stopped")
}

func (r *runner) subscription(ctx context.Context) error {
	toBlock, err := r.getLastBlock(ctx)
	if err != nil {
		return fmt.Errorf("get last block: %w", err)
	}
	if toBlock == nil {
		return nil
	}

	r.log.Debugf("Starting subscription from %d to %d", r.fromBlock, toBlock)
	defer r.log.Debugf("Subscription finished")

	ctx2, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.filterEvents(ctx2, toBlock)
}

func (r *runner) getLastBlock(ctx context.Context) (*uint64, error) {
	ctx2, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	lastBlock, err := r.blockHandler.BlockNumber(ctx2)
	if err != nil {
		return nil, fmt.Errorf("get last block number: %w", err)
	}

	lastBlock -= r.blockWindow

	if lastBlock < r.fromBlock {
		r.log.Infof("Skipping window: start=%d > finish=%d", r.fromBlock, lastBlock)
		return nil, nil
	}

	if r.fromBlock+r.maxBlocks < lastBlock {
		r.log.Debugf("maxBlockPerRequest limit exceeded: setting last block to %d instead of %d", r.fromBlock+r.maxBlocks, lastBlock)
		lastBlock = r.fromBlock + r.maxBlocks
	}

	return &lastBlock, nil
}

func (r *runner) filterEvents(ctx context.Context, toBlock *uint64) error {
	it, err := r.filterer.FilterSBTIdentityProved(&bind.FilterOpts{
		Start:   r.fromBlock,
		End:     toBlock,
		Context: ctx,
	}, nil)
	if err != nil {
		return fmt.Errorf("filter SBTIdentityProved events: %w", err)
	}

	defer func() {
		// https://ethereum.stackexchange.com/questions/8199/are-both-the-eth-newfilter-from-to-fields-inclusive
		// End in FilterLogs is inclusive
		r.fromBlock = *toBlock + 1
		_ = it.Close()
	}()

	for it.Next() {
		evt := it.Event
		if evt == nil {
			r.log.Error("Got nil event")
			continue
		}

		if err = r.handleEvent(*evt); err != nil {
			r.log.WithError(err).Error("Failed to handle event")
			continue
		}
	}

	return nil
}

func (r *runner) handleEvent(evt verifiers.SBTIdentityVerifierSBTIdentityProved) error {
	r.log.WithFields(map[string]any{
		"tx_hash":   evt.Raw.TxHash,
		"tx_index":  evt.Raw.TxIndex,
		"log_index": evt.Raw.Index,
		"block":     evt.Raw.BlockNumber,
	}).Debugf("Got SBTIdentityProved event (identityId=%s)", evt.IdentityId)

	did, err := parseDidFromUint256(evt.IdentityId)
	if err != nil {
		return fmt.Errorf("parse did from uint256 (identityId=%s): %w", evt.IdentityId, err)
	}

	if err = r.createBalanceIfAbsent(did); err != nil {
		return fmt.Errorf("get or create balance (did=%s): %w", did, err)
	}

	poh, err := r.findPohEvent(did)
	if err != nil {
		return fmt.Errorf("find PoH event (did=%s): %w", did, err)
	}
	if poh == nil {
		return nil
	}

	if err = r.fulfillPohEvent(*poh); err != nil {
		return fmt.Errorf("fulfill PoH event: %w", err)
	}

	r.log.Infof("Event %s was fulfilled for DID %s", evtypes.TypeGetPoH, did)
	return nil
}

func (r *runner) createBalanceIfAbsent(did string) error {
	balance, err := r.balancesQ().FilterByDID(did).Get()
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}
	if balance != nil {
		r.log.Debugf("Balance exists for DID %s", did)
		return nil
	}

	r.log.Debugf("Balance not found for DID %s, creating new one", did)
	if err = r.createBalance(did); err != nil {
		return fmt.Errorf("create balance: %w", err)
	}

	return nil
}

func (r *runner) findPohEvent(did string) (*data.Event, error) {
	poh, err := r.eventsQ().
		FilterByUserDID(did).
		FilterByType(evtypes.TypeGetPoH).
		Get()
	if err != nil {
		return nil, fmt.Errorf("get PoH event: %w", err)
	}
	if poh == nil {
		return nil, fmt.Errorf("PoH event was not properly added on balance creation")
	}

	if poh.Status != data.EventOpen {
		r.log.Infof("User %s is not eligible for another PoH event (id=%s status=%s)",
			poh.UserDID, poh.ID, poh.Status)
		return nil, nil
	}

	return poh, nil
}

func (r *runner) fulfillPohEvent(poh data.Event) error {
	getPoh := r.types.Get(evtypes.TypeGetPoH)
	if getPoh == nil {
		return fmt.Errorf("event types were not correctly initialized: missing %s", evtypes.TypeGetPoH)
	}

	_, err := r.eventsQ().FilterByID(poh.ID).Update(data.EventFulfilled, nil, &getPoh.Reward)
	if err != nil {
		return fmt.Errorf("update PoH event status and reward: %w", err)
	}

	return nil
}

func (r *runner) createBalance(did string) error {
	return r.eventsQ().Transaction(func() error {
		err := r.balancesQ().Insert(data.Balance{
			DID:        did,
			ReferralID: referralid.New(did),
		})
		if err != nil {
			return fmt.Errorf("insert balance: %w", err)
		}

		err = r.eventsQ().Insert(r.types.PrepareOpenEvents(did)...)
		if err != nil {
			return fmt.Errorf("insert open events: %w", err)
		}

		return nil
	})
}

func parseDidFromUint256(raw *big.Int) (string, error) {
	id, err := iden3.IDFromInt(raw)
	if err != nil {
		return "", fmt.Errorf("parse iden3 ID from bigInt: %w", err)
	}

	did, err := iden3.ParseDIDFromID(id)
	if err != nil {
		return "", fmt.Errorf("parse DID from iden3 ID (id=%s): %w", id, err)
	}

	return did.String(), nil
}

func (r *runner) balancesQ() data.BalancesQ {
	return pg.NewBalances(r.db.Clone())
}

func (r *runner) eventsQ() data.EventsQ {
	return pg.NewEvents(r.db.Clone())
}
