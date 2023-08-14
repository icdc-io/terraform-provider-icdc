package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVpcRead,
		CreateContext: resourceVpcCreate,
		UpdateContext: resourceVpcUpdate,
		DeleteContext: resourceVpcDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"router": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVpcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("vpcs/%s", d.Get("id"))
	r, err := requestApi("GET", url, nil)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	log.Println(resBody)
	var vpcRequestResponse *VpcRequestResponse

	if err = json.Unmarshal(resBody, &vpcRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	log.Println(PrettyStruct(vpcRequestResponse))
	log.Println(vpcRequestResponse.Vpc.Id)
	d.SetId(vpcRequestResponse.Vpc.Id)
	return diags
}

func resourceVpcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cloudVpcRaw := &VpcCreateBody{
		Vpc: VpcStructBody{
			Name:     d.Get("name").(string),
			TenantId: os.Getenv("ACCOUNT"),
			Router: RouterCreateBody{
				Name: d.Get("name").(string),
			},
		},
	}

	requestBody, err := json.Marshal(cloudVpcRaw)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudVpcRaw))

	url := "vpcs"
	r, err := requestApi("POST", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var vpcRequestResponse *VpcRequestResponse

	if err = json.Unmarshal(resBody, &vpcRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	fmt.Println(PrettyStruct(vpcRequestResponse))
	log.Println(PrettyStruct(vpcRequestResponse))

	vpc_id := vpcRequestResponse.Vpc.Id
	log.Println(PrettyStruct(vpc_id))
	d.SetId(vpc_id)
	return diags
}

func resourceVpcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cloudVpcRaw := &VpcCreateBody{
		Vpc: VpcStructBody{
			Name:     d.Get("name").(string),
			TenantId: os.Getenv("ACCOUNT"),
			Router: RouterCreateBody{
				Name: d.Get("name").(string),
			},
		},
	}

	requestBody, err := json.Marshal(cloudVpcRaw)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudVpcRaw))

	url := fmt.Sprintf("vpcs/%s", d.Get("id").(string))
	r, err := requestApi("PUT", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var vpcRequestResponse *VpcRequestResponse

	if err = json.Unmarshal(resBody, &vpcRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	fmt.Println(PrettyStruct(vpcRequestResponse))
	vpc_id := vpcRequestResponse.Vpc.Id
	d.SetId(vpc_id)
	return diags
}

func resourceVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("vpcs/%s", d.Get("id").(string))
	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")

	return diags
}
