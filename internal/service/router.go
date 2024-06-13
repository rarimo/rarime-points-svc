package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, cfg config.Config) {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxEventTypes(cfg.EventTypes()),
			handlers.CtxBroadcaster(cfg.Broadcaster()),
			handlers.CtxPointPrice(cfg.PointPrice()),
			handlers.CtxLevels(cfg.Levels()),
			handlers.CtxVerifier(cfg.Verifier()),
		),
		handlers.DBCloneMiddleware(cfg.DB()),
	)
	r.Route("/integrations/rarime-points-svc/v1", func(r chi.Router) {
		r.Route("/public", func(r chi.Router) {
			r.Route("/balances/{nullifier}", func(r chi.Router) {
				r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
				r.Get("/", handlers.GetBalance)
				r.Post("/verifypassport", handlers.VerifyPassport)
				r.Get("/withdrawals", handlers.ListWithdrawals)
				r.Post("/withdrawals", handlers.Withdraw)
			})
			r.Route("/events", func(r chi.Router) {
				r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
				r.Get("/", handlers.ListEvents)
				r.Get("/{id}", handlers.GetEvent)
				r.Patch("/{id}", handlers.ClaimEvent)
			})
			r.Route("/balances", func(r chi.Router) {
				r.Get("/", handlers.Leaderboard)
				r.Post("/", handlers.CreateBalance)
			})
			r.Get("/point_price", handlers.GetPointPrice)
			r.Get("/countries_config", handlers.GetCountriesConfig)
			r.Get("/events_config", handlers.GetEventsConfig)
		})
		// must be accessible only within the cluster
		r.Route("/private", func(r chi.Router) {
			r.Patch("/events", handlers.FulfillEvent)
			r.Post("/referrals", handlers.EditReferrals)
		})
	})

	cfg.Log().Info("Service started")
	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
