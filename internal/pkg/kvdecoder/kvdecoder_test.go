package kvdecoder_test

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/goy-ex/order-service/internal/domain"
	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
	"github.com/goy-ex/order-service/internal/pkg/kvdecoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeAssetsJSON(t *testing.T) {
	t.Parallel()

	type testcase struct {
		reader io.Reader
		check  func(t *testing.T, res map[string]*domain.Asset, err error)
	}

	testcases := map[string]testcase{
		"valid": {
			reader: strings.NewReader(`
[
	{
		"symbol": "BTC",
		"name": "Bitcoin",
		"decimals": 8
	}
]`,
			),
			check: func(t *testing.T, res map[string]*domain.Asset, err error) {
				t.Helper()

				require.NoError(t, err)

				asset, ok := res["BTC"]
				require.True(t, ok)
				assert.Equal(t, "BTC", asset.Symbol)
				assert.Equal(t, "Bitcoin", asset.Name)
				assert.Equal(t, 8, asset.Decimals)
			},
		},
		"invalid_syntax": {
			reader: strings.NewReader("lorem ipsum"),
			check: func(t *testing.T, res map[string]*domain.Asset, err error) {
				t.Helper()

				var target *json.SyntaxError
				require.ErrorAs(t, err, &target)
				require.Nil(t, res)
			},
		},
		"invalid_type": {
			reader: strings.NewReader(`
[
	{
		"symbol": "BTC",
		"name": "Bitcoin",
		"decimals": "6"
	}
]`,
			),
			check: func(t *testing.T, res map[string]*domain.Asset, err error) {
				t.Helper()

				var target *json.UnmarshalTypeError
				require.ErrorAs(t, err, &target)
				require.Nil(t, res)
			},
		},
		"invalid_asset": {
			reader: strings.NewReader(`
[
	{
		"symbol": "BTC",
		"name": "Bitcoin",
		"decimals": -1
	}
]`,
			),
			check: func(t *testing.T, res map[string]*domain.Asset, err error) {
				t.Helper()

				var target *domain.NegativeNumberError
				require.ErrorAs(t, err, &target)

				assert.Equal(t, "Decimals", target.FieldName)
				assert.Equal(t, -1, target.Value)
			},
		},
		"empty": {
			reader: strings.NewReader("[]"),
			check: func(t *testing.T, res map[string]*domain.Asset, err error) {
				t.Helper()

				require.NoError(t, err)
				assert.Empty(t, res)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, err := kvdecoder.DecodeAssetsJSON(tc.reader)
			tc.check(t, res, err)
		})
	}
}

func TestDecodePairsJSON(t *testing.T) {
	t.Parallel()

	type testcase struct {
		reader io.Reader
		check  func(t *testing.T, res map[pairpkg.PairKey]*pairpkg.Pair, err error)
	}

	testcases := map[string]testcase{
		"valid": {
			reader: strings.NewReader(`
[
    {
        "base": "BTC",
        "quote": "USDT",
        "priceTick": 1,
        "qtyTick": 1,
        "status": "active"
    }
]`,
			),
			check: func(t *testing.T, res map[pairpkg.PairKey]*pairpkg.Pair, err error) {
				t.Helper()

				require.NoError(t, err)

				pair, ok := res[pairpkg.NewPairKey("BTC", "USDT")]
				require.True(t, ok)
				assert.Equal(t, "BTC", pair.Base)
				assert.Equal(t, "USDT", pair.Quote)
				assert.Equal(t, 1, pair.PriceTick)
				assert.Equal(t, 1, pair.QtyTick)
				assert.Equal(t, "active", pair.Status.String())
			},
		},
		"invalid_syntax": {
			reader: strings.NewReader("lorem ipsum"),
			check: func(t *testing.T, res map[pairpkg.PairKey]*pairpkg.Pair, err error) {
				t.Helper()

				var target *json.SyntaxError
				require.ErrorAs(t, err, &target)
				require.Nil(t, res)
			},
		},
		"invalid_type": {
			reader: strings.NewReader(`
[
    {
        "base": "BTC",
        "quote": "USDT",
        "priceTick": "1",
        "qtyTick": 1,
        "status": "active"
    }
]`,
			),
			check: func(t *testing.T, res map[pairpkg.PairKey]*pairpkg.Pair, err error) {
				t.Helper()

				var target *json.UnmarshalTypeError
				require.ErrorAs(t, err, &target)
				require.Nil(t, res)
			},
		},
		"invalid_asset": {
			reader: strings.NewReader(`
[
    {
        "base": "BTC",
        "quote": "USDT",
        "priceTick": -1,
        "qtyTick": 1,
        "status": "active"
    }
]`,
			),
			check: func(t *testing.T, res map[pairpkg.PairKey]*pairpkg.Pair, err error) {
				t.Helper()

				var target *domain.NotPositiveNumberError
				require.ErrorAs(t, err, &target)

				assert.Equal(t, "PriceTick", target.FieldName)
				assert.Equal(t, -1, target.Value)
			},
		},
		"empty": {
			reader: strings.NewReader("[]"),
			check: func(t *testing.T, res map[pairpkg.PairKey]*pairpkg.Pair, err error) {
				t.Helper()

				require.NoError(t, err)
				assert.Empty(t, res)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, err := kvdecoder.DecodePairsJSON(tc.reader)
			tc.check(t, res, err)
		})
	}
}
