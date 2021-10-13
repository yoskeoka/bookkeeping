package bookkeeping

import (
	"fmt"
)

type Bookkeeping struct {
	db   *DB
	dbGL *DBGeneralLedger
	dbAc *DBAccounts
}

func NewBookkeeping(db *DB) *Bookkeeping {
	return &Bookkeeping{
		db:   db,
		dbGL: NewDBGeneralLedger(db),
		dbAc: NewDBAccounts(db),
	}
}

func (bk *Bookkeeping) Post(gl []GeneralLedger) error {
	if err := balance(gl); err != nil {
		return fmt.Errorf(": %w", err)
	}

	if err := bk.dbGL.Insert(gl...); err != nil {
		return err
	}

	return nil
}

func balance(gl []GeneralLedger) error {
	leftSum, rightSum := 0, 0

	for _, item := range gl {
		leftSum += item.Left
		rightSum += item.Right
	}

	if leftSum != rightSum {
		return fmt.Errorf("credit and debit are not balancing, debit: %v, credit: %v", leftSum, rightSum)
	}

	return nil
}
