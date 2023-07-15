package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSecurityGroupRead,
		Create: resourceSecurityGroupCreate,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,
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
			"direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_range_max": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_range_min": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s", d.Get("id"))
	r, err := requestApi("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error getting api services: %w", err)
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRequestResponse *SecurityGroupRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	log.Println(PrettyStruct(securityGroupRequestResponse))
	d.SetId(securityGroupRequestResponse.Id)
	return nil
}

func resourceSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	cloudGroupRaw := &GroupCreateBody{
		SecurityGroup: SecurityGroupBody{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}

	requestBody, err := json.Marshal(cloudGroupRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudGroupRaw))

	url := fmt.Sprintf("%s/security_groups", d.Get("vpc_id"))
	r, err := requestApi("POST", url, body)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRequestResponse *SecurityGroupRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	log.Println(PrettyStruct(securityGroupRequestResponse))

	group_id := securityGroupRequestResponse.Id
	log.Println(PrettyStruct(group_id))
	d.SetId(group_id)

	/*create rule*/
	/*	cloudGroupRuleRaw := &RuleCreateBody{
			SecurityGroupRule: SecurityGroupRuleBody{
				Direction:       d.Get("direction").(string),
				NetworkProtocol: d.Get("network_protocol").(string),
				Ethertype:       d.Get("ethertype").(string),
				PortRangeMin:    d.Get("port_range_min").(string),
				PortRangeMax:    d.Get("port_range_max").(string),
				Protocol:        d.Get("protocol").(string),
				RemoteGroupId:   d.Get("remote_group_id").(string),
				SecurityGroupId: d.Get("security_group_id").(string),
			},
		}

		requestBody, err = json.Marshal(cloudGroupRuleRaw)
		if err != nil {
			return fmt.Errorf("error marshaling service request: %w", err)
		}

		body = bytes.NewBuffer(requestBody)

		log.Println(PrettyStruct(cloudGroupRuleRaw))
		url = fmt.Sprintf("security_groups/%s/rules", d.Get("id"))
		r, err = requestApi("POST", url, body)

		if err != nil {
			return err
		}
		resBody, err = ioutil.ReadAll(r.Body)
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

		sgroup_id := securityGroupRuleRequestResponse.Id
		log.Println(PrettyStruct(group_id))
		d.SetId(group_id)*/
	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {

	cloudGroupRaw := &GroupCreateBody{
		SecurityGroup: SecurityGroupBody{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}

	requestBody, err := json.Marshal(cloudGroupRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudGroupRaw))

	url := fmt.Sprintf("%s", d.Get("id").(string))
	r, err := requestApi("PUT", url, body)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var securityGroupRequestResponse *SecurityGroupRequestResponse

	if err = json.Unmarshal(resBody, &securityGroupRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(securityGroupRequestResponse))
	vpc_id := securityGroupRequestResponse.Id
	d.SetId(vpc_id)

	/*update rule*/
	/*	cloudGroupRuleRaw := &RuleCreateBody{
			SecurityGroupRule: SecurityGroupRuleBody{
				Direction:       d.Get("direction").(string),
				NetworkProtocol: d.Get("network_protocol").(string),
				Ethertype:       d.Get("ethertype").(string),
				PortRangeMin:    d.Get("port_range_min").(string),
				PortRangeMax:    d.Get("port_range_max").(string),
				Protocol:        d.Get("protocol").(string),
				RemoteGroupId:   d.Get("remote_group_id").(string),
				SecurityGroupId: d.Get("security_group_id").(string),
			},
		}

		requestBody, err = json.Marshal(cloudGroupRuleRaw)
		if err != nil {
			return fmt.Errorf("error marshaling service request: %w", err)
		}

		body = bytes.NewBuffer(requestBody)

		log.Println(PrettyStruct(cloudGroupRuleRaw))
		url = fmt.Sprintf("security_groups/%s/rules/", d.Get("id"))
		_, err = requestApi("POST", url, body)

		if err != nil {
			return err
		}*/

	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {

	url := fmt.Sprintf("%s", d.Get("id").(string))
	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil

}
