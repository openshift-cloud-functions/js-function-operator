# `JSFunction` Operator
## A Kubernetes Operator for JavaScript functions

This project provides a Kubernetes Operator for managing a `JSFunction`
custom resource. A `JSFunction` accepts a JavaScript function inline,
as well as an optional inline package.json.

When a `JSFunction` resource is applied to a cluster,
the controller takes the following steps during reconciliation.

* Check to see if a Knative `Service` for this function exists. If not, create one.
* Create a `ConfigMap` containing with the user supplied data (index.js/package.json)
* Create a `TaskRun`, installing any dependencies and building a runtime image using [`openshift-cloud-functions/faas-js-runtime-image`](https://github.com/openshift-cloud-functions/faas-js-runtime-image)
* Create a `PodSpec` for the `Service` specifying the runtime image just created
* Wires up knative eventing, if `events` is set to `true` in a `JSFunction` custom resource. 
Knative Eventing objects `Subscription` and `Channel` are created and acts as a sink for the function.
Events sent to the `Channel` are passed to the function.

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

```console
./deploy/deploy.sh
```

This will set up a service account with appropriate roles, add the `JSFunction`
type and a build `Task` for the function, and finally deploy the operator to the
current namespace. 

#### Note privileges for the service account
Among the things that happen when a new `JSFunction` is created, is that
a build task involving the creation of runtime images runs. This requires
permissions not available by default on OpenShift. These permissions are added
to the `js-function-operator` service account with the following commands when
you run `./deploy/deploy.sh`.

```sh
oc adm policy add-role-to-user edit -z js-function-operator
oc adm policy add-scc-to-user privileged -z js-function-operator
```

To deploy a function, run 

```sh
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunction_cr.yaml
```

### Testing the Knative Eventing wiring
Set `events` to `true` in a `JSFunction` custom resource and deploy this CR.

To emit events to the automatically created `Channel`, run command specified below.
This will emit event every minute.

```sh
kubectl apply -f hack/cronjob-source.yaml
```
