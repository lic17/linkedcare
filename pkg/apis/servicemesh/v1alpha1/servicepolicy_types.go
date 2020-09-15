package v1alpha1

import (
	"istio.io/api/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServicePolicyConditionType string

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=servicepolicy

type ServicePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServicePolicySpec   `json:"spec"`
	Status ServicePolicyStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=servicepolicys

type ServicePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ServicePolicy `json:"items"`
}

type ServicePolicySpec struct {
	// Label selector for destination rules.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Template used to create a destination rule
	// +optional
	Template DestinationRuleSpecTemplate `json:"template,omitempty"`
}

type DestinationRuleSpecTemplate struct {

	// Metadata of the virtual services created from this template
	// +optional
	metav1.ObjectMeta

	// Spec indicates the behavior of a destination rule.
	// +optional
	Spec v1beta1.DestinationRule `json:"spec,omitempty"`
}

type ServicePolicyStatus struct {
	AppliedConfig string `json:"appliedConfig,omitempty"`
	//Conditions represent the latest available observations of an object's current state:
	//More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#typical-status-properties
	Conditions []ServicePolicyCondition `json:"conditions,omitempty"`
}

type ServicePolicyCondition struct {
	// Type of servicepolicy condition.
	Type ServicePolicyConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition
	Message string `json:"message,omitempty"`
}
