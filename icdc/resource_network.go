package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read:          resourceNetworkRead,
		CreateContext: resourceNetworkCreate,
		Update:        resourceNetworkUpdate,
		Delete:        resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:                  schema.TypeString,
				Required:              true,
				Description:           "name of your vpc network",
				DiffSuppressOnRefresh: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return convertName(old) == convertName(new)
				},
			},
			"mtu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1420,
				Description: "the maximum transmission unit size (default is 1420)",
			},
			"subnet": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dns_nameserver": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	networkName := convertName(d.Get("name").(string))

	subnet := d.Get("subnet").(*schema.Set)
	subnet_list := subnet.List()
	subnet_map := subnet_list[0].(map[string]interface{})
	subnetName := networkName

	//subnetName := (convertName(subnet_map["name"].(string)))
	dns_nameserver := subnet_map["dns_nameserver"].(string)
	cidr := subnet_map["cidr"].(string)

	ipVersion, networkProtocol := fetchNetworkProtocol(cidr)

	network := addNetworkBody{
		Action: "create",
		Name:   networkName,
		Mtu:    d.Get("mtu").(int),
		Subnet: addSubnetBody{
			Name:            subnetName,
			Cidr:            cidr,
			DnsNameservers:  []string{dns_nameserver},
			IpVersion:       ipVersion,
			NetworkProtocol: networkProtocol,
		},
	}

	requestBody, err := json.Marshal(network)

	if err != nil {
		msg := fmt.Errorf("can't marshalling network into json %+v", network)
		return diag.FromErr(msg)
	}

	payload := bytes.NewBuffer(requestBody)

	providerId, err := getProviderId()

	if err != nil {
		msg := fmt.Errorf("can't fetch provider id: %s", err)
		return diag.FromErr(msg)
	}

	//(ahrechushkin): create network request returns information about automation task
	// 				  with fields {success: and message: }
	//                TODO: add also check the result of task for sure
	_, err = requestApi("POST", fmt.Sprintf("api/compute/v1/providers/%s/cloud_networks", providerId), payload)

	if err != nil {
		msg := fmt.Errorf("can't create a new network: %s", err)
		return diag.FromErr(msg)
	}

	//(ahrechushkin): we need to wait for subnet created
	//                cause we use manageiq-api instead of ovn api (in this version of app)

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		f, err := subnetCreated(subnetName, cidr, providerId)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		if f {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cloud subnet is not created"))
	})

	if err != nil {
		return diag.FromErr(err)
	}

	createdNetwork, err := getNetworkObject(networkName, cidr, providerId)

	if err != nil {
		return diag.FromErr(err)
	}

	createdSubnet := map[string]string{
		"id":             createdNetwork.Subnets[0].Id,
		"name":           createdNetwork.Subnets[0].Name,
		"cidr":           createdNetwork.Subnets[0].Cidr,
		"gateway":        createdNetwork.Subnets[0].Gateway,
		"dns_nameserver": createdNetwork.Subnets[0].DnsNameservers[0],
	}

	subnetInterface := make([]interface{}, 1)
	subnetInterface[0] = createdSubnet

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", createdNetwork.Name)
	d.Set("subnet", subnetInterface)
	d.SetId(createdNetwork.Id)

	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	pId, err := getProviderId()

	if err != nil {
		return fmt.Errorf("can't fetch provider_id: %s", err)
	}

	requestUrl := fmt.Sprintf("api/compute/v1/providers/%s/cloud_networks", pId)

	payload := deleteNetworkBody{
		Action: "delete",
		Id:     d.Id(),
	}

	requestBody, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("can't marshalling delete network payload")
	}
	body := bytes.NewBuffer(requestBody)
	_, err = requestApi("POST", requestUrl, body)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	requestUrl := fmt.Sprintf("api/compute/v1/cloud_networks/%s?expand=resources&attributes=cloud_subnets", d.Id())

	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return fmt.Errorf("can't fetch network details [%s]: %s", d.Id(), err)
	}

	var network *Network
	err = responseBody.Decode(&network)

	if err != nil {
		return fmt.Errorf("cant parse response into network %+v object: %s", responseBody, err)
	}
	createdSubnet := map[string]string{
		"id":             network.Subnets[0].Id,
		"name":           network.Subnets[0].Name,
		"cidr":           network.Subnets[0].Cidr,
		"gateway":        network.Subnets[0].Gateway,
		"dns_nameserver": network.Subnets[0].DnsNameservers[0],
	}

	subnetInterface := make([]interface{}, 1)
	subnetInterface[0] = createdSubnet

	err = d.Set("subnet", subnetInterface)

	if err != nil {
		return fmt.Errorf("can't apply changed subnet values: %s", err)
	}

	d.Set("name", network.Name)

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("deprecated feature")
}
