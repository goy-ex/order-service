package domain

import (
	"strconv"

	"github.com/goy-ex/sentinel"
)

type PairStatus byte

const (
	pairStatusInvalid  = 0
	PairStatusInactive = 1
	PairStatusActive   = 2
)

func (s PairStatus) IsValid() bool {
	switch s {
	case PairStatusInactive, PairStatusActive:
		return true
	default:
		return false
	}
}

func (s PairStatus) String() string {
	switch s {
	case PairStatusInactive:
		return "inactive"
	case PairStatusActive:
		return "active"
	default:
		return "invalid(" + strconv.Itoa(int(s)) + ")"
	}
}

func PairStatusFromString(s string) PairStatus {
	switch s {
	case "inactive":
		return PairStatusInactive
	case "active":
		return PairStatusActive
	default:
		return pairStatusInvalid
	}
}

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
