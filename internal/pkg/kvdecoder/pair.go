package kvdecoder

import (
	"encoding/json"
	"fmt"
	"io"

	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
)

type pairJSONDTO struct {
	Base      string `json:"base"`
	Quote     string `json:"quote"`
	Status    string `json:"status"`
	PriceTick int    `json:"priceTick"`
	QtyTick   int    `json:"qtyTick"`
}

func DecodePairsJSON(r io.Reader) (map[pairpkg.PairKey]*pairpkg.Pair, error) {
	dec := json.NewDecoder(r)

	var dtos []pairJSONDTO

	err := dec.Decode(&dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	pairMap := make(map[pairpkg.PairKey]*pairpkg.Pair, len(dtos))

	for _, dto := range dtos {
		pair, err := pairpkg.NewPair(
			dto.Base,
			dto.Quote,
			dto.PriceTick,
			dto.QtyTick,
			pairpkg.PairStatusFromString(dto.Status),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create pair: %w", err)
		}

		pairMap[pairpkg.NewPairKey(dto.Base, dto.Quote)] = pair
	}

	return pairMap, nil
}
