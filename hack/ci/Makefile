## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binary
KUSTOMIZE ?= $(LOCALBIN)/kustomize
## Tool Version
KUSTOMIZE_VERSION ?= v5.0.0

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) --output install_kustomize.sh && bash install_kustomize.sh $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); rm install_kustomize.sh; }

.PHONY: create-k3d
create-k3d: ## Create k3d with kyma CRDs.
	"${KYMA_CLI}" provision k3d -p 8081:80@loadbalancer -p 8443:443@loadbalancer --registry-port ${REGISTRY_PORT} --name ${CLUSTER_NAME} --ci

.PHONY: install-eventing-manager
install-eventing-manager: kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build ./../e2e/config/crd | kubectl apply -f -

.PHONY: install-k3d-tools
install-k3d-tools: ## Create k3d with kyma CRDs.
	curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=${K3D_VERSION} bash
