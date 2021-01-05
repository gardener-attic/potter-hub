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

package fake

import (
	"context"
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"

	"github.com/gardener/potter-hub/pkg/proxy"
)

type Proxy struct {
	Releases []release.Release
}

func (f *Proxy) GetReleaseStatus(ctx context.Context, namespace, relName string, vo proxy.ValidationObject) (release.Status, error) {
	return release.StatusDeployed, nil
}

func (f *Proxy) ResolveManifest(ctx context.Context, namespace, values string, ch *chart.Chart, vo proxy.ValidationObject) (string, error) {
	return "", nil
}

func (f *Proxy) ResolveManifestFromRelease(ctx context.Context, namespace, releaseName string, revision int32, vo proxy.ValidationObject) (string, error) {
	return "", nil
}

func (f *Proxy) ListReleases(ctx context.Context, namespace string, releaseListLimit int, status string, vo proxy.ValidationObject) ([]proxy.AppOverview, error) {
	res := []proxy.AppOverview{}
	for _, r := range f.Releases {
		relStatus := "DEPLOYED" // Default
		if r.Info != nil {
			relStatus = r.Info.Status.String()
		}
		if (namespace == "" || namespace == r.Namespace) &&
			len(res) <= releaseListLimit &&
			(r.Info == nil || strings.EqualFold(status, relStatus)) {
			res = append(res, proxy.AppOverview{
				ReleaseName: r.Name,
				Version:     "",
				Namespace:   r.Namespace,
				Icon:        "",
				Status:      relStatus,
			})
		}
	}
	return res, nil
}

func (f *Proxy) CreateRelease(ctx context.Context, name, namespace, values string, ch *chart.Chart, vo proxy.ValidationObject) (*release.Release, error) {
	for _, r := range f.Releases {
		if r.Name == name {
			return nil, fmt.Errorf("release already exists")
		}
	}
	r := release.Release{
		Name:      name,
		Namespace: namespace,
	}
	f.Releases = append(f.Releases, r)
	return &r, nil
}

func (f *Proxy) UpdateRelease(ctx context.Context, name, namespace, values string, ch *chart.Chart, vo proxy.ValidationObject) (*release.Release, error) {
	for _, r := range f.Releases {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("release %s not found", name)
}

func (f *Proxy) RollbackRelease(ctx context.Context, name, namespace string, revision int32, vo proxy.ValidationObject) (*release.Release, error) {
	for _, r := range f.Releases {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("release %s not found", name)
}

func (f *Proxy) GetRelease(ctx context.Context, name, namespace string, vo proxy.ValidationObject) (*release.Release, error) {
	for _, r := range f.Releases {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("release %s not found", name)
}

func (f *Proxy) DeleteRelease(ctx context.Context, name, namespace string, keepHistory bool, vo proxy.ValidationObject) error {
	for i, r := range f.Releases {
		if r.Name == name {
			if !keepHistory {
				f.Releases[i] = f.Releases[len(f.Releases)-1]
				f.Releases = f.Releases[:len(f.Releases)-1]
			} else {
				r.Info = &release.Info{
					Status: release.StatusUninstalled,
				}
				f.Releases[i] = r
			}
			return nil
		}
	}
	return fmt.Errorf("release %s not found", name)
}
