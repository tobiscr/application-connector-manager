# Runtime Agent

> [!WARNING]
> Runtime Agent is currently not integrated into the Application Connector module. Proceed with caution and consult the respective documentation for the Application Connector module's current configuration and functionality.

Runtime Agent is a Kyma component that connects to [Compass](https://github.com/kyma-incubator/compass). It is an integral part of every Kyma Runtime in the [Compass mode](README.md) and it fetches the latest configuration from Compass. It also provides Runtime-specific information that is displayed in the Compass UI, such as Runtime UI URL, and it provides Compass with Runtime configuration, such as Event Gateway URL, that should be passed to an Application. To learn more, read the section on [configuring the Runtime](technical-reference/07-20-configuring-runtime.md).

The main responsibilities of the component are:

- Establishing a trusted connection between the Kyma Runtime and Compass
- Renewing a trusted connection between the Kyma Runtime and Compass
- Synchronizing with the [Director](https://github.com/kyma-incubator/compass/blob/master/docs/compass/02-01-components.md#director) by fetching new Applications from the Director and creating them in the Runtime, and removing from the Runtime Applications that no longer exist in the Director.

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
