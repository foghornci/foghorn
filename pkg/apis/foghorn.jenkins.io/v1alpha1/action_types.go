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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActionType defines what type of action should be taken
type ActionType string

const (
	// CreatePipelineRun indicates that a Tekton PipelineRun should be created
	// in response to the webhook that triggered this action.
	CreatePipelineRun ActionType = "CreatePipelineRun"
	// CommentOnIssue indicates that we should leave a comment on a specified
	// issue on specified repo on a configured git provider.
	CommentOnIssue ActionType = "CommentOnIssue"
	// CommentOnPR indicates that we should comment on a pull request (or equivalent
	// contstruct) for a repo on a configured git provider.
	CommentOnPR ActionType = "CommentOnPR"
	// LabelIssue indicates that we should label an issue on a repo on a configured
	// git provider.
	LabelIssue ActionType = "LabelIssue"
	// LabelPR indicates that we should label a PR on a repo on a configured git
	// provider
	LabelPR ActionType = "LabelPR"
	// MergePR indicates that a PR is ready to merge in response to an approve
	// or lgtm event.
	MergePR ActionType = "MergePR"
)

// ActionSpec defines the desired state of Action
type ActionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Type        ActionType `json:"type"`
	ParentEvent *GitEvent  `json:"parentEvent"`
}

// ActionStatus defines the observed state of Action
type ActionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +genclient

// Action is the Schema for the actions API
type Action struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActionSpec   `json:"spec,omitempty"`
	Status ActionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ActionList contains a list of Action
type ActionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Action `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Action{}, &ActionList{})
}
