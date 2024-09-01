/*
Copyright 2024.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DuoPodSpec defines the desired state of DuoPod
type DuoPodSpec struct {
	// Container is the container specification for the pods
	Container corev1.Container `json:"container"`
}

// DuoPodStatus defines the observed state of DuoPod
type DuoPodStatus struct {
	// PodNames are the names of the created pods
	PodNames []string `json:"podNames"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DuoPod is the Schema for the duopods API
type DuoPod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DuoPodSpec   `json:"spec,omitempty"`
	Status DuoPodStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DuoPodList contains a list of DuoPod
type DuoPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DuoPod `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DuoPod{}, &DuoPodList{})
}
