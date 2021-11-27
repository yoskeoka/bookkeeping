package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/yoskeoka/bookkeeping"
)

func accountCmd() command {
	fset := flag.NewFlagSet("bk account", flag.ExitOnError)

	subcommands := []command{
		accountListCmd(),
	}

	fset.Usage = func() {
		fmt.Fprintln(fset.Output(), "Subcommands:")
		for _, cmd := range subcommands {
			if cmd.fset == nil || cmd.fn == nil {
				continue // skip not implemented
			}

			fmt.Fprintf(fset.Output(), "  %s:%s%s\n", cmd.name, strings.Repeat(" ", 12-len(cmd.name)), cmd.description)
		}
	}

	return command{
		name:          "account",
		description:   "Manage accounts",
		hasSubcommand: true,
		fset:          fset,
		fn: func(args []string, glOpts *globalOpts) error {
			fset.Parse(args)
			return subcmd("bk account", subcommands, fset.Args(), glOpts)
		},
	}
}

func accountListCmd() command {
	fset := flag.NewFlagSet("bk account list", flag.ExitOnError)
	opts := &accountListOpts{}

	fset.StringVar(&opts.codeFilter, "code", "", "Code filter (wildcard '*' supported)")
	fset.StringVar(&opts.nameFilter, "name", "", "Name filter (wildcard '*' supported)")

	return command{
		name:        "list",
		description: "List accounts",
		fset:        fset,
		fn: func(args []string, glOpts *globalOpts) error {
			fset.Parse(args)
			return accountList(opts, glOpts)
		},
	}
}

type accountListOpts struct {
	codeFilter string
	nameFilter string
}

func accountList(opts *accountListOpts, glOpts *globalOpts) error {

	db, err := bookkeeping.NewDB(filepath.Join(glOpts.dataDir, databaseName))
	if err != nil {
		return err
	}
	bk := bookkeeping.NewBookkeeping(db)

	fetchAcOpts := bookkeeping.FetchAcOpts{
		CodeFilter: opts.codeFilter,
		DescFilter: opts.nameFilter,
	}

	items, err := bk.FetchAc(fetchAcOpts)
	if err != nil {
		return err
	}

	if len(items) > 0 {
		printAccounts(glOpts.output, items)
	} else {
		fmt.Fprintln(glOpts.output, "no accounts found")
	}
	return nil
}

func printAccounts(w io.Writer, items []bookkeeping.Account) {
	fmt.Fprintln(w, "Accounts List")
	fprintLFW(w, "code", 10)
	fprintLFW(w, "name", 40)
	fprintLFW(w, "bs/pl", 6)
	fprintLFW(w, "debit/credit", 14)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 70))

	for _, item := range items {
		fprintLFW(w, fmt.Sprintf("%d", item.Code), 10)
		fprintLFW(w, item.Name, 40)
		bspl := "PL"
		if item.IsBS {
			bspl = "BS"
		}
		fprintLFW(w, bspl, 6)
		dc := "credit"
		if item.IsLeft {
			dc = "debit"
		}

		fprintLFW(w, dc, 14)
		fmt.Fprintln(w)
	}
}
