package bookkeeping

import (
	"database/sql"
)

type GeneralLedger struct {
	ID          int
	Date        sql.NullTime
	Code        int
	Description string
	Left        int
	Right       int
}

type Account struct {
	Code   int
	Name   string
	IsBS   bool
	IsLeft bool
}
