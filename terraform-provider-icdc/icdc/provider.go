package icdc

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "userid (email)",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ICDC_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "user password",
				DefaultFunc: schema.EnvDefaultFunc("ICDC_PASSWORD", nil),
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "operated locations name",
				DefaultFunc: schema.EnvDefaultFunc("ICDC_LOCATION", nil),
			},
			"sso_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "sso url, if your sso url is not `login.icdc.io`",
				Default:     "login.icdc.io/auth",
			},
			"sso_realm": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "input your sso realm (by default: master)",
				Default:     "master",
			},
			"sso_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "input your sso client_id",
				Default:     "insights",
			},
			"auth_group": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "setup auth_group. it contain two parts: ACCOUNT.ROLE",
				DefaultFunc: schema.EnvDefaultFunc("ICDC_AUTH_GROUP", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"icdc_service":        resourceService(),
			"icdc_instance_group": resourceServiceV2(),
			"icdc_network":        resourceNetwork(),
			"icdc_security_group": resourceSecurityGroup(),
			"icdc_dns_zone":       resourceDnsZone(),
			"icdc_dns_record":     resourceDnsRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"icdc_template": dataSourceICDCTemplate(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	location := d.Get("location").(string)
	authGroup := d.Get("auth_group").(string)
	ssoUrl := d.Get("sso_url").(string)
	ssoRealm := d.Get("sso_realm").(string)
	ssoClientId := d.Get("sso_client_id").(string)

	var diags diag.Diagnostics

	jwt, err := getJwt(username, password, ssoUrl, ssoRealm, ssoClientId)

	if err != nil {
		return nil, err
	}

	account := strings.Split(authGroup, ".")[0]
	role := strings.Split(authGroup, ".")[1]

	jwtClaims, diags := jwt.Claims()
	gatewayUrl := jwtClaims.External.Locations[location]

	os.Setenv("API_GATEWAY", gatewayUrl)
	os.Setenv("ROLE", role)
	os.Setenv("AUTH_TOKEN", jwt.AccessToken)
	os.Setenv("LOCATION", location)
	os.Setenv("ACCOUNT", account)
	os.Setenv("AUTH_GROUP", authGroup)

	return nil, diags
}
