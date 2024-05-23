package config

import (
	"fmt"

	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

const proofEventIDValue = "TODO"

func (c *config) Verifier() *zk.Verifier {
	return c.verifier.Do(func() interface{} {
		var cfg struct {
			VerificationKeyPath string `fig:"verification_key_path,required"`
			AllowedAge          int    `fig:"allowed_age,required"`
		}

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "verifier")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out verifier: %w", err))
		}

		v, err := zk.NewPassportVerifier(nil,
			zk.WithVerificationKeyFile(cfg.VerificationKeyPath),
			zk.WithAgeAbove(cfg.AllowedAge),
			zk.WithEventID(proofEventIDValue),
		)

		if err != nil {
			panic(fmt.Errorf("failed to initialize passport verifier: %w", err))
		}

		return v
	}).(*zk.Verifier)
}
