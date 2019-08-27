package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JSFunctionSpec defines the desired state of JSFunction
// +k8s:openapi-gen=true
type JSFunctionSpec struct {
	Func       string `json:"func"`
	Deployment int    `json:"deployment"`
	Package    string `json:"package,omitempty"`
	Events     bool   `json:"events,omitempty"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// JSFunctionStatus defines the observed state of JSFunction
// +k8s:openapi-gen=true
type JSFunctionStatus struct {
	Nodes []string `json:"nodes"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JSFunction is the Schema for the jsfunctions API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type JSFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JSFunctionSpec   `json:"spec,omitempty"`
	Status JSFunctionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JSFunctionList contains a list of JSFunction
type JSFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JSFunction `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JSFunction{}, &JSFunctionList{})
}
