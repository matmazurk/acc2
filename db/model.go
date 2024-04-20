package db

import (
	"time"

	"github.com/matmazurk/acc2/model"
)

type expense struct {
	ID          string
	Description string
	PayerID     uint
	Payer       payer `gorm:"foreignKey:PayerID"`
	CategoryID  uint
	Category    category `gorm:"foreignKey:CategoryID"`
	Amount      string
	Currency    string
	CreatedAt   time.Time
}

func (e expense) toExpense() (model.Expense, error) {
	return model.ExpenseBuilder{
		Id:          e.ID,
		Description: e.Description,
		Payer:       e.Payer.Name,
		Category:    e.Category.Name,
		Amount:      e.Amount,
		Currency:    e.Currency,
		CreatedAt:   e.CreatedAt,
	}.Build()
}

type payer struct {
	ID   uint
	Name string `gorm:"uniqueIndex"`
}

func (p payer) isZero() bool {
	return p == payer{}
}

type category struct {
	ID   uint
	Name string `gorm:"uniqueIndex"`
}

func (c category) isZero() bool {
	return c == category{}
}
