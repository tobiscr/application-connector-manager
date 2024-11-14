# Configuring the Runtime

Runtime Agent periodically requests for the configuration of the Kyma Runtime from UCL.

To fetch the Kyma Runtime configuration, Runtime Agent calls the [`applicationsForRuntime`](https://github.com/kyma-incubator/compass/blob/master/components/director/pkg/graphql/schema.graphql) query offered by the component called UCL Director.

The response for the query contains Applications (an `Application` represents an external system) assigned for the Kyma Runtime.

Each `Application` contains credentials that are only valid and unique for the querying Kyma Runtime (another requestor won't be able to use them). Runtime Agent stores the credentials in a Secrets which will be used by Application Gateway for establishing a trusted outbound communication to the external system.

This data mapping shows how the retrieved configuration of an `Application` from UCL Director is stored in the Kyma Runtime:

| **UCL Director**    | **Kyma Runtime**                    |
|---------------------------|-------------------------------|
| Application               | Application CR                |
| API Bundle                | Service in the Application CR |
| API Definition            | Entry under the service       |


## Application name

The name of the `Application` is used as a key within the Application Connectivity Module and has special requirements:


### Uniqueness

The names of Applications assigned to the Runtime must be unique in Kyma Runtime. If they are not unique, the synchronization will fail.


### Normalization of Application names 

The Runtime Agent can normalize the names of Applications fetched from UCL Director by converting them to lowercase and removing special characters and spaces.

This feature is controlled by the **isNormalized** label, which can be set on the Runtime in UCL.

When the Runtime is initially labeled with `isNormalized=true`, Runtime Agent normalizes the names of Applications. When the Runtime is initially labeled with `isNormalized=false`, or if the Runtime does not contain such a label, Runtime Agent doesn't normalize the names.

The normalization can lead to non-unqiue Application names if names were differentiated only be special characters or just by different lower-upper case letters.


## Reporting Kyma Runtime configuration to UCL

Runtime Agent reports back to the Director the Runtime-specific [LabelDefinitions](https://github.com/kyma-incubator/compass/blob/master/docs/compass/03-04-labels.md#labeldefinitions), which represent Runtime configuration, together with their values.
Runtime-specific LabelDefinitions are Event Gateway URL and Runtime Console URL.
