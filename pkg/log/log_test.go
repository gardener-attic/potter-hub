package log

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
)

func TestPrepareLoggerHandler(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	req.Header.Set(requestIDHeader, "123456")

	// This function will be called by PrepareLoggerHandler as the next function in the middleware chain.
	// We use it for executing our assertions.
	next := func(rw http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r.Context())
		assert.Equal(t, logger.Data["request_id"], "123456", "logger must have expected request id")
	}

	PrepareLoggerHandler(responseRecorder, req, next)
}

func TestPrepareLoggerHandlerWithNoRequestIdHeader(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// This function will be called by PrepareLoggerHandler as the next function in the middleware chain.
	// We use it for executing our assertions.
	next := func(rw http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r.Context())
		assert.Equal(t, logger.Data["request_id"], "", "logger must have expected request id")
	}

	PrepareLoggerHandler(responseRecorder, req, next)
}

func TestGetLoggerPanicsWithEmptyContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	GetLogger(context.TODO())
}
