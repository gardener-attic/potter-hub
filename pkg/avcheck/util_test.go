package avcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	logUtils "github.wdf.sap.corp/kubernetes/hub/pkg/log"
)

func TestCheckPing(t *testing.T) {
	tests := []struct {
		name            string
		httpReturnCode  int
		checkSuccessful bool
	}{
		{
			name:            "successful",
			httpReturnCode:  http.StatusOK,
			checkSuccessful: true,
		},
		{
			name:            "fails due to invalid status code",
			httpReturnCode:  http.StatusInternalServerError,
			checkSuccessful: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.httpReturnCode)
		}))

		nullLogger, _ := test.NewNullLogger()
		ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})
		pingResult := checkPing(ctx, ts.URL)

		ts.Close()

		assert.Equal(t, pingResult.CheckSuccessful, tt.checkSuccessful, "checkSuccessful")
	}
}
