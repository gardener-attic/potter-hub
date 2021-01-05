package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"

	"github.com/gardener/potter-hub/pkg/util"
)

const authHeaderKey = "AUTHORIZATION_HEADER"

func main() {
	glog.Infoln("starting hub-chart-repo")

	glog.Infoln("starting authorization header preparation")
	prepareAuthorizationHeader()
	glog.Infoln("finished authorization header preparation")

	// Disable linting because we do not call this program in a context where os.Args is derived from user input
	// nolint:gosec
	cmd := exec.Command("./chart-repo", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	glog.Infoln("starting chart-repo")
	err := cmd.Run()
	if err != nil {
		glog.Fatalln("chart-repo finished with error: ", err)
	}
	glog.Infoln("chart-repo finished successfully")
	glog.Infoln("hub-chart-repo finished successfully")
}

func prepareAuthorizationHeader() {
	authorizationHeader := os.Getenv(authHeaderKey)
	if len(authorizationHeader) > 0 {
		splittedAuthorizationHeader := strings.Fields(authorizationHeader)
		if len(splittedAuthorizationHeader) != 2 {
			glog.Fatalln("invalid authorization header. expected \"<type> <credentials>\"")
		}

		if splittedAuthorizationHeader[0] == "Basic" {
			username, password, err := util.DecodeBasicAuthCredentials(splittedAuthorizationHeader[1])
			if err != nil {
				glog.Fatalln("Invalid authorization header. ", err)
			}

			if username == "_json_key" {
				glog.Infoln("Starting gcloud oauth flow")
				accessToken, err := util.GetGCloudAccessToken(password)
				if err != nil {
					glog.Fatal("Gcloud oauth flow failed. ", err)
				}
				os.Setenv(authHeaderKey, "Bearer "+accessToken)
				glog.Infoln("Successfully performed gcloud oauth flow and set access token in authorization header")
			}
		}
	} else {
		glog.Infoln("Authorization header not set")
	}
}
