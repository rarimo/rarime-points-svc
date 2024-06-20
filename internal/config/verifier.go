package config

import (
	"fmt"

	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	proofEventIDValue  = "211985299740800702300256033401632392934377086534111448880928528431996790315"
	proofSelectorValue = "23073"
	maxIdentityCount   = 1
)

func (c *config) Verifier() *zk.Verifier {
	return c.verifier.Do(func() interface{} {
		var cfg struct {
			AllowedAge               int    `fig:"allowed_age,required"`
			VerificationKeyPath      string `fig:"verification_key_path,required"`
			AllowedIdentityTimestamp int64  `fig:"allowed_identity_timestamp,required"`
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
			zk.WithIdentityVerifier(c.ProvideVerifier()),
			zk.WithProofSelectorValue(proofSelectorValue),
			zk.WithEventID(proofEventIDValue),
			zk.WithIdentitiesCounter(maxIdentityCount),
			zk.WithIdentitiesCreationTimestampLimit(cfg.AllowedIdentityTimestamp),
		)

		if err != nil {
			panic(fmt.Errorf("failed to initialize passport verifier: %w", err))
		}

		return v
	}).(*zk.Verifier)
}
