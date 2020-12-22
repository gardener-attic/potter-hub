package chart

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

type KubeappsRelease release.Release

type KubeappsChart struct {
	Metadata *chart.Metadata `json:"metadata,omitempty"`
	// LocK is the contents of Chart.lock.
	Lock *chart.Lock `json:"lock,omitempty"`
	// Templates for this chart.
	Templates []*chart.File `json:"templates,omitempty"`
	// Values are default config for this template.
	Values map[string]interface{} `json:"values,omitempty"`
	// Schema is an optional JSON schema for imposing structure on Values
	Schema []byte `json:"schema,omitempty"`
	// Files are miscellaneous files in a chart archive,
	// e.g. README, LICENSE, etc.
	Files []*chart.File `json:"files,omitempty"`

	parent       *chart.Chart
	dependencies []*chart.Chart
}

type KReleaseInfo struct {
	// FirstDeployed is when the release was first deployed.
	FirstDeployed *time.Time `json:"first_deployed,omitempty"`
	// LastDeployed is when the release was last deployed.
	LastDeployed *time.Time `json:"last_deployed,omitempty"`
	// Deleted tracks when this object was deleted.
	Deleted *time.Time `json:"deleted,omitempty"`
	// Description is human-friendly "log entry" about this release.
	Description string `json:"Description,omitempty"`
	// Status is the current state of the release
	Status release.Status `json:"status,omitempty"`
	// Cluster resources as kubectl would print them.
	Resources string `json:"resources,omitempty"`
	// Contains the rendered templates/NOTES.txt if available
	Notes string `json:"notes,omitempty"`
}

type KRelease struct {
	Name string `json:"name,omitempty"`
	// Info provides information about a release
	Info *KReleaseInfo `json:"info,omitempty"`
	// Chart is the chart that was released.
	Chart *KubeappsChart `json:"chart,omitempty"`
	// Config is the set of extra Values added to the chart.
	// These values override the default values inside of the chart.
	Config KConfig `json:"config,omitempty"`
	// Manifest is the string representation of the rendered template.
	Manifest string `json:"manifest,omitempty"`
	// Hooks are all of the hooks declared for this release.
	Hooks []*release.Hook `json:"hooks,omitempty"`
	// Version is an int which represents the version of the release.
	Version int `json:"version,omitempty"`
	// Namespace is the kubernetes namespace of the release.
	Namespace string `json:"namespace,omitempty"`
}

type KConfig struct {
	Raw string `json:"raw,omitempty"`
}

func (rel *KubeappsRelease) MarshalJSON() ([]byte, error) {
	var kChart KubeappsChart

	if rel.Chart == nil {
		kRelease := KRelease{
			Name:      rel.Name,
			Info:      nil,
			Chart:     nil,
			Manifest:  rel.Manifest,
			Hooks:     rel.Hooks,
			Version:   rel.Version,
			Namespace: rel.Namespace,
		}
		jsonBytes, err := json.Marshal(kRelease)

		return jsonBytes, err
	}

	kChart = KubeappsChart{
		Metadata:     rel.Chart.Metadata,
		Lock:         rel.Chart.Lock,
		Templates:    rel.Chart.Templates,
		Values:       rel.Chart.Values,
		Schema:       rel.Chart.Schema,
		Files:        rel.Chart.Files,
		parent:       rel.Chart.Parent(),
		dependencies: rel.Chart.Dependencies(),
	}

	rawConfig, err := yaml.Marshal(rel.Config)

	if err != nil {
		return nil, err
	}

	kReleaseInfo := KReleaseInfo{
		FirstDeployed: &rel.Info.FirstDeployed.Time,
		LastDeployed:  &rel.Info.LastDeployed.Time,
		Deleted:       &rel.Info.Deleted.Time,
		Description:   rel.Info.Description,
		Status:        rel.Info.Status,
		Notes:         rel.Info.Notes,
	}

	// Set Zero Time to nil pointer for nicer marshaling results
	if kReleaseInfo.FirstDeployed != nil && kReleaseInfo.FirstDeployed.IsZero() {
		kReleaseInfo.FirstDeployed = nil
	}
	if kReleaseInfo.LastDeployed != nil && kReleaseInfo.LastDeployed.IsZero() {
		kReleaseInfo.LastDeployed = nil
	}
	if kReleaseInfo.Deleted != nil && kReleaseInfo.Deleted.IsZero() {
		kReleaseInfo.Deleted = nil
	}

	kRelease := KRelease{
		Name:      rel.Name,
		Info:      &kReleaseInfo,
		Chart:     &kChart,
		Config:    KConfig{Raw: string(rawConfig)},
		Manifest:  rel.Manifest,
		Hooks:     rel.Hooks,
		Version:   rel.Version,
		Namespace: rel.Namespace,
	}

	jsonBytes, err := json.Marshal(kRelease)

	return jsonBytes, err
}
