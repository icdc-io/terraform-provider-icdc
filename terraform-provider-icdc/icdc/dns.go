package icdc

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DnsZone struct {
	Name string `json:"name"`
}

type AddDnsZone struct {
	DnsZone `json:"zone"`
}

type AddDnsZoneResponse struct {
	DnsZone `json:"data"`
}

type DnsRecord struct {
	Payload DnsRecordBody `json:"record"`
}

type DnsRecordBody struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	Group    string `json:"group"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Port     int    `json:"port"`
	Ttl      int    `json:"ttl"`
}

type DnsRecordDetails struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Ttl      int    `json:"ttl"`
	Group    string `json:"group"`
	Data     string `json:"data"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Port     int    `json:"port"`
}

type responseListDnsRecords struct {
	Data []DnsRecordDetails `json:"data"`
}

type responseAddDnsRecord struct {
	Data DnsRecordBody `json:"data"`
}

func (r *DnsRecord) setAdditionalFields(d *schema.ResourceData) {
	switch t := strings.ToLower(r.Payload.Type); t {
	case "mx":
		(*r).Payload.Priority = d.Get("priority").(int)
	case "srv":
		(*r).Payload.Priority = d.Get("priority").(int)
		(*r).Payload.Weight = d.Get("weight").(int)
		(*r).Payload.Port = d.Get("port").(int)
	}
}
