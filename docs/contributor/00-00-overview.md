# Overview: Eventing-Publisher-Proxy

The Eventing Publisher Proxy receives legacy and Cloud Event publishing requests from the cluster workloads (microservices or Serverless Functions) and redirects them to the Enterprise Messaging Service Cloud Event Gateway.
It also fetches a list of subscriptions for a connected application.

It is a part of Eventing Manager to process and deliver events in Kyma.
For further information, refer to [Eventing Architecture](https://github.com/kyma-project/eventing-manager/blob/main/docs/user/evnt-architecture.md).
