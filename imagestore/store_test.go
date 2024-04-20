package imagestore_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/matmazurk/acc2/imagestore"
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
