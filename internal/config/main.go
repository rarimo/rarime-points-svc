package config

import (
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/sbtcheck"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	auth.Auther
	evtypes.EventTypeser
	sbtcheck.SbtChecker
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	auth.Auther
	evtypes.EventTypeser
	sbtcheck.SbtChecker

	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:       getter,
		Databaser:    pgdb.NewDatabaser(getter),
		Copuser:      copus.NewCopuser(getter),
		Listenerer:   comfig.NewListenerer(getter),
		Logger:       comfig.NewLogger(getter, comfig.LoggerOpts{}),
		EventTypeser: evtypes.NewConfig(getter),
		Auther:       auth.NewAuther(getter), //nolint:misspell
		SbtChecker:   sbtcheck.NewConfig(getter),
	}
}
