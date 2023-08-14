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

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSecurityGroupRead,
		CreateContext: resourceSecurityGroupCreate,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("vpcs/%s/security_groups/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	tflog.Info(ctx, "Security group read url:", map[string]any{"url": url})

	r, err := requestApi("GET", url, nil)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRequestResponse *SecurityGroupRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	ps, _ := PrettyStruct(securityGroupRequestResponse)
	tflog.Info(ctx, "Security group read response body:", map[string]any{"response": ps})

	d.SetId(securityGroupRequestResponse.SecurityGroup.Id)
	return diags
}

func resourceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cloudGroupRaw := &GroupCreateBody{
		SecurityGroup: SecurityGroupBody{
			Name:        d.Get("name").(string),
			TenantId:    os.Getenv("ACCOUNT"),
			Description: d.Get("description").(string),
		},
	}

	requestBody, err := json.Marshal(cloudGroupRaw)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudGroupRaw))

	url := fmt.Sprintf("vpcs/%s/security_groups", d.Get("vpc_id").(string))
	tflog.Info(ctx, "Security group create url:", map[string]any{"url": url})

	r, err := requestApi("POST", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRequestResponse *SecurityGroupRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	ps, _ := PrettyStruct(securityGroupRequestResponse)
	tflog.Info(ctx, "Security group create response body:", map[string]any{"response": ps})

	d.SetId(securityGroupRequestResponse.SecurityGroup.Id)
	return diags
}

func resourceSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cloudGroupRaw := &GroupCreateBody{
		SecurityGroup: SecurityGroupBody{
			Name:        d.Get("name").(string),
			TenantId:    os.Getenv("ACCOUNT"),
			Description: d.Get("description").(string),
		},
	}

	requestBody, err := json.Marshal(cloudGroupRaw)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudGroupRaw))

	url := fmt.Sprintf("vpcs/%s/security_groups/%s", d.Get("vpc_id").(string), d.Get("id").(string))
	tflog.Info(ctx, "Security group update url:", map[string]any{"url": url})

	r, err := requestApi("PUT", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRequestResponse *SecurityGroupRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	ps, _ := PrettyStruct(securityGroupRequestResponse)
	tflog.Info(ctx, "Security group update response body:", map[string]any{"response": ps})

	d.SetId(securityGroupRequestResponse.SecurityGroup.Id)

	return diags
}

func resourceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("%s", d.Get("id").(string))
	tflog.Info(ctx, "Security group delete url:", map[string]any{"url": url})

	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")

	return diags
}
