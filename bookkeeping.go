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
		return fmt.Errorf("journals are not balancing: %w", err)
	}

	if err := bk.dbJn.Insert(jn...); err != nil {
		return err
	}

	return nil
}

type FetchGLOpts struct {
	AccountIDList []int
}

func (bk *Bookkeeping) FetchGL(opts ...FetchGLOpts) (map[int][]Journal, error) {
	jnFetchOpts := DBJournalsFetchOption{}
	for _, o := range opts {
		jnFetchOpts.Code = append(jnFetchOpts.Code, o.AccountIDList...)
	}
	journals, err := bk.dbJn.Fetch(jnFetchOpts)
	if err != nil {
		return nil, err
	}

	res := make(map[int][]Journal)
	for _, j := range journals {
		if _, ok := res[j.Code]; !ok {
			res[j.Code] = []Journal{}
		}
		res[j.Code] = append(res[j.Code], j)
	}
	return res, err
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
