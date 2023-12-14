package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlbRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlbRouteCreate,
		Read:   resourceAlbRouteRead,
		Update: resourceAlbRouteUpdate,
		Delete: resourceAlbRouteDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloudgw_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "alb-default",
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},
			"target_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
			"insecure": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
			},
			"tls_termination": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"healthcheck": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/",
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "http",
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"follow_redirects": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"method": {
							Type: schema.TypeString,
							Optional: true,
							Default: "get",
						},
					},
				},
			},
		},
	}
}

func resourceAlbRouteCreate(d *schema.ResourceData, m interface{}) error {

	services := servicesByExtId(d.Get("services").([]interface{}))

	healthcheck := d.Get("healthcheck").(*schema.Set)
	hc_list := healthcheck.List()

	hc_enabled := len(hc_list) > 0
	hc := Healthcheck{}
	if hc_enabled {
		hc.assignParams(d)
	}

	routeBody := AlbRoute{
		Name:              d.Get("name").(string),
		Hostname:          d.Get("hostname").(string),
		Path:              d.Get("path").(string),
		TargetPort:        d.Get("target_port").(int),
		Insecure:          d.Get("insecure").(string),
		TlsTermination:    d.Get("tls_termination").(string),
		CloudGatewayId:    cloudGwIdByName(d.Get("cloudgw_name").(string)),
		IpVersion:         strconv.Itoa(d.Get("ip_version").(int)),
		Services:          services,
		HealtcheckEnabled: hc_enabled,
		Healthcheck:       hc,
	}

	payload, err := json.Marshal(routeBody)
	body := bytes.NewBuffer(payload)
	if err != nil {
		return fmt.Errorf("can't marshalling into json %+v", routeBody)
	}

	fmt.Printf("DEBUG PAYLOAD %s", payload)

	requestUrl := "api/traefik_manager/v1/routes"
	responseBody, err := requestApi("POST", requestUrl, body)

	if err != nil {
		return err
	}

	var route AlbRoute

	err = responseBody.Decode(&route)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(route.Id))

	return nil
}

func resourceAlbRouteRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAlbRouteUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAlbRouteDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
