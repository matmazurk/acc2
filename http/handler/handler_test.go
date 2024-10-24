package handler_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matmazurk/acc2/http/handler"
	"github.com/matmazurk/acc2/model"
	"github.com/stretchr/testify/require"
)

func TestExpenses(t *testing.T) {
	pf := newPersistenceFake()
	is := newImagestoreFake()
	h, err := handler.NewHandler(pf, is)
	require.NoError(t, err)

	mux := http.NewServeMux()
	h.Routes(mux)

	t.Run("should_return_400_for_invalid_forms", func(t *testing.T) {
		tcs := []struct {
			name                string
			contentTypeIncluded bool
			formData            map[string]string
		}{
			{
				name:                "no_content_type",
				contentTypeIncluded: false,
			},
			{
				name:                "no_description",
				contentTypeIncluded: true,
				formData:            map[string]string{},
			},
			{
				name:                "no_author",
				contentTypeIncluded: true,
				formData: map[string]string{
					"description": "some description",
				},
			},
			{
				name:                "no_category",
				contentTypeIncluded: true,
				formData: map[string]string{
					"description": "some description",
					"author":      "some payer",
				},
			},
			{
				name:                "no_amount",
				contentTypeIncluded: true,
				formData: map[string]string{
					"description": "some description",
					"payer":       "some payer",
					"category":    "some category",
				},
			},
			{
				name:                "no_currency",
				contentTypeIncluded: true,
				formData: map[string]string{
					"description": "some description",
					"payer":       "some payer",
					"category":    "some category",
					"amount":      "22.22",
				},
			},
		}
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				for name, value := range tc.formData {
					writer.WriteField(name, value)
				}
				require.NoError(t, writer.Close())
				req, err := http.NewRequest("POST", "/expenses", body)
				require.NoError(t, err)
				if tc.contentTypeIncluded {
					req.Header.Set("Content-Type", writer.FormDataContentType())
				}

				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)
				if rr.Result().StatusCode != http.StatusBadRequest {
					t.Logf("body:'%s'", rr.Body.String())
					t.Fatalf("received status code different than expected:\n%d != %d", rr.Result().StatusCode, http.StatusBadRequest)
				}
			})
		}
	})

	t.Run("should_successfully_store_new_expense", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		description := "expensive shopping"
		amount := "22.22"
		for k, v := range map[string]string{
			"description": description,
			"author":      "some payer",
			"category":    "some category",
			"amount":      amount,
			"currency":    "EUR",
		} {
			writer.WriteField(k, v)
		}
		require.NoError(t, writer.Close())
		req, err := http.NewRequest("POST", "/expenses", body)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Result().StatusCode != http.StatusFound {
			t.Logf("body:'%s'", rr.Body.String())
			t.Fatalf("received status code different than expected:\n%d != %d", rr.Result().StatusCode, http.StatusBadRequest)
		}

		require.Len(t, pf.expenses, 1)
		exp := pf.expenses[0]
		require.Equal(t, description, exp.Description())
		require.Equal(t, amount, exp.Amount())
	})
}

type persistenceFake struct {
	expenses   []model.Expense
	payers     []string
	categories []string
}

func newPersistenceFake() *persistenceFake {
	return &persistenceFake{
		expenses:   []model.Expense{},
		payers:     []string{},
		categories: []string{},
	}
}

func (pf *persistenceFake) Insert(e model.Expense) error {
	pf.expenses = append(pf.expenses, e)
	return nil
}

func (pf *persistenceFake) SelectExpenses() ([]model.Expense, error) {
	return pf.expenses, nil
}

func (pf *persistenceFake) CreatePayer(name string) error {
	pf.payers = append(pf.payers, name)
	return nil
}

func (pf *persistenceFake) CreateCategory(name string) error {
	pf.categories = append(pf.categories, name)
	return nil
}

func (pf *persistenceFake) ListPayers() ([]string, error) {
	return pf.payers, nil
}

func (pf *persistenceFake) ListCategories() ([]string, error) {
	return pf.categories, nil
}

func (pf *persistenceFake) RemoveExpense(_ model.Expense) error {
	return nil
}

type imagestoreFake struct {
	photos map[string][]byte
}

func newImagestoreFake() *imagestoreFake {
	return &imagestoreFake{
		photos: map[string][]byte{},
	}
}

func (isf *imagestoreFake) SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	isf.photos[e.ID()+fileExtension] = contents
	return nil
}

func (isf *imagestoreFake) LoadExpensePhoto(e model.Expense) (io.ReadCloser, error) {
	return nil, nil
}

func (isf *imagestoreFake) getPhoto(e model.Expense, fileExtension string) []byte {
	return isf.photos[e.ID()+fileExtension]
}
