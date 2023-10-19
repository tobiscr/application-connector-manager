# CompassConnection

The `applicationconnectors.operator.kyma-project.io` CustomResourceDefinition (CRD)
is a detailed description of the kind of data and the format used to preserve
the configuration and status of the Application Connector Module .
The `ApplicationConnector` custom resource (CR) contains the state and statuses of the module installed on the cluster.

To get the up-to-date CRD and show the output in the `yaml` format, run this command:

```bash
kubectl get crd applicationconnectors.operator.kyma-project.io -o yaml
```
## Sample custom resource

This is a sample resource that instance all parts of 