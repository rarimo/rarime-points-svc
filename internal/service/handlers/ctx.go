package handlers

import (
	"context"
	"net/http"

	"github.com/rarimo/auth-svc/resources"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/saver-grpc-lib/broadcaster"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	eventsQCtxKey
	balancesQCtxKey
	withdrawalsQCtxKey
	eventTypesCtxKey
	userClaimsCtxKey
	broadcasterCtxKey
	pointPriceCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxEventsQ(q data.EventsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, eventsQCtxKey, q)
	}
}

func EventsQ(r *http.Request) data.EventsQ {
	return r.Context().Value(eventsQCtxKey).(data.EventsQ).New()
}

func CtxBalancesQ(q data.BalancesQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, balancesQCtxKey, q)
	}
}

func BalancesQ(r *http.Request) data.BalancesQ {
	return r.Context().Value(balancesQCtxKey).(data.BalancesQ).New()
}

func CtxWithdrawalsQ(q data.WithdrawalsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, withdrawalsQCtxKey, q)
	}
}

func WithdrawalsQ(r *http.Request) data.WithdrawalsQ {
	return r.Context().Value(withdrawalsQCtxKey).(data.WithdrawalsQ).New()
}

func CtxEventTypes(types evtypes.Types) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, eventTypesCtxKey, types)
	}
}

func EventTypes(r *http.Request) evtypes.Types {
	return r.Context().Value(eventTypesCtxKey).(evtypes.Types)
}

func CtxUserClaims(claim []resources.Claim) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userClaimsCtxKey, claim)
	}
}

func UserClaims(r *http.Request) []resources.Claim {
	return r.Context().Value(userClaimsCtxKey).([]resources.Claim)
}

func CtxBroadcaster(broadcaster broadcaster.Broadcaster) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, broadcasterCtxKey, broadcaster)
	}
}

func Broadcaster(r *http.Request) broadcaster.Broadcaster {
	return r.Context().Value(broadcasterCtxKey).(broadcaster.Broadcaster)
}

func CtxPointPrice(price uint64) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pointPriceCtxKey, price)
	}
}

func PointPrice(r *http.Request) uint64 {
	return r.Context().Value(pointPriceCtxKey).(uint64)
}
