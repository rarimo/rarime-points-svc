package connector

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const defaultTimeout = 10 * time.Second

type Pointer interface {
	Points() *Client
}

type points struct {
	once   comfig.Once
	getter kv.Getter
}

func NewPointer(getter kv.Getter) Pointer {
	return &points{getter: getter}
}

func (p *points) Points() *Client {
	return p.once.Do(func() any {
		var cfg struct {
			Addr           *url.URL      `fig:"addr,required"`
			RequestTimeout time.Duration `fig:"request_timeout"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(p.getter, "points")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out points: %s", err))
		}

		if cfg.RequestTimeout == 0 {
			cfg.RequestTimeout = defaultTimeout
		}

		return NewClient(&client{
			base: cfg.Addr,
			http: &http.Client{Timeout: cfg.RequestTimeout},
		})
	}).(*Client)
}
