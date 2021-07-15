package main

import (

	// "context".
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/peterbourgon/ff/v3"
	logger "github.com/sheldonhull/go-semantic-sentences/internal/logger"
)

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail   = 1
	MaxSize    = 10
	MaxBackups = 7
	MaxAge     = 7
)

// Logger contains the package level logger provided from internal logger package that wraps up zerolog.
// var log *logger.Logger //nolint: gochecknoglobals

// main configuration from Matt Ryer with minimal logic, passing to run, to allow easier CLI tests.
func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("no arguments")
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	debug := flag.Bool("debug", false, "sets log level to debug and console pretty output")

	// (&debug,
	// 	"debug",
	// 	false,
	// 	"sets log level to debug and console pretty output")

	if err := ff.Parse(flags, args); // ff.WithEnvVarNoPrefix(),
	// ff.WithConfigFileFlag("config"),
	// ff.WithConfigFileParser(fftoml.Parser),
	err != nil {
		return err
	}

	LogLevel := "info"
	if *debug {
		LogLevel = "debug"
	}

	c := logger.Config{
		Enable:                true,
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJson:      false,
		FileLoggingEnabled:    false,
		Directory:             "",
		Filename:              "",
		MaxSize:               MaxSize,
		MaxBackups:            MaxBackups,
		MaxAge:                MaxAge,
		Level:                 LogLevel,
	}

	_ = logger.InitLogger(c)

	return nil
}

// CountViolations counts the number of lines that would need to be fixed by adding semantic line break. It returns an integer value of the violation count found.
func CountViolations(content []byte) int {
	re := regexp.MustCompile(`(?is)^.*\w\.\s\w.*$`)
	matches := re.FindAllString(string(content), -1)
	logger.Log.Info().Int("ViolationCount", len(matches)).Msg("CountViolations")
	// Log.Logger.Info().Int("ViolationCount", len(matches)).Msg("CountViolations")
	// Log.Info().Int("ViolationCount", len(matches)).Msg("CountViolations")
	//	Log.Info().Msg("CountViolations")

	return len(matches)
}
