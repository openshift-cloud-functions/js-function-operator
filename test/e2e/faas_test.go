package e2e

import (
	goctx "context"
	"github.com/openshift-cloud-functions/js-function-operator/pkg/apis"
	jsfunction "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

const OPERATOR_NAME string = "js-function-operator"

var (
	functionName = "test-function"
	functionSrc  = `
    module.exports = context => {
      return "ok";
    };`
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestJsFaas(t *testing.T) {
	// Add CRDs
	jsFunctionList := &jsfunction.JSFunctionList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, jsFunctionList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	t.Log("Added CRDs")

	// Setup test framework
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
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
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, OPERATOR_NAME, 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Operator ready")

	if err = testRunFunction(t, f, ctx, namespace); err != nil {
		t.Fatal(err)
	}

}

func testRunFunction(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, namespace string) error {
	exampleJs := &jsfunction.JSFunction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      functionName,
			Namespace: namespace,
		},
		Spec: jsfunction.JSFunctionSpec{
			Func: functionSrc,
		},
	}

	err := f.Client.Create(goctx.TODO(), exampleJs, &framework.CleanupOptions{TestContext: ctx})
	if err != nil {
		t.Fatal(err)
	}

}
