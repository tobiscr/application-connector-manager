# Application Connector Module


## What is Application Connectivity in Kyma?

Application Connectivity in Kyma simplifies the interaction between external systems and your Kyma workloads. The main benefits are:


* Smooth and loosely coupled integration of external systems with Kyma workloads using [Kyma Eventing](https://kyma-project.io/#/eventing-manager/user/README)

* Easy consumption of SAP BTP services by supporting the SAP BTP Extensibility approach

* Establishing high-security standards for any interaction between systems by using trusted communication channels and authentication methods

* Reducing configuration changes for Kyma workloads through encapsulating configuration details of external API endpoints

* Monitoring and tracing capabilities to facilitate operational aspects


## Application Connector Module

The Application Connector module bundles all features of Application Connectivity in Kyma. You can install and manage the module using Kyma dashboard.

The module includes Kubernetes operators and is fully configurable over its own Kubernetes custom resources (CRs). For each external system, a dedicated configuration is used. This allows for individual configuration of security aspects (like encryption and authentication) per system.

Besides proxying any ingress and egress requests to external systems and dealing with security concerns, it also includes full integration with SAP BTP Unified Customer Landscape (UCL) to simplify the consumption of SAP BTP services.


## Features

The Application Connector module provides the following features:

* Easy installation of Kyma's Application Connectivity capabilities by enabling the Application Connector module in your Kyma Runtime.

* Simple configuration using Kubernetes CRs and easy management with Kyma dashboard.

* Full integration of BTP's UCL service, which implements the SAP BTP Extensibility concept. This allows for the automated integration of external systems registered in the UCL service.

* Dispatching of incoming requests from external systems to Kyma workloads (for example, a Kyma Serverless Function) by using an Istio Gateway with mTLS and the Kyma Eventing module.

* Proxying outgoing requests to external APIs and transparently covering security requirements like encryption and authentication (like OAuth 2.0 + mTLS, Basic Auth, and Client Certificates).

* Metering of throughput and exposing monitoring metrics.


### Options for integrating external systems

#### Automatically by UCL

If an external systems is registered for the Kyma Runtime in BTP's UCL (Unified Customer Landscape), it's automatically configured by the Application Connector and able to send requests to Kyma workloads. The Application Connector Modules includes a [`Runtime Agent`](components/00-10-runtime-agent.md) and acts as client of the UCL backend. It retrieves automatically the configuration of each external system and integrates it with Kyma.


#### Manually

It is always possible to integrate any exernal system into Kyma by applying the configuration by hand. The steps for configuring and integrating a new external system in your Kyma Runtime are described in [this tutorial](tutorials/README.md).

