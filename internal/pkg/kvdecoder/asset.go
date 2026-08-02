package kvdecoder

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goy-ex/order-service/internal/domain"
)

type assetJSONDTO struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
}

func DecodeAssetsJSON(reader io.Reader) (map[string]*domain.Asset, error) {
	dec := json.NewDecoder(reader)

	var dtos []assetJSONDTO

	err := dec.Decode(&dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	assetMap := make(map[string]*domain.Asset, len(dtos))

	for _, dto := range dtos {
		asset, err := domain.NewAsset(dto.Symbol, dto.Name, dto.Decimals)
		if err != nil {
			return nil, fmt.Errorf("failed to create asset: %w", err)
		}

		assetMap[asset.Symbol] = asset
	}

	return assetMap, nil
}
