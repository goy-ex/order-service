package pair

import (
	"fmt"
	"io"
	"os"

	"github.com/goy-ex/order-service/internal/domain"
)

type fileProvider struct {
	parser   Parser
	filename string
}

func NewFileProvider(filename string, parser Parser) *fileProvider {
	return &fileProvider{
		filename: filename,
		parser:   parser,
	}
}

func (f *fileProvider) Provide() (map[domain.PairKey]*domain.Pair, error) {
	file, err := os.Open(f.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	pairs, err := f.parser.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return pairs, nil
}

type Parser interface {
	Parse(r io.Reader) (map[domain.PairKey]*domain.Pair, error)
}

type ParserFunc func(r io.Reader) (map[domain.PairKey]*domain.Pair, error)

func (f ParserFunc) Parse(r io.Reader) (map[domain.PairKey]*domain.Pair, error) {
	return f(r)
}
