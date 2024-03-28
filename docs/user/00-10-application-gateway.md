# Application Gateway

Application Gateway is an intermediary component between a Function or a microservice and an external API.
It [proxies the requests](technical-reference/04-20-application-gateway.md) from Functions and microservices in Kyma to external APIs based on the configuration stored in Secrets.

Application Gateway also supports redirects for the request flows in which the URL host remains unchanged. For more details, see [Response rewriting](technical-reference/07-10-application-gateway-details.md#response-rewriting).

## Supported API Authentication for Application CR

Application Gateway can call services which are not secured, or are secured with:

- [Basic Authentication](https://tools.ietf.org/html/rfc7617)
- [OAuth](https://tools.ietf.org/html/rfc6750)
- [OAuth 2.0 mTLS](https://datatracker.ietf.org/doc/html/rfc8705)
- Client certificates

Additionally, Application Gateway supports cross-site request forgery (CSRF) tokens as an optional layer of API protection.

Application Gateway calls the registered APIs accordingly, based on the security type specified for the API in the Application CR.

## Provide a Custom Access Token

Application Gateway overrides the registered API's security type if it gets a request which contains the **Access-Token** header. In such a case, Application Gateway rewrites the token from the **Access-Token** header into an OAuth-compliant **Authorization** header and forwards it to the target API.

This mechanism is suited for implementations in which an external application handles user authentication.

See how to [pass an access token in a request header](02-10-pass-access-token-in-request-header.md).
