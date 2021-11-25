package main

import (
	"flag"
	"path/filepath"

	"github.com/yoskeoka/bookkeeping"
)

func deletedbCmd() command {
	fset := flag.NewFlagSet("bk deletedb", flag.ExitOnError)

	return command{
		name: "deletedb",
		fset: fset,
		fn: func(args []string, glOpts *globalOpts) error {
			fset.Parse(args)
			return deletedb(glOpts)
		},
	}
}

func deletedb(glOpts *globalOpts) error {

	db, err := bookkeeping.NewDB(filepath.Join(glOpts.dataDir, databaseName))
	if err != nil {
		return err
	}

	err = db.Delete()
	if err != nil {
		return err
	}

	return nil
}
