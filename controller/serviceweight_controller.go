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

// ServiceWeightReconciler reconciles a ServiceWeight object
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ServiceWeightReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ServiceWeight instance
	serviceWeight := &k8scrdserviceweight.ServiceWeight{}
	err := r.Get(ctx, req.NamespacedName, serviceWeight)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			logger.Info("ServiceWeight resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get ServiceWeight.")
		return ctrl.Result{}, err
	}

	// Create or update the corresponding Service
	// If WeightedBackends is empty, it behaves like a normal Service
	if len(serviceWeight.Spec.WeightedBackends) == 0 {
		// Create/update standard service
		return r.reconcileStandardService(ctx, serviceWeight)
	} else {
		// Create/update weighted service (need to implement the actual weighted traffic logic)
		return r.reconcileWeightedService(ctx, serviceWeight)
	}
}

// reconcileStandardService handles the standard Service behavior when no WeightedBackends are specified
func (r *ServiceWeightReconciler) reconcileStandardService(ctx context.Context, sw *k8scrdserviceweight.ServiceWeight) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Define a new Service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sw.Name,
			Namespace: sw.Namespace,
		},
		Spec: sw.Spec.ServiceSpec,
	}

	// Set ServiceWeight instance as the owner and controller
	if err := controllerutil.SetControllerReference(sw, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Service already exists
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
		err = r.Create(ctx, service)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Update the found Service spec with the desired spec from ServiceWeight
	found.Spec = sw.Spec.ServiceSpec
	logger.Info("Updating Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
	if err := r.Update(ctx, found); err != nil {
		return ctrl.Result{}, err
	}

	// Service updated successfully
	return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
}

// reconcileWeightedService handles the weighted traffic distribution when WeightedBackends are specified
func (r *ServiceWeightReconciler) reconcileWeightedService(ctx context.Context, sw *k8scrdserviceweight.ServiceWeight) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// For this implementation, we'll create a standard Service and log the weighted backends
	// In a real implementation, you would need to integrate with a service mesh (like Istio)
	// or implement custom load balancing logic

	// Define a new Service object
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

	// Set ServiceWeight instance as the owner and controller
	if err := controllerutil.SetControllerReference(sw, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Service already exists
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a weighted Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
		err = r.Create(ctx, service)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Service created successfully
		logger.Info("Created weighted Service with backends", "Backends", sw.Spec.WeightedBackends)
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Update the found Service with the new annotations
	found.Spec = sw.Spec.ServiceSpec
	found.Annotations["example.com/weighted-backends"] = fmt.Sprintf("%v", sw.Spec.WeightedBackends)
	logger.Info("Updating weighted Service", "Service.Name", service.Name, "Service.Namespace", service.Namespace)
	if err := r.Update(ctx, found); err != nil {
		return ctrl.Result{}, err
	}

	// Log the weighted backends for demonstration
	logger.Info("Weighted Service updated with backends", "Backends", sw.Spec.WeightedBackends)

	// Update the status of the ServiceWeight
	return r.updateServiceWeightStatus(ctx, sw)
}

// updateServiceWeightStatus updates the status of the ServiceWeight resource
func (r *ServiceWeightReconciler) updateServiceWeightStatus(ctx context.Context, sw *k8scrdserviceweight.ServiceWeight) (ctrl.Result, error) {
	// Create a copy to avoid modifying the original object
	updatedSW := sw.DeepCopy()

	// Add a condition indicating the ServiceWeight is ready
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             "ServiceCreated",
		Message:            "Service created/updated successfully",
	}

	// Update or add the condition
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

	// Update the status
	if err := r.Status().Update(ctx, updatedSW); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceWeightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8scrdserviceweight.ServiceWeight{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
