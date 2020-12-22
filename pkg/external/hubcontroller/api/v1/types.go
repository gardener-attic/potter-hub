package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Job struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type ReadyRequirements struct {
	Jobs      []Job      `json:"jobs,omitempty"`
	Resources []Resource `json:"resources,omitempty"`
}

type Resource struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	APIVersion    string                 `json:"apiVersion"`
	Resource      string                 `json:"resource"`
	FieldPath     string                 `json:"fieldPath"`
	SuccessValues []runtime.RawExtension `json:"successValues,omitempty"`
}

// ApplicationConfig defines one application to be deployed on a cluster
type ApplicationConfig struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Pattern=^[0-9a-z]*$
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Pattern=^[0-9a-z]*$
	ConfigType string `json:"configType,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	TypeSpecificData runtime.RawExtension `json:"typeSpecificData,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *runtime.RawExtension `json:"values,omitempty"`

	SecretValues *SecretValues `json:"secretValues,omitempty"`

	NamedSecretValues map[string]NamedSecretValues `json:"namedSecretValues,omitempty"`

	NoReconcile bool `json:"noReconcile,omitempty"`

	ReadyRequirements ReadyRequirements `json:"readyRequirements,omitempty"`
}

type SecretValues struct {
	InternalSecretName string `json:"internalSecretName,omitempty"`
	// +kubebuilder:validation:Enum=replace;keep;delete
	Operation string `json:"operation,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	Data *runtime.RawExtension `json:"data,omitempty"`
}

type NamedSecretValues struct {
	InternalSecretName string `json:"internalSecretName,omitempty"`
	// +kubebuilder:validation:Enum=delete
	Operation string `json:"operation,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	StringData map[string]string `json:"data,omitempty"`
}

// DeploymentConfig defines the deployment of one application
type DeploymentConfig struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Pattern=^[0-9a-z]*$
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Pattern=^[0-9a-z]*$
	ConfigType string `json:"configType,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	TypeSpecificData runtime.RawExtension `json:"typeSpecificData,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	Values *runtime.RawExtension `json:"values,omitempty"`

	InternalSecretName string `json:"internalSecretName,omitempty"`

	NamedInternalSecretNames map[string]string `json:"namedInternalSecretNames,omitempty"`

	NoReconcile       bool              `json:"noReconcile,omitempty"`
	ReconcileTime     metav1.Time       `json:"reconcileTime,omitempty"`
	ReadyRequirements ReadyRequirements `json:"readyRequirements,omitempty"`
}

// ApplicationState describes the state of the deployment of an application
type ApplicationState struct {
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:Enum=failed;pending;ok;unknown
	State         string        `json:"state,omitempty"`
	DetailedState DetailedState `json:"detailedState,omitempty"`
}

type Reachability struct {
	Reachable bool        `json:"reachable"`
	Time      metav1.Time `json:"time,omitempty"`
}

type Readiness struct {
	// +kubebuilder:validation:Enum=failed;pending;ok;unknown;notRelevant;finallyFailed
	State string      `json:"state,omitempty"`
	Time  metav1.Time `json:"time,omitempty"`
}

type ErrorEntry struct {
	Description string `json:"description,omitempty"`
	// +kubebulder:validation:Format="date-time"
	Time metav1.Time `json:"time,omitempty"`
}

type ErrorHistory struct {
	ErrorEntries []ErrorEntry `json:"errorEntries,omitempty"`
}

// DetailedState describes the detailed state of a deployment of an application
type DetailedState struct {
	Generation         int64                          `json:"generation,omitempty"`
	DeletionTimestamp  *metav1.Time                   `json:"deletionTimestamp,omitempty"`
	CurrentOperation   CurrentOperation               `json:"currentOperation,omitempty"`
	LastOperation      LastOperation                  `json:"lastOperation,omitempty"`
	Reachability       *Reachability                  `json:"reachability,omitempty"`
	Readiness          *Readiness                     `json:"readiness,omitempty"`
	HdcConditions      []HubDeploymentConfigCondition `json:"hdcConditions,omitempty"`
	TypeSpecificStatus *runtime.RawExtension          `json:"typeSpecificStatus,omitempty"`
}

// CurrentOperation defines the current deployment operation
type CurrentOperation struct {
	// not used anymore
	// +kubebuilder:validation:Enum=install;remove
	Operation string `json:"operation,omitempty"`

	// not used anymore
	Number int32 `json:"number,omitempty"`

	// +kubebulder:validation:Format="date-time"
	Time metav1.Time `json:"time,omitempty"`
}

// LastOperation describes the last deployment operation
type LastOperation struct {
	// +kubebuilder:validation:Enum=install;remove
	Operation string `json:"operation,omitempty"`

	// not used anymore
	Number int32 `json:"number,omitempty"`
	// not used anymore
	SuccessNumber int32 `json:"successNumber,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	SuccessGeneration  int64 `json:"successGeneration,omitempty"`

	// +kubebuilder:validation:Enum=failed;ok
	State         string `json:"state,omitempty"`
	NumberOfTries int32  `json:"numberOfTries,omitempty"`
	// +kubebulder:validation:Format="date-time"
	Time         metav1.Time   `json:"time,omitempty"`
	Description  string        `json:"description,omitempty"`
	ErrorHistory *ErrorHistory `json:"errorHistory,omitempty"`
}
