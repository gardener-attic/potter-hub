package handler

import (
	"net/http"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gardener/potter-hub/pkg/log"
	"github.com/gardener/potter-hub/pkg/util"
)

func NewSystemInfoHandler(config *rest.Config) (*SystemInfoHandler, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	obj := &SystemInfoHandler{
		client: client,
	}
	return obj, nil
}

type SystemInfoHandler struct {
	client *kubernetes.Clientset
}

func (h *SystemInfoHandler) GetControllerVersion(w http.ResponseWriter, r *http.Request) {
	controllerNamespace := util.GetControllerNamespace()

	cm, err := h.client.CoreV1().ConfigMaps(controllerNamespace).Get(r.Context(), util.ControllerSystemInfoConfigMapName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, "cannot get configmap")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}

	controllerVersion := cm.Data["appVersion"]
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(controllerVersion))
	if err != nil {
		err = errors.Wrap(err, "cannot write response body")
		log.GetLogger(r.Context()).Error(err)
	}
}
