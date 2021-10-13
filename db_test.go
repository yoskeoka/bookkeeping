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

	err = tdb.Init()
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
		{Code: 1000, Name: "現金", IsBS: true, IsLeft: true},
		{Code: 3100, Name: "資本金", IsBS: true, IsLeft: false},
	}

	err := a.Insert(testAccounts...)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DBGeneralLedger_Insert(t *testing.T) {
	tdb := NewTestDB(t)
	initAccounts(t, tdb)

	gl := bookkeeping.NewDBGeneralLedger(tdb)

	insertItems := []bookkeeping.GeneralLedger{
		{Date: date(2021, 01, 03), Code: 1000, Description: "現金", Left: 100000, Right: 0},
		{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
	}

	err := gl.Insert(insertItems...)
	if err != nil {
		t.Fatal(err)
	}

	fetchedItems, err := gl.Fetch(bookkeeping.DBGeneralLedgerFetchOption{})
	if err != nil {
		t.Fatal(err)
	}

	if len(fetchedItems) != 2 {
		t.Errorf("DBGeneralLedger.Insert() then Fetch() want 2 items, but got %v", len(fetchedItems))
	}

	if !reflect.DeepEqual(fetchedItems[0].Date, insertItems[0].Date) {
		t.Errorf("DBGeneralLedger.Insert() then Fetch() should return Date = %+v, but got %+v", insertItems[0].Date, fetchedItems[0].Date)
	}
	if fetchedItems[0].Code != insertItems[0].Code {
		t.Errorf("DBGeneralLedger.Insert() then Fetch() should return Code = %v, but got %v", insertItems[0].Code, fetchedItems[0].Code)
	}
	if fetchedItems[0].Description != insertItems[0].Description {
		t.Errorf("DBGeneralLedger.Insert() then Fetch() should return Description = %v, but got %v", insertItems[0].Description, fetchedItems[0].Description)
	}
	if fetchedItems[0].Left != insertItems[0].Left {
		t.Errorf("DBGeneralLedger.Insert() then Fetch() should return Left = %v, but got %v", insertItems[0].Left, fetchedItems[0].Left)
	}
	if fetchedItems[0].Right != insertItems[0].Right {
		t.Errorf("DBGeneralLedger.Insert() then Fetch() should return Right = %v, but got %v", insertItems[0].Right, fetchedItems[0].Right)
	}
}

func date(year int, month time.Month, day int) sql.NullTime {
	return sql.NullTime{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}

func Test_DBGeneralLedger_Fetch(t *testing.T) {
	type args struct {
		items     []bookkeeping.GeneralLedger
		opt       bookkeeping.DBGeneralLedgerFetchOption
		wantItems []bookkeeping.GeneralLedger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"After",
			args{
				[]bookkeeping.GeneralLedger{
					{Date: date(2021, 01, 03), Code: 1000, Description: "現金", Left: 100000, Right: 0},
					{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
					{Date: date(2021, 01, 16), Code: 1000, Description: "現金", Left: 20000, Right: 0},
					{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
				},
				bookkeeping.DBGeneralLedgerFetchOption{After: date(2021, 01, 16)},
				[]bookkeeping.GeneralLedger{
					{Date: date(2021, 01, 16), Code: 1000, Description: "現金", Left: 20000, Right: 0},
					{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
				},
			},
		},
		{
			"Before",
			args{
				[]bookkeeping.GeneralLedger{
					{Date: date(2021, 01, 03), Code: 1000, Description: "現金", Left: 100000, Right: 0},
					{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
					{Date: date(2021, 01, 16), Code: 1000, Description: "現金", Left: 20000, Right: 0},
					{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
				},
				bookkeeping.DBGeneralLedgerFetchOption{Before: date(2021, 01, 03)},
				[]bookkeeping.GeneralLedger{
					{Date: date(2021, 01, 03), Code: 1000, Description: "現金", Left: 100000, Right: 0},
					{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
				},
			},
		},
		{
			"Between (Before-After)",
			args{
				[]bookkeeping.GeneralLedger{
					{Date: date(2021, 01, 03), Code: 1000, Description: "現金", Left: 100000, Right: 0},
					{Date: date(2021, 01, 03), Code: 3100, Description: "資本金", Left: 0, Right: 100000},
					{Date: date(2021, 01, 10), Code: 1000, Description: "現金", Left: 500000, Right: 0},
					{Date: date(2021, 01, 10), Code: 3100, Description: "資本金", Left: 0, Right: 500000},
					{Date: date(2021, 01, 16), Code: 1000, Description: "現金", Left: 20000, Right: 0},
					{Date: date(2021, 01, 16), Code: 3100, Description: "資本金", Left: 0, Right: 20000},
				},
				bookkeeping.DBGeneralLedgerFetchOption{After: date(2021, 01, 10), Before: date(2021, 01, 10)},
				[]bookkeeping.GeneralLedger{
					{Date: date(2021, 01, 10), Code: 1000, Description: "現金", Left: 500000, Right: 0},
					{Date: date(2021, 01, 10), Code: 3100, Description: "資本金", Left: 0, Right: 500000},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := NewTestDB(t)
			initAccounts(t, tdb)
			gl := bookkeeping.NewDBGeneralLedger(tdb)

			err := gl.Insert(tt.args.items...)
			if err != nil {
				t.Fatal(err)
			}

			gotItems, err := gl.Fetch(tt.args.opt)
			if err != nil {
				t.Fatal(err)
			}

			if len(gotItems) != len(tt.args.wantItems) {
				t.Errorf("DBGeneralLedger.Fetch(opt) want %v items, but got %v", len(tt.args.wantItems), len(gotItems))
			}

			if !reflect.DeepEqual(gotItems[0].Date, tt.args.wantItems[0].Date) {
				t.Errorf("DBGeneralLedger.Fetch(opt) should return Date = %+v, but got %+v", tt.args.wantItems[0].Date, gotItems[0].Date)
			}
			if gotItems[0].Code != tt.args.wantItems[0].Code {
				t.Errorf("DBGeneralLedger.Fetch(opt) should return Code = %v, but got %v", tt.args.wantItems[0].Code, gotItems[0].Code)
			}
			if gotItems[0].Description != tt.args.wantItems[0].Description {
				t.Errorf("DBGeneralLedger.Fetch(opt) should return Description = %v, but got %v", tt.args.wantItems[0].Description, gotItems[0].Description)
			}
			if gotItems[0].Left != tt.args.wantItems[0].Left {
				t.Errorf("DBGeneralLedger.Fetch(opt) should return Left = %v, but got %v", tt.args.wantItems[0].Left, gotItems[0].Left)
			}
			if gotItems[0].Right != tt.args.wantItems[0].Right {
				t.Errorf("DBGeneralLedger.Fetch(opt) should return Right = %v, but got %v", tt.args.wantItems[0].Right, gotItems[0].Right)
			}
		})
	}
}
