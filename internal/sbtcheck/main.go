package sbtcheck

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/event"
	iden3 "github.com/iden3/go-iden3-core/v2"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/sbtcheck/verifiers"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

// all retries of runner are done on failures
const abnormalPeriod = 5 * time.Second

type Runner struct {
	networks map[string]network
	// ensure to always call .New() for balancesQ and eventsQ
	balancesQ data.BalancesQ
	eventsQ   data.EventsQ
	types     *evtypes.Types
	log       *logan.Entry
}

type network struct {
	events   *verifiers.SBTIdentityVerifierFilterer
	timeout  time.Duration
	disabled bool
}

func (r *Runner) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	for name, net := range r.networks {
		if net.disabled {
			r.log.Infof("SBT check: network %s disabled", name)
			continue
		}

		r.log.Infof("SBT check: running for network %s", name)
		wg.Add(1)

		if err := r.run(ctx, net, &wg); err != nil {
			return fmt.Errorf("run checker for network %s: %w", name, err)
		}
	}

	wg.Wait()
	r.log.Infof("SBT check: all network checkers stopped")
	return nil
}

func (r *Runner) run(ctx context.Context, net network, wg *sync.WaitGroup) error {
	sink := make(chan *verifiers.SBTIdentityVerifierSBTIdentityProved)

	toCtx, cancel := context.WithTimeout(ctx, net.timeout)
	defer cancel()

	sub, err := net.events.WatchSBTIdentityProved(&bind.WatchOpts{Context: toCtx}, sink, nil)
	if err != nil {
		return fmt.Errorf("subscribe to SBTIdentityProved event: %w", err)
	}

	go running.UntilSuccess(ctx, r.log, "sbt-checker", func(ctx context.Context) (bool, error) {
		err = r.subscribe(ctx, sub, sink)
		if err == nil {
			wg.Done()
		}
		return err == nil, err
	}, abnormalPeriod, abnormalPeriod)

	return nil
}

func (r *Runner) subscribe(
	ctx context.Context,
	sub event.Subscription,
	sink chan *verifiers.SBTIdentityVerifierSBTIdentityProved,
) error {

	for {
		select {
		case <-ctx.Done():
			r.log.Info("SBTIdentityProved subscription stopped")
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("SBTIdentityProved subscription error: %w", err)
		case evt := <-sink:
			if evt == nil {
				r.log.Debug("Got nil SBTIdentityProved event from subscription, continue")
				continue
			}
			if err := r.handleEvent(*evt); err != nil {
				return fmt.Errorf("handle event: %w", err)
			}
		}
	}
}

func (r *Runner) handleEvent(evt verifiers.SBTIdentityVerifierSBTIdentityProved) error {
	did, err := parseDidFromUint256(evt.IdentityId)
	if err != nil {
		return fmt.Errorf("parse did from uint256 (identityId=%s): %w", evt.IdentityId, err)
	}

	balanceID, err := r.getOrCreateBalance(did)
	if err != nil {
		return fmt.Errorf("get or create balance (did=%s): %w", did, err)
	}

	poh, err := r.findPohEvent(balanceID)
	if err != nil {
		return fmt.Errorf("find PoH event (balanceID=%s): %w", balanceID, err)
	}
	if err = r.fulfillPohEvent(*poh); err != nil {
		return fmt.Errorf("update PoH event status to fulfilled: %w", err)
	}

	r.log.Infof("Event %s was fulfilled for DID %s", evtypes.TypeGetPoH, did)
	return nil
}

func (r *Runner) getOrCreateBalance(did string) (string, error) {
	balance, err := r.balancesQ.New().FilterByUserDID(did).Get()
	if err != nil {
		return "", fmt.Errorf("get balance: %w", err)
	}
	if balance != nil {
		r.log.Debugf("Balance exists for DID %s", did)
		return balance.ID, nil
	}

	r.log.Debugf("Balance not found for DID %s, creating new one", did)
	id, err := r.createBalance(did)
	if err != nil {
		return "", fmt.Errorf("create balance: %w", err)
	}

	return id, nil
}

func (r *Runner) findPohEvent(bid string) (*data.Event, error) {
	poh, err := r.eventsQ.New().
		FilterByBalanceID(bid).
		FilterByType(evtypes.TypeGetPoH).
		FilterByStatus(data.EventOpen).
		Get()
	if err != nil {
		return nil, fmt.Errorf("get open PoH event: %w", err)
	}
	if poh == nil {
		return nil, fmt.Errorf("PoH event was not properly added on balance creation")
	}

	return poh, nil
}

func (r *Runner) fulfillPohEvent(poh data.Event) error {
	getPoh := r.types.Get(evtypes.TypeGetPoH)
	if getPoh == nil {
		return fmt.Errorf("event types were not correctly initialized: missing %s", evtypes.TypeGetPoH)
	}

	return r.eventsQ.New().FilterByID(poh.ID).Update(data.Event{
		Status: data.EventFulfilled,
		PointsAmount: sql.NullInt32{
			Int32: getPoh.Reward,
			Valid: true,
		},
	})
}

func (r *Runner) createBalance(did string) (string, error) {
	err := r.balancesQ.New().Insert(data.Balance{DID: did})
	if err != nil {
		return "", fmt.Errorf("insert balance: %w", err)
	}

	balance, err := r.balancesQ.New().FilterByUserDID(did).Get()
	if err != nil {
		return "", fmt.Errorf("get balance back: %w", err)
	}

	err = r.eventsQ.New().Insert(r.types.PrepareOpenEvents(balance.ID)...)
	if err != nil {
		return "", fmt.Errorf("insert open events: %w", err)
	}

	return balance.ID, nil
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
