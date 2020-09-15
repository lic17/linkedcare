package v1alpha1

import (
	"istio.io/api/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StrategyConditionType string

type StrategyType string

const (
	// Canary strategy type
	CanaryType StrategyType = "Canary"

	// BlueGreen strategy type
	BlueGreenType StrategyType = "BlueGreen"

	// Mirror strategy type
	Mirror StrategyType = "Mirror"
)

type StrategyPolicy string

const (
	// apply strategy only until workload is ready
	PolicyWaitForWorkloadReady StrategyPolicy = "WaitForWorkloadReady"

	// apply strategy immediately no matter workload status is
	PolicyImmediately StrategyPolicy = "Immediately"

	// pause strategy
	PolicyPause StrategyPolicy = "Paused"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=strategy

type Strategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StrategySpec   `json:"spec"`
	Status StrategyStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=strategys

type StrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Strategy `json:"items"`
}

type StrategySpec struct {
	// Strategy type
	Type StrategyType `json:"type,omitempty"`

	// Principal version, the one as reference version
	// label version value
	// +optional
	PrincipalVersion string `json:"principal,omitempty"`

	// Governor version, the version takes control of all incoming traffic
	// label version value
	// +optional
	GovernorVersion string `json:"governor,omitempty"`

	// Label selector for virtual services.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Template describes the virtual service that will be created.
	Template VirtualServiceTemplateSpec `json:"template,omitempty"`

	// strategy policy, how the strategy will be applied
	// by the strategy controller
	StrategyPolicy StrategyPolicy `json:"strategyPolicy,omitempty"`
}

// VirtualServiceTemplateSpec
type VirtualServiceTemplateSpec struct {

	// Metadata of the virtual services created from this template
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec indicates the behavior of a virtual service.
	// +optional
	Spec v1beta1.VirtualService `json:"spec,omitempty"`
}

type StrategyStatus struct {
	AppliedConfig string `json:"appliedConfig,omitempty"`
	//Conditions represent the latest available observations of an object's current state:
	//More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#typical-status-properties
	Conditions []StrategyCondition `json:"conditions,omitempty"`
}

type StrategyCondition struct {
	// Type of strategy condition.
	Type StrategyConditionType `json:"type"`
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
