package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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

type ExpenseBuilder struct {
	Id          string
	Description string
	Payer       string
	Category    string
	Amount      string
	Currency    string
	Timestamp   time.Time
}

func (eb ExpenseBuilder) Build() (Expense, error) {
	id, err := parseID(eb.Id)
	if err != nil {
		return Expense{}, errors.Wrapf(err, "could not parse UUID from '%s'", eb.Id)
	}

	if eb.Description == "" {
		return Expense{}, errors.New("description cannot be empty")
	}

	if eb.Payer == "" {
		return Expense{}, errors.New("payer cannot be empty")
	}

	if eb.Category == "" {
		return Expense{}, errors.New("category cannot be empty")
	}

	if eb.Amount == "" {
		return Expense{}, errors.New("amount cannot be empty")
	}

	if eb.Currency == "" {
		return Expense{}, errors.New("currency cannot be empty")
	}

	if eb.Timestamp.IsZero() {
		return Expense{}, errors.New("timestamp cannot be zero value")
	}

	return Expense{
		id:          id,
		description: eb.Description,
		payer:       eb.Payer,
		category:    eb.Category,
		amount:      eb.Amount,
		currency:    eb.Currency,
		timestamp:   eb.Timestamp,
	}, nil
}

func parseID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.NewRandom()
	}
	return uuid.Parse(id)
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
