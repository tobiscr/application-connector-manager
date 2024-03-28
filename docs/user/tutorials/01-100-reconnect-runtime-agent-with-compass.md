# Reconnect Runtime Agent with Compass

> [!WARNING]
> Runtime Agent is currently not integrated into the Application Connector module. Proceed with caution and consult the respective documentation for the Application Connector module's current configuration and functionality.

This tutorial shows how to reconnect Runtime Agent with Compass after the established connection was lost.

## Prerequisites

- [Compass](https://github.com/kyma-incubator/compass)
- [ConfigMap created](../tutorials/01-90-configure-runtime-agent-with-compass.md)

## Steps

To force Runtime Agent to reconnect using the parameters from the Secret, delete the CompassConnection CR:

```bash
kubectl delete compassconnection compass-connection
```

After the Connection CR is removed, Runtime Agent tries to connect to Compass using the token from the Secret.
