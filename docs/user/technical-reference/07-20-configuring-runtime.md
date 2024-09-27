# Configuring the Runtime

Runtime Agent periodically requests for the configuration of its Runtime from Compass.
Changes in the configuration for the Runtime are applied by Runtime Agent on the Runtime.

To fetch the Runtime configuration, Runtime Agent calls the [`applicationsForRuntime`](https://github.com/kyma-incubator/compass/blob/master/components/director/pkg/graphql/schema.graphql) query offered by the Compass component called Director.
The response for the query contains Applications assigned for the Runtime.
Each Application contains only credentials that are valid for the Runtime that called the query. Runtime Agent uses the credentials to create Secrets used by Application Gateway.
Each Runtime Agent can fetch the configurations for Runtimes that belong to its tenant.

The data mapping between Compass Director and Kyma looks as follows:

| **Director (Compass)**    | **Kyma**                      |
|---------------------------|-------------------------------|
| Application               | Application CR                |
| API Bundle                | Service in the Application CR |
| API Definition            | Entry under the service       |

## Application name restriction

The names of Applications assigned to the Runtime must be unique in Kyma Runtime.\
If the name of an Application fetched from Compass is not unique, the synchronization fails.  

## Normalization of Application names 

The Runtime Agent can normalize the names of Applications fetched from Compass by converting them to lowercase and removing special characters and spaces.\
This feature is controlled by the **isNormalized** label, which can be set on the Runtime in Compass\
When the Runtime is initially labeled with `isNormalized=true`, Runtime Agent normalizes the names of Applications.\
When the Runtime is initially labeled with `isNormalized=false`, or if the Runtime does not contain such a label, Runtime Agent doesn't normalize the names.\
If the normalized Application names fetched from Compass are not unique, the synchronization process fails.\
Because of that, it is required to use unique names for Applications with lowercase characters only and not use any special characters or spaces in their names.

## Runtime labeling in Compass

Runtime Agent reports back to the Director the Runtime-specific [LabelDefinitions](https://github.com/kyma-incubator/compass/blob/master/docs/compass/03-04-labels.md#labeldefinitions), which represent Runtime configuration, together with their values.
Runtime-specific LabelDefinitions are Event Gateway URL and Runtime Console URL.
