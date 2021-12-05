package bookkeeping

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"
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

type PL struct {
	NetSales                            int
	CostSales                           int
	GrossProfit                         int
	OperatingExpences                   int
	OperatingIncome                     int
	NonOperatingIncomes                 int
	NonOperatingExpences                int
	ExtraordinaryIncomes                int
	ExtraordinaryExpences               int
	IncomeBeforeProvisionForIncomeTaxes int
	ProvisionForIncomeTaxes             int
	NetIncome                           int
}

type FetchPLOpts struct {
	Start time.Time
	End   time.Time
}

func (bk *Bookkeeping) FetchPL(opt FetchPLOpts) (PL, error) {
	pl := PL{}
	dbOpt := DBJournalsFetchOption{}
	if !opt.Start.IsZero() {
		dbOpt.After = sql.NullTime{Time: opt.Start, Valid: true}
	}
	if !opt.End.IsZero() {
		dbOpt.Before = sql.NullTime{Time: opt.End, Valid: true}
	}
	sales, err := bk.dbJn.Fetch(dbOpt.CodeRange(4000, 4999))
	if err != nil {
		return pl, err
	}
	pl.NetSales = SumJournal(sales)

	costSales, err := bk.dbJn.Fetch(dbOpt.CodeRange(5000, 6999))
	if err != nil {
		return pl, err
	}
	pl.CostSales = SumJournal(costSales)
	pl.GrossProfit = pl.NetSales - pl.CostSales

	operatingExpences, err := bk.dbJn.Fetch(dbOpt.CodeRange(7000, 7999))
	if err != nil {
		return pl, err
	}
	pl.OperatingExpences = SumJournal(operatingExpences)
	pl.OperatingIncome = pl.GrossProfit - pl.OperatingExpences

	nonOperatingIncomes, err := bk.dbJn.Fetch(dbOpt.CodeRange(8100, 8199))
	if err != nil {
		return pl, err
	}
	pl.NonOperatingIncomes = SumJournal(nonOperatingIncomes)

	nonOperatingExpences, err := bk.dbJn.Fetch(dbOpt.CodeRange(8200, 8299))
	if err != nil {
		return pl, err
	}
	pl.NonOperatingExpences = SumJournal(nonOperatingExpences)

	extraordinaryIncomes, err := bk.dbJn.Fetch(dbOpt.CodeRange(8300, 8399))
	if err != nil {
		return pl, err
	}
	pl.ExtraordinaryIncomes = SumJournal(extraordinaryIncomes)

	extraordinaryExpences, err := bk.dbJn.Fetch(dbOpt.CodeRange(8400, 8499))
	if err != nil {
		return pl, err
	}
	pl.ExtraordinaryExpences = SumJournal(extraordinaryExpences)

	pl.IncomeBeforeProvisionForIncomeTaxes = pl.OperatingIncome +
		pl.NonOperatingIncomes - pl.NonOperatingExpences +
		pl.ExtraordinaryIncomes - pl.ExtraordinaryExpences

	provisionForIncomeTaxes, err := bk.dbJn.Fetch(dbOpt.CodeRange(9000, 9999))
	if err != nil {
		return pl, err
	}

	pl.ProvisionForIncomeTaxes = SumJournal(provisionForIncomeTaxes)

	pl.NetIncome = pl.IncomeBeforeProvisionForIncomeTaxes - pl.ProvisionForIncomeTaxes

	return pl, nil
}

type BS struct {
	Date time.Time

	TotalCurrentAssets    int
	TotalNoncurrentAssets int
	TotalAssets           int

	TotalCurrentLiabilities    int
	TotalNoncurrentLiabilities int
	TotalLiabilities           int

	OwnersCapital             int
	RetainedErnings           int
	TotalEquity               int
	TotalLiabilitiesAndEquity int
}

type FetchBSOpts struct {
	Date time.Time
}

func (bk *Bookkeeping) FetchBS(opt FetchBSOpts) (BS, error) {

	bs := BS{}

	dbOpt := DBJournalsFetchOption{}
	if !opt.Date.IsZero() {
		dbOpt.Before = sql.NullTime{Time: opt.Date, Valid: true}
		bs.Date = opt.Date
	} else {
		bs.Date = time.Now()
	}

	currentAssets, err := bk.dbJn.Fetch(dbOpt.CodeRange(1100, 1199))
	if err != nil {
		return bs, err
	}

	bs.TotalCurrentAssets = SumJournal(currentAssets)

	noncurrentAssets, err := bk.dbJn.Fetch(dbOpt.CodeRange(1200, 1299))
	if err != nil {
		return bs, err
	}

	bs.TotalNoncurrentAssets = SumJournal(noncurrentAssets)

	bs.TotalAssets = bs.TotalCurrentAssets + bs.TotalNoncurrentAssets

	currentLiabilities, err := bk.dbJn.Fetch(dbOpt.CodeRange(2100, 2199))
	if err != nil {
		return bs, err
	}

	bs.TotalCurrentLiabilities = SumJournal(currentLiabilities)

	noncurrentLiabilities, err := bk.dbJn.Fetch(dbOpt.CodeRange(2200, 2299))
	if err != nil {
		return bs, err
	}

	bs.TotalNoncurrentLiabilities = SumJournal(noncurrentLiabilities)

	bs.TotalLiabilities = bs.TotalCurrentLiabilities + bs.TotalNoncurrentLiabilities

	ownersCapital, err := bk.dbJn.Fetch(dbOpt.CodeRange(3100, 3199))
	if err != nil {
		return bs, err
	}

	bs.OwnersCapital = SumJournal(ownersCapital)

	retainedEarnings, err := bk.dbJn.Fetch(dbOpt.CodeRange(3200, 3299))
	if err != nil {
		return bs, err
	}

	plOpt := FetchPLOpts{}
	if !opt.Date.IsZero() {
		plOpt.End = opt.Date
	}
	pl, err := bk.FetchPL(plOpt)
	if err != nil {
		return bs, err
	}

	bs.RetainedErnings = SumJournal(retainedEarnings) + pl.NetIncome

	bs.TotalEquity = bs.OwnersCapital + bs.RetainedErnings

	bs.TotalLiabilitiesAndEquity = bs.TotalLiabilities + bs.TotalEquity
	return bs, nil
}

func SumJournal(jnn ...[]Journal) int {
	sum := 0

	for _, jn := range jnn {
		for _, j := range jn {
			if j.Account.IsLeft {
				sum += j.Left - j.Right
			} else {
				sum += j.Right - j.Left
			}
		}
	}
	return sum
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
