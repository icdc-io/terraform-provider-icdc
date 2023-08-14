package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	tflog.Info(ctx, "Vpc read url:", map[string]any{"url": url})

	r, err := requestApi("GET", url, nil)
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

	ps, _ := PrettyStruct(vpcRequestResponse)
	tflog.Info(ctx, "Vpc read response body:", map[string]any{"response": ps})

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
	tflog.Info(ctx, "Vpc create url:", map[string]any{"url": url})
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

	ps, _ := PrettyStruct(vpcRequestResponse)
	tflog.Info(ctx, "Vpc create response body:", map[string]any{"response": ps})

	d.SetId(vpcRequestResponse.Vpc.Id)
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
	tflog.Info(ctx, "Vpc update url:", map[string]any{"url": url})

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

	ps, _ := PrettyStruct(vpcRequestResponse)
	tflog.Info(ctx, "Vpc update response body:", map[string]any{"response": ps})

	d.SetId(vpcRequestResponse.Vpc.Id)
	return diags
}

func resourceVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("vpcs/%s", d.Get("id").(string))
	tflog.Info(ctx, "Vpc delete url:", map[string]any{"url": url})

	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")

	return diags
}
