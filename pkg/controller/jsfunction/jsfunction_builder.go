package jsfunction

import (
	"context"
	"fmt"

	faasv1alpha1 "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1"
	pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

/**
 * Creates a ConfigMap if one doesn't exist, or updates the one that is there
 * with the source function and a possible package.json. If there is a package.json
 * and if this is a new ConfigMap or the contents have been changed, a Tekton
 * TaskRun is created to run a build.
 */
func (r *ReconcileJSFunction) createOrUpdateSourceBuild(function *faasv1alpha1.JSFunction) error {
	// Check if a ConfigMap containing the function code exists yet
	logger := log.WithValues("Function.Namespace", function.Namespace, "Function.Name", function.Name)
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: function.Name, Namespace: function.Namespace}, configMap)
	if err != nil && errors.IsNotFound(err) {
		// No ConfigMap exists yet, create it
		// Create configmap for function code and package.json
		logger.Info("Creating new ConfigMap for function.")
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      function.Name,
				Namespace: function.Namespace,
			},
			Data: mapFunctionData(function),
		}
		if err = controllerutil.SetControllerReference(function, configMap, r.scheme); err != nil {
			logger.Error(err, "Cannot set function controller reference on ConfigMap")
		} else if err = r.client.Create(context.TODO(), configMap); err != nil {
			logger.Error(err, "Cannot create ConfigMap for function")
		}
		logger.Info("Creating TaskRun for function build.")
		err = r.runBuild(function)
	} else if err != nil {
		logger.Error(err, "Error getting ConfigMap for function")
	} else if configMap.Data["index.js"] != function.Spec.Func {
		// Be sure the ConfigMap is updated with the latest code
		logger.Info("Updating ConfigMap with function changes")
		configMap.Data = mapFunctionData(function)
		err = r.client.Update(context.TODO(), configMap)
		if err != nil {
			logger.Error(err, "Error updating ConfigMap for function")
		}
		logger.Info("Creating TaskRun for function build.")
		err = r.runBuild(function)
	}

	return err
}

func (r *ReconcileJSFunction) runBuild(function *faasv1alpha1.JSFunction) error {
	build := buildForFunction(function)
	if err := controllerutil.SetControllerReference(function, build, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), build)
}

func (r *ReconcileJSFunction) configMapWithFunction(f *faasv1alpha1.JSFunction) (*corev1.ConfigMap, error) {
	// Create a config map containing the user code
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: f.Namespace,
		},
		Data: mapFunctionData(f),
	}
	if err := controllerutil.SetControllerReference(f, configMap, r.scheme); err != nil {
		return nil, err
	}
	return configMap, nil
}

func mapFunctionData(f *faasv1alpha1.JSFunction) map[string]string {
	data := map[string]string{"index.js": f.Spec.Func}

	if f.Spec.Package != "" {
		data["package.json"] = f.Spec.Package
	}
	return data
}

func runtimeImageForFunction(f *faasv1alpha1.JSFunction) string {
	return fmt.Sprintf("image-registry.openshift-image-registry.svc:5000/%s/%s-runtime", f.Namespace, f.Name)
}

func buildForFunction(f *faasv1alpha1.JSFunction) *pipeline.TaskRun {
	imageName := runtimeImageForFunction(f)
	return &pipeline.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-build-", f.Name),
			Namespace:    f.Namespace,
		},
		Spec: pipeline.TaskRunSpec{
			ServiceAccount: "js-function-operator",
			TaskRef: &pipeline.TaskRef{
				Name: "js-function-build-runtime",
			},
			Inputs: pipeline.TaskRunInputs{
				Params: []pipeline.Param{{
					Name: "FUNCTION_NAME",
					Value: pipeline.ArrayOrString{
						Type:      "string",
						StringVal: f.Name,
					},
				}},
			},
			Outputs: pipeline.TaskRunOutputs{
				Resources: []pipeline.TaskResourceBinding{
					{
						Name: "image",
						ResourceSpec: &pipeline.PipelineResourceSpec{
							Type: "image",
							Params: []pipeline.ResourceParam{{
								Name:  "url",
								Value: fmt.Sprintf("%s:latest", imageName),
							}},
						},
					},
				},
			},
		},
	}
}
