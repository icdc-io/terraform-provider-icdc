package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
			"ip_version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"dns_nameservers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSubnetRead(d *schema.ResourceData, m interface{}) error {
	responseBody, err := requestApi("GET", fmt.Sprintf("cloud_subnets/%s?expand=resources", d.Id()), nil)

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

	err = d.Set("ip_version", subnet.IpVersion)

	if err != nil {
		return err
	}

	err = d.Set("network_protocol", subnet.NetworkProtocol)

	if err != nil {
		return err
	}

	err = d.Set("dns_nameservers", subnet.DnsNameservers)

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
	responseBody, err := requestApi("GET", "providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager", nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&emsProvider)

	if err != nil {
		return err
	}

	emsProviderId := emsProvider.Resources[0].Id

	/*
		ahrechushkin:
		 workaround
		 https://stackoverflow.com/questions/72402307/interface-conversion-error-while-sending-the-payload-for-post-request-custom-t
	*/

	dnsNameservers := d.Get("dns_nameservers").([]interface{})
	dns := make([]string, len(dnsNameservers))

	for _, dnsNameserver := range dnsNameservers {
		if dnsNameserver != "" {
			dns = append(dns, dnsNameserver.(string))
		}
	}

	// end workaround

	cloudNetworkRaw := &CloudNetworkRequest{
		Action: "create",
		Name:   d.Get("name").(string),
		Subnet: SubnetCreateBody{
			Cidr:            d.Get("cidr").(string),
			IpVersion:       d.Get("ip_version").(int),
			NetworkProtocol: d.Get("network_protocol").(string),
			Name:            d.Get("name").(string),
			DnsNameservers:  dns,
		},
	}

	requestBody, err := json.Marshal(cloudNetworkRaw)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	var response *ServiceRequestResponse

	responseBody, err = requestApi("POST", fmt.Sprintf("providers/%s/cloud_networks", emsProviderId), body)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&response)

	if err != nil {
		return err
	}

	var networkCollection *NetworkCollection

	responseBody, err = requestApi("GET", "cloud_networks?expand=resources&attributes=cloud_subnets", nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&networkCollection)

	if err != nil {
		return err
	}

	time.Sleep(25 * time.Second)

	for _, network := range networkCollection.Resources {
		if network.Name == fmt.Sprintf("%s_%s_%s", os.Getenv("LOCATION"), os.Getenv("ACCOUNT"), d.Get("name").(string)) {
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

	_, err = requestApi("POST", fmt.Sprintf("providers/%s/cloud_networks", d.Get("ems_id").(string)), body)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
