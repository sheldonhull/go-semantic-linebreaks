// Originally source from https://gist.github.com/panta/2530672ca641d953ae452ecb5ef79d7d
// Modified logging package to include more context and Level to be configured
package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

// Configuration for logging.
type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
	// Level is either info or debug output level
	Level string
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func Configure(config Config) *Logger {
	var writers []io.Writer

	// Setting as switch so I can expand later to other layers like trace etc if required
	switch config.Level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	if config.ConsoleLoggingEnabled {
		consoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		consoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s |", i)
		}
		consoleWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		consoleWriter.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}
		writers = append(writers, consoleWriter)
	}

	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	// mw := io.MultiWriter(writers...)
	multi := zerolog.MultiLevelWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// logger := zerolog.New(mw).With().Timestamp().Logger()
	logger := zerolog.New(multi).With().Timestamp().Caller().Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Str("level", config.Level).
		Msg("logging configured")

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0o755); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")

		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
