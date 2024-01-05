# Testing

This document provides an overview of the testing activities used in this project.

## Unit Tests and Env-Tests

To run the unit and integration tests, the following command must be executed. The command ensures that all necessary tooling is executed in case changes to the source code were made, or if this is your first time to execute the tests.

```sh
make generate-and-test
```

## E2E Tests

Because E2E tests need a Kubernetes cluster to run on, they are separate from the remaining tests.
The E2E tests are executed on any PR using [GithubActions](https://github.com/kyma-project/eventing-publisher-proxy/actions/workflows/e2e.yml).
For local execution, you need to follow the steps in the [action](../../.github/workflows/e2e.yml).

As Prerequisites, you need:

- [Docker](https://www.docker.com/) to build the EPP image.
  ```sh
  make docker-build docker-push IMG=<container-registry>/eventing-publisher-proxy:<tag>
  export EPP_IMAGE=$IMG
  ```
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- Access to a Kubernetes cluster (e.g. [k3d](https://k3d.io/) / k8s)  
- [Eventing Manager](https://github.com/kyma-project/eventing-manager/) to execute commands in its MAKEFILE.


