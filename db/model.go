package db

import (
	"time"
)

type expense struct {
	ID          string    `db:"id"`
	Payer       payer     `db:"payer"`
	CategoryID  category  `db:"category"`
	Description string    `db:"description"`
	Amount      string    `db:"amount"`
	Currency    string    `db:"currency"`
	CreatedAt   time.Time `db:"created_at"`
}

type payer struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

func (p payer) isZero() bool {
	return p == payer{}
}

type category struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

func (c category) isZero() bool {
	return c == category{}
}
