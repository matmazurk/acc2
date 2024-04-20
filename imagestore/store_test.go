package imagestore_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/matmazurk/acc2/imagestore"
	"github.com/matmazurk/acc2/model"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	t.Run("should_do_nothing_when_dir_exists", func(t *testing.T) {
		dir, err := os.MkdirTemp(".", "__tmpdir_")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(dir) })

		_, err = imagestore.NewStore(dir)
		require.NoError(t, err)
	})

	t.Run("should_create_dir_when_no_exists", func(t *testing.T) {
		filepath := fmt.Sprintf("./%s%d", "__tmpdir_", time.Now().UnixMilli())
		_, err := imagestore.NewStore(filepath)
		require.NoError(t, err)

		fi, err := os.Stat(filepath)
		require.NoError(t, err)
		require.True(t, fi.IsDir())

		fi, err = os.Stat(filepath + "/" + "photos")
		require.NoError(t, err)
		require.True(t, fi.IsDir())

		err = os.RemoveAll(filepath)
		require.NoError(t, err)
	})
}

func TestSaveExpensePhoto(t *testing.T) {
	filepath := fmt.Sprintf("./%s%d", "__tmpdir_", time.Now().UnixMilli())
	store, err := imagestore.NewStore(filepath)
	require.NoError(t, err)
	defer os.RemoveAll(filepath)

	someDate := time.Date(2024, time.April, 10, 13, 40, 0, 0, time.UTC)
	someExp, err := model.ExpenseBuilder{
		Id:          "57f8ea23-4387-491b-bbb0-7195a0e15127",
		Description: "some expense",
		Payer:       "some payer",
		Category:    "groceries",
		Amount:      "22.22",
		Currency:    "USD",
		CreatedAt:   someDate,
	}.Build()
	require.NoError(t, err)
	fileExtension := "jpeg"

	fileContents := []byte("some contents")
	err = store.SaveExpensePhoto(someExp, fileExtension, io.NopCloser(bytes.NewReader(fileContents)))
	require.NoError(t, err)

	expectedFilepath := filepath + "/photos/100424_1340_57f8ea23-4387-491b-bbb0-7195a0e15127." + fileExtension
	fi, err := os.Stat(expectedFilepath)
	require.NoError(t, err)
	require.False(t, fi.IsDir())

	actualContents, err := os.ReadFile(expectedFilepath)
	require.NoError(t, err)
	require.Equal(t, fileContents, actualContents)
}
