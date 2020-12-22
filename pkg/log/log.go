package log

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	errorUtils "github.wdf.sap.corp/kubernetes/hub/pkg/errors"
)

//nolint:gochecknoglobals
var std *Logger

//nolint:gochecknoinits
func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	loglevel := os.Getenv("LOG_LEVEL")

	if loglevel == "warning" {
		logrus.SetLevel(logrus.WarnLevel)
	} else if loglevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if loglevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	std = &Logger{Entry: logrus.NewEntry(logrus.StandardLogger())}
}

func StandardLogger() *Logger {
	return std
}

const (
	requestIDHeader    = "X-Request-ID"
	LogKeyResponseBody = "response-body"
	LogKeyLoggerName   = "logger"
)

type LoggerKey struct{}

type Logger struct {
	*logrus.Entry
}

// Utility function for logging errors.
// The utility this function provides is that it also logs the stacktrace of the error, if it exists.
func (logger *Logger) Error(err error) {
	var fields = logrus.Fields{}

	httpErrorType, isHTTPError := errorUtils.GetHTTPErrorType(err)
	if isHTTPError {
		fields["error_code"] = int(httpErrorType)
	}

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if stackErr, ok := err.(stackTracer); ok {
		fields["stacktrace"] = fmt.Sprintf("%+v", stackErr.StackTrace())
	} else {
		fields["stacktrace"] = ""
	}

	logger.WithFields(fields).Error(err)
}

func PrepareLoggerHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	reqID := r.Header.Get(requestIDHeader)
	entry := logrus.WithFields(logrus.Fields{
		"request_id": reqID,
	})
	logger := &Logger{entry}
	ctxt := context.WithValue(r.Context(), LoggerKey{}, logger)
	r = r.WithContext(ctxt)
	next(rw, r)
}

func GetLogger(ctx context.Context) *Logger {
	logger := ctx.Value(LoggerKey{})
	return logger.(*Logger)
}

func RequestResponseLogHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log := GetLogger(r.Context())

	entryLogger := log.WithField("request", fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, r.Proto))
	if r.URL.Path == "/ready" || r.URL.Path == "/live" {
		entryLogger.Debug("Process request")
	} else {
		entryLogger.Info("Process request")
	}

	start := time.Now()
	next(rw, r)
	duration := time.Since(start)

	// By default, negroni wraps the standard ResponseWriter with their own implementation.
	// https://github.com/urfave/negroni/blob/master/negroni.go#L106
	// The negroni.ResponseWriter allows us to access certain variables from the HTTP response, e.g. the HTTP status.
	// https://github.com/urfave/negroni/blob/master/response_writer.go
	negroniRw := rw.(negroni.ResponseWriter)

	exitLogger := log.WithFields(logrus.Fields{
		"duration":        duration,
		"status":          strconv.Itoa(negroniRw.Status()),
		"body_bytes_sent": strconv.Itoa(negroniRw.Size()),
	})

	if r.URL.Path == "/ready" || r.URL.Path == "/live" {
		exitLogger.Debug("Finished request")
	} else {
		exitLogger.Info("Finished request")
	}
}
