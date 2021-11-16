package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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
		accountCmd(),
		postCmd(),
		glCmd(),
		deletedbCmd(),
	}

	fset := flag.NewFlagSet("bk", flag.ExitOnError)
	version := fset.Bool("version", false, "Print version")

	glOpts := &globalOpts{
		output: os.Stdout,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Print(err)
		return 1
	}

	glOpts.dataDir = filepath.Join(homeDir, ".bookkeeping")

	fset.Usage = func() {
		fmt.Fprintln(fset.Output(), "Usage: bk <command> [command flags]")
		fset.PrintDefaults()

		fmt.Fprintln(fset.Output())
		fmt.Fprintln(fset.Output(), "Commands:")
		for _, cmd := range commands {
			if cmd.fset == nil || cmd.fn == nil {
				continue // skip not implemented
			}

			fmt.Fprintf(fset.Output(), "  %s:%s%s\n", cmd.name, strings.Repeat(" ", 12-len(cmd.name)), cmd.description)
		}
	}

	fset.Parse(os.Args[1:])

	if *version {
		fmt.Fprintf(fset.Output(), "Version: %s\n", bookkeeping.CommitHash)
		return 0
	}

	args := fset.Args()
	if len(args) == 0 {
		fset.Usage()
		return 1
	}

	err = subcmd("bk", commands, args, glOpts)
	if err != nil {
		fmt.Fprint(fset.Output(), err)
		return 1
	}

	return 0
}

func subcmd(parentCmd string, subCommands []command, args []string, glOpts *globalOpts) error {

	subCmd := args[0]
	for _, cmd := range subCommands {
		if cmd.name == subCmd {

			if len(args) == 1 && cmd.hasSubcommand {
				cmd.fset.Usage()
				return nil
			}

			err := cmd.fn(args[1:], glOpts)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("unknown command: '%s' for '%v'", subCmd, parentCmd)
}

type globalOpts struct {
	dataDir string
	output  io.Writer
}

type command struct {
	name          string
	description   string
	hasSubcommand bool
	fset          *flag.FlagSet
	fn            func(args []string, gOpts *globalOpts) error
}
