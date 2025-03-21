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
	return &worker{
		name: workerName,
		rq:   pg.NewReferrals(cfg.DB().Clone()),
		bq:   pg.NewBalances(cfg.DB().Clone()),
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
	timeStamp := uint64(time.Now().Unix())

	for _, referralEntry := range referralPairs {
		timeStamp += 1
		if err := w.updateReferralPair(referralEntry, timeStamp); err != nil {
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
	w.rq = w.rq.New().WithStatus()
	w.bq = w.bq.New().FilterDisabled().FilterByIsPassportProven()

	withoutPassportScanBalances, err := w.bq.Select()
	if err != nil {
		return referralPairList, errors.Wrap(err, "failed to select without passport scan balances")
	}

	for _, balance := range withoutPassportScanBalances {
		if !balance.ReferredBy.Valid {
			continue
		}
		referral, err := w.rq.Get(balance.ReferredBy.String)
		if err != nil {
			w.log.WithFields(logan.F{
				"error":                           err,
				"referred_by":                     balance.ReferredBy.String,
				"without_passport_scan_nullifier": balance.Nullifier,
			}).Error("failed get referrale by referred_by")
			continue
		}

		if referral.Status == data.StatusConsumed && balance.CreatedAt <= int32(time.Now().Add(-w.exc.CodeLifetime*time.Second).Unix()) {
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

func (w *worker) updateReferralPair(referralEntry referralPair, timeStamp uint64) error {
	if err := w.bq.New().FilterByNullifier(referralEntry.WithoutPassportScanBalance.Nullifier).Update(map[string]any{"referred_by": w.exc.Code}); err != nil {
		return errors.Wrap(err, "failed change referred_by")
	}
	return w.rq.Transaction(func() error {
		refToAdd := handlers.PrepareReferralsToAdd(referralEntry.Referral.Nullifier, 1, timeStamp)
		if err := w.rq.Insert(refToAdd...); err != nil {
			return errors.Wrap(err, "failed to insert referrals")
		}

		if err := w.rq.DeleteByID(referralEntry.Referral.ID); err != nil {
			return errors.Wrap(err, "failed delete referrals")
		}

		return nil
	})
}
