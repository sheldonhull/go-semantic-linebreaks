package main

import (

	// "context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/peterbourgon/ff/v3"

	"github.com/sheldonhull/go-semantic-sentences/internal/logger"
)

var debug bool

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail   = 1
	MaxSize    = 10
	MaxBackups = 7
	MaxAge     = 7
)

// main configuration from Matt Ryer with minimal logic, passing to run, to allow easier CLI tests
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

	flag.BoolVar(&debug,
		"debug",
		false,
		"sets log level to debug and console pretty output")

	if err := ff.Parse(flags, args); // ff.WithEnvVarNoPrefix(),
	// ff.WithConfigFileFlag("config"),
	// ff.WithConfigFileParser(fftoml.Parser),
	err != nil {
		return err
	}

	config := Config{
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      false,
		FileLoggingEnabled:    false,
		Directory:             "",
		Filename:              "",
		MaxSize:               MaxSize,
		MaxBackups:            MaxBackups,
		MaxAge:                MaxAge,
		Level:                 "info",
	}
	Configure(config)

	// l.Logger.Info().Msg("test")
	// Logger.Info().Msg("test")
	// Log.Info().Msg("func run() completed")
	return nil
}

// func main() {
// 	root := &ffcli.Command{
// 		Exec: func(ctx context.Context, args []string) error {
// 			println("hello world")
// 			return nil
// 		},
// 	}

// 	root.ParseAndRun(context.Background(), os.Args[1:])
// }
