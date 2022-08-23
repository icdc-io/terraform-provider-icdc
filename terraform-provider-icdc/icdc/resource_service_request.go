package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"os"

	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceRequestResource struct {
	ServiceName         string `json:"service_name"`
	VmMemory 					  string `json:"vm_memory"`
	NumberOfSockets 	  string `json:"number_of_sockets"`
	CoresPerSocket 		  string `json:"cores_per_socket"`
	Hostname 			      string `json:"hostname"`
	Vlan 				        string `json:"vlan"`
//	EnablePublicAccess  string `json:"enable_public_access"`
	SystemDiskType 		  string `json:"system_disk_type"`
	SystemDiskSize 		  string `json:"system_disk_size"`
	AuthType 			      string `json:"auth_type"`
	Adminpassword 		  string `json:"adminpassword"`
	SshKey 			      	string `json:"ssh_key"`
	ServiceTemplateHref string `json:"service_template_href"`
	RegionNumber 		    string `json:"region_number"`
}

type ServiceRequest struct {
	Action 		string 									 `json:"action"`
	Resources []ServiceRequestResource `json:"resources"`
}

type ServiceRequestResponse struct {
	Results []struct {
		Success            bool `json:"success"`
		Message            string `json:"message"`
		ServiceRequestId 	 string `json:"service_request_id"`
		ServiceRequestHref string `json:"service_request_href"`
		Href               string `json:"href"`
	} `json:"results"` 
}

func resourceServiceRequestCreate (d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}

	serviceRequest := &ServiceRequest{
		Action: "add",
		Resources: []ServiceRequestResource{ServiceRequestResource{
			ServiceName: d.Get("service_name").(string),
			VmMemory: d.Get("vm_memory").(string),
			NumberOfSockets: d.Get("number_of_sockets").(string),
			CoresPerSocket: d.Get("cores_per_socket").(string),
			Hostname: d.Get("hostname").(string),
			Vlan: d.Get("vlan").(string),
			SystemDiskType: d.Get("system_disk_type").(string),
			SystemDiskSize: d.Get("system_disk_size").(string),
			AuthType: d.Get("auth_type").(string),
			Adminpassword: d.Get("adminpassword").(string),
			SshKey: d.Get("ssh_key").(string),
			ServiceTemplateHref: d.Get("service_template_href").(string),
			RegionNumber: d.Get("region_number").(string),
			},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)

	
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/service_orders/cart/service_requests/", os.Getenv("ICDC_API_GATEWAY")), body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ICDC_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("ICDC_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return err
	}


	var response *ServiceRequestResponse

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return err
	}

	// WORKAROUND: This is a fake ID. We need to send a second request to fetch a new service ID for real 
	d.SetId(strconv.Itoa(rand.Intn(10000) + 1))

	return nil
}

func resourceServiceRequestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceRequestUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceRequestDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceRequest() *schema.Resource {
	return &schema.Resource{
		Read: resourceServiceRequestRead,
		Create: resourceServiceRequestCreate,
		Update: resourceServiceRequestUpdate,
		Delete: resourceServiceRequestDelete,
		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_memory": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"number_of_sockets": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cores_per_socket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"system_disk_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"system_disk_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"adminpassword": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_template_href": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

		},
	}
}