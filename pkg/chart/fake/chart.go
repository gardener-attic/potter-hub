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
	"encoding/json"
	"net/http"

	"helm.sh/helm/v3/pkg/chart"

	chartUtils "github.wdf.sap.corp/kubernetes/hub/pkg/chart"
)

type Chart struct{}

func (f *Chart) ParseDetails(data []byte) (*chartUtils.Details, error) {
	details := &chartUtils.Details{}
	err := json.Unmarshal(data, details)
	return details, err
}

func (f *Chart) GetChart(details *chartUtils.Details, netClient chartUtils.HTTPClient) (*chart.Chart, error) {
	valuesMap := make(map[string]interface{})
	valuesMap["values"] = details.Values

	return &chart.Chart{
		Metadata: &chart.Metadata{
			Name: details.ChartName,
		},
		Values: valuesMap,
	}, nil
}

func (f *Chart) InitNetClient(ctx context.Context, details *chartUtils.Details) (chartUtils.HTTPClient, error) {
	return &http.Client{}, nil
}
