package cron

import (
	"gitlab.com/distributed_lab/logan/v3"
)

type logger struct {
	in *logan.Entry
}

func newLogger(in *logan.Entry) *logger {
	return &logger{in: in.WithField("who", "gocron-scheduler")}
}

func (l *logger) Debug(msg string, args ...any) {
	l.in.Debug(l.toLoganArgs(msg, args)...)
}

func (l *logger) Info(msg string, args ...any) {
	l.in.Info(l.toLoganArgs(msg, args)...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.in.Warn(l.toLoganArgs(msg, args)...)
}

func (l *logger) Error(msg string, args ...any) {
	l.in.Error(l.toLoganArgs(msg, args)...)
}

func (l *logger) toLoganArgs(msg string, args []any) []any {
	return append([]any{msg}, args...)
}
