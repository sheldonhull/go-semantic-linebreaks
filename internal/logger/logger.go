// Originally source from https://gist.github.com/panta/2530672ca641d953ae452ecb5ef79d7d
// Modified logging package to include more context and Level to be configured
// duplicating the method: https://github.com/learning-cloud-native-go/myapp/blob/step-6/util/logger/logger.go
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

const DirectoryPermissions = 0755

type Logger struct {
	*zerolog.Logger
}

var Log *Logger //nolint: gochecknoglobals
// Configuration for logging.
type Config struct {
	// Set Enable to false to disable logging
	Enable bool
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

// init logger with disabled state on import to avoid any nil pointers if unused.
func init() { //nolint:gochecknoinits
	Log = InitLogger(Config{
		Enable:                false,
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      false,
		FileLoggingEnabled:    false,
		Directory:             "",
		Filename:              "",
		MaxSize:               10, //nolint: gomnd
		MaxBackups:            10, //nolint: gomnd
		MaxAge:                1,
		Level:                 "info",
	})
}

// InitLogger sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
// To disable usage for unit tests, just pass in the Config propery Disable: true.
func InitLogger(config Config) *Logger {
	var writers []io.Writer

	if !config.Enable {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		// Setting as switch so I can expand later to other layers like trace etc if required
		switch config.Level {
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
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

	return &Logger{&logger}
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, DirectoryPermissions); err != nil {
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

// // Output duplicates the global logger and sets w as its output.
// func (l *Logger) Output(w io.Writer) zerolog.Logger {
// 	return l.logger.Output(w)
// }

// // With creates a child logger with the field added to its context.
// func (l *Logger) With() zerolog.Context {
// 	return l.logger.With()
// }

// // Level creates a child logger with the minimum accepted level set to level.
// func (l *Logger) Level(level zerolog.Level) zerolog.Logger {
// 	return l.logger.Level(level)
// }

// // Sample returns a logger with the s sampler.
// func (l *Logger) Sample(s zerolog.Sampler) zerolog.Logger {
// 	return l.logger.Sample(s)
// }

// // Hook returns a logger with the h Hook.
// func (l *Logger) Hook(h zerolog.Hook) zerolog.Logger {
// 	return l.logger.Hook(h)
// }

// // Debug starts a new message with debug level.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Debug() *zerolog.Event {
// 	return l.logger.Debug()
// }

// // Info starts a new message with info level.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Info() *zerolog.Event {
// 	return l.logger.Info()
// }

// // Warn starts a new message with warn level.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Warn() *zerolog.Event {
// 	return l.logger.Warn()
// }

// // Error starts a new message with error level.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Error() *zerolog.Event {
// 	return l.logger.Error()
// }

// // Fatal starts a new message with fatal level. The os.Exit(1) function
// // is called by the Msg method.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Fatal() *zerolog.Event {
// 	return l.logger.Fatal()
// }

// // Panic starts a new message with panic level. The message is also sent
// // to the panic function.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Panic() *zerolog.Event {
// 	return l.logger.Panic()
// }

// // WithLevel starts a new message with level.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
// 	return l.logger.WithLevel(level)
// }

// // Log starts a new message with no level. Setting zerolog.GlobalLevel to
// // zerolog.Disabled will still disable events produced by this method.
// //
// // You must call Msg on the returned event in order to send the event.
// func (l *Logger) Log() *zerolog.Event {
// 	return l.logger.Log()
// }

// // Print sends a log event using debug level and no extra field.
// // Arguments are handled in the manner of fmt.Print.
// func (l *Logger) Print(v ...interface{}) {
// 	l.logger.Print(v...)
// }

// // Printf sends a log event using debug level and no extra field.
// // Arguments are handled in the manner of fmt.Printf.
// func (l *Logger) Printf(format string, v ...interface{}) {
// 	l.logger.Printf(format, v...)
// }

// // Ctx returns the Logger associated with the ctx. If no logger
// // is associated, a disabled logger is returned.
// func (l *Logger) Ctx(ctx context.Context) *Logger {
// 	return &Logger{logger: zerolog.Ctx(ctx)}
// }
