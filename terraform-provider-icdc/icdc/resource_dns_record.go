package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDnsRecordRead,
		Create: resourceDnsRecordCreate,
		Update: resourceDnsRecordUpdate,
		Delete: resourceDnsRecordDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "ID of related domain zone",
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "One of supported dns types: A, AAAA, CNAME, NS, MX, SRV, TXT",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Name for your dns record without zone",
			},
			"data": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Payload for dns record, e.g. IP-address, hostname, etc.",
			},
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Required for SRV, MX records",
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Required for SRV records",
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Required for SRV records",
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				Description: "Time to live: how long to cache a query before requesting a new one",
			},
			"group": &schema.Schema{
				Type:			schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDnsRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDnsRecordCreate(d *schema.ResourceData, m interface{}) error {
	if strings.ToLower(d.Get("type").(string)) == "srv" {
		var dnsRecordRaw AddSrvRecord	

		dnsRecordRaw.Record.Name = d.Get("name").(string)
		dnsRecordRaw.Record.Type = d.Get("type").(string)
		dnsRecordRaw.Record.Data = d.Get("data").(string)
		dnsRecordRaw.Record.Ttl = d.Get("ttl").(int)
		if d.Get("weight").(int) != 0 {
			dnsRecordRaw.Record.Weight = d.Get("weight").(int)
		}

		if d.Get("port").(int) != 0 {
			dnsRecordRaw.Record.Port = d.Get("port").(int) 
		}

		if d.Get("priority").(int) != 0 {
			dnsRecordRaw.Record.Priority = d.Get("priority").(int)
		}

		requestBody, err := json.Marshal(dnsRecordRaw)

		if err != nil {
			return err
		}
	
		body := bytes.NewBuffer(requestBody)
	
		_, err = requestApi("POST", fmt.Sprintf("api/dns/v1/zones/%s/records", d.Get("zone").(string)), body)
	} else if strings.ToLower(d.Get("type").(string)) == "mx" {
		var dnsRecordRaw AddSrvRecord	
		dnsRecordRaw.Record.Name = d.Get("name").(string)
		dnsRecordRaw.Record.Type = d.Get("type").(string)
		dnsRecordRaw.Record.Data = d.Get("data").(string)
		dnsRecordRaw.Record.Ttl = d.Get("ttl").(int)
		dnsRecordRaw.Record.Priority = d.Get("priority").(int)

		requestBody, err := json.Marshal(dnsRecordRaw)

		if err != nil {
			return err
		}
	
		body := bytes.NewBuffer(requestBody)
	
		_, err = requestApi("POST", fmt.Sprintf("api/dns/v1/zones/%s/records", d.Get("zone").(string)), body)

	} else {
		var dnsRecordRaw AddDnsRecord

		dnsRecordRaw.Record.Name = d.Get("name").(string)
		dnsRecordRaw.Record.Type = d.Get("type").(string)
		dnsRecordRaw.Record.Data = d.Get("data").(string)
		dnsRecordRaw.Record.Ttl = d.Get("ttl").(int)

		requestBody, err := json.Marshal(dnsRecordRaw)

		if err != nil {
			return err
		}
	
		body := bytes.NewBuffer(requestBody)
	
		_, err = requestApi("POST", fmt.Sprintf("api/dns/v1/zones/%s/records", d.Get("zone").(string)), body)
	}

	dnsRecordsList, err := requestApi("GET", fmt.Sprintf("api/dns/v1/zones/%s/records", d.Get("zone").(string)), nil)

	if err != nil {
		return err
	}

	var response *DnsRecordResponse                                     

	err = dnsRecordsList.Decode(&response)

	for i := range response.Data {
		if response.Data[i].Data == d.Get("data").(string) && response.Data[i].Name == d.Get("name") {
			id := fmt.Sprintf("%s.%s", response.Data[i].Id, response.Data[i].Name)
			d.Set("group", response.Data[i].Group)
			d.SetId(id)

			return nil
		}
	}

	return nil
}

func resourceDnsRecordUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}
func resourceDnsRecordDelete(d *schema.ResourceData, m interface{}) error {
	_, err := requestApi("DELETE", fmt.Sprintf("api/dns/v1/zones/%s/records/%s", d.Get("zone").(string), d.Get("id").(string)), nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

