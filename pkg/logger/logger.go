package logger

// import (
// 	"log"
// )

// type Logger interface {
// 	Log(message string)
// }

// type SimpleLogger struct{}

// func New() Logger {
// 	return &SimpleLogger{}
// }

// func (l *SimpleLogger) Log(message string) {
// 	log.Println(message)
// }

import (
	"log"

	"github.com/getsentry/sentry-go"
)

type Logger interface {
	Info(message string, fields map[string]interface{})
	Error(err error, message string, fields map[string]interface{})
}

type logger struct {
	enabled         bool
	localLogEnabled bool
}

func New(SentryEnabled bool, localLogEnabled bool) Logger {
	return &logger{
		enabled:         SentryEnabled,
		localLogEnabled: localLogEnabled,
	}
}

func (l *logger) Info(message string, fields map[string]interface{}) {

	if l.localLogEnabled {
		log.Println(message)
	}

	if !l.enabled {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelInfo)
		if fields != nil {
			scope.SetExtras(fields)
		}
		sentry.CaptureMessage(message)
	})
}

func (l *logger) Error(err error, message string, fields map[string]interface{}) {
	log.Println(err.Error())

	if !l.enabled {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.SetExtra("info", message)
		if fields != nil {
			scope.SetExtras(fields)
		}
		sentry.CaptureException(err)
	})
}
