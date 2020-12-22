package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	appRepoClientSet "github.wdf.sap.corp/kubernetes/hub/cmd/apprepository-controller/pkg/client/clientset/versioned"
	errUtils "github.wdf.sap.corp/kubernetes/hub/pkg/errors"
	"github.wdf.sap.corp/kubernetes/hub/pkg/log"
	"github.wdf.sap.corp/kubernetes/hub/pkg/util"
)

func NewAppRepositoryHandler(config *rest.Config) (*AppRepositoryHandler, error) {
	appRepoClient, err := appRepoClientSet.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	obj := &AppRepositoryHandler{
		client: appRepoClient,
	}
	return obj, nil
}

type AppRepositoryHandler struct {
	client appRepoClientSet.Interface
}

func (h *AppRepositoryHandler) ListAppRepositories(w http.ResponseWriter, r *http.Request) {
	podNamespace := util.GetPodNamespace()
	appRepoList, err := h.client.KubeappsV1alpha1().AppRepositories(podNamespace).List(metav1.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "cannot list apprepositories")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}
	response, err := json.Marshal(appRepoList)
	if err != nil {
		err = errors.Wrap(err, "cannot marshal response")
		util.SendErrResponse(r.Context(), w, errUtils.InternalServerError.New(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		err = errors.Wrap(err, "cannot write response body")
		log.GetLogger(r.Context()).Error(err)
	}
}

func (h *AppRepositoryHandler) GetAppRepository(w http.ResponseWriter, r *http.Request, params Params) {
	podNamespace := util.GetPodNamespace()
	appRepo, err := h.client.KubeappsV1alpha1().AppRepositories(podNamespace).Get(params["name"], metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, "cannot get apprepository")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}
	response, err := json.Marshal(appRepo)
	if err != nil {
		err = errors.Wrap(err, "cannot marshal response")
		util.SendErrResponse(r.Context(), w, errUtils.InternalServerError.New(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		err = errors.Wrap(err, "cannot write response body")
		log.GetLogger(r.Context()).Error(err)
	}
}
