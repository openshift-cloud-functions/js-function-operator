package test

import (
	"github.com/openshift-cloud-functions/js-function-operator/pkg/apis"
	jsfunction "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	metav1errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/serving/pkg/apis/serving/v1alpha1"
	servingv1alpha1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	"net/http"
	"testing"
	"time"
)

const (
	defaultRetryInterval = 1 * time.Second
	defaultTimeout       = defaultRetryInterval * 40
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
	OPERATOR_NAME        = "js-function-operator"
)

func AssertGetRequest(t *testing.T, url string, expectedStatusCode int, expectedBody []byte) {
	res, err := http.Get(url)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	assert.Equal(t, expectedStatusCode, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	assert.Equal(t, expectedBody, b)
}

func WaitForKnativeServiceReadyDefault(t *testing.T, servingClient *servingv1alpha1.ServingV1alpha1Interface, name string, namespace string) *v1alpha1.Service {
	return WaitForKnativeServiceReady(t, servingClient, name, namespace, defaultRetryInterval, defaultTimeout)
}

func WaitForKnativeServiceReady(t *testing.T, servingClient *servingv1alpha1.ServingV1alpha1Interface, name string, namespace string, retryInterval, timeout time.Duration) *v1alpha1.Service {
	var service *v1alpha1.Service

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		service, err = (*servingClient).Services(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if metav1errors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s service\n", name)
				return false, nil
			}
			return false, err
		}

		if service.Status.IsReady() {
			return true, nil
		}
		t.Logf("Waiting for availability of %s service. Actual status: %+v\n", name, service.Status)
		return false, nil
	})

	assert.NoError(t, err)

	return service
}

// This function bootstraps an e2e adding CRDs, initializing the test ctx and deploying the operator
func E2EBootstrap(t *testing.T) (*framework.TestCtx, *framework.Framework, string) {
	// Add CRDs
	jsFunctionList := &jsfunction.JSFunctionList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, jsFunctionList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	t.Log("Added CRDs")

	// Setup test framework
	ctx := framework.NewTestCtx(t)
	t.Log("Test Framework context ready")

	// Initialize cluster resources
	err = ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	f := framework.Global
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, OPERATOR_NAME, 1, defaultRetryInterval, defaultTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Operator ready")

	return ctx, f, namespace
}
