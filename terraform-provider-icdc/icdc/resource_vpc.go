package icdc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		Read:   resourceVpcRead,
		Create: resourceVpcCreate,
		Update: resourceVpcUpdate,
		Delete: resourceVpcDelete,
		Schema: map[string]*schema.Schema{
			/*"id": {
				Type:     schema.TypeString,
				Computed: true,
			},*/
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceVpcRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVpcCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVpcUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVpcDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*

func resourceVpcCreate(d *schema.ResourceData, m interface{}) error {
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

func resourceVpcUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVpcDelete(d *schema.ResourceData, m interface{}) error {
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
