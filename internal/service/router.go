package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-points-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)
	r.Route("/integrations/rarime-points-svc", func(r chi.Router) {
		// configure endpoints here
	})

	return r
}
