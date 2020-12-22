/*

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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

func init() {
	SchemeBuilder.Register(&ClusterBom{}, &ClusterBomList{})
}

// ClusterBomSpec defines the desired state of ClusterBom
type ClusterBomSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name of the secret which contains the target environment data
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=^[0-9a-zA-Z][0-9a-zA-Z_.\-]*[0-9a-zA-Z]$
	SecretRef string `json:"secretRef,omitempty"`

	ApplicationConfigs []ApplicationConfig `json:"applicationConfigs,omitempty"`

	AutoDelete *AutoDelete `json:"autoDelete,omitempty"`
}

// ClusterBomStatus defines the observed state of ClusterBom
type ClusterBomStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ApplicationStates []ApplicationState `json:"applicationStates,omitempty"`
	// +kubebuilder:validation:Enum=failed;pending;ok;unknown
	OverallState       string                `json:"overallState,omitempty"`
	OverallTime        metav1.Time           `json:"overallTime,omitempty"`
	Description        string                `json:"description,omitempty"`
	Conditions         []ClusterBomCondition `json:"conditions,omitempty"`
	ObservedGeneration int64                 `json:"observedGeneration,omitempty"`

	OverallNumOfDeployments      int `json:"overallNumOfDeployments,omitempty"`
	OverallNumOfReadyDeployments int `json:"overallNumOfReadyDeployments,omitempty"`
	OverallProgress              int `json:"overallProgress,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterBom is the Schema for the clusterboms API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CLUSTER",type=string,JSONPath=`.spec.secretRef`
// +kubebuilder:printcolumn:name="OVERALL STATUS",type=string,JSONPath=`.status.overallState`
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type ClusterBom struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterBomSpec   `json:"spec,omitempty"`
	Status ClusterBomStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterBomList contains a list of ClusterBom
type ClusterBomList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterBom `json:"items"`
}

type AutoDelete struct {
	ClusterBomAge int64 `json:"clusterBomAge,omitempty"`
}

// ClusterBomCondition describes the state of a clusterbom at a certain point.
type ClusterBomCondition struct {
	// Type of clusterbom condition; currently only Ready is supported.
	Type ClusterBomConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason ClusterBomConditionReason `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type ClusterBomConditionType string

// These are valid conditions of a clusterbom.
const (
	ClusterBomReady  ClusterBomConditionType = "Ready"
	ClusterReachable ClusterBomConditionType = "ClusterReachable"
)

type ClusterBomConditionReason string

const (

	// Reasons for Ready condition
	ReasonTargetClusterDoesNotExist ClusterBomConditionReason = "TargetClusterDoesNotExist"
	ReasonEmptyClusterBom           ClusterBomConditionReason = "EmptyClusterBom"
	ReasonFailedApps                ClusterBomConditionReason = "FailedApps"
	ReasonFailedAndPendingApps      ClusterBomConditionReason = "FailedAndPendingApps"
	ReasonPendingApps               ClusterBomConditionReason = "PendingApps"
	ReasonAllAppsReady              ClusterBomConditionReason = "AllAppsReady"
	ReasonClusterBomModified        ClusterBomConditionReason = "ReasonClusterBomModified"

	// Reasons for ClusterReachable condition
	ReasonClusterReachable           ClusterBomConditionReason = "ReasonClusterReachable"
	ReasonClusterNotReachable        ClusterBomConditionReason = "ReasonClusterNotReachable"
	ReasonClusterDoesNotExist        ClusterBomConditionReason = "ReasonClusterDoesNotExist"
	ReasonClusterReachabilityUnknown ClusterBomConditionReason = "ReasonClusterReachabilityUnknown"
)
