#!/bin/sh

# This script will set up a minikube cluster with knative serving and eventing
# as well as tekton pipelines.

minikube status

minikube start --memory=8192 --cpus=6 \
  --kubernetes-version=v1.12.0 \
  --vm-driver=kvm2 \
  --disk-size=30g \
  --extra-config=apiserver.enable-admission-plugins="LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"

helm init

# Install Istio
export ISTIO_VERSION=1.1.7
echo "Installing ${ISTIO_VERSION}"
curl -L https://git.io/getLatestIstio | sh -
cd istio-${ISTIO_VERSION}

for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done

echo "Waiting for CRDs to be comitted"
sleep 7

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: istio-system
  labels:
    istio-injection: disabled
EOF

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

kubectl get pods --namespace istio-system

# Install Knative
kubectl apply --selector knative.dev/crd-install=true \
  --filename https://github.com/knative/serving/releases/download/v0.9.0/serving.yaml \
  --filename https://github.com/knative/eventing/releases/download/v0.9.0/release.yaml \
  --filename https://github.com/knative/serving/releases/download/v0.9.0/monitoring.yaml

kubectl apply --filename https://github.com/knative/serving/releases/download/v0.9.0/serving.yaml \
  --filename https://github.com/knative/eventing/releases/download/v0.9.0/release.yaml \
  --filename https://github.com/knative/serving/releases/download/v0.9.0/monitoring.yaml


# Install Tekton
kubectl apply --filename https://storage.googleapis.com/tekton-releases/latest/release.yaml

# Install FaaS
cd ..
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunction_crd.yaml
kubectl apply -f deploy/crds/faas_v1alpha1_jsfunctionbuild_crd.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/build/js-function-build.yaml
kubectl apply -f deploy/operator.yaml

kubectl get namespaces
