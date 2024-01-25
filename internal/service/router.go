package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
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
			handlers.CtxEventsQ(pg.NewEvents(s.cfg.DB())),
			handlers.CtxBalancesQ(pg.NewBalances(s.cfg.DB())),
			handlers.CtxEventTypes(s.cfg.EventTypes()),
		),
	)
	r.Route("/integrations/rarime-points-svc/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(s.cfg.Auth(), s.log))
			r.Get("/balances/{did}", handlers.GetBalance)
			r.Post("/balances", handlers.CreateBalance)
			r.Get("/events", handlers.ListEvents)
			r.Patch("/events/{id}", handlers.ClaimEvent)
		})
		r.Get("/balances", handlers.Leaderboard)
	})

	return r
}
