package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yoskeoka/bookkeeping"
)

func main() {
	exitCode := cli()
	os.Exit(exitCode)
}

var (
	databaseName string = "bookkeeping.db"
)

func cli() int {
	commands := []command{
		postCmd(),
	}

	fset := flag.NewFlagSet("bk", flag.ExitOnError)
	version := fset.Bool("version", false, "Print version")

	glOpts := &globalOpts{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Print(err)
		return 1
	}

	glOpts.dataDir = filepath.Join(homeDir, ".bookkeeping")

	fset.Usage = func() {
		fmt.Fprintln(fset.Output(), "Usage: bk <command> [command flags]")
		fset.PrintDefaults()
		for _, cmd := range commands {
			if cmd.fset == nil || cmd.fn == nil {
				continue // skip not implemented
			}

			fmt.Fprintf(fset.Output(), "\n%s command:\n", cmd.name)
			cmd.fset.SetOutput(fset.Output())
			cmd.fset.PrintDefaults()
		}
	}

	fset.Parse(os.Args[1:])

	if *version {
		fmt.Printf("Version: %s", bookkeeping.CommitHash)
		return 0
	}

	args := fset.Args()
	if len(args) == 0 {
		fset.Usage()
		return 1
	}

	subCmd := args[0]
	for _, cmd := range commands {
		if cmd.name == subCmd {
			err := cmd.fn(args[1:], glOpts)
			if err != nil {
				log.Print(err)
				return 1
			}
			return 0
		}
	}

	log.Printf("Unknown command: %s", subCmd)

	return 1
}

type globalOpts struct {
	dataDir string
}

type command struct {
	name string
	fset *flag.FlagSet
	fn   func(args []string, gOpts *globalOpts) error
}
