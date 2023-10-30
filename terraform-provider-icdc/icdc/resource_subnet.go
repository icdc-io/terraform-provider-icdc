package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
//	"os"
//	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSubnetRead,
		Create: resourceSubnetCreate,
		Update: resourceSubnetUpdate,
		Delete: resourceSubnetDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_name": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ems_ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ems_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_nameserver": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSubnetRead(d *schema.ResourceData, m interface{}) error {
	responseBody, err := requestApi("GET", fmt.Sprintf("api/compute/v1/cloud_subnets/%s?expand=resources", d.Id()), nil)

	if err != nil {
		return err
	}

	var subnet *Subnet

	err = responseBody.Decode(&subnet)

	if err != nil {
		return err
	}

	err = d.Set("name", subnet.Name)

	if err != nil {
		return err
	}

	err = d.Set("ems_ref", subnet.EmsRef)

	if err != nil {
		return err
	}

	err = d.Set("ems_id", subnet.EmsId)

	if err != nil {
		return err
	}

	err = d.Set("cloud_network_id", subnet.CloudNetworkId)

	if err != nil {
		return err
	}

	err = d.Set("cidr", subnet.Cidr)

	if err != nil {
		return err
	}

	err = d.Set("gateway", subnet.Gateway)

	if err != nil {
		return err
	}

	err = d.Set("network_protocol", subnet.NetworkProtocol)

	if err != nil {
		return err
	}

	dnsNameserver := ""

	if len(subnet.DnsNameservers) > 0 {
		dnsNameserver = subnet.DnsNameservers[0]
	}

	err = d.Set("dns_nameserver", dnsNameserver)

	if err != nil {
		return err
	}

	err = d.Set("network_router_id", subnet.NetworkRouterId)

	if err != nil {
		return err
	}

	return nil
}

func resourceSubnetCreate(d *schema.ResourceData, m interface{}) error {
	var emsProvider *EmsProvider
	responseBody, err := requestApi("GET", "api/compute/v1/providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager", nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&emsProvider)

	if err != nil {
		return err
	}

	emsProviderId := emsProvider.Resources[0].Id

	dnsNameserver := d.Get("dns_nameserver").(string)
	dnsNameservers := []string{ dnsNameserver }

	ipVersion := 4

	if d.Get("network_protocol") == "ipv6" {
		ipVersion = 6
	}

	cloudNetworkRaw := &CloudNetworkRequest{
		Action: "create",
		Name:   d.Get("network_name").(string),
		Subnet: SubnetCreateBody{
			Cidr:            d.Get("cidr").(string),
			IpVersion:       ipVersion,
			NetworkProtocol: d.Get("network_protocol").(string),
			Name:            d.Get("network_name").(string),
			DnsNameservers:  dnsNameservers,
		},
	}

	requestBody, err := json.Marshal(cloudNetworkRaw)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	var response *ServiceRequestResponse

	responseBody, err = requestApi("POST", fmt.Sprintf("api/compute/v1/providers/%s/cloud_networks", emsProviderId), body)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&response)

	if err != nil {
		return err
	}

	var networkCollection *NetworkCollection

	responseBody, err = requestApi("GET", "api/compute/v1/cloud_networks?expand=resources&attributes=cloud_subnets", nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&networkCollection)

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	for _, network := range networkCollection.Resources {
		if network.Subnets[0].Cidr == d.Get("cidr") {
			err := d.Set("cloud_network_id", network.Id)

			if err != nil {
				return err
			}
	
			err = d.Set("name", network.Subnets[0].Name)

			if err != nil {
				return err
			}

			err = d.Set("ems_ref", network.Subnets[0].EmsRef)

			if err != nil {
				return err
			}

			err = d.Set("network_router_id", network.Subnets[0].NetworkRouterId)

			if err != nil {
				return err
			}

			err = d.Set("ems_id", network.Subnets[0].EmsId)

			if err != nil {
				return err
			}

			d.SetId(network.Subnets[0].Id)
		}
	}
	return nil
}

func resourceSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubnetDelete(d *schema.ResourceData, m interface{}) error {
	deleteNetworkRequest := &DeleteRequest{
		Action: "delete",
		Id:     d.Get("cloud_network_id").(string),
	}

	requestBody, err := json.Marshal(deleteNetworkRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("api/compute/v1/providers/%s/cloud_networks", d.Get("ems_id").(string)), body)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
