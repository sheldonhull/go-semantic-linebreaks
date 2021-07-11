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
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
	// "github.com/rs/zerolog/log"
	// "github.com/sheldonhull/go-semantic-sentences/internal/logger"
	//_ "github.com/sheldonhull/go-semantic-sentences/internal/logger"
	zl "github.com/sheldonhull/go-semantic-sentences/internal/logger"
)

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail   = 1
	MaxSize    = 10
	MaxBackups = 7
	MaxAge     = 7
)

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

	// var debug bool
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

	config := zl.Config{
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      false,
		FileLoggingEnabled:    false,
		Directory:             "",
		Filename:              "",
		MaxSize:               MaxSize,
		MaxBackups:            MaxBackups,
		MaxAge:                MaxAge,
		Level:                 LogLevel,
	}

	_ = zl.InitLogger(config)
	zl.Logger.Info().Msg("test")
	// log.Info().Msg("test")

	// zl.Logger.Info().Msg("test")
	// zl.Logger.Log().Msg("test2")
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

// CountViolations counts the number of lines that would need to be fixed by adding semantic line break. It returns an integer value of the violation count found.
func CountViolations(content []byte) int {
	// TODO: implement CountViolations()
	re := regexp.MustCompile(`(?is)^.*\w\.\s\w.*$`)
	// fmt.Print(string(content))
	matches := re.FindAllString(string(content), -1)
	// l.Logger.Info().Msg("test")
	// ?zerolog.Info().Msg("test2")
	// log.Debug().Strs("matches", matches).Msg("CountViolations")
	// logger.Logger.Debug().Strs("matches", matches).Msg("CountViolations")
	// zl.Logge?r

	// fmt.Println("==== matches =====")
	// fmt.Printf("%v",matches)
	// totalMatched := 0

	// for _, m := range matches {
	// 	totalMatched = totalMatched + 1
	// 	fmt.Println(m)
	// }

	// fmt.Print(content)
	return len(matches)
}
