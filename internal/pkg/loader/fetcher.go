package loader

import (
	"fmt"
	"io"
	"os"
)

type FileFetcher struct {
	filename string
}

func NewFileFetcher(filename string) *FileFetcher {
	return &FileFetcher{filename: filename}
}

func (f *FileFetcher) Fetch() (io.ReadCloser, error) {
	file, err := os.Open(f.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
