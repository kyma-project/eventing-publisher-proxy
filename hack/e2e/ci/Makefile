.PHONY: install-eventing-manager
install-eventing-manager: ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	kubectl apply -k ./config

.PHONY: install-k3d-tools
install-k3d-tools: ## Create k3d with kyma CRDs.
	curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=${K3D_VERSION} bash
