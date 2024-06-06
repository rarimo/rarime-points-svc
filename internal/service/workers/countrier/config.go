package countrier

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

// Config exists only to Run countrier with provided country list
type Config struct {
	m map[string]country
}

type country struct {
	Code              string `fig:"code,required"`
	ReserveLimit      int64  `fig:"reserve_limit,required"`
	ReserveAllowed    bool   `fig:"reserve_allowed,required"`
	WithdrawalAllowed bool   `fig:"withdrawal_allowed,required"`
}

type Countrier interface {
	Countries() Config
}

type config struct {
	once   comfig.Once
	getter kv.Getter
}

func NewConfig(getter kv.Getter) Countrier {
	return &config{getter: getter}
}

func (c *config) Countries() Config {
	return c.once.Do(func() interface{} {
		var cfg struct {
			Countries []country `fig:"countries,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "countries")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out countries: %s", err))
		}

		countries := make(map[string]country, len(cfg.Countries))
		for _, cou := range cfg.Countries {
			err = validation.Errors{
				"code": validation.Validate(
					cou.Code,
					validation.Required,
					validation.When(cou.Code != data.DefaultCountryCode, is.CountryCode3),
				),
				"reserve_limit": validation.Validate(cou.ReserveLimit, validation.Min(0)),
			}.Filter()

			if err != nil {
				panic(fmt.Errorf("invalid country %s: %w", cou.Code, err))
			}

			countries[cou.Code] = cou
		}

		if _, ok := countries[data.DefaultCountryCode]; !ok {
			panic(fmt.Errorf("default country with code %s is not set", data.DefaultCountryCode))
		}

		return Config{m: countries}
	}).(Config)
}
