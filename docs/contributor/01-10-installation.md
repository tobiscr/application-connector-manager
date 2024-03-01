# Install Application Connector Manager 

- [Install Application Connector Manager](#install-application-connector-manager)
  - [Install Application Connector Manager from the local sources](#install-application-connector-manager-from-the-local-sources)
    - [Prerequisites](#prerequisites)
    - [Procedure](#procedure)
  - [Make targets to run Application Connector Manager locally on k3d](#make-targets-to-run-application-connector-manager-locally-on-k3d)
    - [Run Application Connector Manager on bare k3d](#run-application-connector-manager-on-bare-k3d)
  - [Install Application Connector module on remote Kyma runtime](#install-application-connector-module-on-remote-kyma-runtime)
    - [Prerequisite](#prerequisite)
    - [Procedure](#procedure-1)

Learn how to install the Application Connector module locally (on k3d) or on your remote cluster.

## Install Application Connector Manager from the Local Sources 

### Prerequisites

- Access to a Kubernetes (v1.24 or higher) cluster
- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kubebuilder](https://book.kubebuilder.io/)

### Procedure

You can build and run the Application Connector Manager in the Kubernetes cluster without Kyma.
For the day-to-day development on your machine, you don't always need to have it controlled by Kyma's Lifecycle Manager.

Run the following commands to deploy Application Connector Manager in a target Kubernetes cluster, such as k3d:

1. Clone the project.

   ```bash
   git clone https://github.com/kyma-project/application-connector-manager.git && cd application-connector-manager/
   ```

2. Set the Application Connector Manager image name.

   > NOTE: You can use the local k3d registry or your Docker Hub account to push intermediate images.  
   ```bash
   export IMG=<DOCKER_USERNAME>/custom-application-connector-manager:0.0.1
   ```

3. Test the code.

   ```bash
   make test
   ```
4. Build and push the image to the registry.

   ```bash
   make module-image
   ```
5. Create a target namespace.

   ```bash
   kubectl create ns kyma-system
   ```

6. Deploy Application Connector Manager.

   ```bash
   make deploy
   ```

7. Verify if Application Connector Manager is deployed.

   ```bash
   kubectl get deployments -n kyma-system
   ```

   You should get a result similar to this example:

   ```
   NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
   application-connector-controller-manager   1/1     1            1           20s
   ```

## Make Targets to Run Application Connector Manager Locally on k3d

### Run Application Connector Manager on Bare k3d

When using a local k3d cluster, you can also use the local OCI image registry that comes with it.
Thanks to that, you don't need to push the Application Connector module images to a remote registry and you can test the changes in the Kyma installation set up entirely on your machine.

1. Clone the project.

   ```bash
   git clone https://github.com/kyma-project/application-connector-manager.git && cd application-connector-manager/
   ```
2. Build the manager locally and run it in the k3d cluster.

   ```bash
   make -C hack/local run-without-lifecycle-manager
   ```
3. If you want to clean up the k3d cluster, use the `make -C hack/local stop` make target.

