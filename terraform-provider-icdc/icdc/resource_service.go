package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	//"strconv"
	"time"
	"os"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
  ahrechushkin: 
		- Need to move all http requests to separate function to make code prettier.
	  requestComputeApi(method, endpoint, body)
		- Need to prepare generic type for non-root (service, vm, etc.) Compute objects, i mean RequestResponse, Requests...
		- Need to find a way to pass metada in provider context (now i passed all required info in os environment vars)
		- Need to implement error handling
		- Need to implement logger
*/

type Service struct {
	ID  string `json:"id"`
	Name string `json:"name"`
	SshKey string
	ServiceTemplateId string `json:"service_template_id"`
	Vms [] VmParams `json:"vms"`
			/* ahrechushkin: for sure we can fetch full information about vm we need to make 2 requests.
				1. api/services/:ID?expand=resources&attributes=vms
				2. api/vms/:ID?expand=resources&attributes=hardware
				Maybe make sense a generate object with aggregated information from two endpoints.
		*/
}

type VmParams struct {
	ID string `json:"id"`
	Name string `json:"name"`
	MemoryMb string
	CpuCores string
	StorageType string
	StorageMb string
	Network string
}

func resourceService() *schema.Resource {
	return &schema.Resource{
		Read: resourceServiceRead,
		Create: resourceServiceCreate,
		Update: resourceServiceUpdate,
		Delete: resourceServiceDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vms": &schema.Schema{
				Type:		 schema.TypeList,
				Elem: 	 &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type: schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_mb": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"cpu": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"storage_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"storage_gb": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"network": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"ssh_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type ServiceResources struct {
	ServiceName         string `json:"service_name"`
	VmMemory 					  string `json:"vm_memory"`
	NumberOfSockets 	  string `json:"number_of_sockets"`
	CoresPerSocket 		  string `json:"cores_per_socket"`
	Hostname 			      string `json:"hostname"`
	Vlan 				        string `json:"vlan"`
	SystemDiskType 		  string `json:"system_disk_type"`
	SystemDiskSize 		  string `json:"system_disk_size"`
	AuthType 			      string `json:"auth_type"`
	Adminpassword 		  string `json:"adminpassword"`
	SshKey 			      	string `json:"ssh_key"`
	ServiceTemplateHref string `json:"service_template_href"`
	RegionNumber 		    string `json:"region_number"`
}

type ServiceRequest struct {
	Action 		string 						 `json:"action"`
	Resources []ServiceResources `json:"resources"`
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



type ServiceMiqRequest struct {
	MiqRequestTasks []struct {
		DestinationId string `json:"destination_id"`
		DestinationType string `json:"destination_type"`
	} `json:"miq_request_tasks"`
}



func resourceServiceCreate (d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}

	vlan := fmt.Sprintf("%s (%s)", d.Get("vms.0.network").(string), d.Get("vms.0.network").(string))

	service := Service{
		Name: d.Get("name").(string),
		SshKey: d.Get("ssh_key").(string),
		ServiceTemplateId: d.Get("service_template_id").(string),
		Vms: []VmParams {VmParams{
			MemoryMb: d.Get("vms.0.memory_mb").(string),
			CpuCores: d.Get("vms.0.cpu").(string),
			StorageType: d.Get("vms.0.storage_type").(string),
			StorageMb: d.Get("vms.0.storage_gb").(string),
			Network: vlan,
			},
		},
	}

	serviceRequest := &ServiceRequest{
		Action: "add",
		Resources: []ServiceResources{ServiceResources{
			ServiceName: service.Name,
			VmMemory: service.Vms[0].MemoryMb,
			NumberOfSockets: "1",
			CoresPerSocket: service.Vms[0].CpuCores,
			Hostname: "generated-hostname",
			Vlan: service.Vms[0].Network,
			SystemDiskType: service.Vms[0].StorageType,
			SystemDiskSize: service.Vms[0].StorageMb,
			AuthType: "key",
			SshKey: service.SshKey,
			ServiceTemplateHref: fmt.Sprintf("/api/service_templates/%s", service.ServiceTemplateId),
			RegionNumber: "18",
			},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/service_orders/cart/service_requests/", os.Getenv("API_GATEWAY")), body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return err
	}


	var response *ServiceRequestResponse

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return err
	}

	/*
	ahrechushkin: We need to wait for the service request to be completed.
		To know service id we will make requests to /api/services with filter service_request_id int the loop.
		And setup ID only after creating service in Compute DB.		
		Monkey patching is not the best way to do this, but anyway it works.
	*/

	serviceRequestId := response.Results[0].ServiceRequestId
	
	for {
		serviceId, err := fetchDestinationId(serviceRequestId, "Service")

		if err != nil {
			return err
		}

		if serviceId != "" {
			d.SetId(serviceId)
			break
		}

		time.Sleep(10 * time.Second)
	}


	return nil
}

func fetchDestinationId (serviceRequestId string, destinationType string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/service_requests/%s?expand=resources&attributes=miq_request_tasks", os.Getenv("API_GATEWAY"), serviceRequestId), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var response *ServiceMiqRequest

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	file, _ := json.MarshalIndent(response, "", "  ")
	_ = ioutil.WriteFile("/tmp/miq_request_task_response.json", file, 0644)

	for i := range response.MiqRequestTasks {
		if response.MiqRequestTasks[i].DestinationType == destinationType {
			return response.MiqRequestTasks[i].DestinationId, nil
		}
	}

	return "", nil
}

func resourceServiceRead(d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/services/%s?expand=resources&attributes=vms", os.Getenv("API_GATEWAY"), d.Id()), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	var service *Service

	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		return err
	}

	file, _ := json.MarshalIndent(service, "", "  ")
	file_name := fmt.Sprintf("/tmp/service_%s.json", d.Id())
	_ = ioutil.WriteFile(file_name, file, 0644)

	d.Set("vms", service.Vms)
	d.SetId(d.Id())
	return nil
}

/*
type VmReconfigureRequest struct {
	Action string `json:"action"`
	Resource struct {
		RequestType string `json:"request_type"`
		VmMemory int `json:"vm_memory"`
		NumberOfCpus int `json:"number_of_cpus"`
		NumberOfSockets int `json:"number_of_sockets"`
		CoresPerSocket int `json:"cores_per_socket"` 
	} `json:"resource"`
}
*/

func resourceServiceUpdate(d *schema.ResourceData, m interface{}) error {
	/*
	ahrechushkin: Unfourtunately we can't update service resources.
		We may update only vm resource, but we don't have VM abstraction layer.
		Service -> [VMs -> [Resources -> [VmMemory, NumberOfCpus, NumberOfSockets, CoresPerSocket]]]
		Must be implemented in future.
	*/
	return nil
}

func resourceServiceDelete(d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}

	serviceRequest := &ServiceRequest{
		Action: "request_retire",
	}

	requestBody, err := json.Marshal(serviceRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/services/%s", os.Getenv("API_GATEWAY"), d.Id()), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", os.Getenv("AUTH_GROUP"))

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	var service *Service
	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
