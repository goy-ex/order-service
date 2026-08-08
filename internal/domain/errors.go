package domain

import "fmt"

type NotPositiveNumberError struct {
	Value     any
	FieldName string
}

func (e *NotPositiveNumberError) Error() string {
	return fmt.Sprintf("field %s must be positive, got: %v", e.FieldName, e.Value)
}

type NegativeNumberError struct {
	Value     any
	FieldName string
}

func (e *NegativeNumberError) Error() string {
	return fmt.Sprintf("field %s must be not negative, got: %v", e.FieldName, e.Value)
}

type NoSuchOptionError struct {
	Value     any
	FieldName string
}

func (e NoSuchOptionError) Error() string {
	return fmt.Sprintf("field %s has no such option as %v", e.FieldName, e.Value)
}
