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
	WithoutPassportScanBalance *data.Balance
	Referral                   *data.Referral
}

func newWorker(cfg config.Config, workerName string) *worker {
	db := cfg.DB().Clone()
	return &worker{
		name: workerName,
		rq:   pg.NewReferrals(db),
		bq:   pg.NewBalances(db),
		log:  cfg.Log().WithField("who", workerName),
		exc:  cfg.ExpiredCode(),
	}
}

func (w *worker) job() error {
	referralPairs, err := w.findReferralPairs()
	if err != nil {
		w.log.WithError(err).Error("failed select referral pairs")
		return errors.Wrap(err, "failed select referral pairs")
	}
	w.log.Infof("Find %d pairs", len(referralPairs))

	for _, referralEntry := range referralPairs {
		if err := w.updateReferralPair(referralEntry); err != nil {
			w.log.WithFields(logan.F{
				"error":                           err,
				"without_passport_scan_nullifier": referralEntry.WithoutPassportScanBalance.Nullifier,
				"referral_nullifier":              referralEntry.Referral.Nullifier,
				"referred_by":                     referralEntry.WithoutPassportScanBalance.ReferredBy.String,
			}).Error("failed update pairs")
		}
	}

	return nil
}

func (w *worker) findReferralPairs() (referralPairList []referralPair, err error) {
	w.rq = w.rq.New().WithStatus().WithoutExpiredStatus()
	w.bq = w.bq.New().FilterDisabled().FilterByIsPassportProven(false).FilterByCreatedAtBefore(int(time.Now().Add(-w.exc.CodeLifetime * time.Second).Unix()))

	withoutPassportScanBalances, err := w.bq.Select()
	if err != nil {
		return referralPairList, errors.Wrap(err, "failed to select without passport scan balances")
	}

	for _, balance := range withoutPassportScanBalances {
		referral, err := w.rq.Get(balance.ReferredBy.String)
		if err != nil {
			w.log.WithFields(logan.F{
				"error":                           err,
				"referred_by":                     balance.ReferredBy.String,
				"without_passport_scan_nullifier": balance.Nullifier,
			}).Error("failed get referrale by referred_by")
			continue
		}
		if referral != nil && referral.Status == data.StatusConsumed {
			referralPairList = append(referralPairList, referralPair{
				WithoutPassportScanBalance: &balance,
				Referral:                   referral,
			})

			w.log.WithFields(logan.F{
				"code":                            balance.ReferredBy.String,
				"without_passport_scan_nullifier": balance.Nullifier,
				"referral_nullifier":              referral.Nullifier,
			}).Info("new pair for change")
		}
	}

	return referralPairList, nil
}

func (w *worker) updateReferralPair(referralEntry referralPair) error {
	w.rq = w.rq.New()
	w.bq = w.bq.New()
	return w.bq.Transaction(func() error {
		if err := w.bq.FilterByNullifier(referralEntry.WithoutPassportScanBalance.Nullifier).Update(map[string]any{"referred_by": w.exc.Code}); err != nil {
			return errors.Wrap(err, "failed change referred_by")
		}

		count, err := w.rq.FilterByNullifier(referralEntry.Referral.Nullifier).Count()
		if err != nil {
			return errors.Wrap(err, "failed to get referral count")
		}

		refToAdd := handlers.PrepareReferralsToAdd(referralEntry.Referral.Nullifier, 1, count)
		if err := w.rq.Insert(refToAdd...); err != nil {
			return errors.Wrap(err, "failed to insert referrals")
		}

		if _, err := w.rq.FilterByID(referralEntry.Referral.ID).Update(-1); err != nil {
			return errors.Wrap(err, "failed update referrals")
		}

		return nil
	})
}
