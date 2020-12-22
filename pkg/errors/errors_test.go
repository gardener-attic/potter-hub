package errors

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/pkg/errors"
)

func TestGetHttpErrorType(t *testing.T) {
	tests := []errTest{
		{
			// no http error
			0,
			false,
			errors.New("Dummy error"),
		},
		{
			// bad request
			BadRequest,
			true,
			BadRequest.New(errors.New("Dummy error")),
		},
		{
			// unauthorized
			Unauthorized,
			true,
			errors.Wrap(Unauthorized.New(errors.New("Inner")), "Outer"),
		},
		{
			// forbidden
			Forbidden,
			true,
			errors.Wrap(Forbidden.New(errors.Wrap(errors.New("Inner"), "Outer")), "Outer2"),
		},
		{
			// not found
			NotFound,
			true,
			NotFound.New(errors.Wrap(errors.New("Inner"), "Outer")),
		},
	}

	for _, tt := range tests {
		run(t, tt)
	}
}
func run(t *testing.T, tt errTest) {
	httpErrorType, isHTTPError := GetHTTPErrorType(tt.err)
	assert.Equal(t, isHTTPError, tt.expectedIsHTTPError, "isHTTPError")
	assert.Equal(t, httpErrorType, tt.expectedHTTPErrorType, "httpErrorType")
}

type errTest struct {
	expectedHTTPErrorType HTTPErrorType
	expectedIsHTTPError   bool
	err                   error
}
