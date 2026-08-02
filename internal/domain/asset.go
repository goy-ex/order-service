package domain

type Asset struct {
	Symbol   string
	Name     string
	Decimals int
}

func NewAsset(symbol string, name string, decimals int) (*Asset, error) {
	if decimals < 0 {
		return nil, &FieldNotPositiveError{FieldName: "Decimals", Value: decimals}
	}

	return &Asset{
		Symbol:   symbol,
		Name:     name,
		Decimals: decimals,
	}, nil
}
