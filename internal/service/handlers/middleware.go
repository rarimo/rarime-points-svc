package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AuthMiddleware(auth *auth.Client, log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// in current implementation any error is internal, despite the status
			claims, _, err := auth.ValidateJWT(r.Header)
			if err != nil {
				log.WithError(err).Error("failed to execute auth validate request")
				ape.Render(w, problems.InternalError())
				return
			}

			if len(claims) == 0 {
				log.Debug("No claims returned for user")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}
			if len(claims) > 1 {
				log.Errorf("Expected 1 claim to get user DID from, got %d claims", len(claims))
				ape.RenderErr(w, problems.InternalError())
				return
			}

			claim := claims[0]
			if claim.User == "" {
				log.Debug("No user DID found in claim")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			ctx := CtxUserDID(claim.User)(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
