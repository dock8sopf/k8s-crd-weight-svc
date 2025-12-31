package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// 编辑此文件！这是您拥有的脚手架代码！
// 注意：json tags 是必需的。您添加的任何新字段都必须有 json tags 才能序列化。

// WeightedBackend 定义了带权重的后端服务
type WeightedBackend struct {
	// 后端服务的名称
	Name string `json:"name"`
	// 用于流量分配的权重
	Weight int32 `json:"weight"`
	// 转发流量的端口
	Port int32 `json:"port"`
}

// ServiceWeightSpec 定义了 ServiceWeight 的期望状态
type ServiceWeightSpec struct {
	// 继承 core ServiceSpec 的所有字段
	corev1.ServiceSpec `json:",inline"`

	// WeightedBackends 定义了带权重的后端服务列表
	WeightedBackends []WeightedBackend `json:"weightedBackends,omitempty"`
}

// ServiceWeightStatus 定义了 ServiceWeight 的观察状态
type ServiceWeightStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - 定义集群的观察状态
	// Important: 修改此文件后运行 "make" 重新生成代码

	// Conditions 表示对象状态的最新可用观察结果
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw;svcw

// ServiceWeight 是 serviceweights API 的模式定义
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

// ServiceWeightList 包含一个 ServiceWeight 列表

type ServiceWeightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceWeight `json:"items"`
}

// DeepCopyObject 实现 runtime.Object 接口
func (in *ServiceWeight) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyObject 实现 runtime.Object 接口
func (in *ServiceWeightList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func init() {
	SchemeBuilder.Register(&ServiceWeight{}, &ServiceWeightList{})
}
