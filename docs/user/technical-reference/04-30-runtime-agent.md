# Runtime Agent

Runtime Agent is a client of UCL. Applications which are registered for Kyma in UCL are fetched from the UCL backend and integrated into Kyma.

To allow a birectional communication, it uploads the Kyma runtime configuration (e.g. the Event Gateway URL) to the UCL backend, that should be used by the external system to call workloads in Kyma. To learn more, read the section on [configuring the Runtime](./runtime-agent/07-20-configuring-runtime.md).


Other responsibilities of the Runtime Agent are:

- Establishing (or renewing) a trusted connection between the Kyma Runtime and UCL backend
- Regularly synchronizing with the [UCL Director](https://github.com/kyma-incubator/compass/blob/master/docs/compass/02-01-components.md#director) by fetching new Applications from the UCL Director aor removing those which no longer exist.


## Workflow

![Runtime Agent architecture](../../assets/ra-runtime-agent-workflow.svg)


1. Runtime Agent fetches the certificate from the [Connector](https://github.com/kyma-incubator/compass/blob/master/docs/connector/02-01-connector-service.md) to initialize connection with Compass.

2. Runtime Agent stores the certificate and key for the Connector and the Director in the Secret.

3. Runtime Agent synchronizes the Runtime with the Director. It does so by:

- fetching new Applications from the Director and creating them in the Runtime
- removing from the Runtime the Applications that no longer exist in the Director.

4. Runtime Agent labels the Runtime data in the Director with the Event Gateway URL and the Console URL of the Kyma cluster. These URLs are displayed in the Compass UI.

5. Runtime Agent renews the certificate for the Connector and the Director to maintain connection with Compass. This only happens if the remaining validity period for the certificate passes a certain threshold.


## Useful Links

If you're interested in learning more about Runtime Agent, follow these links to:

- Perform some simple and more advanced tasks:

  - [Enable Kyma with Runtime Agent](02-20-enable-kyma-with-runtime-agent.md)
  - [Establish a secure connection with Compass](tutorials/01-60-establish-secure-connection-with-compass.md)
  - [Maintain a secure connection with Compass](tutorials/01-70-maintain-secure-connection-with-compass.md)
  - [Revoke a client certificate (RA)](tutorials/01-80-revoke-client-certificate.md)
  - [Configure Runtime Agent with Compass](tutorials/01-90-configure-runtime-agent-with-compass.md)
  - [Reconnect Runtime Agent with Compass](tutorials/01-100-reconnect-runtime-agent-with-compass.md)

- Analyze Runtime Agent specification and configuration files:

  - [Compass Connection](resources/06-20-compassconnection.md) custom resource (CR)
  - [Connection with Compass](technical-reference/05-20-connection-with-compass.md)

- Understand technicalities behind the Runtime Agent implementation:

  - [Runtime Agent workflow](technical-reference/04-30-runtime-agent-workflow.md)
  - [Configuring the Runtime](technical-reference/07-20-configuring-runtime.md)

