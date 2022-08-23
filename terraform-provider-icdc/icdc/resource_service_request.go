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

const baseURL = "https://api.ycz.icdc.io/api/compute/v1"
const miqGroup = "icdc.member"

type ServiceRequestResource struct {
	Service_name         string `json:"service_name"`
	Vm_memory 					  string `json:"vm_memory"`
	Number_of_sockets 	  string `json:"number_of_sockets"`
	Cores_per_socket 		  string `json:"cores_per_socket"`
	Hostname 			      string `json:"hostname"`
	Vlan 				        string `json:"vlan"`
	Enable_public_access  string `json:"enable_public_access"`
	System_disk_type 		  string `json:"system_disk_type"`
	System_disk_size 		  string `json:"system_disk_size"`
	Auth_type 			      string `json:"auth_type"`
	Adminpassword 		  string `json:"adminpassword"`
	Ssh_key_list 			    string `json:"ssh_key_list"`
	Ssh_key 			      string `json:"ssh_key"`
	Service_template_href string `json:"service_template_href"`
	Region_number 		    string `json:"region_number"`
}

type ServiceRequest struct {
	Action string `json:"action"`
	Resources []ServiceRequestResource `json:"resources"`
}

type ServiceRequestResponse struct {
	Results []struct {
		Success bool `json:"success"`
		Message string `json:"message"`
		Service_request_id string `json:"service_request_id"`
		Service_request_href string `json:"service_request_href"`
		Href string `json:"href"`
	} `json:"results"` 
}


func check(e error) {
	if e != nil {
			panic(e)
	}
}

func resourceServiceRequestCreate (d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}

	serviceRequest := &ServiceRequest{
		Action: "add",
		Resources: []ServiceRequestResource{ServiceRequestResource{
			Service_name: d.Get("service_name").(string),
			Vm_memory: d.Get("vm_memory").(string),
			Number_of_sockets: d.Get("number_of_sockets").(string),
			Cores_per_socket: d.Get("cores_per_socket").(string),
			Hostname: d.Get("hostname").(string),
			Vlan: d.Get("vlan").(string),
			Enable_public_access: d.Get("enable_public_access").(string),
			System_disk_type: d.Get("system_disk_type").(string),
			System_disk_size: d.Get("system_disk_size").(string),
			Auth_type: d.Get("auth_type").(string),
			Adminpassword: d.Get("adminpassword").(string),
			Ssh_key_list: d.Get("ssh_key_list").(string),
			Ssh_key: d.Get("ssh_key").(string),
			Service_template_href: d.Get("service_template_href").(string),
			Region_number: d.Get("region_number").(string),
			},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)

	
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/service_orders/cart/service_requests/", baseURL), body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")))
	req.Header.Set("X_MIQ_GROUP", miqGroup)

	r, err := client.Do(req)
	if err != nil {
		return err
	}


	var response ServiceRequestResponse

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
			"enable_public_access": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"ssh_key_list": &schema.Schema{
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