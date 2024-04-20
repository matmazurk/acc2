package model

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	id          uuid.UUID
	description string
	payer       string
	category    string
	amount      string
	currency    string
	timestamp   time.Time
}

func NewExpenseWithID(
	id string,
	description string,
	payer string,
	category string,
	amount string,
	currency string,
	timestamp time.Time,
) (Expense, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return Expense{}, err
	}

	return Expense{
		id:          uid,
		description: description,
		payer:       payer,
		category:    category,
		amount:      amount,
		currency:    currency,
		timestamp:   timestamp,
	}, nil
}

func NewExpense(
	description string,
	payer string,
	category string,
	amount string,
	currency string,
	timestamp time.Time,
) (Expense, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return Expense{}, err
	}

	return Expense{
		id:          uid,
		description: description,
		payer:       payer,
		category:    category,
		amount:      amount,
		currency:    currency,
		timestamp:   timestamp,
	}, nil
}

func (e Expense) Equal(other Expense) bool {
	return e.ID() == other.ID() &&
		e.Description() == other.Description() &&
		e.Payer() == other.Payer() &&
		e.Category() == other.Category() &&
		e.Amount() == other.Amount() &&
		e.Currency() == other.Currency() &&
		e.Time().Equal(other.Time())
}

func (e Expense) ID() string {
	return e.id.String()
}

func (e Expense) Description() string {
	return e.description
}

func (e Expense) Payer() string {
	return e.payer
}

func (e Expense) Category() string {
	return e.category
}

func (e Expense) Amount() string {
	return e.amount
}

func (e Expense) Currency() string {
	return e.currency
}

func (e Expense) Time() time.Time {
	return e.timestamp
}
