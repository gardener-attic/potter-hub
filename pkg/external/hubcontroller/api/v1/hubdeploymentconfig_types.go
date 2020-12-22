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
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

func init() {
	SchemeBuilder.Register(&HubDeploymentConfig{}, &HubDeploymentConfigList{})
}

// HubDeploymentConfigSpec defines the desired state of HubDeploymentConfig
type HubDeploymentConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name
	LocalSecretRef string `json:"localSecretRef,omitempty"`

	CurrentOperation CurrentOperation `json:"currentOperation,omitempty"`

	DeploymentConfig DeploymentConfig `json:"hubDeploymentConfig,omitempty"`
}

// HubDeploymentConfigStatus defines the observed state of HubDeploymentConfig
type HubDeploymentConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase              string                         `json:"phase,omitempty"`
	LastOperation      LastOperation                  `json:"lastOperation,omitempty"`
	Reachability       *Reachability                  `json:"reachability,omitempty"`
	Readiness          *Readiness                     `json:"readiness,omitempty"`
	Conditions         []HubDeploymentConfigCondition `json:"conditions,omitempty"`
	TypeSpecificStatus *runtime.RawExtension          `json:"typeSpecificStatus,omitempty"`
}

// +kubebuilder:object:root=true

// HubDeploymentConfig is the Schema for the hubdeploymentconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CONFIG TYPE",type=string,JSONPath=`.spec.hubDeploymentConfig.configType`
// +kubebuilder:printcolumn:name="CURRENT OPERATION NUMBER",type=string,JSONPath=`.spec.currentOperation.number`
// +kubebuilder:printcolumn:name="LAST OPERATION NUMBER",type=string,JSONPath=`.status.lastOperation.number`
// +kubebuilder:printcolumn:name="LAST OPERATION",type=string,JSONPath=`.status.lastOperation.operation`
// +kubebuilder:printcolumn:name="LAST OPERATION STATE",type=string,JSONPath=`.status.lastOperation.state`
// +kubebuilder:printcolumn:name="LAST OPERATION TIME",type=date,JSONPath=`.status.lastOperation.time`
// +kubebuilder:printcolumn:name="LAST OPERATION DESCRIPTION",type=string,JSONPath=`.status.lastOperation.description`
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type HubDeploymentConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HubDeploymentConfigSpec   `json:"spec,omitempty"`
	Status HubDeploymentConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HubDeploymentConfigList contains a list of HubDeploymentConfig
type HubDeploymentConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HubDeploymentConfig `json:"items"`
}

// HubDeploymentConfigCondition describes the state of a hubdeploymentconfig at a certain point.
type HubDeploymentConfigCondition struct {
	// Type of hubdeploymentconfig condition.
	Type HubDeploymentConfigConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason HubDeploymentConfigConditionReason `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type HubDeploymentConfigConditionType string

// These are valid conditions of a hubdeploymentconfig.
const (
	HubDeploymentConfigReady HubDeploymentConfigConditionType = "Ready"
)

type HubDeploymentConfigConditionReason string

const (
	ReasonClusterUnreachable   HubDeploymentConfigConditionReason = "ClusterUnreachable"
	ReasonInitialState         HubDeploymentConfigConditionReason = "InitialState"
	ReasonUpgradePending       HubDeploymentConfigConditionReason = "UpgradePending"
	ReasonRemovePending        HubDeploymentConfigConditionReason = "RemovePending"
	ReasonRunning              HubDeploymentConfigConditionReason = "Running"
	ReasonRemoved              HubDeploymentConfigConditionReason = "Removed"
	ReasonNotRunning           HubDeploymentConfigConditionReason = "NotRunning"
	ReasonFinallyFailed        HubDeploymentConfigConditionReason = "FinallyFailed"
	ReasonNotCurrentGeneration HubDeploymentConfigConditionReason = "NotCurrentGeneration"
)
