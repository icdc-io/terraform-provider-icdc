package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type Network struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Subnets []Subnet `json:"cloud_subnets"`
}
type Subnet struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	EmsRef          string   `json:"ems_ref"`
	EmsId           string   `json:"ems_id"`
	CloudNetworkId  string   `json:"cloud_network_id"`
	Cidr            string   `json:"cidr"`
	Gateway         string   `json:"gateway"`
	IpVersion       int      `json:"ip_version"`
	NetworkProtocol string   `json:"network_protocol"`
	DnsNameservers  []string `json:"dns_nameservers"`
	NetworkRouterId string   `json:"network_router_id"`
}

type SubnetCreateBody struct {
	Cidr            string   `json:"cidr"`
	IpVersion       int      `json:"ip_version"`
	NetworkProtocol string   `json:"network_protocol"`
	Name            string   `json:"name"`
	DnsNameservers  []string `json:"dns_nameservers"`
}

type NetworkCollection struct {
	Resources []Network `json:"resources"`
}

type CloudNetworkRequest struct {
	Action string           `json:"action"`
	Name   string           `json:"name"`
	Subnet SubnetCreateBody `json:"subnet"`
}

type DeleteNetworkRequest struct {
	Action string `json:"action"`
	Id		 string `json:"id"`
}

type EmsProvider struct {
	Resources []struct {
		Id string `json:"id"`
	} `json:"resources"`
}

func resourceSubnetRead(d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 100 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cloud_subnets/%s?expand=resources", os.Getenv("API_GATEWAY"), d.Id()), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	var subnet *Subnet

	err = json.NewDecoder(r.Body).Decode(&subnet)
	
	if err != nil {
		return err
	}

	d.Set("name", subnet.Name)
	d.Set("ems_ref", subnet.EmsRef)
	d.Set("ems_id", subnet.EmsId)
	d.Set("cloud_network_id", subnet.CloudNetworkId)
	d.Set("cidr", subnet.Cidr)
	d.Set("gateway", subnet.Gateway)
	d.Set("ip_version", subnet.IpVersion)
	d.Set("network_protocol", subnet.NetworkProtocol)
	d.Set("dns_nameservers", subnet.DnsNameservers)
	d.Set("network_router_id", subnet.NetworkRouterId)

	return nil
}

func resourceSubnetCreate(d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 100 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager", os.Getenv("API_GATEWAY")), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	var emsProvider *EmsProvider

	err = json.NewDecoder(r.Body).Decode(&emsProvider)
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

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/providers/%s/cloud_networks", os.Getenv("API_GATEWAY"), emsProviderId), body)

	/*
			{
				https://api.ycz.icdc.io/api/compute/v1/providers/18000000000003/cloud_networks/
		    "action": "create",
		    "name": "ahrechushkin",
		    "subnet": {
		        "cidr": "11.15.13.1/24",
		        "ip_version": 4,
		        "network_protocol": "ipv4",
		        "dns_nameservers": [
		            "178.172.238.131"
		        ],
		        "name": "ahrechushkin"
		    }

				{"results":[{"success":true,"message":"Network and subnet created"}]}
			}

			So, it means that we need to fetch provider id and than use it to create subnet. Or we can pass provider id as parameter. \
			Cause we've limitation of one provider per location.
	*/

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err = client.Do(req)
	if err != nil {
		return err
	}

	var response *ServiceRequestResponse

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return err
	}

	/* ahrechushkin:
	 	we need to fetch all networks and find one with name equal to name of network we've just created.
		But we can't fetch it immediately, cause it's not refreshed from ems.
	*/

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/cloud_networks?expand=resources&attributes=cloud_subnets", os.Getenv("API_GATEWAY")), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err = client.Do(req)
	if err != nil {
		return err
	}

	var networkCollection *NetworkCollection

	err = json.NewDecoder(r.Body).Decode(&networkCollection)

	file, _ := json.MarshalIndent(networkCollection, "", "  ")
	_ = ioutil.WriteFile("/tmp/subnet_collection.json", file, 0644)

	time.Sleep(25 * time.Second)

	for _, network := range networkCollection.Resources {
		if network.Name == fmt.Sprintf("%s_%s_%s", os.Getenv("LOCATION"), os.Getenv("ACCOUNT"), d.Get("name").(string)) {
			d.Set("cloud_network_id", network.Id)
			d.Set("name", network.Subnets[0].Name)
			d.Set("ems_ref", network.Subnets[0].EmsRef)
			d.Set("network_router_id", network.Subnets[0].NetworkRouterId)
			d.Set("ems_id", network.Subnets[0].EmsId)
			d.SetId(network.Subnets[0].Id)
		}
	}
	return nil
}

func resourceSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubnetDelete(d *schema.ResourceData, m interface{}) error {
	/*
		POST: providers/18000000000003/cloud_networks
		{"action":"delete","id":"18000000000133"}
	*/


	client := &http.Client{Timeout: 100 * time.Second}

	deleteNetworkRequest := &DeleteNetworkRequest{
		Action: "delete",
		Id: d.Get("cloud_network_id").(string),
	}

	requestBody, err := json.Marshal(deleteNetworkRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/providers/%s/cloud_networks", os.Getenv("API_GATEWAY"), d.Get("ems_id")), body)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
