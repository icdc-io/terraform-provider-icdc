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

func resourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSecurityGroupRuleRead,
		CreateContext: resourceSecurityGroupRuleCreate,
		UpdateContext: resourceSecurityGroupRuleUpdate,
		DeleteContext: resourceSecurityGroupRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_range_max": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_range_min": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityGroupRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cloudGroupRuleRaw := &RuleCreateBody{
		SecurityGroupRule: SecurityGroupRuleBody{
			Direction:       d.Get("direction").(string),
			Ethertype:       d.Get("ethertype").(string),
			PortRangeMin:    d.Get("port_range_min").(string),
			PortRangeMax:    d.Get("port_range_max").(string),
			Protocol:        d.Get("protocol").(string),
			RemoteGroupId:   d.Get("remote_group_id").(string),
			SecurityGroupId: d.Get("security_group_id").(string),
		},
	}

	requestBody, err := json.Marshal(cloudGroupRuleRaw)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudGroupRuleRaw))
	url := fmt.Sprintf("security_groups/%s/rules", d.Get("security_group_id").(string))
	r, err := requestApi("POST", url, body)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRuleRequestResponse *SecurityGroupRuleRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRuleRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	fmt.Println(PrettyStruct(securityGroupRuleRequestResponse))
	log.Println(PrettyStruct(securityGroupRuleRequestResponse))

	//sgroup_rule_id := securityGroupRuleRequestResponse.SecurityGroupRule.Id
	// log.Println(PrettyStruct(sgroup_rule_id))
	d.SetId(securityGroupRuleRequestResponse.SecurityGroupRule.Id)

	return nil
}

func resourceSecurityGroupRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("rules/%s", d.Get("id").(string))
	r, err := requestApi("GET", url, nil)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRuleRequestResponse *SecurityGroupRuleRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRuleRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	log.Println(PrettyStruct(securityGroupRuleRequestResponse))
	log.Println(securityGroupRuleRequestResponse.SecurityGroupRule.Id)
	d.SetId(securityGroupRuleRequestResponse.SecurityGroupRule.Id)

	return nil
}

func resourceSecurityGroupRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSecurityGroupRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	url := fmt.Sprintf("security_groups/%s/rules/%s", d.Get("security_group_id").(string), d.Get("id").(string))
	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")

	return nil
}
