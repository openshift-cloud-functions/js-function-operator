# `JSFunction` Operator
## A Kubernetes Operator for JavaScript functions

This project provides a Kubernetes Operator for managing a `JSFunction`
custom resource. A `JSFunction` accepts a JavaScript function inline,
as well as an optional inline pakcage.json.

When a `JSFunction` resource is applied to a cluster,
the controller takes the following steps during reconciliation.

* Check to see if a Knative `Service` for this function exists. If not, create one.
* Create a `ConfigMap` containing with the user supplied data (index.js/package.json)
* Create a `TaskRun`, installing any dependencies and building a runtime image using [`lanceball/js-runtime`](https://github.com/openshift-cloud-functions/faas-js-runtime-image)
* Create a `PodSpec` for the `Service` specifying the runtime image just created
* Wires up knative eventing (still a work in progress)

This is still in early stages of development and may change rapidly.

### Running the operator
This project has been developed using OpenShift 4.x - both on an AWS cluster,
as well as with Code Ready Containers. To run it in your environment, you can
either use `operator-sdk up local` for testing, or deploy using operator.yaml.

#### Operator dependencies
The project uses [Knative Serving and Eventing](https://knative.dev), as well as 
[OpenShift Pipelines](https://github.com/openshift/tektoncd-pipeline).
Be sure you have these installed in your namespace.

#### Deploy the operator
Be sure you are logged in to your cluster, then run the following commands.

```sh
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunction_crd.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/operator.yaml
kubectl apply -f deploy/build/js-function-build-task.yaml
```

This will set up a service account with appropriate roles, add the `JSFunction`
type, and deploy the operator to the current namespace. To deploy a function,
run `kubectl apply -f deploy/crds/faas_v1alpha1_jsfunction_cr.yaml`.
