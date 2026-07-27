package pair_test

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/goy-ex/order-service/internal/adapter/provider/pair"
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestFileProvider(t *testing.T) {
	t.Parallel()

	type testcase struct {
		parser   pair.Parser
		checkErr func(t *testing.T, err error)
		filename string
	}

	validDummyParse := func(r io.Reader) (map[domain.PairKey]*domain.Pair, error) {
		//nolint:nilnil // dummy
		return nil, nil
	}

	errInvalidDummyParse := errors.New("invalid dummy parse error")

	invalidDummyParse := func(r io.Reader) (map[domain.PairKey]*domain.Pair, error) {
		return nil, errInvalidDummyParse
	}

	testcases := map[string]testcase{
		"file_exist": {
			filename: "./testdata/exist",
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			parser: pair.ParserFunc(validDummyParse),
		},
		"file_not_exist": {
			filename: "./testdata/not_exist",
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, os.ErrNotExist)
			},
			parser: pair.ParserFunc(validDummyParse),
		},
		"invalid": {
			filename: "./testdata/exist",
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, errInvalidDummyParse)
			},
			parser: pair.ParserFunc(invalidDummyParse),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			prov := pair.NewFileProvider(tc.filename, tc.parser)
			_, err := prov.Provide()
			tc.checkErr(t, err)
		})
	}
}

func mustOpenFile(t *testing.T, name string) *os.File {
	t.Helper()

	//#nosec:G304 // test
	file, err := os.Open(name)
	require.NoError(t, err)

	return file
}

func TestParseJSON(t *testing.T) {
	t.Parallel()

	type testcase struct {
		reader   io.Reader
		checkErr func(t *testing.T, err error)
	}

	valid := mustOpenFile(t, "./testdata/valid.json")
	badSyntax := mustOpenFile(t, "./testdata/bad_syntax.txt")
	invalid := mustOpenFile(t, "./testdata/invalid.json")

	testcases := map[string]testcase{
		"valid": {
			reader: valid,
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		"bad_syntax": {
			reader: badSyntax,
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var syntaxError *json.SyntaxError
				require.ErrorAs(t, err, &syntaxError)
				t.Log(err.Error())
			},
		},
		"invalid": {
			reader: invalid,
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "PriceTick", target.FieldName)
				require.Equal(t, -100, target.Value)
				t.Log(err.Error())
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := pair.ParseJSON(tc.reader)
			tc.checkErr(t, err)
		})
	}
}
