package fusionauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FA_DOMAIN", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FA_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fusionauth_lambda":               newLambda(),
			"fusionauth_application":          newApplication(),
			"fusionauth_user":                 newUser(),
			"fusionauth_theme":                newTheme(),
			"fusionauth_tenant":               newTenant(),
			"fusionauth_email":                newEmail(),
			"fusionauth_key":                  newKey(),
			"fusionauth_webhook":              newWebhook(),
			"fusionauth_group":                newGroup(),
			"fusionauth_application_role":     newApplicationRole(),
			"fusionauth_idp_open_id_connect":  newIDPOpenIDConnect(),
			"fusionauth_idp_google":           newIDPGoogle(),
			"fusionauth_registration":         newRegistration(),
			"fusionauth_system_configuration": resourceSystemConfiguration(),
			"fusionauth_form_field":           resourceFormField(),
			"fusionauth_form":                 resourceForm(),
			"fusionauth_idp_apple":            resourceIDPApple(),
			"fusionauth_imported_key":         resourceImportedKey(),
			"fusionauth_idp_external_jwt":     resourceIDPExternalJWT(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"fusionauth_application":      dataSourceApplication(),
			"fusionauth_tenant":           dataSourceTenant(),
			"fusionauth_application_role": dataSourceApplicationRole(),
			"fusionauth_idp":              dataSourceIDP(),
		},
		ConfigureFunc: configureClient,
	}
}
