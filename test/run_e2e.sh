#!/bin/bash
NS=jsfunction-test

dep ensure -v

kubectl create namespace $NS
kubectl create -f deploy/crds/faas_v1alpha1_jsfunction_crd.yaml
kubectl create -f deploy/crds/faas_v1alpha1_jsfunctionbuild_crd.yaml
kubectl create -f deploy/service_account.yaml --namespace $NS
kubectl create -f deploy/role.yaml --namespace $NS
kubectl create -f deploy/role_binding.yaml --namespace $NS
kubectl create -f deploy/build/js-function-build.yaml --namespace $NS
kubectl create -f deploy/operator.yaml --namespace $NS

oc adm policy add-role-to-user edit -z js-function-operator --namespace $NS
oc adm policy add-scc-to-user privileged -z js-function-operator --namespace $NS

operator-sdk test local ./test/e2e --namespace $NS --no-setup --verbose

kubectl delete jsfunction --all --namespace $NS
kubectl delete jsfunctionbuild --all --namespace $NS
kubectl delete serviceaccount js-function-operator --namespace $NS
kubectl delete task js-function-build-runtime --namespace $NS
kubectl delete task js-function-update-service --namespace $NS
kubectl delete deployment js-function-operator --namespace $NS
kubectl delete crd/jsfunctions.faas.redhat.com
kubectl delete crd/jsfunctionbuilds.faas.redhat.com
kubectl delete namespace $NS
