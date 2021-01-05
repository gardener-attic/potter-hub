/*
Copyright (c) 2018 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	_ "expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"github.com/urfave/negroni"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	appRepo "github.com/gardener/potter-hub/cmd/apprepository-controller/pkg/client/clientset/versioned"
	"github.com/gardener/potter-hub/cmd/ui-backend/internal/handler"
	"github.com/gardener/potter-hub/pkg/avcheck"
	chartUtils "github.com/gardener/potter-hub/pkg/chart"
	logUtils "github.com/gardener/potter-hub/pkg/log"
	helmProxy "github.com/gardener/potter-hub/pkg/proxy"
	"github.com/gardener/potter-hub/pkg/util"
)

func main() {
	disableAuth := pflag.Bool("disable-auth", false, "Disable authorization check")
	listLimit := pflag.Int("list-max", 256, "maximum number of releases to fetch")
	userAgentComment := pflag.String("user-agent-comment", "", "UserAgent comment used during outbound requests")
	version := pflag.String("version", "devel", "UserAgent version used during outbound requests")

	hostURL := pflag.String("host-url", "", "URL of the current host address")
	oidcCA := pflag.String("oidc-cluster-ca", "", "CA of the oidc cluster which contains kubeconfig information")
	oidcClusterURL := pflag.String("oidc-cluster-url", "", "URL of the cluster which contains the kubeconfig information")
	pflag.Parse()

	decodedClusterCAData, decodeErr := base64.StdEncoding.DecodeString(*oidcCA)
	if decodeErr != nil {
		logUtils.StandardLogger().Fatalf("Unable to decode oidc cluster CA: %v", decodeErr)
	}

	var authGate negroni.HandlerFunc
	if *oidcClusterURL != "" {
		authGate = handler.KubeconfigAuthorization(*oidcClusterURL, decodedClusterCAData)
		*disableAuth = true
	} else {
		authGate = handler.TokenAuthorization()
	}

	bomHandler := &handler.BomHandler{
		OidcClusterURL: oidcClusterURL,
		OidcClusterCA:  &decodedClusterCAData,
		ClientFactory:  handler.K8sClientFromConfig,
	}

	hp := initHelmProxy(disableAuth, userAgentComment, version, listLimit)
	k8sReverseProxy := handler.NewK8sReverseProxy(*oidcClusterURL, *hostURL, decodedClusterCAData)
	appRepoHandler := initAppRepoHandler()
	systemInfoHandler := initSystemInfoHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Setup avcheck
	avcheckConfig := parseAVCheckConfig()
	helmProxyChecker := avcheck.NewUIBackendChecker("http://localhost" + addr)
	chartServiceChecker := avcheck.NewChartServiceChecker(util.GetEnvOrPanic("CHARTSVC_URL"), avcheckConfig)
	dashboardChecker := avcheck.NewDashboardChecker(util.GetEnvOrPanic("DASHBOARD_URL"))
	go chartServiceChecker.StartChartsAvailableCheckBackgroundJob()

	// Setup routes
	r := mux.NewRouter()
	addHelmProxyRoutes(r, hp, authGate)
	addClusterBomRoutes(r, bomHandler)
	addAvailabilityRoutes(r, avcheckConfig.PathPrefix, helmProxyChecker, chartServiceChecker, dashboardChecker)
	addAppRepoRoutes(r, appRepoHandler)
	addSystemInfoRoutes(r, systemInfoHandler)
	addK8sReverseProxyRoutes(r, k8sReverseProxy)

	n := negroni.New()
	n.UseHandler(r)

	srv := &http.Server{
		Addr:    addr,
		Handler: n,
	}

	go func() {
		logUtils.StandardLogger().WithField("addr", addr).Info("Started Helm Proxy")
		err := srv.ListenAndServe()
		if err != nil {
			logUtils.StandardLogger().Info(err)
		}
	}()

	// Catch SIGINT and SIGTERM
	// Set up channel on which to send signal notifications.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	logUtils.StandardLogger().Debug("Set system to get notified on signals")
	s := <-c
	logUtils.StandardLogger().Infof("Received signal: %v. Waiting for existing requests to finish", s)

	err := shutdownServer(srv)
	if err != nil {
		logUtils.StandardLogger().Fatalf("HTTP server shutdown failed: %+v", err)
	}

	logUtils.StandardLogger().Info("All requests have been served. Exiting")
}

func shutdownServer(server *http.Server) error {
	// Set a timeout value high enough to let k8s terminationGracePeriodSeconds to act
	// accordingly and send a SIGKILL if needed
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Returns the user agent to be used during calls to the chart repositories
// Examples:
// tiller-proxy/devel
// chart-repo/1.0
// tiller-proxy/1.0 (monocular v1.0-beta4)
// More info here https://github.com/kubeapps/kubeapps/issues/767#issuecomment-436835938
func userAgent(userAgentComment, version string) string {
	if version == "" {
		version = "devel"
	}
	ua := "tiller-proxy/" + version
	if userAgentComment != "" {
		ua = fmt.Sprintf("%s (%s)", ua, userAgentComment)
	}
	return ua
}

func parseAVCheckConfig() *avcheck.Configuration {
	avCheckJSONConfig := os.Getenv("AVAILABILITY_CHECK")
	if avCheckJSONConfig == "" {
		logUtils.StandardLogger().Fatal("the environment variable AVAILABILITY_CHECK is not set")
	}

	var config avcheck.Configuration
	err := json.Unmarshal([]byte(avCheckJSONConfig), &config)
	if err != nil {
		logUtils.StandardLogger().Fatalf("cannot unmarshal availability check JSON config: %s", err.Error())
	}

	err = config.Validate()
	if err != nil {
		logUtils.StandardLogger().Fatalf("invalid availability check JSON config: %s", err.Error())
	}

	return &config
}

func addClusterBomRoutes(r *mux.Router, bomHandler *handler.BomHandler) {
	apiv1 := r.PathPrefix("/{clusterNamespace}/{accessData}/v1").Subrouter()

	apiv1.Methods("GET").Path("/clusterboms").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithParams(bomHandler.ListClusterBoms)),
	))

	apiv1.Methods("GET").Path("/clusterboms/{clusterBomName}").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithParams(bomHandler.GetClusterBom)),
	))

	apiv1.Methods("PUT").Path("/clusterboms/{clusterBomName}").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithParams(bomHandler.UpdateClusterBom)),
	))
}

func addHelmProxyRoutes(r *mux.Router, hp *handler.HelmProxy, authGate negroni.HandlerFunc) {
	apiv1 := r.PathPrefix("/{clusterNamespace}/{accessData}/helm/v1").Subrouter()

	apiv1.Methods("GET").Path("/releases").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		authGate,
		negroni.Wrap(handler.WithoutParams(hp.ListAllReleases)),
	))

	apiv1.Methods("GET").Path("/namespaces/{namespace}/releases").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		authGate,
		negroni.Wrap(handler.WithParams(hp.ListReleases)),
	))

	apiv1.Methods("POST").Path("/namespaces/{namespace}/releases").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		authGate,
		negroni.Wrap(handler.WithParams(hp.CreateRelease)),
	))

	apiv1.Methods("GET").Path("/namespaces/{namespace}/releases/{releaseName}").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		authGate,
		negroni.Wrap(handler.WithParams(hp.GetRelease)),
	))

	apiv1.Methods("PUT").Path("/namespaces/{namespace}/releases/{releaseName}").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		authGate,
		negroni.Wrap(handler.WithParams(hp.OperateRelease)),
	))

	apiv1.Methods("DELETE").Path("/namespaces/{namespace}/releases/{releaseName}").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		authGate,
		negroni.Wrap(handler.WithParams(hp.DeleteRelease)),
	))
}

func addAppRepoRoutes(r *mux.Router, appRepoHandler *handler.AppRepositoryHandler) {
	r.Path("/apprepositories").Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithoutParams(appRepoHandler.ListAppRepositories)),
	))

	r.Path("/apprepositories/{name}").Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithParams(appRepoHandler.GetAppRepository)),
	))
}

func addSystemInfoRoutes(r *mux.Router, systemInfoHandler *handler.SystemInfoHandler) {
	r.Path("/controller-version").Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithoutParams(systemInfoHandler.GetControllerVersion)),
	))
}

func addK8sReverseProxyRoutes(r *mux.Router, k8sReverseProxy *handler.K8sReverseProxy) {
	r.PathPrefix("/{clusterNamespace}/{accessData}/reverse-proxy/v1/resourcecluster").
		Methods("GET").
		Handler(negroni.New(
			negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
			negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
			negroni.Wrap(handler.WithParams(k8sReverseProxy.ProxyRequestToResourceCluster)),
		))

	targetClusterPrefix := r.PathPrefix("/{clusterNamespace}/{accessData}/reverse-proxy/v1").Subrouter()

	targetClusterPrefix.Path("/api/v1/namespaces").Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithParams(k8sReverseProxy.ProxyRequestToTargetCluster)),
	))

	targetClusterPrefix.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
		negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
		negroni.Wrap(handler.WithParams(k8sReverseProxy.ProxyRequestToTargetCluster)),
	))
}

func addAvailabilityRoutes(r *mux.Router, pathPrefix string, uiBackendChecker avcheck.IUIBackendChecker, chartServiceChecker avcheck.IChartServiceChecker, dashboardChecker avcheck.IDashboardChecker) {
	r.Methods("GET").Path("/live").Handler(
		negroni.New(
			negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
			negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
			negroni.Wrap(handler.WithoutParams(avcheck.HandlePing)),
		),
	)

	r.Methods("GET").Path("/ready").Handler(
		negroni.New(
			negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
			negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
			negroni.Wrap(handler.WithoutParams(avcheck.HandlePing)),
		),
	)

	// TODO: delete once pipeline avcheck and avs are migrated to the new endpoint
	r.Methods("GET").Path("/{clusterNamespace}/{accessData}/availability").Handler(
		negroni.New(
			negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
			negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
			negroni.Wrap(handler.WithoutParams(avcheck.CreateAVCheckHandler(uiBackendChecker, chartServiceChecker, dashboardChecker))),
		),
	)

	avcheckRoutes := r.PathPrefix(pathPrefix).Subrouter()
	avcheckRoutes.Methods("GET").Path("/").Handler(
		negroni.New(
			negroni.HandlerFunc(logUtils.PrepareLoggerHandler),
			negroni.HandlerFunc(logUtils.RequestResponseLogHandler),
			negroni.Wrap(handler.WithoutParams(avcheck.CreateAVCheckHandler(uiBackendChecker, chartServiceChecker, dashboardChecker))),
		),
	)
}

func getHubClusterConfig() (*rest.Config, bool) {
	var config *rest.Config
	var err error
	isRemoteClusterConfig := false

	remoteClusterKubeconfig := os.Getenv("KUBECONFIG")
	if remoteClusterKubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", remoteClusterKubeconfig)
		isRemoteClusterConfig = true
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		logUtils.StandardLogger().Fatalf("Unable to get cluster config: %v", err)
	}

	return config, isRemoteClusterConfig
}

func initHelmProxy(disableAuth *bool, userAgentComment, version *string, listLimit *int) *handler.HelmProxy {
	var config *rest.Config
	var err error

	config, isRemoteClusterConfig := getHubClusterConfig()
	if isRemoteClusterConfig {
		*disableAuth = true
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logUtils.StandardLogger().Fatalf("Unable to create a kubernetes client: %v", err)
	}

	appRepoClient, err := appRepo.NewForConfig(config)
	if err != nil {
		logUtils.StandardLogger().Fatalf("Unable to create an app repository client: %v", err)
	}

	chartClient := chartUtils.NewClient(kubeClient, appRepoClient, loader.LoadArchive, userAgent(*userAgentComment, *version))

	return &handler.HelmProxy{
		DisableAuth: *disableAuth,
		ListLimit:   *listLimit,
		ChartClient: chartClient,
		ProxyClient: &helmProxy.Proxy{},
	}
}

func initAppRepoHandler() *handler.AppRepositoryHandler {
	config, _ := getHubClusterConfig()
	obj, err := handler.NewAppRepositoryHandler(config)
	if err != nil {
		logUtils.StandardLogger().Fatalf("Unable to create apprepository handler: %v", err)
	}
	return obj
}

func initSystemInfoHandler() *handler.SystemInfoHandler {
	config, _ := getHubClusterConfig()
	obj, err := handler.NewSystemInfoHandler(config)
	if err != nil {
		logUtils.StandardLogger().Fatalf("Unable to create system info handler: %v", err)
	}
	return obj
}
