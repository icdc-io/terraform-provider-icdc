package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceNetworkRead,
		CreateContext: resourceNetworkCreate,
		UpdateContext: resourceNetworkUpdate,
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

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("vpcs/%s/networks/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	tflog.Info(ctx, "Network read url:", map[string]any{"url": url})

	r, err := requestApi("GET", url, nil)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var networtGetResponse *NetworkRequestResponse

	if err = json.Unmarshal(resBody, &networtGetResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	ps, _ := PrettyStruct(networtGetResponse)
	tflog.Info(ctx, "Network read response body:", map[string]any{"response": ps})

	d.SetId(networtGetResponse.Network.Id)
	return diags
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	mtu, _ := strconv.Atoi(d.Get("mtu").(string))
	ipv, _ := strconv.Atoi(d.Get("ip_version").(string))
	enable_dhcp, _ := strconv.ParseBool(d.Get("enable_dhcp").(string))

	cloudNetworkRaw := &CloudNetworkRequest{
		Network: NetworkCreateBody{
			Name:     d.Get("name").(string),
			Mtu:      mtu,
			TenantId: os.Getenv("ACCOUNT"),
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
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudNetworkRaw))

	url := fmt.Sprintf("vpcs/%s/networks", d.Get("vpc_id").(string))
	tflog.Info(ctx, "Network create url:", map[string]any{"url": url})

	r, err := requestApi("POST", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var networkRequestResponse *NetworkRequestResponse

	if err = json.Unmarshal(resBody, &networkRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	ps, _ := PrettyStruct(networkRequestResponse)
	tflog.Info(ctx, "Network create response body:", map[string]any{"response": ps})

	d.SetId(networkRequestResponse.Network.Id)
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	mtu, _ := strconv.Atoi(d.Get("mtu").(string))
	ipv, _ := strconv.Atoi(d.Get("ip_version").(string))
	enable_dhcp, _ := strconv.ParseBool(d.Get("enable_dhcp").(string))

	cloudNetworkRaw := &CloudNetworkRequest{
		Network: NetworkCreateBody{
			Name:     d.Get("name").(string),
			TenantId: os.Getenv("ACCOUNT"),
			Mtu:      mtu,
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
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudNetworkRaw))

	url := fmt.Sprintf("vpcs/%s/networks/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	tflog.Info(ctx, "Network update url:", map[string]any{"url": url})

	r, err := requestApi("PUT", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var networkRequestResponse *NetworkRequestResponse

	if err = json.Unmarshal(resBody, &networkRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	ps, _ := PrettyStruct(networkRequestResponse)
	tflog.Info(ctx, "Network update response body:", map[string]any{"response": ps})

	d.SetId(networkRequestResponse.Network.Id)
	return diags
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("vpcs/%s/networks/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	tflog.Info(ctx, "Network delete url:", map[string]any{"url": url})

	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")

	return diags
}
