.PHONY: build test

build:
	operator-sdk build --verbose docker.io/oscf/js-function-operator:v0.0.1

push:
	docker push docker.io/oscf/js-function-operator:v0.0.1

kube:
	minikube start --memory=16g --cpus=8 --disk-size=30g --vm-driver=kvm2 --kubernetes-version=v1.12.0 --extra-config=apiserver.enable-admission-plugins="LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"

test:
	./test/setup-cluster.sh
	./test/run_e2e.sh