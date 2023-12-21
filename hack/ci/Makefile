PROJECT_ROOT ?= ../..
KYMA_CLI ?= "${PROJECT_ROOT}/${KYMA}"
CLUSTER_NAME ?= kyma
REGISTRY_PORT ?= 5001

.PHONY: create-k3d
create-k3d: ## Create k3d with kyma CRDs.
	"${KYMA_CLI}" provision k3d -p 8081:80@loadbalancer -p 8443:443@loadbalancer --registry-port ${REGISTRY_PORT} --name ${CLUSTER_NAME} --ci
