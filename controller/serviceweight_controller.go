package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8scrdserviceweight "k8s-crd-service-weight/api/v1alpha1"
)

// ServiceWeightReconciler 协调 ServiceWeight 对象
type ServiceWeightReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=example.com,resources=serviceweights,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=serviceweights/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=serviceweights/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile 是主要的 Kubernetes 协调循环的一部分，旨在将集群的当前状态移动到期望状态。
func (r *ServiceWeightReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 获取 ServiceWeight 实例
	serviceWeight := &k8scrdserviceweight.ServiceWeight{}
	err := r.Get(ctx, req.NamespacedName, serviceWeight)
	if err != nil {
		if errors.IsNotFound(err) {
			// 请求对象未找到，可能在协调请求后被删除
			// 拥有的对象会被自动垃圾回收
			logger.Info("ServiceWeight 资源未找到。由于对象必须被删除，忽略此问题。")
			return ctrl.Result{}, nil
		}
		// 读取对象出错 - 重新排队请求
		logger.Error(err, "获取 ServiceWeight 失败。")
		return ctrl.Result{}, err
	}

	// 创建或更新对应的 Service
	// 如果 WeightedBackends 为空，它就像一个普通的 Service
	if len(serviceWeight.Spec.WeightedBackends) == 0 {
		// 创建/更新标准服务
		return r.reconcileStandardService(ctx, serviceWeight)
	} else {
		// 创建/更新加权服务（需要实现实际的加权流量逻辑）
		return r.reconcileWeightedService(ctx, serviceWeight)
	}
}

// reconcileStandardService 在没有指定 WeightedBackends 时处理标准 Service 行为
func (r *ServiceWeightReconciler) reconcileStandardService(ctx context.Context, sw *k8scrdserviceweight.ServiceWeight) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 定义一个新的 Service 对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sw.Name,
			Namespace: sw.Namespace,
		},
		Spec: sw.Spec.ServiceSpec,
	}

	// 将 ServiceWeight 实例设置为 owner 和 controller
	if err := controllerutil.SetControllerReference(sw, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// 检查此 Service 是否已存在
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("创建新 Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
		err = r.Create(ctx, service)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Service 创建成功 - 返回并重新排队
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// 使用 ServiceWeight 中期望的 spec 更新找到的 Service
	found.Spec = sw.Spec.ServiceSpec
	logger.Info("更新 Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
	if err := r.Update(ctx, found); err != nil {
		return ctrl.Result{}, err
	}

	// Service 更新成功
	return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
}

// reconcileWeightedService 在指定了 WeightedBackends 时处理加权流量分发
func (r *ServiceWeightReconciler) reconcileWeightedService(ctx context.Context, sw *k8scrdserviceweight.ServiceWeight) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 在这个实现中，我们将创建一个标准的 Service 并记录加权的后端
	// 在实际实现中，您需要集成服务网格（如 Istio）
	// 或实现自定义的负载均衡逻辑

	// 定义一个新的 Service 对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sw.Name,
			Namespace: sw.Namespace,
			Annotations: map[string]string{
				"example.com/weighted-backends": fmt.Sprintf("%v", sw.Spec.WeightedBackends),
			},
		},
		Spec: sw.Spec.ServiceSpec,
	}

	// 将 ServiceWeight 实例设置为 owner 和 controller
	if err := controllerutil.SetControllerReference(sw, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// 检查此 Service 是否已存在
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("创建加权 Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
		err = r.Create(ctx, service)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Service 创建成功
		logger.Info("创建了带后端的加权 Service", "Backends", sw.Spec.WeightedBackends)
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// 使用新的 annotations 更新找到的 Service
	found.Spec = sw.Spec.ServiceSpec
	found.Annotations["example.com/weighted-backends"] = fmt.Sprintf("%v", sw.Spec.WeightedBackends)
	logger.Info("更新加权 Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
	if err := r.Update(ctx, found); err != nil {
		return ctrl.Result{}, err
	}

	// 记录加权后端用于演示
	logger.Info("加权 Service 已使用后端更新", "Backends", sw.Spec.WeightedBackends)

	// 更新 ServiceWeight 的状态
	return r.updateServiceWeightStatus(ctx, sw)
}

// updateServiceWeightStatus 更新 ServiceWeight 资源的状态
func (r *ServiceWeightReconciler) updateServiceWeightStatus(ctx context.Context, sw *k8scrdserviceweight.ServiceWeight) (ctrl.Result, error) {
	// 创建副本以避免修改原始对象
	updatedSW := sw.DeepCopy()

	// 添加一个表示 ServiceWeight 已就绪的条件
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             "ServiceCreated",
		Message:            "Service 创建/更新成功",
	}

	// 更新或添加条件
	found := false
	for i, c := range updatedSW.Status.Conditions {
		if c.Type == "Ready" {
			updatedSW.Status.Conditions[i] = condition
			found = true
			break
		}
	}

	if !found {
		updatedSW.Status.Conditions = append(updatedSW.Status.Conditions, condition)
	}

	// 更新状态
	if err := r.Status().Update(ctx, updatedSW); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
}

// SetupWithManager 使用 Manager 设置控制器
func (r *ServiceWeightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8scrdserviceweight.ServiceWeight{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
