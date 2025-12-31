// Package v1alpha1 包含 v1alpha1 API 组的 API 架构定义
// +kubebuilder:object:generate=true
// +groupName=example.com
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion 是用于注册这些对象的组版本
	GroupVersion = schema.GroupVersion{Group: "example.com", Version: "v1alpha1"}

	// SchemeBuilder 用于将 Go 类型添加到 GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme 将本组版本中的类型添加到给定的 scheme
	AddToScheme = SchemeBuilder.AddToScheme
)
