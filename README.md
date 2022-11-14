# Application Connector Manager

## Overview

The repository contains following projects:

- **Application Connector Manager** - operator compatible with Lifecycle Manager that manages Application Connector module in Kyma.
- **Application Connector Module build configuration** to deliver the functionality of Application Connector as a Kyma module

> Note: Docker images for Application Connector binaries are build separately from main GitHub [Kyma repository](https://github.com/kyma-project/kyma/)
 
## How it works 
 
The Application Connector Module represents a specific version of Application Connector binaries delivered via Kyma release channel.\
The configuration of released module is described as a `ModuleTemplate` Custom Resource and delivered as OCI image.\
It can be installed on Kyma Cluster managed by Moduletemplates operator.

The installed Application Connector module is represented as ApplicationConnector Kubernetes CR. \
Any update to this CR is intercepted by Application Connector Manager and applied on Application Connector binaries.

> Note: On this stage of development ApplicationConnector CRD contains only single parameter for testing.\
> The ApplicationConnector CRD and will be extended during further development.

See also:
- [lifecycle-manager documetation](https://github.com/kyma-project/lifecycle-manager)
- [Application connector documentation](https://kyma-project.io/docs/kyma/main/01-overview/main-areas/application-connectivity/ac-01-application-connector/) 
- [Modularization of Kyma](https://github.com/kyma-project/community/tree/main/concepts/modularization)

## Prerequisites

- Access to a k8s cluster. You can use [K3d](https://k3d.io) to get a local cluster for testing, or run against a remote cluster.
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kubebuilder](https://book.kubebuilder.io/)


```bash
# you could use one of the following options

# option 1: using brew
brew install kubebuilder

# option 2: fetch sources directly
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

## Manual `application-connector-manager` installation

1. Clone project

```bash
git clone https://github.com/kyma-project/application-connector-manager.git && cd application-connector-manager/
```

2. Set `application-connector-manager` image name

```bash
export IMG=custom-application-connector-manager:0.0.1
export K3D_CLUSTER_NAME=application-connector-manager-demo
```

3. Build project

```bash
make build
```

4. Build image

```bash
make docker-build
```

5. Push image to registry

<div tabs name="Push image" group="application-connector-installation">
  <details>
  <summary label="k3d">
  k3d
  </summary>

   ```bash
   k3d image import $IMG -c $K3D_CLUSTER_NAME
   ```
  </details>
  <details>
  <summary label="Docker registry">
  Globally available Docker registry
  </summary>

   ```bash
   make docker-push
   ```

  </details>
</div>

6. Deploy

```bash
make deploy
```

## Using `application-connector-manager`

- Create ApplicationConnector instance

```bash
kubectl apply -f config/samples/operator_v1alpha1_applicationconnector.yaml
```

- Delete ApplicationConnector instance

```bash
kubectl delete -f config/samples/operator_v1alpha1_applicationconnector.yaml
```

- Update ApplicationConnector properties

TODO: Provide example of CR update

## Build and install of Application Connector module in modular Kyma on the local k3d cluster

1. Setup local k3d cluster and local Docker registry

```bash
k3d cluster create kyma --registry-create registry.localhost:0.0.0.0:5001
```
2. Add the `etc/hosts` entry to register the local Docker registry under the `registry.localhost` name

```
127.0.0.1 registry.localhost
```

3. Export environment variables (ENVs) pointing to module and the module image registries

```bash
export IMG_REGISTRY=registry.localhost:5001/unsigned/operator-images
export MODULE_REGISTRY=registry.localhost:5001/unsigned
```

4. Build Application Connector module 
```bash
make module-build
```

This command builds an OCI image for Application Connector module and pushes it to the registry and path, as defined in `MODULE_REGISTRY`.

5. Build Application Connector manager image
```bash
make module-image
```

This command builds a Docker image for Application Connector Manager and pushes it to the registry and path, as defined in `MODULE_REGISTRY`.

6. Verify if the module and the manager's image are pushed to the local registry

```bash
curl registry.localhost:5001/v2/_catalog
{"repositories":["unsigned/component-descriptors/kyma.project.io/module/keda","unsigned/operator-images/keda-operator"]}
```

7. Inspect the generated module template

The following are temporary workarounds.

Edit the `template.yaml` file and:

- change `target` to `control-plane`
>**NOTE:** This is only required in the single cluster mode

```yaml
spec:
  target: control-plane
```

- change the existing repository context in `spec.descriptor.component`:
>**NOTE:** Because Pods inside the k3d cluster use the docker-internal port of the registry, it tries to resolve the registry against port 5000 instead of 5001. K3d has registry aliases but module-manager is not part of k3d and thus does not know how to properly alias `registry.localhost:5001`

```yaml
repositoryContexts:                                                                           
- baseUrl: registry.localhost:5000/unsigned                                                   
  componentNameMapping: urlPath                                                               
  type: ociRegistry
```


8. Install modular Kyma on the k3d cluster

This installs the latest versions of `module-manager` and `lifecycle-manager`

You can use the `--template` flag to deploy the ApplicationConnector module manifest from the beginning or apply it using kubectl later.

```bash
kyma alpha deploy  --template=./template.yaml

- Kustomize ready
- Lifecycle Manager deployed
- Module Manager deployed
- Modules deployed
- Kyma CR deployed
- Kyma deployed successfully!

Kyma is installed in version:
Kyma installation took:		18 seconds

Happy Kyma-ing! :)
```

Kyma installation is ready, but no module is activated yet
```bash
kubectl get kymas.operator.kyma-project.io -A
NAMESPACE    NAME           STATE   AGE
kcp-system   default-kyma   Ready   71s
```

Keda Module is a known module, but not activated
```bash
kubectl get moduletemplates.operator.kyma-project.io -A 
NAMESPACE    NAME                  AGE
kcp-system   moduletemplate-keda   2m24s
```

9. Give Module Manager permission to install CustomResourceDefinition (CRD) cluster-wide

>**NOTE:** This is a temporary workaround and is only required in the single-cluster mode

Module-manager must be able to apply CRDs to install modules. In the remote mode (with control-plane managing remote clusters) it gets an administrative kubeconfig, targeting the remote cluster to do so. But in local mode (single-cluster mode), it uses Service Account and does not have permission to create CRDs by default.

Run the following to make sure the module manager's Service Account becomes an administrative role:

```bash
kubectl edit clusterrole module-manager-manager-role
```
add
```yaml
- apiGroups:                                                                                                                  
  - "*"                                                                                                                       
  resources:                                                                                                                  
  - "*"                                                                                                                       
  verbs:                                                                                                                      
  - "*"
```

10. Enable Keda in Kyma

Edit Kyma CR ...

```bash
kubectl edit kymas.operator.kyma-project.io -n kcp-system default-kyma
```
..to add ApplicationConnector module

```yaml
spec:
  modules:
  - name: applicationconnector
```

