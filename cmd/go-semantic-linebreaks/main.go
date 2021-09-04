package main

import (

	// "context".
	"errors"
	"flag"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/ff"
	"github.com/pterm/pterm"

	//	"github.com/sheldonhull/go-semantic-linebreaks/internal/logger"
	"github.com/sheldonhull/go-semantic-linebreaks/pkg/linter"
)

const (
	exitOK     = 0
	exitFail   = 1
	MaxSize    = 10
	MaxBackups = 7
	MaxAge     = 7
	filter     = ".md"
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

// Files contains the matching results for processing.
type FileData struct {
	FullName string
}

// RelativeDepth returns an int reflecting how deeply nested the file is from the source directory.

func (f *FileData) RelativeDepth(fullPathSource string) int {
	return strings.Count(f.RelativePath(fullPathSource), "/")
}

// Parts returns a string slice of the file path split by the separator.
func (f *FileData) Parts(fullPathSource string) []string {
	return strings.Split(f.RelativePath(fullPathSource), string(filepath.Separator))
}

func (f *FileData) RelativePath(fullPathSource string) string {
	return strings.Trim(strings.Replace(f.FullName, fullPathSource, "", -1), string(filepath.Separator))
}

// Directory returns the relative directory of the file.
// For example: /home/user/workdir/nested/file.md would have a relative directory of: `nested` when the fullPathSource directory is `/home/user/workdir`.
func (f *FileData) RelDirectory(fullPathSource string) string {
	return filepath.Dir(f.RelativePath(fullPathSource))
}

func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("no arguments")
	}
	// wd, err := os.Getwd()
	// if err != nil {
	// 	pterm.Error.Println("%v", err)
	// }
	ApplicationHeader()
	pterm.EnableDebugMessages()

	ffs := flag.NewFlagSet("", flag.ExitOnError)

	var (
		source = ffs.String("source", "", "source directory or file")
		debug  = ffs.Bool("debug", false, "sets log level to debug and console pretty output")
		write  = ffs.Bool("write", false, "default to stdout, otherwise replace contents of the file")
	)

	if err := ff.Parse(ffs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		pterm.Error.Printf("ff.Parse: %v\n", err)

		return err
	}

	if *source == "" {
		pterm.Error.Println("no source provided")

		return errors.New("no source provided")
	}

	if *debug {
		pterm.EnableDebugMessages()
		pterm.Error.ShowLineNumber = true
		pterm.DefaultSection.Println("Debug Output")
	}

	if !*write {
		pterm.Info.Println("write not set so no changes will be written")
	}

	pterm.Debug.Printf("source: %10v\n", *source)
	pterm.Debug.Printf("write: %10v\n", *write)
	pterm.Debug.Printf("debug: %10v\n", *debug)

	filedata := []FileData{}

	fullPathSource, err := filepath.Abs(*source)
	if err != nil {
		pterm.Error.Printf("filepath.Abs(%s): %v\n", *source, err)

		return err
	}

	fileInfo, err := os.Stat(fullPathSource)
	if err != nil {
		pterm.Error.Printf("os.Stat(%s): %v\n", fullPathSource, err)

		return err
	}

	if fileInfo.IsDir() {
		// d, err := os.ReadDir(fullPathSource)
		// if err != nil {
		// 	pterm.Error.Printf("os.ReadDir(%s): [%v]\n", fullPathSource, err)

		// 	return err
		// }

		// if path describes a directory, then it will pull the matching files recursively.
		// if a single file, then this will be bypassed

		err = filepath.WalkDir(fullPathSource, func(path string, info fs.DirEntry, err error) error {
			if info.IsDir() {
				return nil
			}

			if filepath.Ext(info.Name()) == filter {
				// wd, _ := os.Getwd()

				fp, err := filepath.Abs(filepath.Join(fullPathSource, info.Name()))
				if err != nil {
					pterm.Error.Printf("filepath.Abs(%s): %v\n", info.Name(), err)

					return err
				}

				// pterm.Debug.Printf("filepath.Abs(filepath.Join(wd, f.Name())): [%s]\n", fp)

				filedata = append(filedata, FileData{
					FullName: fp,
				})
			}
			return nil
		})
	} else {
		pterm.Success.Printf("single file input provided: [%v]\n", fullPathSource)
		filedata = append(filedata, FileData{
			FullName: fullPathSource,
		})
	}

	// Parse the file tree for debug output into a nice tree format
	if *debug {
		pterm.DefaultSection.Println("Files Parsed")

		leveledList := pterm.LeveledList{}

		for _, f := range filedata {
			// relpath := strings.Trim(strings.Replace(f, wd, "", -1), string(filepath.Separator))
			// parts := strings.Split(relpath, string(filepath.Separator))
			// pterm.Debug.Printf("file: [%s] relpath: [%s] parts: [%v]\n", f, relpath, parts)

			leveledList = append(leveledList, pterm.LeveledListItem{
				Level: 0,
				Text:  f.RelativePath(fullPathSource),
			})

			pterm.Debug.Println("==== leveledlist ====")
			pterm.Debug.Printf("FullName: 		[%s]\n", f.FullName)
			pterm.Debug.Printf("fullPathSource:	[%s]\n", fullPathSource)
			pterm.Debug.Printf("RelativeDepth: 	[%s]\n", f.RelativeDepth(fullPathSource))
			pterm.Debug.Printf("RelativePath:   [%s]\n", f.RelativePath(fullPathSource))
			pterm.Debug.Printf("RelDirectory:   [%s]\n", f.RelDirectory(fullPathSource))
			pterm.Debug.Printf("Parts:          [%v]\n", f.Parts(fullPathSource))
		}

		// Generate tree from LeveledList.
		root := pterm.NewTreeFromLeveledList(leveledList)

		// Render TreePrinter
		pterm.DefaultTree.WithRoot(root).Render()
	}

	pterm.Info.Printf("%s: %d files\n", fullPathSource, len(filedata))

	for _, f := range filedata {
		b, err := ioutil.ReadFile(f.FullName)
		if err != nil {
			pterm.Error.Printf("ioutil.ReadFile(%s): [%v]\n", f, err)

			return err
		}

		// TODO: use temp file to ensure reliability on failure then replace original file
		count := linter.CountViolations(b)
		formatted := linter.FormatSemanticLineBreak(b)

		if *write {
			err = ioutil.WriteFile(f.FullName, []byte(formatted), os.ModeDevice)
		}

		if err != nil {
			pterm.Error.Printf("ioutil.WriteFile(%s): [%v]\n", f, err)

			return err
		}

		if *write {
			pterm.Success.Printf("✔️ %s [violation count: %d]\n", f, count)
		} else {
			pterm.Info.Printf("⏩  %s [violation count: %d]\n", f, count)
		}
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
