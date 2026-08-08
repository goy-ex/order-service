package usecase

import "errors"

// PairNotFoundError indicates that no trading pair is registered under the
// given key.
type PairNotFoundError struct {
	// PairKey is the key that was looked up.
	PairKey string
}

// Error implements the error interface.
func (e *PairNotFoundError) Error() string {
	return "pair '" + e.PairKey + "' not found"
}

var ErrNilInput = errors.New("usecase input is nil")
