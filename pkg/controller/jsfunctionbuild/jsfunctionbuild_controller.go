package jsfunctionbuild

import (
	"context"
	"fmt"

	faasv1alpha1 "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1"
	pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
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

var log = logf.Log.WithName("controller_jsfunctionbuild")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new JSFunctionBuild Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileJSFunctionBuild{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("jsfunctionbuild-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource JSFunctionBuild
	err = c.Watch(&source.Kind{Type: &faasv1alpha1.JSFunctionBuild{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource PipelineResources and requeue the owner JSFunctionBuild
	err = c.Watch(&source.Kind{Type: &pipeline.PipelineResource{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &faasv1alpha1.JSFunctionBuild{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileJSFunctionBuild implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileJSFunctionBuild{}

// ReconcileJSFunctionBuild reconciles a JSFunctionBuild object
type ReconcileJSFunctionBuild struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a JSFunctionBuild object and makes changes based on the state read
// and what is in the JSFunctionBuild.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileJSFunctionBuild) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling JSFunctionBuild")

	// Fetch the JSFunctionBuild instance
	functionBuild := &faasv1alpha1.JSFunctionBuild{}
	err := r.client.Get(context.TODO(), request.NamespacedName, functionBuild)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	/// PipelineResource section start
	// Create a new PipelineResource
	pipelineResource := newPipelineResourceForBuild(functionBuild)

	// Set this JSFunctionBuild as the owner and controller
	if err = controllerutil.SetControllerReference(functionBuild, pipelineResource, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this PipelineResource already exists
	found := &pipeline.PipelineResource{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pipelineResource.Name, Namespace: pipelineResource.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new PipelineResource", "PipelineResource.Namespace", pipelineResource.Namespace, "PipelineResource.Name", pipelineResource.Name)
		if err = controllerutil.SetControllerReference(functionBuild, pipelineResource, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.client.Create(context.TODO(), pipelineResource)
		if err != nil {
			return reconcile.Result{}, err
		}

		// PipelineResource created successfully - run a build
		reqLogger.Info("Creating a new PipelineRun for this build")
		if err = r.RunPipeline(functionBuild); err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else if found.Spec.Params[0].Value != pipelineResource.Spec.Params[0].Value {
		// PipelineResource for this build exists but is out dated.
		// Update the PipelineResource with the new build image
		reqLogger.Info("Updating existing PipelineResource for build.")
		found.Spec.Params[0].Value = pipelineResource.Spec.Params[0].Value
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}

		// PipelineResource updated successfully - trigger a build
		reqLogger.Info("Creating a new PipelineRun for this build")
		if err = r.RunPipeline(functionBuild); err != nil {
			return reconcile.Result{}, err
		}
	}
	/// PipelineResource section end

	return reconcile.Result{}, nil
}

// RunPipeline creates a tekton PipelineRun resource
func (r *ReconcileJSFunctionBuild) RunPipeline(build *faasv1alpha1.JSFunctionBuild) error {
	pipelineRun := newPipelineRunForBuild(build)

	// Set this JSFunctionBuild as the owner and Controller
	if err := controllerutil.SetControllerReference(build, pipelineRun, r.scheme); err != nil {
		return err
	}
	if err := r.client.Create(context.TODO(), pipelineRun); err != nil {
		return err
	}
	return nil
}

func newPipelineResourceForBuild(build *faasv1alpha1.JSFunctionBuild) *pipeline.PipelineResource {
	labels := map[string]string{
		"app": build.Name,
	}
	return &pipeline.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pipeline-resource", build.Name),
			Namespace: build.Namespace,
			Labels:    labels,
		},
		Spec: pipeline.PipelineResourceSpec{
			Type: "image",
			Params: []pipeline.ResourceParam{{
				Name:  "url",
				Value: build.Spec.Image,
			}},
		},
	}
}

func newPipelineRunForBuild(build *faasv1alpha1.JSFunctionBuild) *pipeline.PipelineRun {
	return &pipeline.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pipeline-run-%d", build.Name, build.Spec.Revision),
			Namespace: build.Namespace,
		},
		Spec: pipeline.PipelineRunSpec{
			ServiceAccount: "js-function-operator",
			PipelineRef: pipeline.PipelineRef{
				Name: "js-function-build-pipeline",
			},
			Resources: []pipeline.PipelineResourceBinding{{
				Name: "image",
				ResourceRef: pipeline.PipelineResourceRef{
					Name: fmt.Sprintf("%s-pipeline-resource", build.Name),
				},
			}},
			Params: []pipeline.Param{{
				Name: "FUNCTION_NAME",
				Value: pipeline.ArrayOrString{
					Type:      "string",
					StringVal: build.Name,
				},
			}},
		},
	}
}
