package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
	"log"
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
			"ems_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"egress": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ems_ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"source_ip_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ingress": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ems_ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"source_ip_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {

	/*

		Request URL: https://api.ycz.icdc.io/api/compute/v1/providers/18000000000003/security_groups
		{"action":"create","name":"ycz_icdc_tf-1"}

		{"results": [
			{"success":true,
			"message":"Creating security group",
			"task_id":"18000000108557",
			"task_href":"https://compute.ycz.icdc.io/api/tasks/18000000108557"
			}
			]
		}

	*/

	var emsProvider *EmsProvider
	responseBody, err := requestApi("GET", "providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager", nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&emsProvider)

	if err != nil {
		return err
	}

	emsProviderId := emsProvider.Resources[0].Id

	securityGroupCreateRequest := SecurityGroupCreateRequest{
		Name:   d.Get("name").(string),
		Action: "create",
	}

	requestBody, err := json.Marshal(securityGroupCreateRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	responseBody, err = requestApi("POST", fmt.Sprintf("providers/%s/security_groups", emsProviderId), body)

	if err != nil {
		return err
	}

	var taskResponse TaskResponse

	err = responseBody.Decode(&taskResponse)

	if err != nil {
		return err
	}

	log.Println(taskResponse)

	if !taskResponse.Results[0].Success {
		return fmt.Errorf("Error creating security group: %s", taskResponse.Results[0].Message)
	}
 
	taskId := taskResponse.Results[0].TaskId

	// Wait for task to complete
	time.Sleep(30 * time.Second)

	taskResultResponse, err := requestApi("GET", fmt.Sprintf("tasks/%s?expand=resources&attributes=task_results", taskId), nil)

	if err != nil {
		return err
	}

	var securityGroupTaskResult SecurityGroupTaskResult

	err = taskResultResponse.Decode(&securityGroupTaskResult)

	log.Println(securityGroupTaskResult)

	if err != nil {
		return err
	}

	securityGroupEmsRef := securityGroupTaskResult.TaskResults.SecurityGroups.EmsRef

	// Wait for completely ems refreshing

	//time.Sleep(45 * time.Second)

	securityGroupCollectionResponse, err := requestApi("GET", fmt.Sprintf("security_groups?expand=resources&filter[]=ems_ref=%s&attributes=firewall_rules", securityGroupEmsRef), nil)

	if err != nil {
		return err
	}

	var securityGroupCollection SecurityGroupCollection

	err = securityGroupCollectionResponse.Decode(&securityGroupCollection)

	log.Println(securityGroupCollection)

	if err != nil {
		return err
	}

	//err = d.Set("Name", securityGroupCollection.Resources[0].Name)

	if err != nil {
		return err
	}

	d.SetId(securityGroupCollection.Resources[0].Id)

	return nil

}

func resourceSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {

	var emsProvider *EmsProvider
	responseBody, err := requestApi("GET", "providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager", nil)

	if err != nil {
		return err
	}

	err = responseBody.Decode(&emsProvider)

	if err != nil {
		return err
	}

	emsProviderId := emsProvider.Resources[0].Id

	securityGroupDeleteRequest := &SecurityGroupDeleteRequest{
		Action: "delete",
		Id:     d.Id(),
		Name:   d.Get("name").(string),
	}

	requestBody, err := json.Marshal(securityGroupDeleteRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("providers/%s/security_groups", emsProviderId), body)

	if err != nil {
		return err
	}

	return nil
}
