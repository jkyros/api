package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfiguration provides information to configure an operator to manage openshift-machine-config-operator.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec is the specification of the desired behavior of the OpenShift API Server.
	// +kubebuilder:validation:Required
	// +required
	Spec MachineConfigurationSpec `json:"spec"`

	// status defines the observed status of the OpenShift API Server.
	// +optional
	Status MachineConfigurationStatus `json:"status"`
}

type MachineConfigurationSpec struct {
	OperatorSpec `json:",inline"`
}

type MachineConfigurationStatus struct {
	OperatorStatus `json:",inline"`

	// latestAvailableRevision is the latest revision used as suffix of revisioned
	// secrets like encryption-config. A new revision causes a new deployment of
	// pods.
	// +optional
	// +kubebuilder:validation:Minimum=0
	LatestAvailableRevision int32 `json:"latestAvailableRevision,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigurationList is a collection of items
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items contains the items
	Items []MachineConfiguration `json:"items"`
}
