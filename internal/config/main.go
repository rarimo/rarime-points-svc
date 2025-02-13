package config

import (
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/countrier"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/sbtcheck"
	zk "github.com/rarimo/zkverifier-kit"
	"github.com/rarimo/zkverifier-kit/root"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	auth.Auther //nolint:misspell
	evtypes.EventTypeser
	sbtcheck.SbtChecker
	countrier.Countrier
	Broadcasterer

	Levels() Levels
	Verifier() *zk.Verifier
	PointPrice() PointsPrice
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	auth.Auther
	root.VerifierProvider
	evtypes.EventTypeser
	sbtcheck.SbtChecker
	countrier.Countrier
	Broadcasterer

	levels     comfig.Once
	verifier   comfig.Once
	pointPrice comfig.Once
	getter     kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:           getter,
		Databaser:        pgdb.NewDatabaser(getter),
		Listenerer:       comfig.NewListenerer(getter),
		Logger:           comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Auther:           auth.NewAuther(getter), //nolint:misspell
		VerifierProvider: root.NewVerifierProvider(getter, root.PoseidonSMT),
		EventTypeser:     evtypes.NewConfig(getter),
		SbtChecker:       sbtcheck.NewConfig(getter),
		Countrier:        countrier.NewConfig(getter),
		Broadcasterer:    NewBroadcaster(getter),
	}
}
