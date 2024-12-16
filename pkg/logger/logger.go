package logger

import (
	"log"
)

type Logger interface {
	Log(message string)
}

type SimpleLogger struct{}

func New() Logger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Log(message string) {
	log.Println(message)
}
