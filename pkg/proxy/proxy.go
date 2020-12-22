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

package proxy

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	grpcStatus "google.golang.org/grpc/status"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	logUtils "github.wdf.sap.corp/kubernetes/hub/pkg/log"
)

// nolint
var (
	appMutex map[string]*sync.Mutex
)

// nolint
func init() {
	appMutex = make(map[string]*sync.Mutex)
}

// Proxy contains all the elements to contact Tiller and the K8s API
type Proxy struct {
}

// NewProxy creates a Proxy
func NewProxy() *Proxy {
	return &Proxy{}
}

// AppOverview represents the basics of a release
type AppOverview struct {
	ReleaseName   string         `json:"releaseName"`
	Description   string         `json:"description"`
	Version       string         `json:"version"`
	Namespace     string         `json:"namespace"`
	Icon          string         `json:"icon,omitempty"`
	Status        string         `json:"status"`
	Chart         string         `json:"chart"`
	ChartMetadata chart.Metadata `json:"chartMetadata"`
}

func (p *Proxy) getRelease(vo ValidationObject, name, namespace string) (*release.Release, error) {
	config := vo.initActionConfig(namespace)

	rls, err := action.NewGet(config).Run(name)

	if err != nil {
		return nil, errors.New(prettyError(err).Error())
	}

	// We check that the release found is from the provided namespace.
	// If `namespace` is an empty string we do not do that check
	// This check check is to prevent users of for example updating releases that might be
	// in namespaces that they do not have access to.
	if namespace != "" && rls.Namespace != namespace {
		return nil, errors.Errorf("Release %q not found in namespace %q", name, namespace)
	}

	return rls, err
}

// GetReleaseStatus prints the status of the given release if exists
func (p *Proxy) GetReleaseStatus(ctx context.Context, namespace, relName string, vo ValidationObject) (release.Status, error) {
	config := vo.initActionConfig(namespace)

	stats, err := action.NewStatus(config).Run(relName)

	if err == nil {
		if stats.Info != nil {
			return stats.Info.Status, nil
		}
	}
	return release.StatusUnknown, errors.Wrapf(err, "Unable to fetch release status for %s", relName)
}

// ResolveManifest returns a manifest given the chart parameters
func (p *Proxy) ResolveManifest(ctx context.Context, namespace, values string, ch *chart.Chart, vo ValidationObject) (string, error) {
	// We use the release returned after running a dry-run to know the elements to install

	config := vo.initActionConfig(namespace)

	install := action.NewInstall(config)
	install.DryRun = true
	install.ReleaseName = "dummyrlsname"
	install.Namespace = namespace
	// TODO: add values override (?)
	valuesMap := make(map[string]interface{})

	resDry, err := install.Run(ch, valuesMap)

	if err != nil {
		return "", errors.Wrap(err, "Could not run install dry run")
	}
	// The manifest returned has some extra new lines at the beginning
	return strings.TrimLeft(resDry.Manifest, "\n"), nil
}

// ResolveManifestFromRelease returns a manifest given the release name and revision
func (p *Proxy) ResolveManifestFromRelease(ctx context.Context, namespace, releaseName string, revision int32, vo ValidationObject) (string, error) {
	rel, err := p.GetRelease(ctx, releaseName, namespace, vo)

	if err != nil {
		return "", err
	}
	// The manifest returned has some extra new lines at the beginning
	return strings.TrimLeft(rel.Manifest, "\n"), nil
}

// Apply the same filtering than helm CLI
// Ref: https://github.com/helm/helm/blob/d3b69c1fc1ac62f1cc40f93fcd0cba275c0596de/cmd/helm/list.go#L173
func filterList(rels []*release.Release) []*release.Release {
	idx := map[string]int{}

	for _, r := range rels {
		name, version := r.Name, r.Version
		if max, ok := idx[name]; ok {
			// check if we have a greater version already
			if max > version {
				continue
			}
		}
		idx[name] = version
	}

	uniq := make([]*release.Release, 0, len(idx))
	for _, r := range rels {
		if idx[r.Name] == r.Version {
			uniq = append(uniq, r)
		}
	}
	return uniq
}

// ListReleases list releases in a specific namespace if given
func (p *Proxy) ListReleases(ctx context.Context, namespace string, releaseListLimit int, status string, vo ValidationObject) ([]AppOverview, error) {
	config := vo.initActionConfig(namespace)
	listCommand := action.NewList(config)
	if status == "all" {
		listCommand.All = true
		listCommand.SetStateMask()
	}

	releases, err := listCommand.Run()

	if err != nil {
		return []AppOverview{}, errors.Wrapf(err, "Unable to list helm releases")
	}
	appList := []AppOverview{}
	if releases != nil {
		filteredReleases := filterList(releases)
		for _, r := range filteredReleases {
			if namespace == "" || namespace == r.Namespace {
				appList = append(appList, AppOverview{
					ReleaseName:   r.Name,
					Description:   r.Info.Description,
					Version:       r.Chart.Metadata.Version,
					Namespace:     r.Namespace,
					Icon:          r.Chart.Metadata.Icon,
					Status:        r.Info.Status.String(),
					Chart:         r.Chart.Metadata.Name,
					ChartMetadata: *r.Chart.Metadata,
				})
			}
		}
	}
	return appList, nil
}

func lock(name string) {
	if appMutex[name] == nil {
		appMutex[name] = &sync.Mutex{}
	}
	appMutex[name].Lock()
}

func unlock(name string) {
	appMutex[name].Unlock()
}

// CreateRelease creates a tiller release
func (p *Proxy) CreateRelease(ctx context.Context, name, namespace, values string, ch *chart.Chart, vo ValidationObject) (*release.Release, error) {
	lock(name)
	defer unlock(name)

	log := logUtils.GetLogger(ctx)

	log.Printf("Installing release %s into namespace %s", name, namespace)

	err := ensureNamespace(ctx, namespace, vo)
	if err != nil {
		return nil, err
	}

	config := vo.initActionConfig(namespace)
	log.Printf("Got action config")

	install := action.NewInstall(config)
	install.Namespace = namespace
	install.ReleaseName = name

	valOpts, err := getValueMap(values)
	if err != nil {
		return nil, err
	}

	log.Printf("Installing chart %s", name)
	res, err := install.Run(ch, valOpts)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create the release")
	}

	hubsec, enabled := os.LookupEnv("HUBSEC_DOCKERCONFIGJSON")
	log.Printf("Secret configured: %t", enabled)
	if enabled {
		log.Printf("Secret configured: %s", hubsec)
		imageSecret := newImageSecret(ctx, res, hubsec, vo)
		err = imageSecret.createOrUpdateImageSecret(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to create image pull secret")
		}
	}

	log.Printf("%s successfully installed in %s", name, namespace)

	return res, err
}

func getValueMap(values string) (map[string]interface{}, error) {
	valuesMap := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(values), &valuesMap); err != nil {
		return valuesMap, errors.Wrapf(err, "Failed to parse values")
	}
	return valuesMap, nil
}

// UpdateRelease upgrades a tiller release
func (p *Proxy) UpdateRelease(ctx context.Context, name, namespace, values string, ch *chart.Chart, vo ValidationObject) (*release.Release, error) {
	lock(name)
	defer unlock(name)

	log := logUtils.GetLogger(ctx)

	// Check if the release already exists
	_, err := p.getRelease(vo, name, namespace)
	if err != nil {
		return nil, err
	}
	log.Printf("Updating release %s", name)

	config := vo.initActionConfig(namespace)

	upgrade := action.NewUpgrade(config)
	upgrade.Namespace = namespace
	upgrade.MaxHistory = 10

	valOpts, err := getValueMap(values)
	if err != nil {
		return nil, err
	}

	rel, err := upgrade.Run(name, ch, valOpts)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to update the release")
	}
	return rel, err
}

// RollbackRelease rolls back to a specific revision
func (p *Proxy) RollbackRelease(ctx context.Context, name, namespace string, revision int32, vo ValidationObject) (*release.Release, error) {
	lock(name)
	defer unlock(name)
	// Check if the release already exists
	config := vo.initActionConfig(namespace)

	_, err := p.getRelease(vo, name, namespace)
	if err != nil {
		return nil, err
	}

	rollbackCommand := action.NewRollback(config)
	rollbackCommand.Version = int(revision)
	err = rollbackCommand.Run(name)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to rollback the release")
	}
	return p.getRelease(vo, name, namespace)
}

// GetRelease returns the info of a release
func (p *Proxy) GetRelease(ctx context.Context, name, namespace string, vo ValidationObject) (*release.Release, error) {
	lock(name)
	defer unlock(name)
	return p.getRelease(vo, name, namespace)
}

// DeleteRelease deletes a release
func (p *Proxy) DeleteRelease(ctx context.Context, name, namespace string, keepHistory bool, vo ValidationObject) error {
	lock(name)
	defer unlock(name)

	log := logUtils.GetLogger(ctx)

	log.Printf("Deleting release %s in namespace %s", name, namespace)
	// Validate that the release actually belongs to the namespace
	_, err := p.getRelease(vo, name, namespace)
	if err != nil {
		return err
	}

	config := vo.initActionConfig(namespace)
	uninstall := action.NewUninstall(config)
	uninstall.KeepHistory = keepHistory

	rel, uninstallErr := uninstall.Run(name)
	if uninstallErr != nil {
		return errors.Wrap(uninstallErr, "Unable to delete the release")
	}

	hubsec, enabled := os.LookupEnv("HUBSEC_DOCKERCONFIGJSON")
	log.Printf("Secret configured: %t", enabled)
	if enabled {
		log.Printf("Secret configured: %s", hubsec)
		overviews, listErr := p.ListReleases(ctx, namespace, 0, "all", vo)
		if listErr != nil {
			return errors.Wrap(listErr, "Unable to list release to check if imagepullsecret has to be deleted")
		}
		if len(overviews) == 0 {
			imageSecrets := newImageSecret(ctx, rel.Release, hubsec, vo)
			err = imageSecrets.deleteImageSecret(ctx)
			if err != nil {
				return errors.Wrap(err, "Unable to delete image pull secret")
			}
		}
	}

	log.Printf("%s successfully deleted in %s", name, namespace)

	return err
}

// ensureNamespace make sure we create a namespace in the cluster in case it does not exist
func ensureNamespace(ctx context.Context, namespace string, vo ValidationObject) error {
	log := logUtils.GetLogger(ctx)

	log.Printf("Ensuring namespace %s exists", namespace)
	// The namespace might not exist yet, so we create the clientset with the default namespace
	clientset, err := vo.getClientSet("default")
	if err != nil {
		return errors.Wrapf(err, "Error creating kubernetes client for namespace %s.", namespace)
	}
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return nil
		}
		if k8sErrors.IsNotFound(err) {
			_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}, metav1.CreateOptions{})
			if err != nil {
				return errors.Wrapf(err, "Could not create namespace %s.", namespace)
			}
			return nil
		}
		return errors.Wrapf(err, "Unable to fetch namespace %s", namespace)
	}
	return nil
}

// extracted from https://github.com/helm/helm/blob/master/cmd/helm/helm.go#L227
// prettyError unwraps or rewrites certain errors to make them more user-friendly.
func prettyError(err error) error {
	// Add this check can prevent the object creation if err is nil.
	if err == nil {
		return nil
	}
	// If it's grpc's error, make it more user-friendly.
	if s, ok := grpcStatus.FromError(err); ok {
		return fmt.Errorf(s.Message())
	}
	// Else return the original error.
	return err
}

// TillerClient for exposed funcs
type TillerClient interface {
	GetReleaseStatus(ctx context.Context, namespace string, relName string, vo ValidationObject) (release.Status, error)
	ResolveManifest(ctx context.Context, namespace, values string, ch *chart.Chart, vo ValidationObject) (string, error)
	ResolveManifestFromRelease(ctx context.Context, namespace string, releaseName string, revision int32, vo ValidationObject) (string, error)
	ListReleases(ctx context.Context, namespace string, releaseListLimit int, status string, vo ValidationObject) ([]AppOverview, error)
	CreateRelease(ctx context.Context, name, namespace, values string, ch *chart.Chart, vo ValidationObject) (*release.Release, error)
	UpdateRelease(ctx context.Context, name, namespace string, values string, ch *chart.Chart, vo ValidationObject) (*release.Release, error)
	RollbackRelease(ctx context.Context, name, namespace string, revision int32, vo ValidationObject) (*release.Release, error)
	GetRelease(ctx context.Context, name, namespace string, vo ValidationObject) (*release.Release, error)
	DeleteRelease(ctx context.Context, name, namespace string, keepHistory bool, vo ValidationObject) error
}
