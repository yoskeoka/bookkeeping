package bookkeeping

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed _embed/sql/*
var sqlFiles embed.FS

type DB struct {
	dbFilePath string
	dbConn     *sql.DB
}

func NewDB(path string) (*DB, error) {
	dir := filepath.Dir(path)
	_, statDirErr := os.Stat(dir)
	if errors.Is(statDirErr, os.ErrNotExist) {
		err := os.MkdirAll(dir, 0744)
		if err != nil {
			return nil, err
		}
	}

	var initRequired = false

	_, dbStatErr := os.Stat(path)
	if errors.Is(dbStatErr, os.ErrNotExist) {
		initRequired = true
	}

	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open database file %s: %v", path, err)
	}

	db := &DB{
		dbFilePath: path,
		dbConn:     sqlDB,
	}

	if initRequired {
		fmt.Println("initializing database...")
		if err := db.InitSchema(); err != nil {
			return nil, fmt.Errorf("database init schema error: %v", err)
		}

		if err := db.InitAccounts(); err != nil {
			return nil, fmt.Errorf("database init accounts data error: %v", err)
		}
		fmt.Println("database initialized:", path)
	}

	return db, nil
}

func (d *DB) InitSchema() error {
	sb, err := sqlFiles.ReadFile("_embed/sql/schema.sql")
	if err != nil {
		return err
	}

	_, err = d.dbConn.Exec(string(sb))
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) InitAccounts() error {
	acc, err := sqlFiles.ReadFile("_embed/sql/accounts_ja.sql")
	if err != nil {
		return err
	}

	_, err = d.dbConn.Exec("delete from accounts")
	if err != nil {
		return err
	}

	_, err = d.dbConn.Exec(string(acc))
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Close() error {
	return d.dbConn.Close()
}

func (d *DB) Delete() error {
	return os.Remove(d.dbFilePath)
}

type DBAccounts struct {
	db *DB
}

func NewDBAccounts(db *DB) *DBAccounts {
	return &DBAccounts{db}
}

func (a *DBAccounts) Insert(items ...Account) error {
	tx, err := a.db.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into accounts(code, name, is_bs, is_left) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(item.Code, item.Name, item.IsBS, item.IsLeft)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type DBJournals struct {
	db *DB
}

func NewDBJournals(db *DB) *DBJournals {
	return &DBJournals{db}
}

func (jn *DBJournals) Insert(items ...Journal) error {
	tx, err := jn.db.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into journals(code, date, description, left, right) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(item.Code, item.Date, item.Description, item.Left, item.Right)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type DBJournalsFetchOption struct {
	After  sql.NullTime
	Before sql.NullTime
	Code   []int
}

func (jn *DBJournals) Fetch(opt DBJournalsFetchOption) ([]Journal, error) {
	q := []string{
		`
		SELECT jn.id, jn.date, jn.code, jn.description, jn.left, jn.right, 
				a.code, a.name, a.is_bs, a.is_left
		FROM journals AS jn
		INNER JOIN accounts AS a ON a.code = jn.code
		`,
	}
	w := []string{}
	args := []interface{}{}

	if opt.After.Valid {
		w = append(w, "? <= jn.date")
		args = append(args, opt.After)
	}
	if opt.Before.Valid {
		w = append(w, "? >= jn.date")
		args = append(args, opt.Before)
	}
	if len(opt.Code) > 0 {
		w = append(w, "jn.code IN ("+strings.Repeat("?,", len(opt.Code)-1)+"?)")
		for _, c := range opt.Code {
			args = append(args, c)
		}
	}

	if len(w) > 0 {
		q = append(q, "WHERE", strings.Join(w, " AND"))
	}

	query := strings.Join(q, " ")
	stmt, err := jn.db.dbConn.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	items := []Journal{}
	for rows.Next() {
		item := Journal{}
		err := rows.Scan(
			&item.ID, &item.Date, &item.Code, &item.Description, &item.Left, &item.Right,
			&item.Account.Code, &item.Account.Name, &item.Account.IsBS, &item.Account.IsLeft,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
