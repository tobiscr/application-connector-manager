# Application Connector Module


## What is Application Connectivity in Kyma?

The Application Connectivity in Kyma simplifies the interaction between external systems and your Kyma workloads.

Main benefits are:

* Smooth and loosely coupled integration of external systems with Kyma workloads via Kyma eventing
* Easy consumption of BTP services by supporting the BTP Extensibility approach.
* Establishing high security standards for any interaction between systens by using trusted communication channels and authentication methods
* Encapsulating configuration details of external API endpoints reduces confiugration changes for Kyma workloads
* Offering monitoring and tracing capabilities to facilitate operational aspects


## Application Connector Module

The Application Connector module bundles all features for Application Connectivity in Kyma. The installation and confiugration of the Application Connectivity module is managed over the Kyma Dashboard.

The module includes its own Kuberentes operators and is fully confiugrable over its own Kuberentes Custome Resources. For each external system, a dedicated configuration is used. This allows an individual configuration of security and aspects (like encryption and authentication) per system.

Beside proxing any ingress- / egress-request to external systems and dealing with security concerns, it also contains a full integration with BTP's UCL (unified customer landscape) to simplify the consumption of BTP services.


## Features

In summary, the Application Connector module provides the following features:

* Easy installation of Kyma's Application Connectivity capabiliites by enabling the Kyma Application Connector module on your Kyma cluster
* Simple configuration over Kuberentes Custom Resources and managable over the Kyma Dashboard
* Full integration of BTP's UCL service (Unified Customer Landscape) which implements the BTP Extensibility Concept. This allows an automated integration of external systems which are registed in the UCL service.
* Dispatching of incoming requests from external systems to Kyma workloads (e.g. a [Kyma Serverless Function](https://kyma-project.io/#/serverless-manager/user/resources/06-10-function-cr)) by using an Istio Gateway with mTSL and the [Kyma Eventing Module](https://kyma-project.io/#/eventing-manager/user/README)
* Proxying outgoing requests to external APIs and transparently covering security requirements like encryption and authentication (like OAuth 2.0 + mTLS, Basic Auth, Client Certificates).
* Metering of throughput and exposing monitoring metrics.


### Integration of external systems

#### Automatically by UCL

If an external systems is registered for Kyma in BTP's UCL, its automatically configured by the Application Connector and able to send requests to Kyma workloads. The Application Connector Modules includes a [`Runtime Agent`](components/00-10-runtime-agent.md) which acts as client and communicates with the UCL backend. It retrieves the configuration of each external system and integrates it in Kyma.

An example how a system can be registered in UCL and gets integrated into Kyma is provided in [this tutorial](tutorials/mode-ucl/README.md).


#### Manually

It is always possible to integrate any exernal system into Kyma by applying the configuration by hand. The steps for configuring a new external system are described in [this tutrorial](tutorials/mode-manual/README.md).


## Architecture

The diagram shows all components of the Application Connector Module and describes the flow, how an external application is integrated into Kyma.

![Application Connector Architecture](./assets/ac-architecture.png)


### Components

The diagram shows all componentes included in a common Application Connector Module use case. To get more details of a particular component, please follow the links.

|Component|Purpose|
|--|--|
|External Application|The remote system which wants to interact with a Kyma workload or should be called by it.|
|UCL|The UCL (Unified Customer Landscape) implements the BTP Extensibility Concept and administrates system formations.|
|[Runtime Agent](./technical-reference/04-30-runtime-agent.md)|The Runtime Agent is a client of the UCL system and integrates system formations automatically in the Kyma cluster.|
|Certificate Secret|Stores the UCL managed certificates used for trusted communications.|
|[Application Customer Resource (CR)](./resources/06-10-application.md)|Used to store metadata of the external system (like endpoints, authentication method etc.). Each Application custom resource corresponds to a single external system.|
|Application Credentials Secret|Stores endpoint/API credentials of the external system.|
|[Istio Ingress Gateway](./technical-reference/04-10-istio-gateway.md)|The Application Connector Module uses an Istio Gateway as endpoint for incoming requests from external system. The Gateway supports mTLS for establishing trusted connections between the external system and the Kyma cluster.|
|Application Connectivity Validator|Verifies incoming requests (e.g. by analying the certificate subject) and passes the request to the Kyma eventing module.|
|Eventing Module|The Kyma Eventing Module is used for dispatching incoming requests from external systems to Kyma worloads.|
|Kyma Workload|Can be a customer workload (e.g. deployed applications) or Kyma hosted serverless function.|
|[Application Gateway](./technical-reference/04-20-application-gateway.md)|This component acts as an proxy for outgoing communication from a Kyma workload to an external system. It supports various types of authentication methods.|
|Application Connector Manager|It takes care of installation and configuration of all the Application Connector module components on your cluster. It manages the lifecycle of the Application Connector module based on the dedicated ApplicationConnector custom resource (CR).|


### Workflow

The diagram includes also the flow from registring an external system at UCL until it is interacting with Kyma workloads. Some flow steps are only relevant when the integartion of the external system is managed manually, but all of them are applied when the integration happens automatically via UCL.

<table>
    <tr>
        <th>Step</th>
        <th>Description</th>
        <th>Manual integration</th>
        <th>Integration by UCL</th>
    </tr>
    <tr>
        <td colspan="4"><strong>Inbound communciation</strong></td>
    </tr>
    <tr>
        <td>1</td>
        <td>The external system is registered in the UCL system (Unified Customer Landscape).</td>
        <td></td>
        <td>X</td>
    </tr>
    <tr>
        <td>2</td>
        <td>Runtime Agent connects to the UCL system and gathers all registered applications for this Kyma runtime.</td>
        <td></td>
        <td>X</td>
    </tr>
    <tr>
        <td>3</td>
        <td>The agent retrieves from UCL the certificate used for securing the communication with the external system. The certificate is stored in a Kubernetes secret.</td>
        <td></td>
        <td>X</td>
    </tr>
    <tr>
        <td>4</td>
        <td>Additionally, the application metadata of the external system (e.g. authentication type, API endpoints etc.) are stored in an Application custom resource (CR).</td>
        <td>X<br/>This is a manual step applied by an administrator</td>
        <td>X<br/>Step is automatically applied by Runtime Agent</td>
    </tr>
    <tr>
        <td>5</td>
        <td>Finally, the Runtime Agents stores the credentials for accessing the API of an external system in a Kubernetes secret.
These credentials are used for outbound communication by the Application Gateway.</td>
        <td></td>
        <td>X</td>
    </tr>
    <tr>
        <td>6</td>
        <td>The external systems communicates with the Kyma cluster via an Istio Ingress Gateway.</td>
        <td>X</td>
        <td>X</td>
    </tr>
    <tr>
        <td>7</td>
        <td>Istio gateway is using the provided certificate for securing the communication.</td>
        <td></td>
        <td>X</td>
    </tr>
    <tr>
        <td>8</td>
        <td>Incoming calls are verified by the Connectivity Validator by investigating the subject of the certificate. It finally forwards the request to the Eventing component.</td>
        <td></td>
        <td>X</td>
    </tr>
    <tr>
        <td>9</td>
        <td>Kyma workloads can receive incoming events and process them.</td>
        <td>X</td>
        <td>X</td>
    </tr>
    <tr>
        <td colspan="4"><strong>Onbound communciation</strong></td>
    </tr>
    <tr>
        <td>a</td>
        <td>The Application Gateway acts as an proxy for any outbound communication from a Kyma workload to an external system API.</td>
        <td>X</td>
        <td>X</td>
    </tr>
    <tr>
        <td>b</td>
        <td>The Application Gateway takes care about security and authentication during the communication with the external system API.</td>
        <td>X</td>
        <td>X</td>
    </tr>
</table>

