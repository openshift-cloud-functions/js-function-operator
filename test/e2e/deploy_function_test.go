package e2e

import (
	goctx "context"
	jsfunction "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1"
	"github.com/openshift-cloud-functions/js-function-operator/test"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knservingclient "knative.dev/serving/pkg/client/clientset/versioned"
	"net/http"
	"testing"
)

const OPERATOR_NAME string = "js-function-operator"

var (
	functionName = "test-function"
	functionSrc  = `
    module.exports = context => {
      return "ok";
    };`
)

func TestDeployFunction(t *testing.T) {
	ctx, f, namespace := test.E2EBootstrap(t)

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

	// Try to do a request
	res, err := http.Get(serviceURL.String())

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.Equal(t, "ok", string(body))

}
