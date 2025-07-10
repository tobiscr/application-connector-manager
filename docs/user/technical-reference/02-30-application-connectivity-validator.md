# Application Connectivity Validator

Application Connectivity Validator is an intermediary component processing incoming events from a connected External System.

It verifies incoming requests for all connected External Systems by comparing the certificate subjects with the registered applications. If the subject matches with an application name , the request is accepted and the validator passes it to the Kyma Eventing module.

The following diagram shows how Application Connectivity Validator service interacts with other Kyma modules and external system, within Application Connector components.

![Application Connectivity Validator Diagram](../assets/app-conn-validator.png)

## Workflow
1. An external system sends an event to Kyma with MTLS protocol using certificate issued for the application.
2. Istio Ingress Gateway verifies the client certificate, following MTLS protocol, allows for connection and proxies the event to Application Connectivity Validator.
3. Application Connectivity Validator verifies if the subject of a client certificate attached to the event matches the application name.
4. Application Connectivity Validator forwards the request to the Eventing Module.
5. Eventing module proxies the event to the Kyma Workload subscribed for the event.


> [!NOTE]
> As the Istio Ingress Gateway provides essential parts of Application Connectivity Validator workflow \
> the Istio Service Mesh is required to be installed to receive events. \
> For more information, see the [Istio Ingress Gateway](02-10-istio-gateway.md)