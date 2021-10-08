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
		// TODO: remove test db file
	})

	return tdb
}

func Test_DBGeneralLedger_Insert(t *testing.T) {
	tdb := NewTestDB(t)
	gl := bookkeeping.NewDBGeneralLedger(tdb)

	insertItems := []bookkeeping.GeneralLedger{
		{Date: sql.NullTime{Time: time.Date(2021, 01, 03, 0, 0, 0, 0, time.UTC), Valid: true}, Code: 1000, Description: "現金", Left: 100000, Right: 0},
		{Date: sql.NullTime{Time: time.Date(2021, 01, 03, 0, 0, 0, 0, time.UTC), Valid: true}, Code: 3100, Description: "資本金", Left: 0, Right: 100000},
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
