package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/goy-ex/order-service/internal/domain"
)

func badRequestToAPIError(err error) *APIError {
	switch {
	case errors.Is(err, ErrDecodeFail):
		return decodeFailToAPIError(err)
	case errors.Is(err, ErrValidationFail):
		return validationFailToAPIError(err)
	}

	if e, ok := errors.AsType[*domain.NotPositiveNumberError](err); ok {
		return &APIError{Status: http.StatusBadRequest, Message: e.Error()}
	}

	if e, ok := errors.AsType[*domain.NoSuchOptionError](err); ok {
		return &APIError{Status: http.StatusBadRequest, Message: e.Error()}
	}

	return &APIError{Status: http.StatusBadRequest, Message: "Bad Request"}
}

func decodeFailToAPIError(err error) *APIError {
	if e, ok := errors.AsType[*json.SyntaxError](err); ok {
		message := "json syntax error near " + strconv.FormatInt(e.Offset, 10)

		return &APIError{Status: http.StatusBadRequest, Message: message}
	}

	if e, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		message := fmt.Sprintf("field %q expects type %s, but got JSON value of type %s", e.Field, e.Type, e.Value)

		return &APIError{Status: http.StatusBadRequest, Message: message}
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return &APIError{Status: http.StatusBadRequest, Message: io.ErrUnexpectedEOF.Error()}
	}

	return &APIError{Status: http.StatusBadRequest, Message: "Bad Request"}
}

func validationFailToAPIError(err error) *APIError {
	errs, ok := errors.AsType[validator.ValidationErrors](err)

	if !ok || len(errs) == 0 {
		return &APIError{Status: http.StatusBadRequest, Message: "Bad Request"}
	}

	fe := errs[0]
	message := fmt.Sprintf("field %q failed validation with tag %q", fe.Field(), fe.Tag())

	return &APIError{Status: http.StatusBadRequest, Message: message}
}
