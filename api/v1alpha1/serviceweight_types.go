package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WeightedBackend defines the backend service with weight
type WeightedBackend struct {
	// Name of the backend service
	Name string `json:"name"`
	// Weight for traffic distribution
	Weight int32 `json:"weight"`
	// Port to forward traffic
	Port int32 `json:"port"`
}

// ServiceWeightSpec defines the desired state of ServiceWeight
type ServiceWeightSpec struct {
	// Inherit all fields from core ServiceSpec
	corev1.ServiceSpec `json:",inline"`

	// WeightedBackends defines the list of backend services with weights
	WeightedBackends []WeightedBackend `json:"weightedBackends,omitempty"`
}

// ServiceWeightStatus defines the observed state of ServiceWeight
type ServiceWeightStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw;svcw

// ServiceWeight is the Schema for the serviceweights API
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Selector",type="string",JSONPath=".spec.selector"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"

type ServiceWeight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceWeightSpec   `json:"spec,omitempty"`
	Status ServiceWeightStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceWeightList contains a list of ServiceWeight

type ServiceWeightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceWeight `json:"items"`
}

// DeepCopyObject implements runtime.Object
func (in *ServiceWeight) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyObject implements runtime.Object
func (in *ServiceWeightList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func init() {
	SchemeBuilder.Register(&ServiceWeight{}, &ServiceWeightList{})
}
