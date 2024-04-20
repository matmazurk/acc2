package model_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/matmazurk/acc2/model"
	"github.com/stretchr/testify/require"
)

func TestExpenseBuilder(t *testing.T) {
	tcs := []struct {
		name        string
		in          model.ExpenseBuilder
		errContains string
	}{
		{
			name: "invalid_uuid",
			in: model.ExpenseBuilder{
				Id: "invalid-uuid",
			},
			errContains: "could not parse UUID from 'invalid-uuid'",
		},
		{
			name: "invalid_description",
			in: model.ExpenseBuilder{
				Description: "",
			},
			errContains: "description cannot be empty",
		},
		{
			name: "invalid_payer",
			in: model.ExpenseBuilder{
				Description: "some description",
				Payer:       "",
			},
			errContains: "payer cannot be empty",
		},
		{
			name: "invalid_category",
			in: model.ExpenseBuilder{
				Description: "some description",
				Payer:       "some payer",
				Category:    "",
			},
			errContains: "category cannot be empty",
		},
		{
			name: "invalid_amount",
			in: model.ExpenseBuilder{
				Description: "some description",
				Payer:       "some payer",
				Category:    "some category",
				Amount:      "",
			},
			errContains: "amount cannot be empty",
		},
		{
			name: "invalid_currency",
			in: model.ExpenseBuilder{
				Description: "some description",
				Payer:       "some payer",
				Category:    "some category",
				Amount:      "22.22",
				Currency:    "",
			},
			errContains: "currency cannot be empty",
		},
		{
			name: "invalid_timestamp",
			in: model.ExpenseBuilder{
				Description: "some description",
				Payer:       "some payer",
				Category:    "some category",
				Amount:      "22.22",
				Currency:    "USD",
				Timestamp:   time.Time{},
			},
			errContains: "timestamp cannot be zero value",
		},
	}

	for _, tc := range tcs {
		t.Run("Build_fails_"+tc.name, func(t *testing.T) {
			expense, err := tc.in.Build()
			fmt.Println(err.Error())
			require.Empty(t, expense)
			require.ErrorContains(t, err, tc.errContains)
		})
	}

	t.Run("Build_succeeded", func(t *testing.T) {
		eb := model.ExpenseBuilder{
			Id:          "57f8ea23-4387-491b-bbb0-7195a0e15127",
			Description: "some description",
			Payer:       "some payer",
			Category:    "some category",
			Amount:      "2.22",
			Currency:    "USD",
			Timestamp:   time.Now(),
		}
		expense, err := eb.Build()
		require.NoError(t, err)
		require.Equal(t, eb.Id, expense.ID())
		require.Equal(t, eb.Description, expense.Description())
		require.Equal(t, eb.Payer, expense.Payer())
		require.Equal(t, eb.Category, expense.Category())
		require.Equal(t, eb.Amount, expense.Amount())
		require.Equal(t, eb.Currency, expense.Currency())
		require.True(t, eb.Timestamp.Equal(expense.Time()))
	})
}

