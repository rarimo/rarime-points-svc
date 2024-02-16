package cron

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"gitlab.com/distributed_lab/logan/v3"
)

// how many services should call Start to do the actual start
const countOfServices = 2

var sin = struct {
	scheduler gocron.Scheduler
	mu        sync.Mutex
	init      bool
	counter   int
	log       *logan.Entry
}{}

func Init(log *logan.Entry) {
	sin.mu.Lock()
	defer sin.mu.Unlock()

	if sin.init {
		return
	}

	var err error
	sin.scheduler, err = gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLogger(newLogger(log)),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize scheduler: %w", err))
	}
	sin.log = log.WithField("who", "cron-scheduler")
	sin.init = true
}

func NewJob(jobDef gocron.JobDefinition, task gocron.Task, opts ...gocron.JobOption) (gocron.Job, error) {
	sin.mu.Lock()
	defer sin.mu.Unlock()

	if !sin.init {
		panic("scheduler not initialized")
	}

	job, err := sin.scheduler.NewJob(jobDef, task, opts...)
	if err != nil {
		return nil, fmt.Errorf("add new job: %w", err)
	}

	return job, nil
}

// Start must be called several times by asynchronous services to start the scheduler only once.
func Start(ctx context.Context) {
	sin.mu.Lock()
	defer sin.mu.Unlock()

	if !sin.init {
		panic("scheduler not initialized")
	}

	sin.counter++
	if sin.counter < countOfServices {
		sin.log.Debugf("Waiting for %d services to start", countOfServices-sin.counter)
		return
	}

	sin.scheduler.Start()
	logJobs()
	<-ctx.Done() // all cron jobs are shut down when ctx is canceled

	if err := sin.scheduler.Shutdown(); err != nil {
		sin.log.WithError(err).Error("Scheduler shutdown failed")
		return
	}
	sin.log.Info("Scheduler shutdown succeeded")
}

func logJobs() {
	// mutex lock must be already acquired in the caller
	jobs := sin.scheduler.Jobs()
	logged := make([]string, 0, len(jobs))

	for _, job := range jobs {
		nextRun, err := job.NextRun()
		if err != nil {
			sin.log.WithError(err).
				Errorf("Failed to get next run time: name=%s uuid=%s", job.Name(), job.ID())
			continue
		}

		logged = append(logged, fmt.Sprintf("(name=%s next_run=%s)",
			job.Name(), nextRun.Format(time.RFC3339)))
	}

	if len(logged) == 0 {
		sin.log.Warn("No jobs successfully scheduled")
	}

	sin.log.Infof("Scheduled jobs: %s", strings.Join(logged, ", "))
}
