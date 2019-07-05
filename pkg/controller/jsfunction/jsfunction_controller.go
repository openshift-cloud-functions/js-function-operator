package jsfunction

import (
	"context"

	faasv1alpha1 "github.com/lance/js-function-operator/pkg/apis/faas/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
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

	// Watch for changes to secondary resource Deployment and requeue the owner JSFunction
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
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
			reqLogger.Info("Function resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Function.")
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: function.Name, Namespace: function.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// No deployment for this function exists. Create a new one
		deployment := r.deploymentForFunction(function)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// TODO ensure deployment size is the same as the spec
	// TODO update the JSFunction status with the pod names
	// TODO update status nodes if necessary

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func (r *ReconcileJSFunction) deploymentForFunction(f *faasv1alpha1.JSFunction) *appsv1.Deployment {
	replicas := f.Spec.Size
	labels := map[string]string{"app": "jsfunction", "jsfunction_cr": f.Name}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: f.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "node:12-alpine",
						Name:    "Node 12",
						Command: []string{"node", "-e", f.Spec.Func},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          f.Name,
						}},
					}},
				},
			},
		},
	}

	// Set JSFunction instance as the owner and controller
	if err := controllerutil.SetControllerReference(f, deployment, r.scheme); err != nil {
		log.Error(err, "Failed to set controller reference for function deployment")
	}

	return deployment
}
