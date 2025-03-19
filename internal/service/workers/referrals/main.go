package referrals

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/cron"
)

const workerName = "refChecker"

func Run(ctx context.Context, cfg config.Config, sig chan struct{}) {
	cron.Init(cfg.Log())

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Errorf("%v: failed to create new scheduler: %w", workerName, err))
	}

	worker := newWorker(cfg, workerName)
	task := gocron.NewTask(worker.job)
	jobType := gocron.CronJob(
		"0 6 */7 * *",
		false,
	)

	_, err = s.NewJob(jobType, task)
	if err != nil {
		panic(fmt.Errorf("%v: failed to create new job: %w", workerName, err))
	}
	sig <- struct{}{}

	s.Start()

	select {
	case <-ctx.Done():
		if err := s.Shutdown(); err != nil {
			panic(fmt.Errorf("%v: failed to shutdown scheduler: %w", workerName, err))
		}
	}
}
