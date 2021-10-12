package bookkeeping

import (
	"database/sql"
	"embed"
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
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	return nil, err
	// }
	// f := filepath.Join(homeDir, "bookkeeping.db")

	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db := &DB{
		dbFilePath: path,
		dbConn:     sqlDB,
	}

	return db, nil
}

func (d *DB) Init() error {
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

type DBGeneralLedger struct {
	db *DB
}

func NewDBGeneralLedger(db *DB) *DBGeneralLedger {
	return &DBGeneralLedger{db}
}

func (gl *DBGeneralLedger) Insert(items ...GeneralLedger) error {
	tx, err := gl.db.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into general_ledger(code, date, description, left, right) values(?, ?, ?, ?, ?)")
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

type DBGeneralLedgerFetchOption struct {
	After  sql.NullTime
	Before sql.NullTime
}

func (gl *DBGeneralLedger) Fetch(opt DBGeneralLedgerFetchOption) ([]GeneralLedger, error) {
	q := []string{"select id, date, code, description, left, right from general_ledger"}
	w := []string{}
	args := []interface{}{}

	if opt.After.Valid {
		w = append(w, "? <= date")
		args = append(args, opt.After)
	}
	if opt.Before.Valid {
		w = append(w, "? >= date")
		args = append(args, opt.Before)
	}

	if len(w) > 0 {
		q = append(q, "WHERE", strings.Join(w, " AND"))
	}

	query := strings.Join(q, " ")
	stmt, err := gl.db.dbConn.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	items := []GeneralLedger{}
	for rows.Next() {
		item := GeneralLedger{}
		err := rows.Scan(&item.ID, &item.Date, &item.Code, &item.Description, &item.Left, &item.Right)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
