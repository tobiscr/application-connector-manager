# Application Connector Manager

## Overview

This repository contains the following projects:

- **Application Connector Manager** - operator compatible with Lifecycle Manager that manages the Application Connector module in Kyma.
- **Application Connector Module build configuration** - delivers the functionality of Application Connector as a Kyma module.

> **NOTE:** Docker images for the Application Connector binaries are built separately from the main GitHub [Kyma repository](https://github.com/kyma-project/kyma/).
 
## How it works 
 
The Application Connector module represents a specific version of the Application Connector binaries.
The configuration of the released module is described as a ModuleTemplate custom resource (CR) and delivered as an OCI image.
You can install it on the Kyma cluster managed by Template Operator.

The installed Application Connector module is represented as the ApplicationConnector Kubernetes CR.

```yaml
apiVersion: operator.kyma-project.io/v1alpha1
kind: ApplicationConnector
metadata:
  labels:
    app.kubernetes.io/name: applicationconnector
    app.kubernetes.io/instance: applicationconnector-sample
    app.kubernetes.io/part-of: application-connector-manager
    app.kuberentes.io/managed-by: kustomize
    app.kubernetes.io/created-by: application-connector-manager
  name: applicationconnector-sample
spec:
  disableLegacyConnectivity : "false"
```

Any update to this CR is intercepted by Application Connector Manager and applied to the Application Connector binaries.

> **NOTE:** At this stage of development, the ApplicationConnector Custom Resource Definition (CRD) contains only one parameter for testing.
> The ApplicationConnector CRD will be extended during further development.

See also:
- [Lifecycle Manager documentation](https://github.com/kyma-project/lifecycle-manager#lifecycle-manager)
- [Application Connector module documentation](docs/user/README.md) 
- [Modularization of Kyma](https://kyma-project.io/docs/kyma/latest#kyma-modules)