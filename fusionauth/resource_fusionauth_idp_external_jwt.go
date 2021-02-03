package fusionauth

import (
	"encoding/json"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type IDPExternalJWTProviderBody struct {
	IdentityProvider fusionauth.ExternalJWTIdentityProvider `json:"identityProvider"`
}

type IDPExternalJWTAppConfig struct {
	CreateRegistration bool `json:"createRegistration"`
	Enabled            bool `json:"enabled,omitempty"`
}

func resourceIDPExternalJWT() *schema.Resource {
	return &schema.Resource{
		Create: createIDPExternalJWT,
		Read:   readIDPExternalJWT,
		Update: updateIDPExternalJWT,
		Delete: deleteIdentityProvider,
		Schema: map[string]*schema.Schema{
			"application_configuration": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "The configuration for each Application that the identity provider is enabled for.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"create_registration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging intres.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines if this identity provider is enabled for the Application specified by the applicationId key.",
						},
					},
				},
			},
			"claim_map": {
				Type:         schema.TypeMap,
				Optional:     true,
				Description:  "A map of incoming claims to User fields, User data or Registration data. The key of the map is the incoming claim name from the configured identity provider.",
				ValidateFunc: validateClaimMap,
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.",
			},
			"domains": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of domains that are managed by this Identity Provider.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"header_key_parameter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name header claim that identifies the public key used to verify the signature. In most cases this be kid or x5t.",
			},
			// TODO: Not implemented in client
			// "keys": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "A map of public keys used to verify JWT signatures issued from the configured Identity Provider. The key is the key identifier, this may be referred to as the kid or for X.509 certificates the x5t claim may be used.\nThe map may contain one entry with an empty map key. When provided this key will be used when no header claim is provided to indicate which public key should be used to verify the signature. Generally speaking this will only be used when the Identity Provider issues JWTs without a key identifier in the header.",
			// },
			"lambda_reconcile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Identity Provider.",
			},
			"oauth2_authorization_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authorization endpoint for this Identity Provider. This value is not utilized by FusionAuth is only provided to be returned by the Lookup Identity Provider API response. During integration you may then utilize this value to perform the browser redirect to the OAuth2 authorize endpoint.",
			},
			"oauth2_token_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The token endpoint for this Identity Provider. This value is not utilized by FusionAuth is only provided to be returned by the Lookup Identity Provider API response. During integration you may then utilize this value to complete the OAuth2 grant workflow.",
			},
			"unique_identity_claim": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the claim that represents the unique identify of the User. This will generally be email or the name of the claim that provides the email address.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func createIDPExternalJWT(data *schema.ResourceData, i interface{}) error {
	o := buildIDPExternalJWT(data)

	b, err := json.Marshal(o)
	if err != nil {
		return err
	}

	client := i.(Client)

	bb, err := createIdenityProvider(b, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return err
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func readIDPExternalJWT(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	b, err := readIdenityProvider(data.Id(), client)
	if err != nil {
		return err
	}

	var ipb IDPExternalJWTProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceDataFromIDPExternalJWT(data, ipb.IdentityProvider)
}

func updateIDPExternalJWT(data *schema.ResourceData, i interface{}) error {
	o := buildIDPExternalJWT(data)

	b, err := json.Marshal(o)
	if err != nil {
		return err
	}

	client := i.(Client)
	bb, err := updateIdenityProvider(b, data.Id(), client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return err
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func buildIDPExternalJWT(data *schema.ResourceData) IDPExternalJWTProviderBody {
	idp := fusionauth.ExternalJWTIdentityProvider{
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Name: data.Get("name").(string),
			Type: fusionauth.IdentityProviderType_ExternalJWT,
		},
		Domains:            handleStringSlice("domains", data),
		HeaderKeyParameter: data.Get("header_key_parameter").(string),
		// TODO: handle keys
		Oauth2: fusionauth.IdentityProviderOauth2Configuration{
			AuthorizationEndpoint: data.Get("oauth2_authorization_endpoint").(string),
			TokenEndpoint:         data.Get("oauth2_token_endpoint").(string),
		},
		UniqueIdentityClaim: data.Get("unique_identity_claim").(string),
	}

	if x, ok := data.GetOk("claim_map"); ok {
		m := make(map[string]string)
		cm := x.(map[string]interface{})
		for k, v := range cm {
			m[k] = v.(string)
		}
		idp.ClaimMap = m
	}

	ac := buildIDPExternalJWTAppConfig("application_configuration", data)
	idp.ApplicationConfiguration = ac
	return IDPExternalJWTProviderBody{IdentityProvider: idp}
}

func buildIDPExternalJWTAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	s := data.Get(key)
	set, ok := s.(*schema.Set)
	if !ok {
		return m
	}
	l := set.List()
	for _, x := range l {
		ac := x.(map[string]interface{})
		aid := ac["application_id"].(string)
		oc := IDPExternalJWTAppConfig{
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
		}
		m[aid] = oc
	}
	return m
}

func buildResourceDataFromIDPExternalJWT(data *schema.ResourceData, res fusionauth.ExternalJWTIdentityProvider) error {
	if err := data.Set("claim_map", res.ClaimMap); err != nil {
		return fmt.Errorf("idpExternalJwt.claim_map: %s", err.Error())
	}
	if err := data.Set("debug", res.Debug); err != nil {
		return fmt.Errorf("idpExternalJwt.debug: %s", err.Error())
	}
	if err := data.Set("domains", res.Domains); err != nil {
		return fmt.Errorf("idpExternalJwt.domains: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return fmt.Errorf("idpExternalJwt.enabled: %s", err.Error())
	}
	if err := data.Set("header_key_parameter", res.HeaderKeyParameter); err != nil {
		return fmt.Errorf("idpExternalJwt.header_key_parameter: %s", err.Error())
	}

	// TODO: get keys
	// if err := data.Set("keys", res.); err != nil {
	// 	return fmt.Errorf("idpExternalJwt.keys: %s", err.Error())
	// }
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return fmt.Errorf("idpExternalJwt.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return fmt.Errorf("idpExternalJwt.name: %s", err.Error())
	}
	if err := data.Set("oauth2_authorization_endpoint", res.Oauth2.AuthorizationEndpoint); err != nil {
		return fmt.Errorf("idpExternalJwt.oauth2_authorization_endpoint: %s", err.Error())
	}
	if err := data.Set("oauth2_token_endpoint", res.Oauth2.TokenEndpoint); err != nil {
		return fmt.Errorf("idpExternalJwt.oauth2_token_endpoint: %s", err.Error())
	}
	if err := data.Set("unique_identity_claim", res.UniqueIdentityClaim); err != nil {
		return fmt.Errorf("idpExternalJwt.unique_identity_claim: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]IDPExternalJWTAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(res.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return fmt.Errorf("idpExternalJwt.application_configuration: %s", err.Error())
	}
	return nil
}

func validateClaimMap(i interface{}, k string) (warnings []string, errors []error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be map", k))
		return
	}

	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected value for %q to be string", k))
			return
		}
		switch s {
		case "birthDate", "firstName", "lastName", "fullName", "middleName", "mobilePhone", "imageUrl", "timezone", "UserData", "RegistrationData":
			continue
		default:
			errors = append(errors, fmt.Errorf("valid values for %q are: %q", k, []string{"birthDate", "firstName", "lastName", "fullName", "middleName", "mobilePhone", "imageUrl", "timezone", "UserData", "RegistrationData"}))
		}
	}
	return warnings, errors
}
