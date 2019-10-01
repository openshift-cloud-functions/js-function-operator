#!/bin/sh

# This script will set up a minikube cluster with kn serving, kn eventing,
# kn monitoring, tekton pipelines, istio, and helm. It requires a significant
# amount of system resources and runs minikube with 20g ram, 6 cpus, and
# 30g disk.
#
# The script assumes that you have minikube and helm installed, and will create
# a jsfunction-operator-test profile within minikube for testing. It does not
# modify in any way any other minikube profiles. When you are finished, you can
# delete this profile with the following commands.
#
# minikube stop -p jsfunction-operator-test
# minikube delete -p jsfunction-operator-test
#
# Once the script completes, it will take up to 30 minutes for the cluster to
# be fully stabilized.

echo "This script creates a kubernetes cluster using minikube with the profile jsfunction-operator-test."
echo "To continue, type 'y'. Any other input will stop execution if this script."
read CONT

if [ "y" != "$CONT" ] ; then exit 1 ; fi

set -x

minikube status

minikube start -p jsfunction-operator-test --memory=20g --cpus=6 \
  --kubernetes-version=v1.12.0 \
  --vm-driver=kvm2 \
  --disk-size=30g \
  --extra-config=apiserver.enable-admission-plugins="LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"

echo "Installing helm tiller"
helm init

# Install Istio
export ISTIO_VERSION=1.1.7
echo "Installing ${ISTIO_VERSION}"
curl -L https://git.io/getLatestIstio | sh -
cd istio-${ISTIO_VERSION}

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: istio-system
  labels:
    istio-injection: disabled
EOF

for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done
sleep 5

# A lighter template, with just pilot/gateway.
# Based on install/kubernetes/helm/istio/values-istio-minimal.yaml
helm template --namespace=istio-system \
  --set prometheus.enabled=false \
  --set mixer.enabled=false \
  --set mixer.policy.enabled=false \
  --set mixer.telemetry.enabled=false \
  `# Pilot doesn't need a sidecar.` \
  --set pilot.sidecar=false \
  --set pilot.resources.requests.memory=128Mi \
  `# Disable galley (and things requiring galley).` \
  --set galley.enabled=false \
  --set global.useMCP=false \
  `# Disable security / policy.` \
  --set security.enabled=false \
  --set global.disablePolicyChecks=true \
  `# Disable sidecar injection.` \
  --set sidecarInjectorWebhook.enabled=false \
  --set global.proxy.autoInject=disabled \
  --set global.omitSidecarInjectorConfigMap=true \
  `# Set gateway pods to 1 to sidestep eventual consistency / readiness problems.` \
  --set gateways.istio-ingressgateway.autoscaleMin=1 \
  --set gateways.istio-ingressgateway.autoscaleMax=1 \
  `# Set pilot trace sampling to 100%` \
  --set pilot.traceSampling=100 \
  install/kubernetes/helm/istio \
  > ./istio-lean.yaml

kubectl apply -f istio-lean.yaml

echo "Waiting for 1 minute for Istio system to initialize"
sleep 60
kubectl get pods --namespace istio-system

# Install Knative
echo "Preparing the knative installation"
kubectl apply --selector knative.dev/crd-install=true \
  --filename https://github.com/knative/serving/releases/download/v0.8.0/serving.yaml \
  --filename https://github.com/knative/eventing/releases/download/v0.8.0/release.yaml \
  --filename https://github.com/knative/serving/releases/download/v0.8.0/monitoring.yaml

# Apparently there is a known race condition on install, so just to be safe do this again
kubectl apply --selector knative.dev/crd-install=true \
  --filename https://github.com/knative/serving/releases/download/v0.8.0/serving.yaml \
  --filename https://github.com/knative/eventing/releases/download/v0.8.0/release.yaml \
  --filename https://github.com/knative/serving/releases/download/v0.8.0/monitoring.yaml

echo "Applying knative resources"
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.8.0/serving.yaml \
  --filename https://github.com/knative/eventing/releases/download/v0.8.0/release.yaml \
  --filename https://github.com/knative/serving/releases/download/v0.8.0/monitoring.yaml

echo "Waiting for 1 minute for Knative system to initialize"
sleep 60
kubectl get pods -n knative-monitoring
kubectl get pods -n knative-eventing
kubectl get pods -n knative-serving


# Install Tekton
kubectl apply --filename https://github.com/tektoncd/pipeline/releases/download/v0.5.2/release.yaml

echo "Waiting for 1 minute for Tekton system to initialize"
sleep 60
kubectl get pods -n tekton-pipelines


# Install FaaS
cd ..
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunction_crd.yaml
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunctionbuild_crd.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/build/js-function-build.yaml
kubectl apply -f deploy/operator.yaml

kubectl get pods -n istio-system
kubectl get pods -n tekton-pipelines
kubectl get pods -n knative-monitoring
kubectl get pods -n knative-eventing
kubectl get pods -n knative-serving
kubectl get pods -n default

echo "Cluster set up complete. Wait for up to 10 minutes for the cluster to be fully functional."