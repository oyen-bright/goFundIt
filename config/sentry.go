package config

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func InitSentry(environment string, isDebug bool) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://589962e0c2695cb8afe6a0743ce33f0b@o4508531503923200.ingest.de.sentry.io/4508531543638096",
		Environment:      environment,
		Debug:            isDebug,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		AttachStacktrace: true,
	})

	if err != nil {
		log.Fatalf("Sentry initialization failed: %v", err)
		return err
	}

	return nil
}
