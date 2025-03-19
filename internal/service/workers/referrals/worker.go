package referrals

import (
	"time"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type worker struct {
	name string
	rq   data.ReferralsQ
	bq   data.BalancesQ
	log  *logan.Entry
	exc  config.ExpiredCode
}

type referralPair struct {
	GuestBalance *data.Balance
	Referral     *data.Referral
}

func newWorker(cfg config.Config, workerName string) *worker {
	return &worker{
		name: workerName,
		rq:   pg.NewReferrals(cfg.DB().Clone()),
		bq:   pg.NewBalances(cfg.DB().Clone()),
		log:  cfg.Log().WithField("who", workerName),
		exc:  cfg.ExpiredCode(),
	}
}

func (w *worker) job() error {
	w.rq = w.rq.New().WithStatus()
	w.bq = w.bq.New().FilterDisabled()

	referralPairs, codesForDelete, err := w.findReferralPairs()
	if err != nil {
		w.log.WithError(err).Error("failed select referral pairs")
		panic(err)
	}
	w.log.Infof("Find %d pairs", len(referralPairs))
	refCodesCountMap := w.updateReferredBys(referralPairs)
	if err = w.rq.DeleteByID(codesForDelete...); err != nil {
		w.log.WithError(err).Error("failed delete referrals")
		panic(err)
	}
	w.insertNewReferralCodes(refCodesCountMap)

	return nil
}

func (w *worker) findReferralPairs() (referralPairList []referralPair, codesForDelete []string, err error) {
	guestBalances, err := w.bq.Select()
	if err != nil {
		return referralPairList, nil, errors.Wrap(err, "failed to select guests balances")
	}

	for _, balance := range guestBalances {
		if !balance.ReferredBy.Valid {
			continue
		}
		referral, err := w.rq.Get(balance.ReferredBy.String)
		if err != nil {
			w.log.WithFields(logan.F{
				"error":           err,
				"referred_by":     balance.ReferredBy.String,
				"guest_nullifier": balance.Nullifier,
			}).Error("failed get referrale by referred_by")
			continue
		}
		if referral.Status == data.StatusConsumed && balance.CreatedAt <= int32(time.Now().AddDate(0, 0, -7).Unix()) {
			referralPairList = append(referralPairList, referralPair{
				GuestBalance: &balance,
				Referral:     referral,
			})
			codesForDelete = append(codesForDelete, referral.ID)

			w.log.WithFields(logan.F{
				"code":               balance.ReferredBy.String,
				"guest_nullifier":    balance.Nullifier,
				"referral_nullifier": referral.Nullifier,
			}).Info("new pair for change")
		}
	}

	return referralPairList, codesForDelete, nil
}

func (w *worker) updateReferredBys(pairs []referralPair) map[string]int {
	refCodesCountMap := make(map[string]int)
	for _, referralEntry := range pairs {
		w.bq = w.bq.New()
		if err := w.bq.FilterByNullifier(referralEntry.GuestBalance.Nullifier).Update(map[string]any{"referred_by": w.exc.Code}); err != nil {
			w.log.WithFields(logan.F{
				"error":              err,
				"guest_nullifier":    referralEntry.GuestBalance.Nullifier,
				"referral_nullifier": referralEntry.Referral.Nullifier,
				"new_code":           referralEntry.Referral.ID,
			}).Error("failed change referred_by")
		}
		refCodesCountMap[referralEntry.Referral.Nullifier] = refCodesCountMap[referralEntry.Referral.Nullifier] + 1
	}
	return refCodesCountMap
}

func (w *worker) insertNewReferralCodes(refCodesCountMap map[string]int) {
	for nullifier, refsCount := range refCodesCountMap {
		timeStamp := uint64(time.Now().Unix())
		refToAdd := handlers.PrepareReferralsToAdd(nullifier, uint64(refsCount), timeStamp)
		if err := w.rq.New().Insert(refToAdd...); err != nil {
			w.log.WithFields(logan.F{
				"error":     err,
				"nullifier": nullifier,
				"timestamp": timeStamp,
			}).Error("failed to insert referrals")
			continue
		}

		w.log.WithFields(logan.F{
			"timestamp": timeStamp,
			"nullifier": nullifier,
			"count":     refsCount,
		}).Info("create new referrals")
	}
}
