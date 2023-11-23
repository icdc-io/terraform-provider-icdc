package icdc

import (
	"context"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Your username(email) into ICDC platform",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ICDC_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Your password",
				DefaultFunc: schema.EnvDefaultFunc("ICDC_PASSWORD", nil),
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Name of operated location",
				DefaultFunc: schema.EnvDefaultFunc("ICDC_LOCATION", nil),
			},
			"sso_url": &schema.Schema{
				Type:				schema.TypeString,
				Optional: 	true,
				Sensitive: 	true,
				Description: "input sso url if your sso is not login.icdc.io",
				Default: 		"login.icdc.io/auth",
			},
			"sso_realm": &schema.Schema{
				Type:				schema.TypeString,
				Optional: 	true,
				Sensitive: 	true,
				Description: "input your sso realm",
				Default: 		"master",
			},
			"sso_client_id": &schema.Schema{
				Type:				schema.TypeString,
				Optional: 	true,
				Sensitive: 	true,
				Description: "input your sso client id",
				Default: 		"insights",
			},
			"auth_group": &schema.Schema{
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
			"icdc_subnet":         resourceSubnet(),
			"icdc_security_group": resourceSecurityGroup(),
			"icdc_dns_zone":			 resourceDnsZone(),
			"icdc_dns_record":		 resourceDnsRecord(),
		},
		DataSourcesMap:       map[string]*schema.Resource{
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

	var url = fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token", ssoUrl, ssoRealm)
	var buf = []byte("username=" + username + "&password=" + password + "&client_id=" + ssoClientId + "&grant_type=password")
	var jwt JwtToken
	resp, _ := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(buf))
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal([]byte(body), &jwt)

	defer resp.Body.Close()

	if err != nil {
		return nil, diags
	}

	account := strings.Split(authGroup, ".")[0]
	role := strings.Split(authGroup, ".")[1]

	gatewayUrl := findGatewayUrl(jwt.AccessToken, location)

	os.Setenv("API_GATEWAY", gatewayUrl)
	os.Setenv("ROLE", role)
	os.Setenv("AUTH_TOKEN", jwt.AccessToken)
	os.Setenv("LOCATION", location)
	os.Setenv("ACCOUNT", account)

	return nil, diags
}

type IcdcClaims struct {
	External struct {
		Locations map[string]string `json:"locations"`
	} `json:"external"`
}

func findGatewayUrl(token string, location string) string {

	base64TokenClaims := strings.Split(token, ".")[1]

	base64TokenClaims += strings.Repeat("=", ((4 - len(base64TokenClaims)%4) % 4))

	fmt.Println(base64TokenClaims)
	rawClaims, err := base64.StdEncoding.DecodeString(base64TokenClaims)

	if err != nil {
		panic(err)
	}

	var icdcClaims IcdcClaims

	err = json.Unmarshal(rawClaims, &icdcClaims)

	if err != nil {
		panic(err)
	}

	return icdcClaims.External.Locations[location]
}
