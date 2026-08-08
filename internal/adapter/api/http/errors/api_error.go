package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/goy-ex/sentinel"
)

type APIError struct {
	Message string
	Status  int
}

func ToAPIError(err error) *APIError {
	switch {
	case errors.Is(err, sentinel.ErrBadRequest):
		return badRequestToAPIError(err)
	default:
		return &APIError{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
}

func WriteAPIError(w http.ResponseWriter, apiErr *APIError) error {
	buf := bytes.NewBuffer(make([]byte, 0))
	body := map[string]any{
		"error": apiErr.Message,
	}

	err := json.NewEncoder(buf).Encode(&body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return fmt.Errorf("failed to encode error: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(apiErr.Status)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write error: %w", err)
	}

	return nil
}
