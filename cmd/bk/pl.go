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

func plCmd() command {
	fset := flag.NewFlagSet("bk pl", flag.ExitOnError)
	opts := &plOpts{}
	fset.Var(&dateFlag{&opts.startDate}, "start", "start date of P&L time period. (format: yyyymmdd)")
	fset.Var(&dateFlag{&opts.endDate}, "end", "end date of P&L time period. (format: yyyymmdd)")

	return command{
		name:        "pl",
		description: "Show profit and loss statement (P&L)",
		fset:        fset,
		fn: func(args []string, plOpts *globalOpts) error {
			fset.Parse(args)
			return pl(opts, plOpts)
		},
	}
}

type plOpts struct {
	startDate time.Time
	endDate   time.Time
}

func pl(opts *plOpts, glOpts *globalOpts) error {

	db, err := bookkeeping.NewDB(filepath.Join(glOpts.dataDir, databaseName))
	if err != nil {
		return err
	}
	bk := bookkeeping.NewBookkeeping(db)

	fetchPLOpts := bookkeeping.FetchPLOpts{
		Start: opts.startDate,
		End:   opts.endDate,
	}

	items, err := bk.FetchPL(fetchPLOpts)
	if err != nil {
		return err
	}

	printPL(glOpts.output, items)

	return nil
}

func printPL(w io.Writer, pl bookkeeping.PL) {
	fmt.Fprintln(w, "Profit and Loss Statement:")
	fmt.Fprintln(w)

	fprintLFW(w, "description", 45)
	fprintRFW(w, "amount", 20)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 70))

	fprintLFW(w, "Net Sales", 45)
	fprintRFW(w, pl.NetSales, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Cost Sales", 45)
	fprintRFW(w, pl.CostSales, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Gross Profit", 45)
	fprintRFW(w, pl.GrossProfit, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Operating Expences", 45)
	fprintRFW(w, pl.OperatingExpences, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Operating Income", 45)
	fprintRFW(w, pl.OperatingIncome, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Non Operating Incomes", 45)
	fprintRFW(w, pl.NonOperatingIncomes, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Non Operating Expences", 45)
	fprintRFW(w, pl.NonOperatingExpences, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Extraordinary Incomes", 45)
	fprintRFW(w, pl.ExtraordinaryIncomes, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Extraordinary Expences", 45)
	fprintRFW(w, pl.ExtraordinaryExpences, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Income Before Provision For Income Taxes", 45)
	fprintRFW(w, pl.IncomeBeforeProvisionForIncomeTaxes, 20)
	fmt.Fprintln(w)

	fprintLFW(w, "Net Income", 45)
	fprintRFW(w, pl.NetIncome, 20)
	fmt.Fprintln(w)
}
