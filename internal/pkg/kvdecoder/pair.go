package kvdecoder

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goy-ex/order-service/internal/domain"
)

type pairJSONDTO struct {
	Base      string `json:"base"`
	Quote     string `json:"quote"`
	Status    string `json:"status"`
	PriceTick int    `json:"priceTick"`
	QtyTick   int    `json:"qtyTick"`
}

func DecodePairsJSON(r io.Reader) (map[domain.PairKey]*domain.Pair, error) {
	dec := json.NewDecoder(r)

	var dtos []pairJSONDTO

	err := dec.Decode(&dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	pairMap := make(map[domain.PairKey]*domain.Pair, len(dtos))

	for _, dto := range dtos {
		pair, err := domain.NewPair(
			dto.Base,
			dto.Quote,
			dto.PriceTick,
			dto.QtyTick,
			domain.PairStatusFromString(dto.Status),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create pair: %w", err)
		}

		pairMap[domain.NewPairKey(dto.Base, dto.Quote)] = pair
	}

	return pairMap, nil
}
