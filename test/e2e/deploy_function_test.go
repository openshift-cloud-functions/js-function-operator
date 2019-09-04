package e2e

import (
	goctx "context"
	jsfunction "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1"
	"github.com/openshift-cloud-functions/js-function-operator/test"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knservingclient "knative.dev/serving/pkg/client/clientset/versioned"
	"testing"
)

var (
	functionName = "test-function"
	functionSrc  = `
    module.exports = context => {
      return "ok";
    };`
)

func TestDeployFunction(t *testing.T) {
	ctx, f, namespace := test.E2EBootstrap(t)

	t.Logf("Deploying function in namespace %v", namespace)

	exampleJs := &jsfunction.JSFunction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      functionName,
			Namespace: namespace,
		},
		Spec: jsfunction.JSFunctionSpec{
			Func: functionSrc,
		},
	}

	// Create the service
	err := f.Client.Create(goctx.TODO(), exampleJs, &framework.CleanupOptions{TestContext: ctx})
	assert.NoError(t, err)

	knClient := knservingclient.NewForConfigOrDie(f.KubeConfig).ServingV1alpha1()

	// Wait for knative service to be ready
	knService := test.WaitForKnativeServiceReadyDefault(t, &knClient, functionName, namespace)
	serviceURL := knService.Status.Address.URL

	test.AssertGetRequest(t, serviceURL.String(), 200, []byte("ok"))

}
