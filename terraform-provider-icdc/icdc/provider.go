package icdc

import (
	"context"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_PASSWORD", nil),
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_LOCATION", nil),
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_ROLE", nil),
			},
			"platform": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_NAME", nil),
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CPV_ACCOUNT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"icdc_service": resourceService(),
			"icdc_subnet":  resourceSubnet(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	location := d.Get("location").(string)
	account := d.Get("account").(string)
	role := d.Get("role").(string)

	var diags diag.Diagnostics

	var url = fmt.Sprintf("https://login.%s.io/auth/realms/master/protocol/openid-connect/token", d.Get("platform").(string))
	var buf = []byte("username=" + username + "&password=" + password + "&client_id=insights&grant_type=password")
	var jwt JwtToken
	resp, _ := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(buf))
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal([]byte(body), &jwt)

	if err != nil {
		return err
	}

	gateway_url := fmt.Sprintf("https://api.%s.icdc.io/api/compute/v1", location)

	os.Setenv("API_GATEWAY", gateway_url)
	os.Setenv("ROLE", role)
	os.Setenv("AUTH_TOKEN", jwt.AccessToken)
	os.Setenv("LOCATION", location)
	os.Setenv("ACCOUNT", account)

	return nil, diags
}
