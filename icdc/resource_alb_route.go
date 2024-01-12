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
							Type:     schema.TypeString,
							Optional: true,
							Default:  "get",
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
	hcList := healthcheck.List()

	hcEnabled := len(hcList) > 0
	hc := Healthcheck{}
	if hcEnabled {
		hc.assignParams(d)
	}

	routeBody := AlbRoute{
		Name:               d.Get("name").(string),
		Hostname:           d.Get("hostname").(string),
		Path:               d.Get("path").(string),
		TargetPort:         d.Get("target_port").(int),
		Insecure:           d.Get("insecure").(string),
		TlsTermination:     d.Get("tls_termination").(string),
		CloudGatewayId:     cloudGwIdByName(d.Get("cloudgw_name").(string)),
		IpVersion:          strconv.Itoa(d.Get("ip_version").(int)),
		Services:           services,
		HealthcheckEnabled: hcEnabled,
		Healthcheck:        hc,
	}

	payload, err := json.Marshal(routeBody)
	body := bytes.NewBuffer(payload)
	if err != nil {
		return fmt.Errorf("can't marshalling into json %+v", routeBody)
	}

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
	requestUrl := fmt.Sprintf("api/traefik_manager/v1/routes/%s", d.Id())

	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return fmt.Errorf("error fetching alb route: %s", err)
	}

	var routeResponse AlbRouteApi

	err = responseBody.Decode(&routeResponse)

	route := routeResponse.Route

	if err != nil {
		return fmt.Errorf("error decoding alb route response: %s", err)
	}

	err = d.Set("name", route.Name)

	if err != nil {
		return fmt.Errorf("error setting value: %s", err)
	}

	err = d.Set("hostname", route.Hostname)

	if err != nil {
		return fmt.Errorf("error setting value: %s", err)
	}

	err = d.Set("path", route.Path)

	if err != nil {
		return fmt.Errorf("error setting value: %s", err)
	}

	err = d.Set("target_port", route.TargetPort)

	if err != nil {
		return fmt.Errorf("error setting value: %s", err)
	}

	err = d.Set("insecure", route.Insecure)

	if err != nil {
		return fmt.Errorf("error setting value: %s", err)
	}

	err = d.Set("tls_termination", route.TlsTermination)

	if err != nil {
		return fmt.Errorf("error setting value: %s", err)
	}

	// ahrechushkin: in v1.0.0 we doesn't support in-place update healthcheck

	return nil
}

func resourceAlbRouteUpdate(d *schema.ResourceData, m interface{}) error {

	/* ahrechushkin: update action will be implemented after implementing PATCH action in alb-api
	var route AlbRoute

	if d.HasChange("name") {
		route.Name = d.Get("name").(string)
	}

	if d.HasChange("hostname") {
		route.Hostname = d.Get("hostname").(string)
	}

	if d.HasChange("path") {
		route.Path = d.Get("path").(string)
	}

	if d.HasChange("target_port") {
		route.TargetPort = d.Get("target_port").(int)
	}

	if d.HasChange("insecure") {
		route.Insecure = d.Get("insecure").(string)
	}

	if d.HasChange("tls_termination") {
		route.TlsTermination = d.Get("tls_termination").(string)
	}

	if d.HasChange("services") {
		route.Services = servicesByExtId(d.Get("services").([]interface{}))
	}

	if d.HasChange("ip_version") {
		route.IpVersion = strconv.Itoa(d.Get("ip_version").(int))
	}

	fmt.Printf("[---DEBUG---] route changes %+v", route)

	route.CloudGatewayId = cloudGwIdByName(d.Get("cloudgw_name").(string))
	if d.Get("healthcheck_enabled") != nil {
		hc := Healthcheck{}
		hc.assignParams(d)
		route.Healthcheck = hc
	}

	routeBody := AlbRouteApi{Route: route}

	payload, err := json.Marshal(routeBody)
	body := bytes.NewBuffer(payload)
	if err != nil {
		return fmt.Errorf("can't marshalling into json %+v", routeBody)
	}

	requestUrl := fmt.Sprintf("api/traefik_manager/v1/routes/%s", d.Id())
	_, err = requestApi("PUT", requestUrl, body)

	if err != nil {
		return err
	}
	*/

	return nil
}

func resourceAlbRouteDelete(d *schema.ResourceData, m interface{}) error {
	requestUrl := fmt.Sprintf("api/traefik_manager/v1/routes/%s", d.Id())
	_, err := requestApi("DELETE", requestUrl, nil)

	if err != nil {
		return fmt.Errorf("error deleting alb_route resource, %s", err)
	}

	d.SetId("")
	return nil
}
