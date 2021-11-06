package bookkeeping

import (
	"fmt"
)

type Bookkeeping struct {
	db   *DB
	dbJn *DBJournals
	dbAc *DBAccounts
}

func NewBookkeeping(db *DB) *Bookkeeping {
	return &Bookkeeping{
		db:   db,
		dbJn: NewDBJournals(db),
		dbAc: NewDBAccounts(db),
	}
}

func (bk *Bookkeeping) Post(jn []Journal) error {
	if err := balance(jn); err != nil {
		return fmt.Errorf(": %w", err)
	}

	if err := bk.dbJn.Insert(jn...); err != nil {
		return err
	}

	return nil
}

type FetchGLOpts struct {
	accountIDList []int
}

func (bk *Bookkeeping) FetchGL(opts ...FetchGLOpts) ([]Journal, error) {
	return nil, nil
}

func balance(jn []Journal) error {
	leftSum, rightSum := 0, 0

	for _, item := range jn {
		leftSum += item.Left
		rightSum += item.Right
	}

	if leftSum != rightSum {
		return fmt.Errorf("credit and debit are not balancing, debit: %v, credit: %v", leftSum, rightSum)
	}

	return nil
}
