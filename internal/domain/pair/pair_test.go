package pair_test

import (
	"testing"

	"github.com/goy-ex/order-service/internal/domain"
	"github.com/goy-ex/order-service/internal/domain/pair"
	"github.com/goy-ex/order-service/internal/pkg/fixture"
	"github.com/stretchr/testify/require"
)

func TestNewPair(t *testing.T) {
	t.Parallel()

	type testcase struct {
		checkErr func(t *testing.T, err error)
		pair     pair.Pair
	}

	testcases := map[string]testcase{
		"valid": {
			pair: fixture.NewValidPair(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		"invalid_status": {
			pair: func() pair.Pair {
				valid := fixture.NewValidPair()
				valid.Status = 42

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.NoSuchOptionError

				require.ErrorAs(t, err, &target)
				require.Equal(t, "Status", target.FieldName)
				require.Equal(t, "invalid(42)", target.Value)
			},
		},
		"bad_price_tick": {
			pair: func() pair.Pair {
				valid := fixture.NewValidPair()
				valid.PriceTick = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.NotPositiveNumberError

				require.ErrorAs(t, err, &target)
				require.Equal(t, "PriceTick", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_qty_tick": {
			pair: func() pair.Pair {
				valid := fixture.NewValidPair()
				valid.QtyTick = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.NotPositiveNumberError

				require.ErrorAs(t, err, &target)
				require.Equal(t, "QtyTick", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := pair.NewPair(
				tc.pair.Base,
				tc.pair.Quote,
				tc.pair.PriceTick,
				tc.pair.QtyTick,
				tc.pair.Status,
			)
			tc.checkErr(t, err)
		})
	}
}

func TestPairStatusFromString(t *testing.T) {
	t.Parallel()

	type testcase struct {
		str  string
		want string
	}

	tcs := map[string]testcase{
		"active": {
			str:  pair.PairStatusActive.String(),
			want: pair.PairStatusActive.String(),
		},
		"inactive": {
			str:  pair.PairStatusInactive.String(),
			want: pair.PairStatusInactive.String(),
		},
		"invalid": {
			str:  "invalid",
			want: "invalid",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ps := pair.PairStatusFromString(tc.str)
			require.Equal(t, tc.want, ps.String())
		})
	}
}

func TestNewPairKey(t *testing.T) {
	t.Parallel()

	pk := pair.NewPairKey("BTC", "USDT")
	require.Equal(t, "BTC/USDT", string(pk))
}
