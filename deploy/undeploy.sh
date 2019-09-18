#!/bin/sh

kubectl delete --all jsfunctions
kubectl delete deployment js-function-operator
kubectl delete role js-function-operator
kubectl delete rolebinding js-function-operator
kubectl delete rolebinding js-function-operator-cluster
kubectl delete serviceaccount js-function-operator
kubectl delete task js-function-build-runtime
kubectl delete task js-function-update-service
kubectl delete pipeline js-function-build-pipeline

kubectl delete customresourcedefinition.apiextensions.k8s.io/jsfunctions.faas.redhat.com
kubectl delete customresourcedefinition.apiextensions.k8s.io/jsfunctionbuilds.faas.redhat.com
