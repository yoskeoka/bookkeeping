package bookkeeping_test

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/yoskeoka/bookkeeping"
)

func NewTestDB(t *testing.T) *bookkeeping.DB {
	t.Helper()
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "bookkeeping_test.db")
	t.Logf("test db: %v", f)

	tdb, err := bookkeeping.NewDB(f)
	if err != nil {
		t.Fatal(err)
	}

	err = tdb.InitSchema()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		tdb.Close()
		tdb.Delete()
	})

	return tdb
}

func initAccounts(t *testing.T, tdb *bookkeeping.DB) {
	a := bookkeeping.NewDBAccounts(tdb)
	testAccounts := []bookkeeping.Account{
		{Code: 1110, Name: "現金及び預金", IsBS: true, IsLeft: true},
		{Code: 1120, Name: "売掛金", IsBS: true, IsLeft: true},
		{Code: 1130, Name: "商品", IsBS: true, IsLeft: true},
		{Code: 1210, Name: "有形固定資産", IsBS: true, IsLeft: true},
		{Code: 1211, Name: "機械装置", IsBS: true, IsLeft: true},
		{Code: 2100, Name: "買掛金", IsBS: true, IsLeft: false},
		{Code: 2101, Name: "短期借入金", IsBS: true, IsLeft: false},
		{Code: 2102, Name: "未払い法人税等", IsBS: true, IsLeft: false},
		{Code: 2103, Name: "預り金", IsBS: true, IsLeft: false},
		{Code: 2200, Name: "長期借入金", IsBS: true, IsLeft: false},
		{Code: 3100, Name: "資本金", IsBS: true, IsLeft: false},
		{Code: 3200, Name: "資本剰余金", IsBS: true, IsLeft: false},
		{Code: 4100, Name: "商品売上高", IsBS: false, IsLeft: false},
		{Code: 5100, Name: "期首商品棚卸高", IsBS: false, IsLeft: true},
		{Code: 5200, Name: "商品仕入高", IsBS: false, IsLeft: true},
		{Code: 5300, Name: "期末商品棚卸高", IsBS: false, IsLeft: true},
		{Code: 7200, Name: "給与・賞与", IsBS: false, IsLeft: true},
		{Code: 7300, Name: "経費", IsBS: false, IsLeft: true},
		{Code: 8100, Name: "営業外収益", IsBS: false, IsLeft: false},
		{Code: 8200, Name: "営業外費用", IsBS: false, IsLeft: true},
		{Code: 8300, Name: "特別利益", IsBS: false, IsLeft: false},
		{Code: 8400, Name: "特別損失", IsBS: false, IsLeft: true},
		{Code: 9000, Name: "法人税等", IsBS: false, IsLeft: true},
	}

	err := a.Insert(testAccounts...)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DBJournals_Insert(t *testing.T) {
	tdb := NewTestDB(t)
	initAccounts(t, tdb)

	jn := bookkeeping.NewDBJournals(tdb)

	insertItems := []bookkeeping.Journal{
		{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
		{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
	}

	err := jn.Insert(insertItems...)
	if err != nil {
		t.Fatal(err)
	}

	fetchedItems, err := jn.Fetch(bookkeeping.DBJournalsFetchOption{})
	if err != nil {
		t.Fatal(err)
	}

	if len(fetchedItems) != 2 {
		t.Errorf("DBJournals.Insert() then Fetch() want 2 items, but got %v", len(fetchedItems))
	}

	if !reflect.DeepEqual(fetchedItems[0].Date, insertItems[0].Date) {
		t.Errorf("DBJournals.Insert() then Fetch() should return Date = %+v, but got %+v", insertItems[0].Date, fetchedItems[0].Date)
	}
	if fetchedItems[0].Code != insertItems[0].Code {
		t.Errorf("DBJournals.Insert() then Fetch() should return Code = %v, but got %v", insertItems[0].Code, fetchedItems[0].Code)
	}
	if fetchedItems[0].Description != insertItems[0].Description {
		t.Errorf("DBJournals.Insert() then Fetch() should return Description = %v, but got %v", insertItems[0].Description, fetchedItems[0].Description)
	}
	if fetchedItems[0].Left != insertItems[0].Left {
		t.Errorf("DBJournals.Insert() then Fetch() should return Left = %v, but got %v", insertItems[0].Left, fetchedItems[0].Left)
	}
	if fetchedItems[0].Right != insertItems[0].Right {
		t.Errorf("DBJournals.Insert() then Fetch() should return Right = %v, but got %v", insertItems[0].Right, fetchedItems[0].Right)
	}
}

func date(year int, month time.Month, day int) sql.NullTime {
	return sql.NullTime{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}

func Test_DBJournals_Fetch(t *testing.T) {
	type args struct {
		opt bookkeeping.DBJournalsFetchOption
	}
	tests := []struct {
		name      string
		args      args
		seedItems []bookkeeping.Journal
		wantItems []bookkeeping.Journal
	}{
		{
			"After",
			args{bookkeeping.DBJournalsFetchOption{After: date(2021, 01, 16)}},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
				{Date: date(2021, 01, 16), Code: 1110, Description: "現金及び預金", Left: 20000, Right: 0},
				{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
			},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 16), Code: 1110, Description: "現金及び預金", Left: 20000, Right: 0},
				{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
			},
		},
		{
			"Before", args{bookkeeping.DBJournalsFetchOption{Before: date(2021, 01, 03)}},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
				{Date: date(2021, 01, 16), Code: 1110, Description: "現金及び預金", Left: 20000, Right: 0},
				{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
			},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
			},
		},
		{
			"Between (Before-After)",
			args{bookkeeping.DBJournalsFetchOption{After: date(2021, 01, 10), Before: date(2021, 01, 10)}},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
				{Date: date(2021, 01, 10), Code: 1110, Description: "現金及び預金", Left: 500000, Right: 0},
				{Date: date(2021, 01, 10), Code: 3100, Description: "資本金", Left: 0, Right: 500000},
				{Date: date(2021, 01, 16), Code: 1110, Description: "現金及び預金", Left: 20000, Right: 0},
				{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
			},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 10), Code: 1110, Description: "現金及び預金", Left: 500000, Right: 0},
				{Date: date(2021, 01, 10), Code: 3100, Description: "資本金", Left: 0, Right: 500000},
			},
		},
		{
			"Code",
			args{bookkeeping.DBJournalsFetchOption{Code: []int{1110}}},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
			},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
			},
		},
		{
			"CodeRangeFrom",
			args{bookkeeping.DBJournalsFetchOption{CodeRangeFrom: 3000}},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
			},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
			},
		},
		{
			"CodeRangeTo",
			args{bookkeeping.DBJournalsFetchOption{CodeRangeTo: 3000}},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
				{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
			},
			[]bookkeeping.Journal{
				{Date: date(2021, 01, 03), Code: 1110, Description: "現金及び預金", Left: 100000, Right: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := NewTestDB(t)
			initAccounts(t, tdb)
			jn := bookkeeping.NewDBJournals(tdb)

			err := jn.Insert(tt.seedItems...)
			if err != nil {
				t.Fatal(err)
			}

			gotItems, err := jn.Fetch(tt.args.opt)
			if err != nil {
				t.Fatal(err)
			}

			if len(gotItems) != len(tt.wantItems) {
				t.Errorf("DBJournals.Fetch(opt) want %v items, but got %v", len(tt.wantItems), len(gotItems))
			}

			if !reflect.DeepEqual(gotItems[0].Date, tt.wantItems[0].Date) {
				t.Errorf("DBJournals.Fetch(opt) should return Date = %+v, but got %+v", tt.wantItems[0].Date, gotItems[0].Date)
			}
			if gotItems[0].Code != tt.wantItems[0].Code {
				t.Errorf("DBJournals.Fetch(opt) should return Code = %v, but got %v", tt.wantItems[0].Code, gotItems[0].Code)
			}
			if gotItems[0].Description != tt.wantItems[0].Description {
				t.Errorf("DBJournals.Fetch(opt) should return Description = %v, but got %v", tt.wantItems[0].Description, gotItems[0].Description)
			}
			if gotItems[0].Left != tt.wantItems[0].Left {
				t.Errorf("DBJournals.Fetch(opt) should return Left = %v, but got %v", tt.wantItems[0].Left, gotItems[0].Left)
			}
			if gotItems[0].Right != tt.wantItems[0].Right {
				t.Errorf("DBJournals.Fetch(opt) should return Right = %v, but got %v", tt.wantItems[0].Right, gotItems[0].Right)
			}
		})
	}
}
