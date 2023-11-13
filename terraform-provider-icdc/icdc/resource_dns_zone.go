package icdc

import (
	"fmt"
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDnsZone() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDnsZoneRead,
		Create: resourceDnsZoneCreate,
		Update: resourceDnsZoneUpdate,
		Delete: resourceDnsZoneDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Description: "dns zone name",
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDnsZoneRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDnsZoneCreate(d *schema.ResourceData, m interface{}) error {
	var dnsZoneRaw AddDnsZone
	dnsZoneRaw.Zone.Name = d.Get("name").(string)

	requestBody, err := json.Marshal(dnsZoneRaw)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	var response *AddDnsZoneResponse                                     

	responseBody, err := requestApi("POST", "api/dns/v1/zones", body)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&response)

	if err != nil {
		return err
	}

	d.SetId(response.Data.Name)

	return nil
}

func resourceDnsZoneUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDnsZoneDelete(d *schema.ResourceData, m interface{}) error {

	_, err := requestApi("DELETE", fmt.Sprintf("api/dns/v1/zones/%s", d.Get("id").(string)), nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

