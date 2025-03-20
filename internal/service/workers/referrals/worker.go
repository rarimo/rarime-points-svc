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

	referralPairs, err := w.findReferralPairs()
	if err != nil {
		w.log.WithError(err).Error("failed select referral pairs")
		return errors.Wrap(err, "failed select referral pairs")
	}
	w.log.Infof("Find %d pairs", len(referralPairs))

	refCodesMap := make(map[string][]string)
	for _, referralEntry := range referralPairs {
		if err := w.updateReferredBys(referralEntry.GuestBalance.Nullifier); err != nil {
			w.log.WithFields(logan.F{
				"error":              err,
				"guest_nullifier":    referralEntry.GuestBalance.Nullifier,
				"referral_nullifier": referralEntry.Referral.Nullifier,
				"referred_by":        referralEntry.GuestBalance.ReferredBy.String,
			}).Error("failed change referred_by")
			continue
		}
		refCodesMap[referralEntry.Referral.Nullifier] = append(refCodesMap[referralEntry.Referral.Nullifier], referralEntry.GuestBalance.ReferredBy.String)
	}

	var codesForDelete []string
	for nullifier, codes := range refCodesMap {
		timeStamp, err := w.insertNewReferralCodes(nullifier, len(codes))
		if err != nil {
			w.log.WithFields(logan.F{
				"error":     err,
				"nullifier": nullifier,
				"timestamp": timeStamp,
			}).Error("failed to insert referrals")
			continue
		}
		codesForDelete = append(codesForDelete, codes...)
	}

	if err = w.rq.DeleteByID(codesForDelete...); err != nil {
		w.log.WithError(err).Error("failed delete referrals")
		return errors.Wrap(err, "failed delete referrals")
	}

	return nil
}

func (w *worker) findReferralPairs() (referralPairList []referralPair, err error) {
	guestBalances, err := w.bq.Select()
	if err != nil {
		return referralPairList, errors.Wrap(err, "failed to select guest balances")
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

		if referral.Status == data.StatusConsumed && balance.CreatedAt <= int32(time.Now().Add(-w.exc.CodeLifetime*time.Second).Unix()) {
			referralPairList = append(referralPairList, referralPair{
				GuestBalance: &balance,
				Referral:     referral,
			})

			w.log.WithFields(logan.F{
				"code":               balance.ReferredBy.String,
				"guest_nullifier":    balance.Nullifier,
				"referral_nullifier": referral.Nullifier,
			}).Info("new pair for change")
		}
	}

	return referralPairList, nil
}

func (w *worker) updateReferredBys(guestNullifier string) error {
	return w.bq.Transaction(func() error {
		if err := w.bq.New().FilterByNullifier(guestNullifier).Update(map[string]any{"referred_by": w.exc.Code}); err != nil {
			return err
		}

		return nil
	})
}

func (w *worker) insertNewReferralCodes(nullifier string, refsCount int) (uint64, error) {
	timeStamp := uint64(time.Now().Unix())
	return timeStamp, w.rq.Transaction(func() error {
		refToAdd := handlers.PrepareReferralsToAdd(nullifier, uint64(refsCount), timeStamp)
		if err := w.rq.New().Insert(refToAdd...); err != nil {
			return errors.Wrap(err, "failed to insert referrals")
		}

		return nil
	})
}
