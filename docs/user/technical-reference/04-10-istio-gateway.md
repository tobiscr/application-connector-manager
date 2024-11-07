# Istio Ingress Gateway

The Application Connector Module relies on an Istio Ingress Gateway as endpoint for incoming requests from external systems.

On a managed SAP Kyma cluster (SKR) is the Kyma Istio Module per deault pre-installed. Alterantively, it's also possible to install Istio manually.

The Istio gateway will be created during the Application Connector Module installation.


## DNS name

The DNS name of the Ingress Gatway is cluster-dependent. fir SKR clsuters it follows the format `gateway.{cluster-dns}`.


## Security

### Client Certificates

For external systems which were automatically integrated by UCL, the Application Connector uses the mutual TLS protocol with Client Authentication enabled. As a result, anyone attempting to connect to Application Connector must present a valid client certificate, which is dedicated to a specific Application. In this way, the traffic is fully encrypted, and the client has a valid identity.

### TLS Certificate Verification

By default, the TLS certificate verification is enabled when sending data and requests to every application.
You can [disable the TLS certificate verification](tutorials/01-50-disable-tls-certificate-verification.md) in the communication between Kyma and an application to allow Kyma to send requests and data to an unsecured application. Disabling the certificate verification can be useful in certain testing scenarios.
