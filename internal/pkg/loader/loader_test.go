package loader_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/goy-ex/order-service/internal/pkg/loader"
	"github.com/stretchr/testify/require"
)

func validFetch() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("hello")), nil
}

func validDecode[T any](io.Reader) (T, error) {
	var zero T

	return zero, nil
}

var errInvalidFetch = errors.New("invalid fetch")

func invalidFetch() (io.ReadCloser, error) {
	return nil, errInvalidFetch
}

var errInvalidDecode = errors.New("invalid parse")

func invalidDecode[T any](io.Reader) (T, error) {
	var zero T

	return zero, errInvalidDecode
}

var errInvalidClose = errors.New("invalid close")

type invalidReadCloser struct {
	io.Reader
}

func (invalidReadCloser) Close() error {
	return errInvalidClose
}

func invalidCloseFetch() (io.ReadCloser, error) {
	return &invalidReadCloser{}, nil
}

func TestLoader(t *testing.T) {
	type testcase[T any] struct {
		loader *loader.Loader[T]
		check  func(t *testing.T, err error)
	}

	t.Parallel()

	testcases := map[string]testcase[int]{
		"valid": {
			loader: func() *loader.Loader[int] {
				return loader.NewLoader(
					loader.FetcherFunc(validFetch),
					loader.DecoderFunc[int](validDecode[int]),
				)
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		"invalid_fetch": {
			loader: func() *loader.Loader[int] {
				return loader.NewLoader(
					loader.FetcherFunc(invalidFetch),
					loader.DecoderFunc[int](validDecode[int]),
				)
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, errInvalidFetch)
			},
		},
		"invalid_parse": {
			loader: func() *loader.Loader[int] {
				return loader.NewLoader(
					loader.FetcherFunc(validFetch),
					loader.DecoderFunc[int](invalidDecode[int]),
				)
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, errInvalidDecode)
			},
		},
		"invalid_close": {
			loader: func() *loader.Loader[int] {
				return loader.NewLoader(
					loader.FetcherFunc(invalidCloseFetch),
					loader.DecoderFunc[int](validDecode[int]),
				)
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, errInvalidClose)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.loader.Load()
			tc.check(t, err)
		})
	}
}

// func TestFilePairProvider(t *testing.T) {
// 	t.Parallel()

// 	type testcase struct {
// 		parser   pair.Parser
// 		checkErr func(t *testing.T, err error)
// 		filename string
// 	}

// 	testcases := map[string]testcase{
// 		"file_exist": {
// 			filename: "./testdata/exist",
// 			checkErr: func(t *testing.T, err error) {
// 				t.Helper()
// 				require.NoError(t, err)
// 			},
// 			parser: pair.ParserFunc(validParse),
// 		},
// 		"file_not_exist": {
// 			filename: "./testdata/not_exist",
// 			checkErr: func(t *testing.T, err error) {
// 				t.Helper()
// 				require.ErrorIs(t, err, os.ErrNotExist)
// 			},
// 			parser: pair.ParserFunc(validParse),
// 		},
// 		"invalid": {
// 			filename: "./testdata/exist",
// 			checkErr: func(t *testing.T, err error) {
// 				t.Helper()
// 				require.ErrorIs(t, err, errInvalidParse)
// 			},
// 			parser: pair.ParserFunc(invalidParse),
// 		},
// 	}

// 	for name, tc := range testcases {
// 		t.Run(name, func(t *testing.T) {
// 			t.Parallel()

// 			prov := pair.NewFilePairProvider(tc.filename, tc.parser)
// 			_, err := prov.Provide()
// 			tc.checkErr(t, err)
// 		})
// 	}
// }



// func TestParseJSON(t *testing.T) {
// 	t.Parallel()

// 	type testcase struct {
// 		reader   io.Reader
// 		checkErr func(t *testing.T, err error)
// 	}

// 	valid := mustOpenFile(t, "./testdata/valid.json")
// 	badSyntax := mustOpenFile(t, "./testdata/bad_syntax.txt")
// 	invalid := mustOpenFile(t, "./testdata/invalid.json")

// 	testcases := map[string]testcase{
// 		"valid": {
// 			reader: valid,
// 			checkErr: func(t *testing.T, err error) {
// 				t.Helper()
// 				require.NoError(t, err)
// 			},
// 		},
// 		"bad_syntax": {
// 			reader: badSyntax,
// 			checkErr: func(t *testing.T, err error) {
// 				t.Helper()

// 				var syntaxError *json.SyntaxError
// 				require.ErrorAs(t, err, &syntaxError)
// 				t.Log(err.Error())
// 			},
// 		},
// 		"invalid": {
// 			reader: invalid,
// 			checkErr: func(t *testing.T, err error) {
// 				t.Helper()

// 				var target *domain.FieldNotPositiveError
// 				require.ErrorAs(t, err, &target)
// 				require.Equal(t, "PriceTick", target.FieldName)
// 				require.Equal(t, -100, target.Value)
// 				t.Log(err.Error())
// 			},
// 		},
// 	}

// 	for name, tc := range testcases {
// 		t.Run(name, func(t *testing.T) {
// 			t.Parallel()

// 			_, err := pair.ParseJSON(tc.reader)
// 			tc.checkErr(t, err)
// 		})
// 	}
// }
