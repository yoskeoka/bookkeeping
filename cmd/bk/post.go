package main

import (
	"database/sql"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yoskeoka/bookkeeping"
)

func postCmd() command {
	fset := flag.NewFlagSet("bk post", flag.ExitOnError)
	opts := &postOpts{}
	fset.Var(&dateFlag{&opts.date}, "date", "Journal post date. (format: yyyymmdd)")
	fset.Func("left", "Journal debit item. (format: <account code>/<amount>[/<description>]", func(v string) error {
		opts.left = append(opts.left, v)
		return nil
	})
	fset.Func("right", "Journal credit item. (format: <account code>/<amount>[/<description>]", func(v string) error {
		opts.right = append(opts.right, v)
		return nil
	})

	return command{
		name: "post",
		fset: fset,
		fn: func(args []string, glOpts *globalOpts) error {
			fset.Parse(args)
			return post(opts, glOpts)
		},
	}
}

type postOpts struct {
	left  []string
	right []string
	date  time.Time
}

func post(opts *postOpts, glOpts *globalOpts) error {

	journalItems := make([]bookkeeping.Journal, 0, len(opts.left)+len(opts.right))

	for _, s := range opts.left {
		code, amnt, desc, err := parseJournalItem(s)
		if err != nil {
			return err
		}
		jn := bookkeeping.Journal{Code: code, Left: amnt, Description: desc, Date: sql.NullTime{Time: opts.date, Valid: true}}

		journalItems = append(journalItems, jn)
	}

	for _, s := range opts.right {
		code, amnt, desc, err := parseJournalItem(s)
		if err != nil {
			return err
		}
		jn := bookkeeping.Journal{Code: code, Right: amnt, Description: desc, Date: sql.NullTime{Time: opts.date, Valid: true}}

		journalItems = append(journalItems, jn)
	}

	db, err := bookkeeping.NewDB(filepath.Join(glOpts.dataDir, databaseName))
	if err != nil {
		return err
	}
	bk := bookkeeping.NewBookkeeping(db)

	err = bk.Post(journalItems)
	if err != nil {
		return err
	}
	return nil
}

func parseJournalItem(s string) (accCode int, amount int, desc string, err error) {
	cols := strings.Split(s, "/")
	if len(cols) < 2 || len(cols) > 3 {
		return 0, 0, "", fmt.Errorf("cannot parse '%s' as journal item format, format: <account code>/<amount>[/<description>]", s)
	}

	code, err := strconv.Atoi(cols[0])
	if err != nil {
		return 0, 0, "", fmt.Errorf("cannot parse '%s' as account code: %w", cols[0], err)
	}

	a, err := strconv.Atoi(cols[1])
	if err != nil {
		return 0, 0, "", fmt.Errorf("cannot parse '%s' as amount: %w", cols[1], err)
	}

	d := ""
	if len(cols) >= 3 {
		d = cols[2]
	}

	return code, a, d, nil
}
