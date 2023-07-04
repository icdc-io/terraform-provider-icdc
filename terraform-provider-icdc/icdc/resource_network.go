package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceNetworkRead,
		Create:        resourceNetworkCreate,
		Update:        resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
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
			"ipv6_address_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_nameservers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

/*
	   test := "111111111111111111"
	   fmt.Println(test)
	   fmt.Println(d)
	   return nil
	   --
	   	//tflog.Info(ctx, "11111111111111111111111111111111")
	//tflog.Debug(ctx, "222222222222222222")
	//log.Printf("[DEBUG] 333333333333333333333")
*/

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	// Warning or errors can be collected in a slice type
	//var diags diag.Diagnostics
	fmt.Println("--resourceNetworkCreate--")
	fmt.Println(d)
	/*&{map[cidr:0xc000336b40 dns_nameservers:0xc000337180 enable_dhcp:0xc0003368c0 gateway_ip:0xc000336c80 id:0xc0003363c0
	ip_version:0xc000336f00 ipv6_address_mode:0xc000336dc0 mtu:0xc000336640 name:0xc000336500 network_id:0xc000336a00 vpc_id:0xc000336780]*/

	dnsNameservers := d.Get("dns_nameservers").([]interface{})
	dns := make([]string, len(dnsNameservers))

	for _, dnsNameserver := range dnsNameservers {
		if dnsNameserver != "" {
			dns = append(dns, dnsNameserver.(string))
		}
	}

	mtu, _ := strconv.Atoi(d.Get("mtu").(string))
	ipv, _ := strconv.Atoi(d.Get("ip_version").(string))
	enable_dhcp, _ := strconv.ParseBool(d.Get("enable_dhcp").(string))

	cloudNetworkRaw := &CloudNetworkRequest{
		Network: NetworkCreateBody{
			Name:     d.Get("name").(string),
			Mtu:      mtu,
			TenantId: os.Getenv("ACCOUNT"),
			Subnet: SubnetParams{
				IpVersion:       ipv,
				Cidr:            d.Get("cidr").(string),
				GatewayIp:       d.Get("gateway_ip").(string),
				Ipv6AddressMode: d.Get("ipv6_address_mode").(string),
				EnableDhcp:      enable_dhcp,
				DnsNameservers:  dns,
			},
		},
	}

	requestBody, err := json.Marshal(cloudNetworkRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudNetworkRaw))

	responseBody, err := requestApi("POST", d.Get("vpc_id").(string), body)
	if err != nil {
		return fmt.Errorf("error creating network: %w", err)
	}
	// fmt.Println(networkRequestResponse)
	var networkRequestResponse *NetworkRequestResponse
	if err = responseBody.Decode(&networkRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(networkRequestResponse)
	log.Println(PrettyStruct(networkRequestResponse))

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	// Warning or errors can be collected in a slice type
	//var diags diag.Diagnostics
	fmt.Println("--resourceNetworkUpdate--")
	fmt.Println(d)
	/*&{map[cidr:0xc000336b40 dns_nameservers:0xc000337180 enable_dhcp:0xc0003368c0 gateway_ip:0xc000336c80 id:0xc0003363c0
	ip_version:0xc000336f00 ipv6_address_mode:0xc000336dc0 mtu:0xc000336640 name:0xc000336500 network_id:0xc000336a00 vpc_id:0xc000336780]*/

	dnsNameservers := d.Get("dns_nameservers").([]interface{})
	dns := make([]string, len(dnsNameservers))

	for _, dnsNameserver := range dnsNameservers {
		if dnsNameserver != "" {
			dns = append(dns, dnsNameserver.(string))
		}
	}

	mtu, _ := strconv.Atoi(d.Get("mtu").(string))
	ipv, _ := strconv.Atoi(d.Get("ip_version").(string))
	enable_dhcp, _ := strconv.ParseBool(d.Get("enable_dhcp").(string))

	cloudNetworkRaw := &CloudNetworkRequest{
		Network: NetworkCreateBody{
			Name:     d.Get("name").(string),
			Mtu:      mtu,
			TenantId: os.Getenv("ACCOUNT"),
			Subnet: SubnetParams{
				IpVersion:       ipv,
				Cidr:            d.Get("cidr").(string),
				GatewayIp:       d.Get("gateway_ip").(string),
				Ipv6AddressMode: d.Get("ipv6_address_mode").(string),
				EnableDhcp:      enable_dhcp,
				DnsNameservers:  dns,
			},
		},
	}

	requestBody, err := json.Marshal(cloudNetworkRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudNetworkRaw))

	responseBody, err := requestApi("PUT", d.Get("vpc_id").(string)+"/networks"+d.Get("network_id").(string), body)
	if err != nil {
		return fmt.Errorf("error creating network: %w", err)
	}
	// fmt.Println(networkRequestResponse)
	var networUpdateResponse *NetworkUpdateResponse
	if err = responseBody.Decode(&networUpdateResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	//fmt.Println(networkRequestResponse)
	log.Println(PrettyStruct(networUpdateResponse))

	return nil
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourceSubnetRead(d *schema.ResourceData, m interface{}) error {
	fmt.Println("resourceSubnetRead")
	responseBody, err := requestApi("GET", fmt.Sprintf("cloud_subnets/%s?expand=resources", d.Id()), nil)
	fmt.Println("resourceSubnetRead")
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

	fmt.Println("resourceSubnetCreate")
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
/*
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
	fmt.Println("resourceSubnetUpdate")
	return nil
}

func resourceSubnetDelete(d *schema.ResourceData, m interface{}) error {
	fmt.Println("resourceSubnetDelete")
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
*/
