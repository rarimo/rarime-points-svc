package handlers

import (
	"context"
	"net/http"

	"github.com/rarimo/points-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	eventsQCtxKey
	balancesQCtxKey
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
