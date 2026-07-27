package domain

import "github.com/goy-ex/sentinel"

type PairStatus byte

const (
	PairStatusHalted  = 0
	PairStatusTrading = 1
)

func (ps PairStatus) IsValid() bool {
	switch ps {
	case PairStatusHalted, PairStatusTrading:
		return true
	default:
		return false
	}
}

func (ps PairStatus) String() string {
	switch ps {
	case PairStatusHalted:
		return "halted"
	case PairStatusTrading:
		return "trading"
	default:
		return "?"
	}
}

type Price int64
type Quantity int

type Pair struct {
	Base      string
	Quote     string
	PriceTick int
	QtyTick   int
	Status    PairStatus
}

func NewPair(base, quote string, priceTick, qtyTick int, status PairStatus) (*Pair, error) {
	if priceTick <= 0 {
		return nil, sentinel.BadRequest(&FieldNotPositiveError{
			FieldName: "PriceTick",
			Value:     priceTick,
		})
	}

	if qtyTick <= 0 {
		return nil, sentinel.BadRequest(&FieldNotPositiveError{
			FieldName: "QtyTick",
			Value:     qtyTick,
		})
	}

	if !status.IsValid() {
		return nil, sentinel.BadRequest(&NoSuchOptionError{
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

type PairKey string

func NewPairKey(base, quote string) PairKey {
	return PairKey(base + "/" + quote)
}
