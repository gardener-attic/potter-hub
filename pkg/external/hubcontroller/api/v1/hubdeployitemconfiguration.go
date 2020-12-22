package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=true
type HubDeployItemConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	LocalSecretRef    string           `json:"localSecretRef,omitempty"`
	DeploymentConfig  DeploymentConfig `json:"hubDeploymentConfig,omitempty"`
}

// +kubebuilder:object:root=true
type HubDeployItemProviderStatus struct {
	metav1.TypeMeta    `json:",inline"`
	metav1.ObjectMeta  `json:"metadata,omitempty"`
	LastOperation      LastOperation         `json:"lastOperation,omitempty"`
	Reachability       *Reachability         `json:"reachability,omitempty"`
	Readiness          *Readiness            `json:"readiness,omitempty"`
	TypeSpecificStatus *runtime.RawExtension `json:"typeSpecificStatus,omitempty"`
}
