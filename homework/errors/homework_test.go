package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	var stringBuilder strings.Builder
	if len(e.errors) == 0 {
		return "No errors occurred"
	}
	stringBuilder.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errors)))
	for _, err := range e.errors {
		stringBuilder.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}
	stringBuilder.WriteString("\n")
	return stringBuilder.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError
	if !errors.As(err, &multiErr) {
		multiErr = &MultiError{}
	}
	multiErr.errors = append(multiErr.errors, errs...)
	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
