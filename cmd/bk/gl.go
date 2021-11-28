package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/yoskeoka/bookkeeping"
)

func glCmd() command {
	fset := flag.NewFlagSet("bk gl", flag.ExitOnError)
	opts := &glOpts{}
	fset.Func("code", "Account code.", func(s string) error {
		code, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("cannot parse '%s' as account code: %w", s, err)
		}

		opts.code = append(opts.code, code)
		return nil
	})

	return command{
		name:        "gl",
		description: "Show general ledger",
		fset:        fset,
		fn: func(args []string, glOpts *globalOpts) error {
			fset.Parse(args)
			return gl(opts, glOpts)
		},
	}
}

type glOpts struct {
	code []int
}

func gl(opts *glOpts, glOpts *globalOpts) error {

	db, err := bookkeeping.NewDB(filepath.Join(glOpts.dataDir, databaseName))
	if err != nil {
		return err
	}
	bk := bookkeeping.NewBookkeeping(db)

	fetchGLOpts := bookkeeping.FetchGLOpts{
		AccountIDList: append(make([]int, 0, len(opts.code)), opts.code...),
	}

	items, err := bk.FetchGL(fetchGLOpts)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("no journal records found")
	}

	printGL(glOpts.output, items)

	return nil
}

func printGL(w io.Writer, items map[int][]bookkeeping.Journal) {
	fmt.Fprintln(w, "General Ledger:")

	for code, item := range items {
		first := item[0]
		fmt.Fprintln(w)
		printGLItem(w, code, first.Account.Name, item)
	}
}

func printGLItem(w io.Writer, code int, name string, items []bookkeeping.Journal) {

	fmt.Fprintf(w, "Account code %d: '%s'\n", code, name)
	fprintLFW(w, "date", 20)
	fprintLFW(w, "description", 40)
	fprintLFW(w, "debit", 20)
	fprintLFW(w, "credit", 20)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 100))

	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	sort.SliceStable(items, func(i, j int) bool { return items[i].Date.Time.Before(items[j].Date.Time) })

	for _, item := range items {
		fprintLFW(w, item.Date.Time.Format("2006/01/02"), 20)
		fprintLFW(w, item.Description, 40)
		fprintLFW(w, strconv.Itoa(item.Left), 20)
		fprintLFW(w, strconv.Itoa(item.Right), 20)
		fmt.Fprintln(w)
	}
}
