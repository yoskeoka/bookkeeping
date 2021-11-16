package bookkeeping

import (
	"fmt"
	"regexp"
	"strconv"
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

	for _, j := range jn {
		if err := bk.validateJournalRecord(j); err != nil {
			return err
		}
	}

	if err := bk.dbJn.Insert(jn...); err != nil {
		return err
	}

	return nil
}

func (bk *Bookkeeping) validateJournalRecord(j Journal) error {
	accs, err := bk.dbAc.Fetch(DBAccountsFetchOption{CodePattern: strconv.Itoa(j.Code)})
	if err != nil {
		return err
	}

	if len(accs) != 1 {
		desc := ""
		if len(j.Description) > 0 {
			desc = "/" + j.Description
		}
		norm := "debit"
		if j.Right > 0 {
			norm = "credit"
		}
		return fmt.Errorf("code '%d' is not available (in journal %s record '%d/%d%s')", j.Code, norm, j.Code, j.Left+j.Right, desc)
	}
	return nil
}

type FetchAcOpts struct {
	CodeFilter string
	DescFilter string
}

var accountCodePattern = regexp.MustCompile("[*0-9]+")

func (bk *Bookkeeping) FetchAc(opt FetchAcOpts) ([]Account, error) {

	if len(opt.CodeFilter) > 0 && !accountCodePattern.MatchString(opt.CodeFilter) {
		return nil, fmt.Errorf("code filter may contain numbers or '*' for wildcard")
	}

	acFetchOpt := DBAccountsFetchOption{
		CodePattern:        opt.CodeFilter,
		DescriptionPattern: opt.DescFilter,
	}
	return bk.dbAc.Fetch(acFetchOpt)
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

	if leftSum == 0 || rightSum == 0 {
		return fmt.Errorf("credit or debit is zero-amount")
	}

	if leftSum != rightSum {
		return fmt.Errorf("credit and debit are not balancing, debit: %v, credit: %v", leftSum, rightSum)
	}

	return nil
}
