package main

import (

	// "context".
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/peterbourgon/ff/v3"

	//	"github.com/sheldonhull/go-semantic-sentences/internal/logger"
	"github.com/sheldonhull/go-semantic-sentences/internal/logger"
	"github.com/sheldonhull/go-semantic-sentences/pkg/linter"
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
	if err := Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("no arguments")
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	debug := flag.Bool("debug", false, "sets log level to debug and console pretty output")
	source := flag.String("source", "", "source file")
	write := flag.Bool("write", false, "default to stdout, otherwise replace contents of the file")

	// (&debug,
	// 	"debug",
	// 	false,
	// 	"sets log level to debug and console pretty output")

	// ff.WithEnvVarNoPrefix(),

	// ff.WithConfigFileFlag("config"),
	// ff.WithConfigFileParser(fftoml.Parser),
	if err := ff.Parse(flags, args); err != nil {
		return err
	}

	pterm.
		pterm.Info.Println("debug", *debug).
		pterm.Info.Println("source", *source).
		pterm.Info.Println("write", *write)

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
	ApplicationHeader()
	if files, err := os.ReadDir(*source); err != nil {
		logger.Log.Error().Err(err).Str("source", *source).Msg("ReadDir")
		os.Exit(exitFail)
	}
	leveledList := pterm.LeveledList{}
	for _, f := range files {
leveledList = append(leveledList, pterm.LeveledListItem{Level: 1, Text: f.})

	for _, file := range files {
		if file.IsDir() {
			logger.Log.Error().Err(err).Str("source", *source).Msg("Lint")
			continue
		}
		if err := linter.Lint(file.Name()); err != nil {
			logger.Log.Error().Err(err).Str("source", *source).Msg("Lint")
			os.Exit(exitFail)
		}
	}

	filename := *source

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Log.Error().Err(err).Str("filename", filename).Msg("ReadFile")
		os.Exit(exitFail)
	}

	// for _, file := range files{ }

	logger.Log.Info().Int("ViolationCount", len(matches)).Msg("CountViolations")
	formatted := FormatSemanticLineBreak(b)
	ioutil.WriteFile(filename, []byte(formatted), os.ModeDevice)

	return nil
}

// clear console output
func clear() {
    print("\033[H\033[2J")
}

// ApplicationHeader is pterm formatted output to make things look fancy
func ApplicationHeader() *pterm.TextPrinter {
    return pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println(
        "Go Semantic Linebreaks")
}