package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type LogFormat string

func (logFormat LogFormat) String() string {
	return string(logFormat)
}

const (
	JSONFormat    LogFormat = "json"
	ConsoleFormat LogFormat = "console"
)

// NewLogger returns a zerolog.Logger
func NewLogger(logLevel string, logFormat LogFormat, opts ...Option) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		return &zerolog.Logger{}, err
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log := zerolog.Logger{}

	switch logFormat {
	case JSONFormat:
		log = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	case ConsoleFormat:
		log = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger().Output(zerolog.ConsoleWriter{TimeFormat: time.RFC3339Nano, Out: os.Stdout})
	default:
		err = fmt.Errorf("not support log format [%s]", logFormat)
	}

	for _, opt := range opts {
		log = opt.Apply(log)
	}

	return &log, err
}
