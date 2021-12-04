package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/yoskeoka/bookkeeping"
)

func bsCmd() command {
	fset := flag.NewFlagSet("bk bs", flag.ExitOnError)
	opts := &bsOpts{}
	fset.Var(&dateFlag{&opts.Date}, "date", "date of Balance Sheet. (format: yyyymmdd)")

	return command{
		name:        "bs",
		description: "Show Balance Sheet",
		fset:        fset,
		fn: func(args []string, bsOpts *globalOpts) error {
			fset.Parse(args)
			return bs(opts, bsOpts)
		},
	}
}

type bsOpts struct {
	Date time.Time
}

func bs(opts *bsOpts, glOpts *globalOpts) error {

	db, err := bookkeeping.NewDB(filepath.Join(glOpts.dataDir, databaseName))
	if err != nil {
		return err
	}
	bk := bookkeeping.NewBookkeeping(db)

	fetchBsOpts := bookkeeping.FetchBSOpts{
		Date: opts.Date,
	}

	bs, err := bk.FetchBS(fetchBsOpts)
	if err != nil {
		return err
	}

	printBS(glOpts.output, bs)

	return nil
}

func printBS(w io.Writer, bs bookkeeping.BS) {
	fmt.Fprintln(w, "Balance Sheet:")
	fmt.Fprintln(w)

	indent := strings.Repeat(" ", 10)

	fmt.Fprintln(w, "Assets:")
	fprintLFW(w, "description", 45)
	fprintRFW(w, "amount", 20)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 65))

	fprintLFW(w, indent+"Total Current Assets", 45)
	fprintRFW(w, bs.TotalCurrentAssets, 20)
	fmt.Fprintln(w)

	fprintLFW(w, indent+"Total Noncurrent Assets", 45)
	fprintRFW(w, bs.TotalNoncurrentAssets, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Total Assets", 45)
	fprintRFW(w, bs.TotalAssets, 20)
	fmt.Fprintln(w)

	fmt.Fprintln(w)

	fmt.Fprintln(w, "Liabilities:")
	fprintLFW(w, "description", 45)
	fprintRFW(w, "amount", 20)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 65))

	fprintLFW(w, "Total Liabilities", 45)
	fprintRFW(w, bs.TotalLiabilities, 20)
	fmt.Fprintln(w)

	fmt.Fprintln(w)

	fmt.Fprintln(w, "Equity:")
	fprintLFW(w, "description", 45)
	fprintRFW(w, "amount", 20)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 65))

	fprintLFW(w, indent+"Owner's Capital", 45)
	fprintRFW(w, bs.OwnersCapital, 20)
	fmt.Fprintln(w)

	fprintLFW(w, indent+"Retained Earnings", 45)
	fprintRFW(w, bs.RetainedErnings, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Total Equity", 45)
	fprintRFW(w, bs.TotalEquity, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Total Liabilities and Equity", 45)
	fprintRFW(w, bs.TotalLiabilitiesAndEquity, 20)
	fmt.Fprintln(w)
}
