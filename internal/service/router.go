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
			handlers.CtxWithdrawalsQ(pg.NewWithdrawals(s.cfg.DB())),
			handlers.CtxEventTypes(s.cfg.EventTypes()),
			handlers.CtxBroadcaster(s.cfg.Broadcaster()),
			handlers.CtxPointPrice(s.cfg.PointPrice()),
		),
	)
	r.Route("/integrations/rarime-points-svc/v1", func(r chi.Router) {
		r.Route("/balances", func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(s.cfg.Auth(), s.log))
			r.Post("/", handlers.CreateBalance)
			r.Route("/{did}", func(r chi.Router) {
				r.Get("/", handlers.GetBalance)
				r.Get("/withdrawals", handlers.ListWithdrawals)
				r.Post("/withdrawals", handlers.Withdraw)
			})
		})
		r.Route("/events", func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(s.cfg.Auth(), s.log))
			r.Get("/", handlers.ListEvents)
			r.Patch("/{id}", handlers.ClaimEvent)
		})
		r.Get("/balances", handlers.Leaderboard)
	})

	return r
}
