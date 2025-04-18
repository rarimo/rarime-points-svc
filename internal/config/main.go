package config

import (
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/countrier"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/sbtcheck"
	"github.com/rarimo/saver-grpc-lib/broadcaster"
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
	broadcaster.Broadcasterer
	evtypes.EventTypeser
	sbtcheck.SbtChecker
	countrier.Countrier
	FaceVerifierer

	Levels() Levels
	Verifier() *zk.Verifier
	PointPrice() PointsPrice
	Maintenance() Maintenance
	ExpiredCode() ExpiredCode
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	auth.Auther
	broadcaster.Broadcasterer
	root.VerifierProvider
	evtypes.EventTypeser
	sbtcheck.SbtChecker
	countrier.Countrier
	FaceVerifierer

	levels      comfig.Once
	verifier    comfig.Once
	pointPrice  comfig.Once
	maintenance comfig.Once
	expiredCode comfig.Once

	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:           getter,
		Databaser:        pgdb.NewDatabaser(getter),
		Listenerer:       comfig.NewListenerer(getter),
		Logger:           comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Auther:           auth.NewAuther(getter), //nolint:misspell
		Broadcasterer:    broadcaster.New(getter),
		VerifierProvider: root.NewVerifierProvider(getter, root.PoseidonSMT),
		EventTypeser:     evtypes.NewConfig(getter),
		SbtChecker:       sbtcheck.NewConfig(getter),
		Countrier:        countrier.NewConfig(getter),
		FaceVerifierer:   NewFaceVerifier(getter),
	}
}
