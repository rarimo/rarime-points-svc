package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/points-svc/internal/data/pg"
	"github.com/rarimo/points-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxEventsQ(pg.NewEvents(s.cfg.DB())),
			handlers.CtxBalancesQ(pg.NewBalances(s.cfg.DB())),
		),
	)
	r.Route("/integrations/points-svc", func(r chi.Router) {
		r.Get("/balance", handlers.GetBalance)
		r.Get("/leaderboard", handlers.Leaderboard)
		r.Get("/events", handlers.ListEvents)
		r.Put("/events/{id}", handlers.ClaimEvent)
	})

	return r
}
