package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read:   resourceNetworkRead,
		Create: resourceNetworkCreate,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mtu": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_dhcp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_nameservers": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("vpcs/%s/networks/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	r, err := requestApi("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error getting api services: %w", err)
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var networtGetResponse *NetworkRequestResponse

	if err = json.Unmarshal(resBody, &networtGetResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(networtGetResponse))
	log.Println(PrettyStruct(networtGetResponse))
	d.SetId(networtGetResponse.Network.Id)
	return nil
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	mtu, _ := strconv.Atoi(d.Get("mtu").(string))
	ipv, _ := strconv.Atoi(d.Get("ip_version").(string))
	enable_dhcp, _ := strconv.ParseBool(d.Get("enable_dhcp").(string))

	cloudNetworkRaw := &CloudNetworkRequest{
		Network: NetworkCreateBody{
			Name: d.Get("name").(string),
			Mtu:  mtu,
			Subnet: SubnetParams{
				Name:           d.Get("name").(string),
				IpVersion:      ipv,
				Cidr:           d.Get("cidr").(string),
				GatewayIp:      d.Get("gateway_ip").(string),
				EnableDhcp:     enable_dhcp,
				DnsNameservers: d.Get("dns_nameservers").(string),
			},
		},
	}

	requestBody, err := json.Marshal(cloudNetworkRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudNetworkRaw))

	url := fmt.Sprintf("vpcs/%s/networks", d.Get("vpc_id").(string))
	r, err := requestApi("POST", url, body)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var networkRequestResponse *NetworkRequestResponse

	if err = json.Unmarshal(resBody, &networkRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(networkRequestResponse))
	NetworkId := networkRequestResponse.Network.Id
	d.SetId(NetworkId)
	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	mtu, _ := strconv.Atoi(d.Get("mtu").(string))
	ipv, _ := strconv.Atoi(d.Get("ip_version").(string))
	enable_dhcp, _ := strconv.ParseBool(d.Get("enable_dhcp").(string))

	cloudNetworkRaw := &CloudNetworkRequest{
		Network: NetworkCreateBody{
			Name: d.Get("name").(string),
			Mtu:  mtu,
			Subnet: SubnetParams{
				Name:           d.Get("name").(string),
				IpVersion:      ipv,
				Cidr:           d.Get("cidr").(string),
				GatewayIp:      d.Get("gateway_ip").(string),
				EnableDhcp:     enable_dhcp,
				DnsNameservers: d.Get("dns_nameservers").(string),
			},
		},
	}

	requestBody, err := json.Marshal(cloudNetworkRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudNetworkRaw))

	url := fmt.Sprintf("vpcs/%s/networks/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	r, err := requestApi("PUT", url, body)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var networkRequestResponse *NetworkRequestResponse

	if err = json.Unmarshal(resBody, &networkRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(networkRequestResponse))
	NetworkId := networkRequestResponse.Network.Id
	d.SetId(NetworkId)
	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {

	url := fmt.Sprintf("vpcs/%s/networks/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil

}
