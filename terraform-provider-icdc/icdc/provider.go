package icdc

import (
	"context"

	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"os"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "username": &schema.Schema{
        Type:        schema.TypeString,
        Required:		true,
        DefaultFunc: schema.EnvDefaultFunc("ICDC_USERNAME", nil),
      },
      "password": &schema.Schema{
        Type:        schema.TypeString,
        Required:		true,
        Sensitive:   true,
        DefaultFunc: schema.EnvDefaultFunc("ICDC_PASSWORD", nil),
      },
			"location": &schema.Schema{
        Type:        schema.TypeString,
        Required:		true,
        Sensitive:   true,
        DefaultFunc: schema.EnvDefaultFunc("ICDC_LOCATION", nil),
      },
			"group": &schema.Schema{
        Type:        schema.TypeString,
        Required:		true,
        Sensitive:   true,
        DefaultFunc: schema.EnvDefaultFunc("ICDC_GROUP", nil),
      },
    },
		ResourcesMap: map[string]*schema.Resource{
			"icdc_service_request": resourceServiceRequest(),
			//"icdc_service": resourceService(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

type IcdcToken struct {
	ApiGateway string
	Group string
	Jwt string
}

type JwtToken struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
 	username := d.Get("username").(string)
  password := d.Get("password").(string)
	location := d.Get("location").(string)
	group := d.Get("group").(string)

	var diags diag.Diagnostics

	var url = "https://login.icdc.io/auth/realms/master/protocol/openid-connect/token"
	var buf = []byte("username=" + username + "&password=" + password + "&client_id=insights&grant_type=password")
	var jwt JwtToken
	resp, _ := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(buf))
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &jwt)

	gateway_url := fmt.Sprintf("https://api.%s.icdc.io/api/compute/v1", location)

	os.Setenv("ICDC_API_GATEWAY", gateway_url)
	os.Setenv("ICDC_GROUP", group)
	os.Setenv("ICDC_TOKEN", jwt.AccessToken)

  return nil, diags
}