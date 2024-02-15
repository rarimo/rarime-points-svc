package config

import (
	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/sbtcheck"
	"github.com/rarimo/saver-grpc-lib/broadcaster"
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
	broadcaster.Broadcasterer
	evtypes.EventTypeser
	sbtcheck.SbtChecker

	PointPrice() uint64
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	auth.Auther
	broadcaster.Broadcasterer
	evtypes.EventTypeser
	sbtcheck.SbtChecker

	pointPrice comfig.Once
	getter     kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:        getter,
		Databaser:     pgdb.NewDatabaser(getter),
		Copuser:       copus.NewCopuser(getter),
		Listenerer:    comfig.NewListenerer(getter),
		Logger:        comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Auther:        auth.NewAuther(getter), //nolint:misspell
		Broadcasterer: broadcaster.New(getter),
		EventTypeser:  evtypes.NewConfig(getter),
		SbtChecker:    sbtcheck.NewConfig(getter),
	}
}
