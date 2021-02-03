# External JWT Identity Provider Resource

This is a special type of identity provider that is only used via the JWT Reconcile API. This identity provider defines the claims inside the incoming JWT and how they map to fields in the FusionAuth User object.

In order for this identity provider to use the JWT, it also needs the public key or HMAC secret that the JWT was signed with. FusionAuth will verify that the JWT is valid and has not expired. Once the JWT has been validated, FusionAuth will reconcile it to ensure that the User exists and is up-to-date.



[External JWT Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/external-jwt/)

## Example Useage

```hcl
resource "fusionauth_idp_external_jwt" "jwt" {
  claim_map = {
    first_name = "firstName"
    last_name  = "lastName"
    dept       = "RegistrationData"
  }
  debug                         = false
  enabled                       = true
  header_key_parameter          = "kid"
  name                          = "Acme Corp. ADFS"
  oauth2_authorization_endpoint = "https://acme.com/adfs/oauth2/authorize?client_id=cf3b00da-9551-460a-ad18-33232e6cbff0&response_type=code&redirect_uri=https://acme.com/oauth2/redirect"
  oauth2_token_endpoint         = "https://acme.com/adfs/oauth2/token"
  unique_identity_claim         = "email"
}
```

## Argument Reference

* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the Application to apply this configuration to.
    - `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `domains` - (Optional) An array of domains that are managed by this Identity Provider.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `header_key_parameter` - (Required) The name header claim that identifies the public key used to verify the signature. In most cases this be kid or x5t.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `name` - (Required) The name of the Identity Provider.
* `oauth2_authorization_endpoint` - (Optional) The authorization endpoint for this Identity Provider. This value is not utilized by FusionAuth is only provided to be returned by the Lookup Identity Provider API response. During integration you may then utilize this value to perform the browser redirect to the OAuth2 authorize endpoint.
* `oauth2_token_endpoint` - (Optional) TThe token endpoint for this Identity Provider. This value is not utilized by FusionAuth is only provided to be returned by the Lookup Identity Provider API response. During integration you may then utilize this value to complete the OAuth2 grant workflow.
* `unique_identity_claim` - (Required) The name of the claim that represents the unique identify of the User. This will generally be email or the name of the claim that provides the email address.