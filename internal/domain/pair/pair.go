package pair

import (
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/goy-ex/sentinel"
)

type Pair struct {
	Base      string
	Quote     string
	PriceTick int
	QtyTick   int
	Status    PairStatus
}

func NewPair(base, quote string, priceTick, qtyTick int, status PairStatus) (*Pair, error) {
	if priceTick <= 0 {
		return nil, sentinel.BadRequest(&domain.FieldNotPositiveError{
			FieldName: "PriceTick",
			Value:     priceTick,
		})
	}

	if qtyTick <= 0 {
		return nil, sentinel.BadRequest(&domain.FieldNotPositiveError{
			FieldName: "QtyTick",
			Value:     qtyTick,
		})
	}

	if !status.IsValid() {
		return nil, sentinel.BadRequest(&domain.NoSuchOptionError{
			FieldName: "Status",
			Value:     status.String(),
		})
	}

	return &Pair{
		Base:      base,
		Quote:     quote,
		PriceTick: priceTick,
		QtyTick:   qtyTick,
		Status:    status,
	}, nil
}
