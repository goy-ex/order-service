package domain_test

import (
	"testing"

	"github.com/goy-ex/order-service/internal/domain"
	"github.com/goy-ex/order-service/internal/pkg/fixture"
	"github.com/stretchr/testify/require"
)

func TestOrder_New(t *testing.T) {
	t.Parallel()

	type testcase struct {
		checkErr func(t *testing.T, err error)
		order    domain.Order
	}

	testcases := map[string]testcase{
		"valid": {
			order: fixture.NewValidOrder(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
		"bad_price": {
			order: func() domain.Order {
				valid := fixture.NewValidOrder()
				valid.Price = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "Price", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_amount": {
			order: func() domain.Order {
				valid := fixture.NewValidOrder()
				valid.Amount = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "Amount", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_remaining": {
			order: func() domain.Order {
				valid := fixture.NewValidOrder()
				valid.Remaining = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "Remaining", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := domain.NewOrder(
				tc.order.ID,
				tc.order.UserID,
				tc.order.Pair,
				tc.order.Side,
				tc.order.Price,
				tc.order.Amount,
				tc.order.Remaining,
				tc.order.CreatedAt,
			)
			tc.checkErr(t, err)
		})
	}
}

//nolint:funlen // tabletest
func TestPair_New(t *testing.T) {
	t.Parallel()

	type testcase struct {
		checkErr func(t *testing.T, err error)
		pair     domain.Pair
	}

	testcases := map[string]testcase{
		"valid": {
			pair: fixture.NewValidPair(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		"bad_status": {
			pair: func() domain.Pair {
				valid := fixture.NewValidPair()
				valid.Status = 42

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.NoSuchOptionError

				require.ErrorAs(t, err, &target)
				require.Equal(t, "Status", target.FieldName)
				require.Equal(t, "?", target.Value)
			},
		},
		"bad_price_tick": {
			pair: func() domain.Pair {
				valid := fixture.NewValidPair()
				valid.PriceTick = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError

				require.ErrorAs(t, err, &target)
				require.Equal(t, "PriceTick", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_qty_tick": {
			pair: func() domain.Pair {
				valid := fixture.NewValidPair()
				valid.QtyTick = 0

				return valid
			}(),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError

				require.ErrorAs(t, err, &target)
				require.Equal(t, "QtyTick", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := domain.NewPair(
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
