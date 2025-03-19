package referrals

import (
	"time"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type worker struct {
	name string
	rq   data.ReferralsQ
	bq   data.BalancesQ
	log  *logan.Entry
	exc  config.ExterminatedCode
}

type referralPair struct {
	Referred *data.Balance
	Referrer *data.Referral
}

func newWorker(cfg config.Config, workerName string) *worker {
	return &worker{
		name: workerName,
		rq:   pg.NewReferrals(cfg.DB().Clone()),
		bq:   pg.NewBalances(cfg.DB().Clone()),
		log:  cfg.Log().WithField("who", workerName),
		exc:  cfg.ExterminatedCode(),
	}
}

func (w *worker) job() error {
	w.rq = w.rq.New().WithStatus()
	w.bq = w.bq.New().FilterDisabled()

	referralPairs, refToUsageLeftMap, err := w.findReferralPairs()
	if err != nil {
		w.log.WithFields(logan.F{
			"error": err,
		}).Errorf("failed select referral pairs")
		panic(err)
	}
	w.log.Infof("Find %d pairs", len(referralPairs))
	w.updateReferreds(referralPairs)
	w.updateUsageLeft(refToUsageLeftMap)

	return nil
}

func (w *worker) findReferralPairs() (referralPairList []referralPair, refToUsageLeftMap map[string]int32, err error) {
	refToUsageLeftMap = make(map[string]int32)

	referredBalances, err := w.bq.Select()
	if err != nil {
		return referralPairList, nil, errors.Wrap(err, "failed to select referred balances")
	}

	for _, referred := range referredBalances {
		if !referred.ReferredBy.Valid {
			continue
		}
		referee, err := w.rq.Get(referred.ReferredBy.String)
		if err != nil {
			w.log.WithFields(logan.F{
				"error":              err,
				"referred_by":        referred.ReferredBy.String,
				"referred_nullifier": referred.Nullifier,
			}).Errorf("failed get referrale by referred_by")
			continue
		}
		if referee.Status == data.StatusConsumed && referred.CreatedAt <= int32(time.Now().AddDate(0, 0, -7).Unix()) {
			referralPairList = append(referralPairList, referralPair{
				Referred: &referred,
				Referrer: referee,
			})

			usageLeft, ok := refToUsageLeftMap[referee.Nullifier]
			if !ok {
				usageLeft = referee.UsageLeft
			}
			refToUsageLeftMap[referee.Nullifier] = usageLeft + 1
			w.log.WithFields(logan.F{
				"code":               referred.ReferredBy.String,
				"referred_nullifier": referred.Nullifier,
				"referee_nullifier":  referee.Nullifier,
			}).Info("new pair for change")
		}
	}

	return referralPairList, refToUsageLeftMap, nil
}

func (w *worker) updateReferreds(pairs []referralPair) {
	for _, referralEntry := range pairs {
		w.bq = w.bq.New()
		if err := w.bq.FilterByNullifier(referralEntry.Referred.Nullifier).Update(map[string]any{"referred_by": w.exc.Code}); err != nil {
			w.log.WithFields(logan.F{
				"error":              err,
				"referred_nullifier": referralEntry.Referred.Nullifier,
				"referee_nullifier":  referralEntry.Referrer.Nullifier,
				"new_code":           referralEntry.Referrer.ID,
			}).Errorf("failed change referred_by")
		}
	}
}

func (w *worker) updateUsageLeft(refMap map[string]int32) {
	for nullifier, usageLeft := range refMap {
		w.rq = w.rq.New()
		_, err := w.rq.FilterByNullifier(nullifier).Update(int(usageLeft))
		if err != nil {
			w.log.WithFields(logan.F{
				"error":      err,
				"nullifier":  nullifier,
				"usage_left": usageLeft,
			}).Errorf("failed change usage left")
		}
	}
}
