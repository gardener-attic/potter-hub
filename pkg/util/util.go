package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kubeapps/common/response"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	errorUtils "github.com/gardener/potter-hub/pkg/errors"
	logUtils "github.com/gardener/potter-hub/pkg/log"
)

const (
	StrTrue                           = "true"
	defaultPodNamespace               = metav1.NamespaceSystem
	defaultControllerNamespace        = "potter-controller"
	ControllerSystemInfoConfigMapName = "system-info"
)

func SendErrResponse(ctx context.Context, w http.ResponseWriter, err error) {
	logUtils.GetLogger(ctx).Error(err)
	httpErrorType, isHTTPError := errorUtils.GetHTTPErrorType(err)
	if isHTTPError {
		response.NewErrorResponse(int(httpErrorType), err.Error()).Write(w)
	} else {
		response.NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
	}
}

func GetEnvOrPanic(envKey string) string {
	env := os.Getenv(envKey)
	if env == "" {
		panic(fmt.Sprintf("Environment variable %s was not found", envKey))
	}
	return env
}

func AreAllItemsTrue(list []bool) bool {
	areAllItemsTrue := true
	for _, item := range list {
		if !item {
			areAllItemsTrue = false
			break
		}
	}
	return areAllItemsTrue
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	header := strings.Split(authHeader, " ")
	if len(header) != 2 {
		err := errors.New("Invalid authorization header")
		return "", err
	}
	if header[0] != "Bearer" {
		err := errors.New("Only Bearer access allowed")
		return "", err
	}
	token := header[1]
	if token == "" {
		err := errors.New("Token missing in authorization header")
		return "", err
	}

	return token, nil
}

func CheckAndSendK8sError(ctx context.Context, w http.ResponseWriter, err error) {
	logUtils.GetLogger(ctx).Error(err)

	if apierrors.IsBadRequest(errors.Cause(err)) {
		response.NewErrorResponse(http.StatusBadRequest, err.Error()).Write(w)
	} else if apierrors.IsUnauthorized(errors.Cause(err)) {
		response.NewErrorResponse(http.StatusUnauthorized, err.Error()).Write(w)
	} else if apierrors.IsForbidden(errors.Cause(err)) {
		response.NewErrorResponse(http.StatusForbidden, err.Error()).Write(w)
	} else if apierrors.IsNotFound(errors.Cause(err)) {
		response.NewErrorResponse(http.StatusNotFound, err.Error()).Write(w)
	} else if apierrors.IsConflict(errors.Cause(err)) {
		response.NewErrorResponse(http.StatusConflict, err.Error()).Write(w)
	} else {
		response.NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
	}
}

func GetPodNamespace() string {
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = defaultPodNamespace
	}

	return namespace
}

func GetControllerNamespace() string {
	namespace := os.Getenv("CONTROLLER_NAMESPACE")
	if namespace == "" {
		namespace = defaultControllerNamespace
	}

	return namespace
}

// DecodeBasicAuthCredentials Decodes basic auth credential string and returns username, password or an error
func DecodeBasicAuthCredentials(base64EncodedBasicAuthCredentials string) (string, string, error) {
	decodedCredentials, err := base64.StdEncoding.DecodeString(base64EncodedBasicAuthCredentials)
	if err != nil {
		return "", "", errors.Wrap(err, "Couldn't decode basic auth credentials")
	}
	splittedCredentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(splittedCredentials) < 2 {
		return "", "", errors.New("Password missing in credential string. Could not split by colon ':'")
	}

	username := splittedCredentials[0]
	password := splittedCredentials[1]
	return username, password, nil
}

// GetGCloudAccessToken receives and returns an access token from a google service account json string. Returns the access token or error.
func GetGCloudAccessToken(gcloudServiceAccountJSON string) (string, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(gcloudServiceAccountJSON), "https://www.googleapis.com/auth/devstorage.read_only")
	if err != nil {
		return "", errors.Wrap(err, "Couldn't create Google Service Account object")
	}
	tokenSource := jwtConfig.TokenSource(context.TODO())
	token, err := tokenSource.Token()
	if err != nil {
		return "", errors.Wrap(err, "Couldn't fetch token from token source")
	}
	return token.AccessToken, nil
}
