package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	BadRequest          = HTTPErrorType(400)
	Unauthorized        = HTTPErrorType(401)
	Forbidden           = HTTPErrorType(403)
	NotFound            = HTTPErrorType(404)
	Conflict            = HTTPErrorType(409)
	UnprocessableEntity = HTTPErrorType(422)
	InternalServerError = HTTPErrorType(500)
)

type HTTPErrorType uint

type HTTPError struct {
	httpErrorType HTTPErrorType
	cause         error
}

func (e HTTPError) Cause() error {
	return e.cause
}

func (e HTTPError) HTTPErrorType() HTTPErrorType {
	return e.httpErrorType
}

func (e HTTPError) Error() string {
	return e.cause.Error()
}

func (e HTTPError) StackTrace() errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	stackTrace := errors.StackTrace{}
	if stackErr, ok := e.cause.(stackTracer); ok {
		stackTrace = stackErr.StackTrace()
	}
	return stackTrace
}

func (et HTTPErrorType) New(err error) error {
	return &HTTPError{et, err}
}

func (et HTTPErrorType) NewError(cause string) error {
	return &HTTPError{et, errors.New(cause)}
}

func (et HTTPErrorType) NewErrorf(format string, args ...interface{}) error {
	return &HTTPError{et, errors.New(fmt.Sprintf(format, args...))}
}

func GetHTTPErrorType(err error) (HTTPErrorType, bool) {
	type causer interface {
		Cause() error
	}

	type httperror interface {
		HTTPErrorType() HTTPErrorType
	}

	var httpErrorType HTTPErrorType
	isHTTPError := false

	for err != nil {
		httpError, ok := err.(httperror)
		if ok {
			httpErrorType = httpError.HTTPErrorType()
			isHTTPError = true
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return httpErrorType, isHTTPError
}
