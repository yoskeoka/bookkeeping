package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"

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
		name: "gl",
		fset: fset,
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

	fmt.Fprintln(glOpts.output, items)

	return nil
}
