#!/bin/sh

kubectl apply -f deploy/crds/faas_v1alpha1_jsfunction_crd.yaml
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunctionbuild_crd.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/build/js-function-build.yaml
kubectl apply -f deploy/operator.yaml

oc adm policy add-role-to-user edit -z js-function-operator
oc adm policy add-scc-to-user privileged -z js-function-operator
