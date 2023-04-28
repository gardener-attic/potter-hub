module github.com/gardener/potter-hub

go 1.16

replace (
	github.com/Masterminds/squirrel => github.com/Masterminds/squirrel v1.5.3
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.13
	github.com/docker/distribution => github.com/distribution/distribution v2.7.1+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.5+incompatible
	github.com/emicklei/go-restful v2.9.5+incompatible => github.com/emicklei/go-restful v2.16.0+incompatible
)

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/Microsoft/hcsshim v0.8.14 // indirect
	github.com/arschles/assert v2.0.0+incompatible
	github.com/bshuster-repo/logrus-logstash-hook v1.0.2 // indirect
	github.com/bugsnag/bugsnag-go v2.1.2+incompatible // indirect
	github.com/bugsnag/panicwrap v1.3.4 // indirect
	github.com/containerd/continuity v0.0.0-20201208142359-180525291bb7 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-test/deep v1.0.7
	github.com/gofrs/uuid v4.1.0+incompatible // indirect
	github.com/golang/glog v1.0.0
	github.com/google/go-cmp v0.5.9
	github.com/googleapis/gnostic v0.5.5
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kubeapps/common v0.0.0-20200304064434-f6ba82e79f47
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	github.com/unrolled/render v1.4.0 // indirect
	github.com/urfave/negroni v1.0.1-0.20200608235619-7de0dfc1ff79
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/gorelic v0.0.7 // indirect
	github.com/yvasiyarov/newrelic_platform_go v0.0.0-20160601141957-9c099fbc30e9 // indirect
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b
	google.golang.org/grpc v1.49.0
	gopkg.in/yaml.v2 v2.4.0

	// If you update helm you need to update the kubernetes libs as well
	helm.sh/helm/v3 v3.11.1
	k8s.io/api v0.26.0
	k8s.io/apimachinery v0.26.0
	k8s.io/cli-runtime v0.26.0
	k8s.io/client-go v0.26.0
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/yaml v1.3.0
)
