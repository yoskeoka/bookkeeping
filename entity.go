package bookkeeping

import (
	"database/sql"
)

type Journal struct {
	ID          int
	Date        sql.NullTime
	Code        int
	Description string
	Left        int
	Right       int

	Account Account
}

type Account struct {
	Code   int
	Name   string
	IsBS   bool
	IsLeft bool
}
