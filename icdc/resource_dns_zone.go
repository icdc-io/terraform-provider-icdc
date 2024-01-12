package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDnsZone() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDnsZoneRead,
		Create: resourceDnsZoneCreate,
		Update: resourceDnsZoneUpdate,
		Delete: resourceDnsZoneDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "dns zone name",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceDnsZoneRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDnsZoneCreate(d *schema.ResourceData, m interface{}) error {
	dnsZone := AddDnsZone{
		DnsZone{
			Name: d.Get("name").(string),
		},
	}

	requestBody, err := json.Marshal(dnsZone)

	if err != nil {
		return fmt.Errorf("cant serealize dns zone %+v", dnsZone)
	}

	payload := bytes.NewBuffer(requestBody)
	responseAddZoneBody, err := requestApi("POST", "api/dns/v1/zones", payload)

	if err != nil {
		return fmt.Errorf("the probrem occurs creating dns zone, %s", err)
	}

	var addDnsZoneResponse *AddDnsZoneResponse

	err = responseAddZoneBody.Decode(&addDnsZoneResponse)

	if err != nil {
		return fmt.Errorf("the problem occurs decoding dns zone response body: %s", err)
	}

	d.SetId(addDnsZoneResponse.Name)
	return nil
}

func resourceDnsZoneUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("unsupported action: update dns zone")
}

func resourceDnsZoneDelete(d *schema.ResourceData, m interface{}) error {

	_, err := requestApi("DELETE", fmt.Sprintf("api/dns/v1/zones/%s", d.Get("id").(string)), nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
