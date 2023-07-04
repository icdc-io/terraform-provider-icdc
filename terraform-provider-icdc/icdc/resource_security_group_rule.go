package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupRuleCreate,
		Read:   resourceSecurityGroupRuleRead,
		Update: resourceSecurityGroupRuleUpdate,
		Delete: resourceSecurityGroupRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ems_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_ip_range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	// Get all existing security group rules
	responseBody, err := requestApi("GET", fmt.Sprintf("security_groups/%s?expand=resources&attributes=firewall_rules", d.Get("resource_id").(string)), nil)

	if err != nil {
		return err
	}

	var securityGroupRulesCollection *SecurityGroupRulesCollection

	err = responseBody.Decode(&securityGroupRulesCollection)

	if err != nil {
		return err
	}

	existingRulesIlds := make([]string, len(securityGroupRulesCollection.Rules))

	for i, rule := range securityGroupRulesCollection.Rules {
		existingRulesIlds[i] = rule.Id
	}

	if err != nil {
		return err
	}

	var securityGroupRule AddSecurityGroupRule
	securityGroupRule.Action = "add_firewall_rule"
	securityGroupRule.Direction = d.Get("direction").(string)
	securityGroupRule.NetworkProtocol = d.Get("network_protocol").(string)
	securityGroupRule.PortRangeMin = d.Get("port").(string)
	securityGroupRule.PortRangeMax = d.Get("end_port").(string)
	securityGroupRule.Protocol = d.Get("protocol").(string)
	//securityGroupRule.RemoteGroupId = d.Get("remote_group_ip").(string)
	securityGroupRule.SourceIpRange = d.Get("source_ip_range").(string)
	securityGroupRule.SecurityGroupId = d.Get("security_group_id").(string)

	requestBody, err := json.Marshal(securityGroupRule)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(requestBody), &securityGroupRule)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("security_groups/%s", d.Get("resource_id").(string)), body)

	if err != nil {
		return err
	}

	responseBody, err = requestApi("GET", fmt.Sprintf("security_groups/%s?expand=resources&attributes=firewall_rules", d.Get("resource_id").(string)), nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&securityGroupRulesCollection)

	if err != nil {
		return err
	}

	for _, rule := range securityGroupRulesCollection.Rules {
		if !contains(existingRulesIlds, rule.Id) {
			d.SetId(rule.Id)

			err = d.Set("ems_ref", rule.EmsRef)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func resourceSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	responseBody, err := requestApi("GET", fmt.Sprintf("security_groups/%s?expand=resources&attributes=firewall_rules", d.Get("resource_id").(string)), nil)

	if err != nil {
		return err
	}

	var securityGroupRulesCollection *SecurityGroupRulesCollection
	err = responseBody.Decode(&securityGroupRulesCollection)

	if err != nil {
		return err
	}

	for i, rule := range securityGroupRulesCollection.Rules {
		if rule.Id == d.Id() {
			err = d.Set("ems_ref", rule.EmsRef)

			if err != nil {
				return err
			}

			direction := directionMapper(rule.Direction)

			err = d.Set("direction", direction)

			if err != nil {
				return err
			}

			err = d.Set("network_protocol", rule.NetworkProtocol)

			if err != nil {
				return err
			}

			err = d.Set("port", rule.PortRangeMin)

			if err != nil {
				return err
			}

			err = d.Set("end_port", rule.PortRangeMax)

			if err != nil {
				return err
			}

			err = d.Set("protocol", rule.Protocol)

			if err != nil {
				return err
			}

			err = d.Set("source_ip_range", rule.SourceIpRange)

			if err != nil {
				return err
			}

			err = d.Set("security_group_id", rule.SecurityGroupId)

			if err != nil {
				return err
			}

			err = d.Set("resource_id", securityGroupRulesCollection.Rules[i].SecurityGroupId)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func directionMapper(direction string) string {
	if direction == "inbound" {
		return "ingress"
	}

	return "egress"
}

func resourceSecurityGroupRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {

	/*
		{"action":"remove_firewall_rule","id":"3d88adc1-04c0-451f-9bb9-f9596a8e91fc"}
	*/
	deleteSecurityGroupRule := &DeleteRequest{
		Action: "remove_firewall_rule",
		Id:     d.Get("ems_ref").(string),
	}

	requestBody, err := json.Marshal(deleteSecurityGroupRule)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("security_groups/%s", d.Get("resource_id").(string)), body)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
