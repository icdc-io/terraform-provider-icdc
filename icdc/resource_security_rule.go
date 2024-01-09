package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

func resourceSecurityRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityRuleCreate,
		ReadContext:   resourceSecurityRuleRead,
		UpdateContext: resourceSecurityRuleUpdate,
		DeleteContext: resourceSecurityRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ems_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"direction": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"egress", "ingress"}, true)),
			},
			"port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"protocol": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"", "icmp", "tcp", "udp"}, true)),
			},
			"network_protocol": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "ipv4",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ipv4", "ipv6"}, true)),
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"remote_ip_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	defer ctx.Done()
	var diags diag.Diagnostics

	var rangeMin, rangeMax string

	ranges := strings.Split(d.Get("port_range").(string), "-")

	if len(ranges) == 1 {
		rangeMin = ranges[0]
	} else {
		rangeMin = ranges[0]
		rangeMax = ranges[1]
	}

	securityGroup, err := fetchSecurityGroup(d.Get("group_id").(string))

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	rule := SecurityRule{
		Action:          "add_firewall_rule",
		Direction:       d.Get("direction").(string),
		PortRangeMin:    rangeMin,
		PortRangeMax:    rangeMax,
		Protocol:        d.Get("protocol").(string),
		NetworkProtocol: d.Get("network_protocol").(string),
		RemoteGroupId:   d.Get("remote_group_id").(string),
		SourceIpRange:   d.Get("remote_ip_subnet").(string),
		SecurityGroupId: securityGroup.EmsRef,
	}

	requestBody, err := json.Marshal(rule)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	existedRules, err := rulesListSnapshot(d.Get("group_id").(string))

	er := make(map[string]int)

	for ndx, r := range existedRules {
		er[r.Id] = ndx
	}

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	requestUrl := fmt.Sprintf("api/compute/v1/security_groups/%s", d.Get("group_id").(string))
	responseBody, err := requestApi("POST", requestUrl, bytes.NewBuffer(requestBody))

	var miqTask MiqTask
	err = responseBody.Decode(&miqTask)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if miqTask.Success != "true" {
		err = fmt.Errorf(miqTask.Message)
		return append(diags, diag.FromErr(err)...)
	}

	securityGroup, err = fetchSecurityGroup(d.Get("group_id").(string))

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	var nr SecurityRule

	for _, r := range securityGroup.SecurityGroupRules {
		_, ok := er[r.Id]

		if !ok {
			nr = r
			break
		}
	}
	err = d.Set("ems_ref", nr.EmsRef)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(nr.Id)

	return nil
}

func resourceSecurityRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	defer ctx.Done()
	var diags diag.Diagnostics

	err := fmt.Errorf("method does not supported")
	return append(diags, diag.FromErr(err)...)
}

func resourceSecurityRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSecurityRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	defer ctx.Done()
	var diags diag.Diagnostics

	r := SecurityRule{
		Action: "remove_firewall_rule",
		Id:     d.Get("ems_ref").(string),
	}

	fmt.Printf("[---DEBUG---] rule %+v", r)

	requestBody, err := json.Marshal(r)
	requestUrl := fmt.Sprintf("api/compute/v1/security_groups/%s", d.Get("group_id").(string))
	responseBody, err := requestApi("POST", requestUrl, bytes.NewBuffer(requestBody))

	var miqTask MiqTaskDelete
	err = responseBody.Decode(&miqTask)

	fmt.Printf("[---DEBUG---] miqTask result %+v", miqTask)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if !miqTask.Success {
		err = fmt.Errorf(miqTask.Message)
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")
	return nil
}
