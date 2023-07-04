package icdc

import (
	"context"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_PASSWORD", nil),
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_LOCATION", nil),
			},
			"location_number": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_LOCATION_NUMBER", nil),
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_ROLE", nil),
			},
			"platform": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_NAME", nil),
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_ACCOUNT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"icdc_service":        resourceService(),
			"icdc_network":        resourceNetwork(),
			"icdc_security_group": resourceSecurityGroup(),
			"icdc_vpc":            resourceVPC(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	location := d.Get("location").(string)
	location_number := d.Get("location_number").(string)
	account := d.Get("account").(string)
	role := d.Get("role").(string)

	var diags diag.Diagnostics

	var url = fmt.Sprintf("https://login.%s.io/auth/realms/master/protocol/openid-connect/token", d.Get("platform").(string))
	var buf = []byte("username=" + username + "&password=" + password + "&client_id=insights&grant_type=password")
	var jwt JwtToken
	resp, _ := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(buf))
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal([]byte(body), &jwt)

	defer resp.Body.Close()

	if err != nil {
		return nil, diags
	}

	gateway_url := "http://10.207.1.79:3000/api/v1/vpcs"

	os.Setenv("API_GATEWAY", gateway_url)
	os.Setenv("ROLE", role)
	os.Setenv("AUTH_TOKEN", jwt.AccessToken)
	os.Setenv("LOCATION", location)
	os.Setenv("LOCATION_NUMBER", location_number)
	os.Setenv("ACCOUNT", account)

	return nil, diags
}
