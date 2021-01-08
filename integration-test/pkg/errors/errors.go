package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	NoType   = ErrorType(iota)
	NotFound = ErrorType(404)
)

type ErrorType uint

type CustomError struct {
	ErrorType     ErrorType
	OriginalError error
	ContextInfo   map[string]string
}

func (err CustomError) Error() string {
	return err.OriginalError.Error()
}
func (et ErrorType) New(msg string) error {
	return CustomError{ErrorType: et, OriginalError: errors.New(msg)}
}

func (et ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)

	return CustomError{ErrorType: et, OriginalError: err}
}

func (et ErrorType) Wrap(err error, msg string) error {
	return et.Wrapf(err, msg)
}

func (et ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return CustomError{ErrorType: et, OriginalError: newErr}
}
