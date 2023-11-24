package icdc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDnsRecordRead,
		Create: resourceDnsRecordCreate,
		Update: resourceDnsRecordUpdate,
		Delete: resourceDnsRecordDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of related domain zone",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of supported dns types: A, AAAA, CNAME, NS, MX, SRV, TXT",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for your dns record without zone",
			},
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Payload for dns record, e.g. IP-address, hostname, etc.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Required for SRV, MX records",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Required for SRV records",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Required for SRV records",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Time to live: how long to cache a query before requesting a new one",
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDnsRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDnsRecordCreate(d *schema.ResourceData, m interface{}) error {
	record := DnsRecord{
		Payload: DnsRecordBody{
			Type: d.Get("type").(string),
			Name: d.Get("name").(string),
			Data: d.Get("data").(string),
			Ttl:  d.Get("ttl").(int),
		},
	}

	recordPtr := &record
	recordPtr.setAdditionalFields(d)

	requestBody, err := json.Marshal(record)

	if err != nil {
		return fmt.Errorf("can't serealize record payload: %+v", record)
	}

	payload := bytes.NewBuffer(requestBody)
	requestAddRecordUrl := fmt.Sprintf("api/dns/v1/zones/%s/records", d.Get("zone").(string))
	responseAddRecordBody, err := requestApi("POST", requestAddRecordUrl, payload)

	if err != nil {
		return fmt.Errorf("the problem occurs creating DNS record: %s", err)
	}

	var addRecordResponse *responseAddDnsRecord
	err = responseAddRecordBody.Decode(addRecordResponse)

	if err != nil {
		return fmt.Errorf("the problem occurs deserealizing response of create dns record %s", requestBody)
	}

	err = d.Set("group", addRecordResponse.Data.Group)

	if err != nil {
		return errors.New("can't set 'group' property to a new dns record")
	}

	requestRecordsListUrl := fmt.Sprintf("api/dns/v1/zones/%s/records", d.Get("zone").(string))
	responseRecordsListBody, err := requestApi("GET", requestRecordsListUrl, nil)

	if err != nil {
		return fmt.Errorf("the problem occurs retrieve DNS record: %s", err)
	}

	var recordsList *responseListDnsRecords
	err = responseRecordsListBody.Decode(recordsList)

	if err != nil {
		return fmt.Errorf("the problem occurs deserealizing response of create dns records list of %s", d.Get("zone").(string))
	}

	for _, dnsRecord := range recordsList.Data {
		if dnsRecord.Name == d.Get("name").(string) &&
			dnsRecord.Data == d.Get("data").(string) {
			id := fmt.Sprintf("%s.%s", dnsRecord.Id, dnsRecord.Name)
			d.SetId(id)

			return nil
		}
	}

	return fmt.Errorf("the records wasn't created successfully")
}

func resourceDnsRecordUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("unsupported action: update dns record")
}
func resourceDnsRecordDelete(d *schema.ResourceData, m interface{}) error {
	_, err := requestApi("DELETE", fmt.Sprintf("api/dns/v1/zones/%s/records/%s", d.Get("zone").(string), d.Get("id").(string)), nil)

	if err != nil {
		return fmt.Errorf("the problem occure delete dns record %+v", d)
	}

	d.SetId("")

	return nil
}
