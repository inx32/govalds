package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"forge.pi.home.arpa/govalds/bot/pkg/config"
	"forge.pi.home.arpa/govalds/bot/pkg/systemd"
	"forge.pi.home.arpa/govalds/bot/pkg/terminal"
	"github.com/rs/zerolog"
)

var logPartsOrder = []string{
	zerolog.LevelFieldName,
	zerolog.CallerFieldName,
	"from",
	zerolog.MessageFieldName,
}

var logTime bool

func init() {
	if !systemd.IsSystemd() {
		logTime = true
		logPartsOrder = append([]string{zerolog.TimestampFieldName}, logPartsOrder...)
	}
}

func createConsoleWriter(writer io.Writer, module string) *zerolog.ConsoleWriter {
	w := &zerolog.ConsoleWriter{
		Out:           writer,
		NoColor:       !terminal.UseColors(nil),
		TimeFormat:    "Jan 02 15:04:05",
		PartsOrder:    logPartsOrder,
		FieldsExclude: []string{"from"},
	}
	ApplyFormatters(w, module)
	return w
}

func CreateDefaultLogger(module string) *zerolog.Logger {
	if os.Getenv("LOG_DISABLE") != "" {
		logger := zerolog.New(io.Discard).Level(zerolog.Disabled)
		return &logger
	}

	var w io.Writer

	if os.Getenv("LOG_FORCE_JSON") != "" {
		w = os.Stdout
	} else {
		stdout := createConsoleWriter(os.Stdout, module)
		stderr := createConsoleWriter(os.Stderr, module)
		w = NewSplit(stdout, stderr, zerolog.WarnLevel)
	}

	logger := zerolog.New(w).Level(zerolog.InfoLevel)
	if logTime {
		logger = logger.With().Timestamp().Logger()
	}

	if levelString := os.Getenv("LOG_LEVEL"); levelString != "" {
		level, err := zerolog.ParseLevel(levelString)
		if err == nil {
			logger = logger.Level(level)
		}
	}

	if logger.GetLevel() == zerolog.TraceLevel {
		logger = logger.With().Caller().Logger()
	}
	return &logger
}

func NewLogger(conf *config.Logging, module string) (*zerolog.Logger, error) {
	if conf == nil {
		return CreateDefaultLogger(module), nil
	}
	if conf.Disabled {
		logger := zerolog.New(io.Discard).Level(zerolog.Disabled)
		return &logger, nil
	}

	var w io.Writer = os.Stdout

	if conf.File != nil {
		path := conf.File.Path.Value()
		if conf.File.Placeholder.Value() != "" {
			timefmt := time.Now().Format(conf.File.Placeholder.Value())
			path = strings.Replace(path, "%s", timefmt, 1)
		}
		var err error
		w, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			return nil, err
		}
	} else {
		w = NewSplit(os.Stdout, os.Stderr, zerolog.WarnLevel)
	}

	enableTimestamp := !systemd.IsSystemd()
	if conf.Timestamp != nil {
		enableTimestamp = *conf.Timestamp
	}

	if !conf.JSON {
		logPartsOrderOverride := []string{
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			"from",
			zerolog.MessageFieldName,
		}
		if enableTimestamp {
			logPartsOrderOverride = append(
				[]string{zerolog.TimestampFieldName},
				logPartsOrderOverride...,
			)
		}

		w = &zerolog.ConsoleWriter{
			Out:           w,
			NoColor:       !terminal.UseColors(conf.Colors),
			TimeFormat:    time.Stamp,
			PartsOrder:    logPartsOrderOverride,
			FieldsExclude: []string{"from"},
		}
		ApplyFormatters(w.(*zerolog.ConsoleWriter), module)
	}

	logger := zerolog.New(w).Level(zerolog.InfoLevel)
	if enableTimestamp {
		logger = logger.With().Timestamp().Logger()
	}

	if conf.Level != nil {
		level, err := zerolog.ParseLevel(conf.Level.Value())
		if err != nil {
			return nil, err
		}
		logger = logger.Level(level)
	}

	if logger.GetLevel() == zerolog.TraceLevel {
		logger = logger.With().Caller().Logger()
	}

	return &logger, nil
}
