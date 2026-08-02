package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goy-ex/order-service/internal/pkg/loader"
	"github.com/stretchr/testify/require"
)

func TestFileFetcher(t *testing.T) {
	type testcase struct {
		fetcher *loader.FileFetcher
		check   func(t *testing.T, err error)
	}

	t.Parallel()

	testcases := map[string]testcase{
		"valid": {
			fetcher: loader.NewFileFetcher(filepath.Join("testdata", "exist")),
			check: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
		"invalid_filename": {
			fetcher: loader.NewFileFetcher(filepath.Join("testdata", "not_exist")),
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, os.ErrNotExist)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.fetcher.Fetch()
			tc.check(t, err)
		})
	}
}
