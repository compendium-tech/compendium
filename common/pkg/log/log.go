package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

var LoggerKey loggerKey = loggerKey{}

// Logger set by logger middleware, stored in a context.
func L(ctx context.Context) *logrus.Entry {
	if logger, ok := ctx.Value(LoggerKey).(*logrus.Entry); ok {
		return logger
	}

	return logrus.NewEntry(logrus.StandardLogger())
}
