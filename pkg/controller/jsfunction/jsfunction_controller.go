package jsfunction

import (
	"context"

	knv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	knv1beta1 "github.com/knative/serving/pkg/apis/serving/v1beta1"

	faasv1alpha1 "github.com/lance/js-function-operator/pkg/apis/faas/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_jsfunction")

// Add creates a new JSFunction Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileJSFunction{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("jsfunction-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource JSFunction
	err = c.Watch(&source.Kind{Type: &faasv1alpha1.JSFunction{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner JSFunction
	err = c.Watch(&source.Kind{Type: &knv1alpha1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &faasv1alpha1.JSFunction{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileJSFunction implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileJSFunction{}

// ReconcileJSFunction reconciles a JSFunction object
type ReconcileJSFunction struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a JSFunction object and makes changes based on the state read
// and what is in the JSFunction.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileJSFunction) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling JSFunction")

	// Fetch the JSFunction instance
	function := &faasv1alpha1.JSFunction{}
	err := r.client.Get(context.TODO(), request.NamespacedName, function)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Function resource not found. Reconciled object must have been deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Function. Requeing the request.")
		return reconcile.Result{}, err
	}

	// Check if a Service for this JSFunction already exists, if not create a new one
	found := &knv1alpha1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: function.Name, Namespace: function.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// No service for this function exists. Create a new one
		service := r.serviceForFunction(function)
		reqLogger.Info("Creating a new knative Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return reconcile.Result{}, err
		}

		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service for JSFunction")
		return reconcile.Result{}, err
	}

	// TODO update the JSFunction status with the pod names
	// TODO update status nodes if necessary

	reqLogger.Info("JSFunction Service exists. Requeueing.", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	return reconcile.Result{Requeue: true}, nil
}

func (r *ReconcileJSFunction) serviceForFunction(f *faasv1alpha1.JSFunction) *knv1alpha1.Service {
	// replicas := f.Spec.Size
	// labels := map[string]string{"app": "jsfunction", "jsfunction_cr": f.Name}
	service := &knv1alpha1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "serving.knative.dev/v1beta1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: f.Namespace,
		},
		Spec: knv1alpha1.ServiceSpec{
			ConfigurationSpec: knv1alpha1.ConfigurationSpec{
				Template: &knv1alpha1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{"sidecar.istio.io/inject": "false"},
					},
					Spec: knv1alpha1.RevisionSpec{
						RevisionSpec: knv1beta1.RevisionSpec{
							PodSpec: corev1.PodSpec{
								Containers: []corev1.Container{{
									Image:   "node:12-alpine",
									Name:    "node-12",
									Command: []string{"node", "-e", f.Spec.Func},
									Ports: []corev1.ContainerPort{{
										ContainerPort: 8080,
										Name:          f.Name,
									}},
								}},
							},
						},
					},
				},
			},
			RouteSpec: knv1alpha1.RouteSpec{},
		},
	}

	// Set JSFunction instance as the owner and controller
	if err := controllerutil.SetControllerReference(f, service, r.scheme); err != nil {
		log.Error(err, "Failed to set controller reference for function deployment")
	}

	return service
}
