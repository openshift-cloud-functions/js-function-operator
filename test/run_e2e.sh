#!/bin/bash

kubectl create namespace operator-test
kubectl create -f deploy/crds/faas_v1alpha1_jsfunction_crd.yaml
kubectl create -f deploy/service_account.yaml --namespace operator-test
kubectl create -f deploy/role.yaml --namespace operator-test
kubectl create -f deploy/role_binding.yaml --namespace operator-test
kubectl create -f deploy/operator.yaml --namespace operator-test

oc adm policy add-role-to-user edit -z js-function-operator
oc adm policy add-scc-to-user privileged -z js-function-operator

operator-sdk test local ./test/e2e --namespace operator-test --no-setup --verbose

kubectl delete namespace operator-test
