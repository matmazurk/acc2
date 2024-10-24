package db_test

import (
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/matmazurk/acc2/db"
	"github.com/matmazurk/acc2/model"
	"github.com/stretchr/testify/require"
)

const dbFile = ".test.db"

func TestDB(t *testing.T) {
	c, err := db.New(dbFile)
	require.NoError(t, err)

	payer := uuid.NewString()
	category := uuid.NewString()

	t.Run("should_properly_insert_list_payers", func(t *testing.T) {
		err := c.CreatePayer(payer)
		require.NoError(t, err)
		err = c.CreatePayer(payer)
		require.Error(t, err)

		payers, err := c.ListPayers()
		require.NoError(t, err)
		idx := slices.IndexFunc(payers, func(name string) bool { return name == payer })
		require.NotEqual(t, -1, idx)
	})

	t.Run("should_properly_insert_list_categories", func(t *testing.T) {
		err := c.CreateCategory(category)
		require.NoError(t, err)
		err = c.CreateCategory(category)
		require.Error(t, err)

		categories, err := c.ListCategories()
		require.NoError(t, err)
		idx := slices.IndexFunc(categories, func(name string) bool { return name == category })
		require.NotEqual(t, -1, idx)
	})

	t.Run("should_properly_insert_select_expenses", func(t *testing.T) {
		now := time.Now()
		exp1, err := model.ExpenseBuilder{
			Description: "shopping",
			Payer:       payer,
			Category:    category,
			Amount:      "10.22",
			Currency:    "EUR",
			CreatedAt:   now,
		}.Build()
		require.NoError(t, err)
		err = c.Insert(exp1)
		require.NoError(t, err)

		exp2, err := model.ExpenseBuilder{
			Description: "some other shopping",
			Payer:       payer,
			Category:    category,
			Amount:      "22.22",
			Currency:    "EUR",
			CreatedAt:   now.Add(time.Minute),
		}.Build()
		require.NoError(t, err)
		err = c.Insert(exp2)
		require.NoError(t, err)

		exp3, err := model.ExpenseBuilder{
			Description: "yet some other shopping",
			Payer:       payer,
			Category:    category,
			Amount:      "21.37",
			Currency:    "EUR",
			CreatedAt:   now.Add(time.Hour),
		}.Build()
		require.NoError(t, err)
		err = c.Insert(exp3)
		require.NoError(t, err)

		expectedOrder := []model.Expense{
			exp3,
			exp2,
			exp1,
		}

		exps, err := c.SelectExpenses()
		require.NoError(t, err)

		filteredExps := filterExpenses(exps, exp1.ID(), exp2.ID(), exp3.ID())
		if len(filteredExps) != 3 {
			t.Fatalf("not all expenses found during filtering: %d != %d", len(filteredExps), 3)
		}
		for i, e := range filteredExps {
			expensesEqual(t, expectedOrder[i], e)
		}
	})

	t.Run("should_properly_insert_delete_expense", func(t *testing.T) {
		now := time.Now()
		exp, err := model.ExpenseBuilder{
			Description: "shopping",
			Payer:       payer,
			Category:    category,
			Amount:      "10.22",
			Currency:    "EUR",
			CreatedAt:   now,
		}.Build()
		require.NoError(t, err)
		err = c.Insert(exp)
		require.NoError(t, err)

		exps, err := c.SelectExpenses()
		require.NoError(t, err)
		idx := slices.IndexFunc(exps, func(e model.Expense) bool { return e.ID() == exp.ID() })
		require.Positive(t, idx)

		err = c.RemoveExpense(exp)
		require.NoError(t, err)

		exps, err = c.SelectExpenses()
		require.NoError(t, err)
		idx = slices.IndexFunc(exps, func(e model.Expense) bool { return e.ID() == exp.ID() })
		require.Equal(t, -1, idx)
	})
}

func expensesEqual(t *testing.T, e1, e2 model.Expense) {
	t.Helper()

	if e1.ID() != e2.ID() {
		t.Errorf("ids not matching '%s' != '%s'", e1.ID(), e2.ID())
	}

	if e1.Description() != e2.Description() {
		t.Errorf("descriptions not matching '%s' != '%s'", e1.Description(), e2.Description())
	}

	if e1.Payer() != e2.Payer() {
		t.Errorf("payers not matching '%s' != '%s'", e1.Payer(), e2.Payer())
	}

	if e1.Category() != e2.Category() {
		t.Errorf("categories not matching '%s' != '%s'", e1.Category(), e2.Category())
	}

	if e1.Amount() != e2.Amount() {
		t.Errorf("amounts not matching '%s' != '%s'", e1.Amount(), e2.Amount())
	}

	if e1.Currency() != e2.Currency() {
		t.Errorf("currencies not matching '%s' != '%s'", e1.Currency(), e2.Currency())
	}

	if !e1.CreatedAt().Equal(e2.CreatedAt()) {
		t.Errorf("currencies not matching '%s' != '%s'", e1.CreatedAt().Format(time.RFC3339), e2.CreatedAt().Format(time.RFC3339))
	}
}

func filterExpenses(exps []model.Expense, ids ...string) []model.Expense {
	var ret []model.Expense
	for _, exp := range exps {
		for _, id := range ids {
			if exp.ID() == id {
				ret = append(ret, exp)
			}
		}

	}

	return ret
}
