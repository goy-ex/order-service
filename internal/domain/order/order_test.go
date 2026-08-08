package order_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain"
	orderpkg "github.com/goy-ex/order-service/internal/domain/order"
	"github.com/goy-ex/order-service/internal/pkg/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	t.Parallel()

	type testcase struct {
		check func(t *testing.T, err error)
		order orderpkg.Order
	}

	testcases := map[string]testcase{
		"valid": {
			order: fixture.NewValidOrder(),
			check: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
		"bad_price": {
			order: func() orderpkg.Order {
				valid := fixture.NewValidOrder()
				valid.Price = 0

				return valid
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "Price", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_amount": {
			order: func() orderpkg.Order {
				valid := fixture.NewValidOrder()
				valid.Qty = 0

				return valid
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "Amount", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_remaining": {
			order: func() orderpkg.Order {
				valid := fixture.NewValidOrder()
				valid.RemainingQty = 0

				return valid
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.FieldNotPositiveError
				require.ErrorAs(t, err, &target)
				require.Equal(t, "Remaining", target.FieldName)
				require.Equal(t, 0, target.Value)
			},
		},
		"bad_side": {
			order: func() orderpkg.Order {
				valid := fixture.NewValidOrder()
				valid.Side = domain.SideFromString("invalid")

				return valid
			}(),
			check: func(t *testing.T, err error) {
				t.Helper()

				var target *domain.NoSuchOptionError
				require.ErrorAs(t, err, &target)
				assert.Equal(t, "Side", target.FieldName)
				assert.Equal(t, "invalid", target.Value)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := orderpkg.NewOrder(
				tc.order.ID,
				tc.order.UserID,
				tc.order.Pair,
				tc.order.Side,
				tc.order.Price,
				tc.order.Qty,
				tc.order.RemainingQty,
				tc.order.CreatedAt,
			)
			tc.check(t, err)
		})
	}
}

func TestNewOrderCreated(t *testing.T) {
	t.Parallel()

	id, err := uuid.NewV7()
	require.NoError(t, err)

	createdAt := time.Now()

	valid := fixture.NewValidOrder()
	order, err := orderpkg.NewOrder(
		id,
		valid.UserID,
		valid.Pair,
		valid.Side,
		valid.Price,
		valid.Qty,
		valid.RemainingQty,
		createdAt,
	)
	require.NoError(t, err)

	oc := orderpkg.NewOrderCreated(id, order)

	assert.Equal(t, oc.ID, id)
	assert.Equal(t, oc.Order, order)
	assert.Equal(t, oc.CreatedAt, createdAt)
}
