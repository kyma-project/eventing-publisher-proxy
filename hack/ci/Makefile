KYMA_CLI ?= "/home/runner/work/eventing-publisher-proxy/eventing-publisher-proxy/eventing-manager/bin/kyma-unstable"
CLUSTER_NAME ?= kyma
REGISTRY_PORT ?= 5001
K3D_VERSION ?= 1.28.7

.PHONY: create-k3d
create-k3d: ## Create k3d with kyma CRDs.
	"${KYMA_CLI}" provision k3d -p 8081:80@loadbalancer -p 8443:443@loadbalancer --registry-port ${REGISTRY_PORT} --name ${CLUSTER_NAME} -k ${K3D_VERSION} --ci
