package domain

import "fmt"

type FieldNotPositiveError struct {
	Value     any
	FieldName string
}

func (e *FieldNotPositiveError) Error() string {
	return fmt.Sprintf("field %s must be positive, got: %v", e.FieldName, e.Value)
}

type NoSuchOptionError struct {
	Value     any
	FieldName string
}

func (e NoSuchOptionError) Error() string {
	return fmt.Sprintf("field %s has no such option as %v", e.FieldName, e.Value)
}
