package bookkeeping_test

import (
	"testing"

	"github.com/yoskeoka/bookkeeping"
)

func insertTransactionData(t *testing.T, tdb *bookkeeping.DB) {

	bk := bookkeeping.NewBookkeeping(tdb)
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 1), Code: 1110, Left: 500000, Description: "会社設立"},
		{Date: date(2020, 5, 1), Code: 3100, Right: 500000, Description: "会社設立"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 2), Code: 1110, Left: 1000000, Description: "設備導入資金"},
		{Date: date(2020, 5, 2), Code: 2200, Right: 1000000, Description: "設備導入資金"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 3), Code: 7300, Left: 50000, Description: "事務用品"},
		{Date: date(2020, 5, 3), Code: 1110, Right: 50000, Description: "事務用品"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 3), Code: 1211, Left: 500000, Description: "パソコン"},
		{Date: date(2020, 5, 3), Code: 1110, Right: 500000, Description: "パソコン"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 5), Code: 5200, Left: 100000, Description: "おもちゃ仕入"},
		{Date: date(2020, 5, 5), Code: 1110, Right: 100000, Description: "おもちゃ仕入"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 7), Code: 1110, Left: 200000, Description: "おもちゃ販売"},
		{Date: date(2020, 5, 7), Code: 4100, Right: 200000, Description: "おもちゃ販売"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 10), Code: 1110, Left: 1000000, Description: "運転資金"},
		{Date: date(2020, 5, 10), Code: 2101, Right: 1000000, Description: "運転資金"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 11), Code: 5200, Left: 2000000, Description: "おもちゃ仕入"},
		{Date: date(2020, 5, 11), Code: 2100, Right: 2000000, Description: "おもちゃ仕入"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 12), Code: 1120, Left: 4000000, Description: "おもちゃ販売"},
		{Date: date(2020, 5, 12), Code: 4100, Right: 4000000, Description: "おもちゃ販売"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 15), Code: 2100, Left: 2000000, Description: "買掛金清算"},
		{Date: date(2020, 5, 15), Code: 1110, Right: 2000000, Description: "買掛金清算"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 16), Code: 1110, Left: 3000000, Description: "売掛金回収"},
		{Date: date(2020, 5, 16), Code: 1120, Right: 3000000, Description: "売掛金回収"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 20), Code: 7200, Left: 300000, Description: "事務員A給与"},
		{Date: date(2020, 5, 20), Code: 1110, Right: 290000, Description: "給与"},
		{Date: date(2020, 5, 20), Code: 2103, Right: 10000, Description: "源泉所得税"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 21), Code: 2101, Left: 1000000, Description: "返済"},
		{Date: date(2020, 5, 21), Code: 8200, Left: 100000, Description: "支払利息"},
		{Date: date(2020, 5, 21), Code: 1110, Right: 1100000, Description: "返済"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 22), Code: 7300, Left: 200000, Description: "旅費交通費"},
		{Date: date(2020, 5, 22), Code: 1110, Right: 200000, Description: "旅費交通費"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 31), Code: 1130, Left: 100000, Description: "繰越商品"},
		{Date: date(2020, 5, 31), Code: 5300, Right: 100000, Description: "繰越商品"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 31), Code: 7300, Left: 100000, Description: "パソコン減価償却"},
		{Date: date(2020, 5, 31), Code: 1211, Right: 100000, Description: "パソコン減価償却"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bk.Post([]bookkeeping.Journal{
		{Date: date(2020, 5, 31), Code: 9000, Left: 450000, Description: "法人税"},
		{Date: date(2020, 5, 31), Code: 2102, Right: 450000, Description: "法人税"},
	}); err != nil {
		t.Fatal(err)
	}

}

func Test_FetchGL(t *testing.T) {
	tdb := NewTestDB(t)
	initAccounts(t, tdb)
	insertTransactionData(t, tdb)

	bk := bookkeeping.NewBookkeeping(tdb)
	gl, err := bk.FetchGL()
	if err != nil {
		t.Fatal(err)
	}

	if bookkeeping.SumJournal(gl[1110]) != 1460000 {
		t.Errorf("code 1110 balance must be 1460000, but got %v", bookkeeping.SumJournal(gl[1110]))
	}

	if bookkeeping.SumJournal(gl[3100]) != 500000 {
		t.Errorf("code 3100 balance must be 500000, but got %v", bookkeeping.SumJournal(gl[3100]))
	}
}

func Test_FetchPL(t *testing.T) {
	tdb := NewTestDB(t)
	initAccounts(t, tdb)
	insertTransactionData(t, tdb)

	bk := bookkeeping.NewBookkeeping(tdb)
	pl, err := bk.FetchPL(bookkeeping.FetchPLOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if pl.NetSales != 4200000 {
		t.Errorf("pl.NetSales must be 4200000, but got %v", pl.NetSales)
	}

	if pl.CostSales != 2000000 {
		t.Errorf("pl.CostSales must be 2000000, but got %v", pl.CostSales)
	}

	if pl.NetIncome != 1000000 {
		t.Errorf("pl.NetIncome must be 1000000, but got %v", pl.NetIncome)
	}
}

func Test_FetchBS(t *testing.T) {
	tdb := NewTestDB(t)
	initAccounts(t, tdb)
	insertTransactionData(t, tdb)

	bk := bookkeeping.NewBookkeeping(tdb)
	bs, err := bk.FetchBS(bookkeeping.FetchBSOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if bs.TotalAssets != 2960000 {
		t.Errorf("bs.TotalAssets must be 2960000, but got %v", bs.TotalAssets)
	}

	if bs.TotalLiabilities != 1460000 {
		t.Errorf("bs.TotalLiabilities must be 1460000, but got %v", bs.TotalLiabilities)
	}

	if bs.TotalEquity != 1500000 {
		t.Errorf("bs.TotalEquity must be 1500000, but got %v", bs.TotalEquity)
	}

	if bs.TotalLiabilitiesAndEquity != 2960000 {
		t.Errorf("bs.TotalLiabilitiesAndEquity must be 2960000, but got %v", bs.TotalLiabilitiesAndEquity)
	}
}
