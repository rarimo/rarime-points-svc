package reopener

import (
	"reflect"

	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"gitlab.com/distributed_lab/logan/v3"
)

// Reflect was better option than re-implementing methods, because gocron library
// is adjusted to its own logger, despite supporting custom ones.
func getLogLevel(log *logan.Entry) gocron.LogLevel {
	var (
		val   = reflect.ValueOf(log).Elem()
		entry = val.FieldByName("entry")
		level = entry.FieldByName("Level").Interface().(logrus.Level)
	)

	switch level {
	case logrus.DebugLevel:
		return gocron.LogLevelDebug
	case logrus.InfoLevel:
		return gocron.LogLevelInfo
	case logrus.WarnLevel:
		return gocron.LogLevelWarn
	default:
		return gocron.LogLevelError
	}
}
