package domain_test

import (
	"testing"

	"github.com/goy-ex/order-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAsset(t *testing.T) {
	t.Parallel()

	type testcase struct {
		symbol   string
		name     string
		decimals int
		check    func(t *testing.T, asset *domain.Asset, err error)
	}

	testcases := map[string]testcase{
		"decimals_positive": {
			symbol:   "BTC",
			name:     "Bitcoin",
			decimals: 8,
			check: func(t *testing.T, asset *domain.Asset, err error) {
				t.Helper()

				require.NoError(t, err)
				assert.Equal(t, "BTC", asset.Symbol)
				assert.Equal(t, "Bitcoin", asset.Name)
				assert.Equal(t, 8, asset.Decimals)
			},
		},
		"decimals_zero": {
			symbol:   "BTC",
			name:     "Bitcoin",
			decimals: 0,
			check: func(t *testing.T, asset *domain.Asset, err error) {
				t.Helper()

				require.NoError(t, err)
				assert.Equal(t, "BTC", asset.Symbol)
				assert.Equal(t, "Bitcoin", asset.Name)
				assert.Equal(t, 0, asset.Decimals)
			},
		},
		"decimals_negative": {
			symbol:   "BTC",
			name:     "Bitcoin",
			decimals: -1,
			check: func(t *testing.T, asset *domain.Asset, err error) {
				t.Helper()

				var target *domain.NegativeNumberError
				require.ErrorAs(t, err, &target)
				assert.Equal(t, "Decimals", target.FieldName)
				assert.Equal(t, -1, target.Value)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			asset, err := domain.NewAsset(tc.symbol, tc.name, tc.decimals)
			tc.check(t, asset, err)
		})
	}
}

func TestEventTypeFromString(t *testing.T) {
	t.Parallel()

	type testcase struct {
		str  string
		want bool
	}

	testcases := map[string]testcase{
		"valid": {
			str:  domain.EventTypeOrderCreated.String(),
			want: true,
		},
		"invalid": {
			str:  "invalid",
			want: false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			et := domain.EventTypeFromString(tc.str)
			require.Equal(t, et.IsValid(), tc.want)
		})
	}
}

func TestEventType_String(t *testing.T) {
	t.Parallel()

	type testcase struct {
		et   domain.EventType
		want string
	}

	testcases := map[string]testcase{
		"valid": {
			et:   domain.EventTypeOrderCreated,
			want: domain.EventTypeOrderCreated.String(),
		},
		"invalid_from_string": {
			et:   domain.EventTypeFromString("invalid"),
			want: "invalid",
		},
		"invalid_from_number": {
			et:   domain.EventType(100),
			want: "invalid(100)",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, tc.et.String())
		})
	}
}

func TestSideFromString(t *testing.T) {
	t.Parallel()

	type testcase struct {
		str  string
		want string
	}

	tcs := map[string]testcase{
		"bid": {
			str:  domain.SideBid.String(),
			want: domain.SideBid.String(),
		},
		"ask": {
			str:  domain.SideAsk.String(),
			want: domain.SideAsk.String(),
		},
		"invalid": {
			str:  "invalid",
			want: "invalid",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := domain.SideFromString(tc.str)
			require.Equal(t, tc.want, s.String())
		})
	}
}

func TestSide_String(t *testing.T) {
	t.Parallel()

	type testcase struct {
		side domain.Side
		want string
	}

	tcs := map[string]testcase{
		"invalid_from_number": {
			side: domain.Side(100),
			want: "invalid(100)",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, tc.side.String())
		})
	}
}
