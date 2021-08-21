package main

import (

	// "context".
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"

	//	"github.com/sheldonhull/go-semantic-sentences/internal/logger"
	"github.com/sheldonhull/go-semantic-sentences/pkg/linter"
)

const (
	exitOK     = 0
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
		// fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
	os.Exit(exitOK)
}

func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("no arguments")
	}

	ApplicationHeader()
	pterm.EnableDebugMessages()

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var (
		debug  = fs.Bool("debug", false, "sets log level to debug and console pretty output")
		source = fs.String("source", "", "source directory or file")
		write  = fs.Bool("write", false, "default to stdout, otherwise replace contents of the file")
	)
	// ff.Parse(fs, args, ff.WithEnvVarNoPrefix())
	if err := fs.Parse(args); err != nil {
		pterm.Error.Printf("ff.Parse: %v\n", err)

		return err
	}

	if *debug {
		pterm.EnableDebugMessages()
		pterm.Error.ShowLineNumber = true
	}
	pterm.Debug.Println("")
	pterm.Debug.Printf("source: %10v\n", *source)
	pterm.Debug.Printf("write: %10v\n", *write)
	pterm.Debug.Printf("debug: %10v\n", *debug)

	fullpath, err := filepath.Abs(*source)
	if err != nil {
		pterm.Error.Printf("filepath.Abs(%s): %v\n", *source, err)
		return err
	}

	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		pterm.Error.Printf("os.Stat(%s): %v\n", fullpath, err)
		return err
	}
	var files []string

	if fileInfo.Mode().IsDir() {

		d, err := os.ReadDir(fullpath)
		for _, f := range d {
			files = append(files, f.Name())
		}
		if err != nil {
			pterm.Error.Printf("os.ReadDir(%s): [%v]\n", fullpath, err)
			return err
		}
	}
	if fileInfo.Mode().IsRegular() {
		files = append(files, fullpath)
	}
	// 	if err != nil {

	// 	pterm.Error.Printf("ReadDir [%s]: %v\n", *source, err)
	// 	return err
	// }
	// leveledList := pterm.LeveledList{}
	pterm.Info.Printf("%s: %d files\n", fullpath, len(files))
	for _, f := range files {
		// leveledList = append(leveledList, pterm.LeveledListItem{Level: 1, Text: f.})
		// if f.IsDir() {
		// 	pterm.Info.Println("🔁 skipping since directory object", f.Name())

		// 	continue
		// }
		b, err := ioutil.ReadFile(f)
		if err != nil {
			pterm.Error.Printf("ioutil.ReadFile(%s): [%v]\n", f, err)
			return err
		}

		// for _, file := range files{ }
		count := linter.CountViolations(b)

		formatted := linter.FormatSemanticLineBreak(b)
		err = ioutil.WriteFile(f, []byte(formatted), os.ModeDevice)
		if err != nil {
			pterm.Error.Printf("ioutil.WriteFile(%s): [%v]\n", f, err)
			return err
		}
		pterm.Success.Printf("✔️ %s [violation count: %d]\n", f, count)
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
