# Testing

The Eventing Publisher Proxy uses the following testing activities:

## Unit Tests and Integration Tests

To run the unit and integration tests, you must run the following command:

```sh
make generate-and-test
```
The command ensures that all necessary tooling is executed if the source code changed, or if this is your first time to execute the tests

## E2E Tests

Because E2E tests need a Kubernetes cluster to run on, they are separate from the remaining tests.
The E2E tests are executed on any PR using [GithubActions](https://github.com/kyma-project/eventing-publisher-proxy/actions/workflows/e2e.yml).
For local execution, follow the steps in the [action](../../.github/workflows/e2e.yml).

As prerequisites, you need:

- [Docker](https://www.docker.com/) to build the EPP image.
  ```sh
  make docker-build docker-push IMG=<container-registry>/eventing-publisher-proxy:<tag>
  export EPP_IMAGE=$IMG
  ```
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- Access to a Kubernetes cluster (like [k3d](https://k3d.io/) or k8s)
- [Eventing Manager](https://github.com/kyma-project/eventing-manager/) to execute commands in its MAKEFILE.


