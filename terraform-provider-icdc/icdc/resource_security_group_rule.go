package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSecurityGroupRuleRead,
		Create: resourceSecurityGroupRuleCreate,
		Update: resourceSecurityGroupRuleUpdate,
		Delete: resourceSecurityGroupRuleDelete,
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

func resourceSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudGroupRuleRaw))
	url := fmt.Sprintf("security_groups/%s/rules", d.Get("security_group_id").(string))
	r, err := requestApi("POST", url, body)

	if err != nil {
		return err
	}
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRuleRequestResponse *SecurityGroupRuleRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRuleRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(securityGroupRuleRequestResponse))
	log.Println(PrettyStruct(securityGroupRuleRequestResponse))

	sgroup_rule_id := securityGroupRuleRequestResponse.SecurityGroupRule.Id
	log.Println(PrettyStruct(sgroup_rule_id))
	d.SetId(sgroup_rule_id)

	return nil
}

func resourceSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("rules/%s", d.Get("id").(string))
	r, err := requestApi("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error getting api services: %w", err)
	}

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRuleRequestResponse *SecurityGroupRuleRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRuleRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	log.Println(PrettyStruct(securityGroupRuleRequestResponse))
	log.Println(securityGroupRuleRequestResponse.SecurityGroupRule.Id)
	d.SetId(securityGroupRuleRequestResponse.SecurityGroupRule.Id)

	return nil
}

func resourceSecurityGroupRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {

	url := fmt.Sprintf("security_groups/%s/rules/%s", d.Get("security_group_id").(string), d.Get("id").(string))
	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
