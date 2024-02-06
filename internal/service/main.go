package service

import (
	"context"
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
)

type service struct {
	log   *logan.Entry
	copus types.Copus
	cfg   config.Config
}

func (s *service) run(ctx context.Context) error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return fmt.Errorf("cop failed: %w", err)
	}

	ape.Serve(ctx, r, s.cfg, ape.ServeOpts{})
	return nil
}

func newService(cfg config.Config) *service {
	return &service{
		log:   cfg.Log(),
		copus: cfg.Copus(),
		cfg:   cfg,
	}
}

func Run(ctx context.Context, cfg config.Config) {
	if err := newService(cfg).run(ctx); err != nil {
		panic(err)
	}
}
