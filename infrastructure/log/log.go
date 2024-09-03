package log

import (
	"context"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const CorrelationID string = "X-Correlation-ID"

func Init(logFormat, logLevel string) {
	switch strings.ToLower(logFormat) {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}
	switch strings.ToLower(logLevel) {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(os.Stdout)
}

func getEntry(ctx context.Context, ctxName string) *log.Entry {
	return log.WithFields(log.Fields{
		"context":       ctxName,
		"correlationId": ctx.Value(CorrelationID),
	})
}

func Info(ctx context.Context, ctxName string, format string, args ...interface{}) {
	getEntry(ctx, ctxName).Infof(format, args...)
}

func Warn(ctx context.Context, ctxName string, format string, args ...interface{}) {
	getEntry(ctx, ctxName).Warnf(format, args...)
}

func Debug(ctx context.Context, ctxName string, format string, args ...interface{}) {
	getEntry(ctx, ctxName).Debugf(format, args...)
}

func Error(ctx context.Context, ctxName string, format string, args ...interface{}) {
	getEntry(ctx, ctxName).Errorf(format, args...)
}
