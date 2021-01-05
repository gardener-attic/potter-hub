package util

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	errorUtils "github.com/gardener/potter-hub/pkg/errors"
	logUtils "github.com/gardener/potter-hub/pkg/log"
)

func TestSendErrResponse(t *testing.T) {
	tests := []struct {
		name            string
		expectedCode    int
		expectedMessage string
		err             error
	}{
		{
			"error without explicit status code",
			500,
			`{"code":500,"message":"Dummy error"}`,
			errors.New("Dummy error"),
		},
		{
			"bad request",
			400,
			`{"code":400,"message":"Outer: Inner"}`,
			errors.Wrap(errorUtils.BadRequest.New(errors.New("Inner")), "Outer"),
		},
		{
			"unauthorized",
			401,
			`{"code":401,"message":"Dummy message"}`,
			errorUtils.Unauthorized.New(errors.New("Dummy message")),
		},
		{
			"forbidden",
			403,
			`{"code":403,"message":"Outer2: Outer: Inner"}`,
			errors.Wrap(errorUtils.Forbidden.New(errors.Wrap(errors.New("Inner"), "Outer")), "Outer2"),
		},
		{
			"not found",
			404,
			`{"code":404,"message":"Outer: Inner"}`,
			errorUtils.NotFound.New(errors.Wrap(errors.New("Inner"), "Outer")),
		},
		{
			"conflict",
			409,
			`{"code":409,"message":"Dummy error"}`,
			errorUtils.Conflict.New(errors.New("Dummy error")),
		},
		{
			"unprocessable entity",
			422,
			`{"code":422,"message":"Dummy error"}`,
			errorUtils.UnprocessableEntity.New(errors.New("Dummy error")),
		},
		{
			"internal server error",
			500,
			`{"code":500,"message":"Dummy error"}`,
			errorUtils.InternalServerError.New(errors.New("Dummy error")),
		},
	}

	nullLogger, _ := test.NewNullLogger()
	ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			SendErrResponse(ctx, responseRecorder, tt.err)

			assert.Equal(t, responseRecorder.Code, tt.expectedCode, "code")
			assert.Equal(t, responseRecorder.Body.String(), tt.expectedMessage, "message")
		})
	}
}

func TestSendErrResponseWithNoLoggerInContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// Code under test
	err := errors.New("Dummy error")
	ctx := context.TODO()
	responseRecorder := httptest.NewRecorder()
	SendErrResponse(ctx, responseRecorder, err)
}

func TestGetTokenFromRequest(t *testing.T) {
	tests := []struct {
		name                string
		authorizationHeader string
		expectedErrMsg      string
		expectedToken       string
	}{
		{
			name:                "Authorization header set to empty string",
			authorizationHeader: "",
			expectedErrMsg:      "Invalid authorization header",
			expectedToken:       "",
		},
		{
			name:                "Invalid authorization type",
			authorizationHeader: "Basic YWRtaW46cGFzc3dvcmQ=",
			expectedErrMsg:      "Only Bearer access allowed",
			expectedToken:       "",
		},
		{
			name:                "Invalid header format",
			authorizationHeader: "BearerPleaseProvideASpaceCharacter",
			expectedErrMsg:      "Invalid authorization header",
			expectedToken:       "",
		},
		{
			name:                "Token equals empty string",
			authorizationHeader: "Bearer ",
			expectedErrMsg:      "Token missing in authorization header",
			expectedToken:       "",
		},
		{
			name:                "Token equals empty string",
			authorizationHeader: "Bearer TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCAuLi4=",
			expectedErrMsg:      "",
			expectedToken:       "TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCAuLi4=",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/test", nil)
			request.Header.Add("Authorization", tt.authorizationHeader)

			token, err := GetTokenFromRequest(request)

			assert.Equal(t, token, tt.expectedToken, "token")
			if tt.expectedErrMsg != "" {
				assert.ExistsErr(t, err, "error")
				assert.Equal(t, err.Error(), tt.expectedErrMsg, "errMsg")
			} else {
				assert.NoErr(t, err)
			}
		})
	}
}

func Test_DecodeBasicAuthCredentials_Successfully(t *testing.T) {
	credentials := "aHVtYW46c2VjcmV0"
	expectedUsername := "human"
	expectedPassword := "secret"

	actualUsername, actualPassword, err := DecodeBasicAuthCredentials(credentials)

	if err != nil {
		t.Errorf("Unexpected error when Decoding basic authentication %v", err)
	}
	if actualUsername != expectedUsername {
		t.Errorf("Expecting %s to be resolved as %s", actualUsername, expectedUsername)
	}
	if actualPassword != expectedPassword {
		t.Errorf("Expecting %s to be resolved as %s", actualPassword, expectedPassword)
	}
}

func Test_DecodeBasicAuthCredentials_WithError(t *testing.T) {
	credentials := "a"

	_, _, err := DecodeBasicAuthCredentials(credentials)

	if err == nil {
		t.Errorf("Expected decoding error but no error was thrown")
	}
}

func Test_DecodeBasicAuthCredentials_WithSplitError(t *testing.T) {
	credentials := "aHVtYW4="

	_, _, err := DecodeBasicAuthCredentials(credentials)

	if err == nil {
		t.Errorf("Expected an error when input credentials do not contain a colon")
	}
}
