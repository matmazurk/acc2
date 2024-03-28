package expense

import (
	"time"

	"github.com/Rhymond/go-money"
)

type Expense struct {
	description string
	payer       string
	category    string
	money       string
	timestamp   time.Time
}

func NewExpense(
	description string,
	payer string,
	category string,
	money money.Money,
	timestamp time.Time,
) Expense {
	return Expense{}
}
