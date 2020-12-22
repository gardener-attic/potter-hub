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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

func init() {
	SchemeBuilder.Register(&ClusterBomSync{}, &ClusterBomSyncList{})
}

// ClusterBomSyncSpec defines a block for a ClusterBom and its HDCs
type ClusterBomSyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ID        string      `json:"id,omitempty"`
	Timestamp metav1.Time `json:"timestamp,omitempty"`
	Until     metav1.Time `json:"until,omitempty"`
}

// ClusterBomSyncStatus defines the observed state of ClusterBomSync
type ClusterBomSyncStatus struct {
}

// +kubebuilder:object:root=true

// ClusterBomSync is the Schema for the clusterbomsync API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="UNTIL",type="string",JSONPath=".spec.until"
type ClusterBomSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterBomSyncSpec   `json:"spec,omitempty"`
	Status ClusterBomSyncStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterBomList contains a list of ClusterBom
type ClusterBomSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterBomSync `json:"items"`
}
