package http

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/usecase"
)

type CreateOrderDTO struct {
	Base         string    `json:"base"         validate:"required"`
	Quote        string    `json:"quote"        validate:"required"`
	Side         string    `json:"side"         validate:"required"`
	Price        string    `json:"price"        validate:"required"`
	Qty          string    `json:"qty"          validate:"required"`
	RemainingQty string    `json:"remainingQty" validate:"required"`
	UserID       uuid.UUID `json:"userId"       validate:"required"`
}

func (dto CreateOrderDTO) toInput() (*usecase.CreateOrderInput, error) {
	errs := make([]error, 0, 3)

	price, err := strconv.ParseInt(dto.Price, 10, 64)
	errs = append(errs, err)

	qty, err := strconv.ParseInt(dto.Qty, 10, 64)
	errs = append(errs, err)

	remaining, err := strconv.ParseInt(dto.RemainingQty, 10, 64)
	errs = append(errs, err)

	err = errors.Join(errs...)
	if err != nil {
		return nil, err
	}

	return &usecase.CreateOrderInput{
		UserID:       dto.UserID,
		Base:         dto.Base,
		Quote:        dto.Quote,
		Side:         dto.Side,
		Price:        price,
		Qty:          qty,
		RemainingQty: remaining,
	}, nil
}
