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
		t.Run(tc.name, func(t *testing.T) {
			expense, err := tc.in.Build()
			fmt.Println(err.Error())
			require.Empty(t, expense)
			require.ErrorContains(t, err, tc.errContains)
		})
	}
}

