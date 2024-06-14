package cli

import (
	"context"
	"sync"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/service"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/countrier"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/expirywatch"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/nooneisforgotten"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/reopener"
)

// runServices manages service's dependencies and runs them in the correct order
func runServices(ctx context.Context, cfg config.Config, wg *sync.WaitGroup) {
	// signals indicate the finished initialization of each worker
	var (
		reopenerSig         = make(chan struct{})
		countrierSig        = make(chan struct{})
		expiryWatchSig      = make(chan struct{})
		noOneIsForgottenSig = make(chan struct{})
	)

	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	// these services can safely run in parallel and don't have dependencies
	run(func() { reopener.Run(ctx, cfg, reopenerSig) })
	run(func() { expirywatch.Run(ctx, cfg, expiryWatchSig) })
	run(func() { countrier.Run(cfg, countrierSig) })

	// these two depend on reopener, because events must be opened before they can be
	// fulfilled, then both services do not overlap each other
	<-reopenerSig
	run(func() { nooneisforgotten.Run(cfg, noOneIsForgottenSig) })
	//run(func() { sbtcheck.Run(ctx, cfg) }) // see deprecation notice

	// service depends on all the workers for good UX, except sbtcheck, as it has
	// long catchup period and users are fine to wait
	<-countrierSig
	<-expiryWatchSig
	<-noOneIsForgottenSig
	run(func() { service.Run(ctx, cfg) })
}
