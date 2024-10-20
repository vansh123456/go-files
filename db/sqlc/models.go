// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

func (e *Currency) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Currency(s)
	case string:
		*e = Currency(s)
	default:
		return fmt.Errorf("unsupported scan type for Currency: %T", src)
	}
	return nil
}

type NullCurrency struct {
	Currency Currency
	Valid    bool // Valid is true if Currency is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCurrency) Scan(value interface{}) error {
	if value == nil {
		ns.Currency, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Currency.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCurrency) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Currency), nil
}

type Account struct {
	ID          int64
	Owner       pgtype.Text
	Balance     pgtype.Int8
	Currency    pgtype.Text
	CreatedAt   pgtype.Timestamptz
	CountryCode pgtype.Int4
}

type Entry struct {
	ID        int64
	AccountID pgtype.Int8
	Amount    int64
	CreatedAt pgtype.Timestamptz
}

type Transfer struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	CreatedAt     pgtype.Timestamptz
}
