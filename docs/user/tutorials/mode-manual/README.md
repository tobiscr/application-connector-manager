<span class="fd-avatar--thumbnail fd-avatar fd-avatar--40 fd-avatar--accent-color-10" style="background-image: url('https://avatars.wdf.sap.corp/avatar/I539990')" role="img"></span> [Tobias Schuhmacher](https://people.wdf.sap.corp/profiles/I539990) - Product Owner of Team Framefrog, Kyma runtime

Last update: Oct 2024

# Integrate Kyma with an external system

## Table of Contents

- [Integrate Kyma with an external system](#integrate-kyma-with-an-external-system)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
    - [Prerequisites](#prerequisites)
  - [Initial Setup](#initial-setup)
    - [Create a Kyma instance in the BTP Cockpit and Add the Required Modules](#create-a-kyma-instance-in-the-btp-cockpit-and-add-the-required-modules)

## Introduction

The document describes the steps for connecting an external system (e.g. HTTPBin) with a Kyma instance.

In this example, Kyma will send an authenticated requests to the external process.


### Prerequisites

Besides the Kyma default modules like Istio and API Gateway, you must enable the following Kyma modules:

* `Application Connector` acting as a client for the BTP extensibility mechanism
* `Serverless` to run a Function that is processing SAP CX events
* `Eventing` to receive and dispatch incoming events from SAP Commerce Cloud to the Serverless Function
* `NATS` as an extension of the Eventing module, which includes an in-memory eventing backend.

## Initial Setup

### Create a Kyma instance in the BTP Cockpit and Add the Required Modules

1. [Create a Kyma instance](https://help.sap.com/docs/btp/sap-business-technology-platform/create-kyma-environment-instance) in the BTP Cockpit

2. Open Kyma dashboard

   ![BTP Kyma Open Dashboard](assets/btp-kyma-open-dashboard.png)

3. Add the following Kyma modules using Kyma dashboard: `Application Connector`, `Eventing`, `NATS`, and `Serverless`. See the [Add and Delete a Kyma Module](https://help.sap.com/docs/btp/sap-business-technology-platform/enable-and-disable-kyma-module?#add-and-delete-a-kyma-module-using-kyma-dashboard) tutorial, to learn how to do it.

4. Configure NATS as a backend for the Eventing module. This will change the status of the Eventing module from  `Warning` to `Ready`.

   ![Kyma configure NATS eventing backend](assets/kyma-nats-backend.png)

