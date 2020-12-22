module github.wdf.sap.corp/kubernetes/hub

go 1.15

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.1.1+incompatible
	github.com/docker/docker => github.com/moby/moby v17.12.1-ce+incompatible
	github.com/docker/docker@v0.7.3-0.20190327010347-be7ac8be2ae0 => github.com/moby/moby v17.12.1-ce+incompatible
	github.com/docker/docker@v1.4.2-0.20200203170920-46ec8731fbce => github.com/moby/moby v17.12.1-ce+incompatible
	// replace needed for hub-controller dependency
	github.com/moby/moby@v0.7.3-0.20190826074503-38ab9da00309 => github.com/moby/moby v17.12.1-ce+incompatible
)

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/arschles/assert v1.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-test/deep v1.0.7
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-cmp v0.5.3
	github.com/googleapis/gnostic v0.5.3
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kubeapps/common v0.0.0-20200304064434-f6ba82e79f47
	github.com/opencontainers/selinux v1.6.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/unrolled/render v1.0.3 // indirect
	github.com/urfave/negroni v1.0.1-0.20200608235619-7de0dfc1ff79
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	google.golang.org/grpc v1.33.2
	gopkg.in/yaml.v2 v2.3.0
	// If you update helm you need to update the kubernetes libs as well
	helm.sh/helm/v3 v3.4.1
	k8s.io/api v0.19.5
	k8s.io/apimachinery v0.19.5
	k8s.io/cli-runtime v0.19.5
	k8s.io/client-go v0.19.5
	rsc.io/letsencrypt v0.0.3 // indirect

	// add controller runtime kubernetes client
	sigs.k8s.io/controller-runtime v0.6.4
	sigs.k8s.io/yaml v1.2.0
)
