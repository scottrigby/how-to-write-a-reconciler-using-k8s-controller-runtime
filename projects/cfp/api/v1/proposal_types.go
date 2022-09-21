/*
Copyright 2022.

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

// ProposalSpec defines the desired state of Proposal
type ProposalSpec struct {
	// Title of the proposal
	// +kubebuilder:validation:MaxLength=50
	// +kubebuilder:validation:MinLength=1
	// +required
	Title string `json:"title"`

	// Abstract on what the proposal is about
	// +kubebuilder:validation:MaxLength=50
	// +kubebuilder:validation:MinLength=1
	// +required
	Abstract string `json:"abstract"`

	// Type of talk the proposal is on.
	// +kubebuilder:validation:Enum=talk;tutorial;keynote;lightning
	// +kubebuilder:default=talk
	Type string `json:"type"`

	// +required
	Draft bool `json:"draft"`

	// speaker submitting this talk
	// +required
	SpeakerRef *SpeakerRef `json:"speakerRef"`
}

type SpeakerRef struct {
	// Name of speaker custom resource
	// +kubebuilder:validation:Type=string
	Name string `json:"name"`

	// Namespace of speaker ref
	// +kubebuilder:validation:Type=string
	Namespace string `json:"namespace"`
}

type SubmissionStatus struct {
	// Status represents the current status of the proposal
	// It should move from draft to pending and then either accepted or rejected.
	// +kubebuilder:validation:Enum=draft;pending;accepted;rejected
	Status string `json:"status"`
}

// ProposalStatus defines the observed state of Proposal
type ProposalStatus struct {
	// ObservedGeneration is the last observed generation of the Speaker object.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// The time at which the proposal was submitted
	Timestamp metav1.Timestamp `json:"timestamp:omitempty"`

	SubmissionStatus SubmissionStatus `json:"submissionStatus:omitempty"`

	Conditions []metav1.Condition `json:"conditions:omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// Proposal is the Schema for the proposals API
type Proposal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProposalSpec   `json:"spec,omitempty"`
	Status ProposalStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProposalList contains a list of Proposal
type ProposalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Proposal `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Proposal{}, &ProposalList{})
}
