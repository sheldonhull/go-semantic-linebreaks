package main

import (

	// "context".
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pterm/pterm"

	//	"github.com/sheldonhull/go-semantic-sentences/internal/logger"
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

	ApplicationHeader()
	pterm.EnableDebugMessages()

	fs := flag.NewFlagSet("", flag.ExitOnError)

	var (
		debug  = fs.Bool("debug", false, "sets log level to debug and console pretty output")
		source = fs.String("source", "test.md", "source directory or file")
		write  = fs.Bool("write", false, "default to stdout, otherwise replace contents of the file")
	)
	// ff.Parse(fs, args, ff.WithEnvVarNoPrefix())
	if err := fs.Parse(args); err != nil {
		pterm.Error.Println("ff.Parse: %v", err)

		return err
	}

	if *debug {
		pterm.EnableDebugMessages()
		pterm.Error.ShowLineNumber = true
	}

	pterm.Debug.Printf("source: %10v\n", *source)
	pterm.Debug.Printf("write: %10v\n", *write)
	pterm.Debug.Printf("debug: %10v\n", *debug)

	// files := []os.FileInfo{}

	// if os.FileInfo.IsDir(*source) {
	// 	files, err := os.ReadDir(*source)
	// 	if err != nil {
	// 		pterm.Error.Println("ReadDir: %v", err)
	// 		os.Exit(exitFail)
	// 	}
	// } else {
	// 	files := os.ReadFile(*source)
	// }

	// if err != nil {
	// 	pterm.Error.Printf("ReadDir: %v\n", err)
	// 	os.Exit(exitFail)
	// }

	// leveledList := pterm.LeveledList{}
	files := *source
	for _, f := range files {
		// leveledList = append(leveledList, pterm.LeveledListItem{Level: 1, Text: f.})

		if os.FileInfo.IsDir(f) {
			pterm.Info.Println("🔁 skipping since directory object", f.Name())

			continue
		}

		b, err := ioutil.ReadFile(f.Name())
		if err != nil {
			pterm.Error.Printf("ReadFile: [%v]\n", err)
			os.Exit(exitFail)
		}

		// for _, file := range files{ }
		formatted := linter.FormatSemanticLineBreak(b)
		ioutil.WriteFile(f.Name(), []byte(formatted), os.ModeDevice)
	}
	return nil
}

// clear console output
// func clear() {
//     fmt.Fprintf(os.Stdout, "\033[H\033[2J")
// }

// ApplicationHeader is pterm formatted output to make things look fancy.
func ApplicationHeader() *pterm.TextPrinter {
	return pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println(
		"Go Semantic Linebreaks")
}
