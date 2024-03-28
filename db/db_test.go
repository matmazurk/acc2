package db_test

import (
	"os"
	"testing"
	"time"

	"github.com/matmazurk/acc2/db"
	"github.com/matmazurk/acc2/expense"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	dir, err := os.MkdirTemp(".", "_testing_bin_")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	dbFilepath := dir + "/exps.db"
	database, err := db.New(dbFilepath)
	require.NoError(t, err)

	const payer = "some-payer"
	const groceries = "groceries"

	t.Run("should_properly_insert_list_payers", func(t *testing.T) {
		err := database.CreatePayer(payer)
		require.NoError(t, err)

		payers, err := database.ListPayers()
		require.NoError(t, err)

		require.Equal(t, []string{payer}, payers)
	})

	t.Run("should_properly_insert_list_categories", func(t *testing.T) {
		err = database.CreateCategory(groceries)
		require.NoError(t, err)

		categories, err := database.ListCategories()
		require.NoError(t, err)

		require.Equal(t, []string{groceries}, categories)
	})

	t.Run("should_properly_insert_select_expenses", func(t *testing.T) {
		now := time.Now()
		exp1, err := expense.NewExpense(
			"shopping",
			payer,
			groceries,
			"10.22",
			"EUR",
			now,
		)
		require.NoError(t, err)
		err = database.Insert(exp1)
		require.NoError(t, err)

		exp2, err := expense.NewExpense(
			"some other shopping",
			payer,
			groceries,
			"22.22",
			"EUR",
			now.Add(time.Minute),
		)
		require.NoError(t, err)
		err = database.Insert(exp2)
		require.NoError(t, err)

		exp3, err := expense.NewExpense(
			"yet some other shopping",
			payer,
			groceries,
			"21.37",
			"EUR",
			now.Add(time.Hour),
		)
		require.NoError(t, err)
		err = database.Insert(exp3)
		require.NoError(t, err)

		expectedOrder := []expense.Expense{
			exp3,
			exp2,
			exp1,
		}

		exps, err := database.SelectExpenses()
		require.NoError(t, err)
		require.Len(t, exps, 3)
		for i, e := range exps {
			expensesEqual(t, expectedOrder[i], e)
		}
	})
}

func expensesEqual(t *testing.T, e1, e2 expense.Expense) {
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

	if !e1.Time().Equal(e2.Time()) {
		t.Errorf("currencies not matching '%s' != '%s'", e1.Time().Format(time.RFC3339), e2.Time().Format(time.RFC3339))
	}
}
