package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
  ahrechushkin:
		- Need to find a way to pass metada in provider context (now i passed all required info in os environment vars)
		- Need to implement error handling
		- Need to implement logger
*/

func resourceService() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServiceRead,
		Create: resourceServiceCreate,
		Update: resourceServiceUpdate,
		Delete: resourceServiceDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vms": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
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
						"cpu_cores": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"storage_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"storage_mb": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ipaddresses": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func resourceServiceCreate(d *schema.ResourceData, m interface{}) error {
	vlan := fmt.Sprintf("%s (%s)", d.Get("vms.0.subnet").(string), d.Get("vms.0.subnet").(string))

	service := Service{
		Name:              d.Get("name").(string),
		SshKey:            d.Get("ssh_key").(string),
		ServiceTemplateId: d.Get("service_template_id").(string),
		Vms: []VmParams{VmParams{
			MemoryMb:    d.Get("vms.0.memory_mb").(string),
			CpuCores:    d.Get("vms.0.cpu_cores").(string),
			StorageType: d.Get("vms.0.storage_type").(string),
			StorageMb:   d.Get("vms.0.storage_mb").(string),
			Network:     vlan,
		},
		},
	}

	serviceRequest := &ServiceRequest{
		Action: "add",
		Resources: []ServiceResources{ServiceResources{
			ServiceName:         service.Name,
			VmMemory:            service.Vms[0].MemoryMb,
			NumberOfSockets:     "1",
			CoresPerSocket:      service.Vms[0].CpuCores,
			Hostname:            "generated-hostname",
			Vlan:                service.Vms[0].Network,
			SystemDiskType:      service.Vms[0].StorageType,
			SystemDiskSize:      service.Vms[0].StorageMb,
			AuthType:            "key",
			SshKey:              service.SshKey,
			ServiceTemplateHref: fmt.Sprintf("/api/service_templates/%s", service.ServiceTemplateId),
			RegionNumber:        "18",
		},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)
	responseBody, err := requestApi("POST", "service_orders/cart/service_requests/", body)

	if err != nil {
		return err
	}

	var serviceRequestResponse *ServiceRequestResponse

	err = responseBody.Decode(&serviceRequestResponse)

	if err != nil {
		return err
	}

	/*
		ahrechushkin: We need to wait for the service request to be completed.
			To know service id we will make requests to /api/services with filter service_request_id int the loop.
			And setup ID only after creating service in Compute DB.
			Monkey patching is not the best way to do this, but anyway it works.
	*/

	serviceRequestId := serviceRequestResponse.Results[0].ServiceRequestId

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

func fetchDestinationId(serviceRequestId string, destinationType string) (string, error) {

	responseBody, err := requestApi("GET", fmt.Sprintf("service_requests/%s?expand=resources&attributes=miq_request_tasks", serviceRequestId), nil)

	if err != nil {
		return "", err
	}

	var response *ServiceMiqRequest

	err = responseBody.Decode(&response)
	if err != nil {
		return "", err
	}

	for i := range response.MiqRequestTasks {
		if response.MiqRequestTasks[i].DestinationType == destinationType {
			return response.MiqRequestTasks[i].DestinationId, nil
		}
	}

	return "", nil
}

func resourceServiceRead(d *schema.ResourceData, m interface{}) error {
	responseBody, err := requestApi("GET", fmt.Sprintf("services/%s?expand=resources&attributes=vms", d.Id()), nil)

	if err != nil {
		return err
	}

	var service *Service

	err = responseBody.Decode(&service)
	if err != nil {
		return err
	}

	vms := flattenVms(service.Vms)
	if err := d.Set("vms", vms); err != nil {
		return err
	}

	d.SetId(d.Id())
	return nil
}

func flattenVms(vmsList []VmParams) []interface{} {
	if vmsList != nil {
		vms := make([]interface{}, len(vmsList))

		for i, vm := range vmsList {

			var remoteVm Vm
			responseBody, err := requestApi("GET", fmt.Sprintf("vms/%s?expand=resources&attributes=hardware,disks,lans,ipaddresses", vm.ID), nil)

			if err != nil {
				return nil
			}

			err = responseBody.Decode(&remoteVm)

			if err != nil {
				return nil
			}

			vml := make(map[string]interface{})
			vml["id"] = remoteVm.Id
			vml["name"] = remoteVm.Name
			vml["memory_mb"] = strconv.Itoa(remoteVm.Hardware.MemoryMb)
			vml["cpu_cores"] = strconv.Itoa(remoteVm.Hardware.CpuCores)
			vml["subnet"] = remoteVm.Network[0].Name
			vml["storage_type"] = "nvme"
			vml["storage_mb"] = strconv.Itoa(remoteVm.Disks[0].Size / (1 << 30))
			vml["ipaddresses"] = remoteVm.Ipaddresses

			vms[i] = vml
		}

		return vms
	}

	return make([]interface{}, 0)
}

func resourceServiceUpdate(d *schema.ResourceData, m interface{}) error {
	/*
		ahrechushkin: Unfourtunately we can't update service resources.
			We may update only vm resource, but we don't have VM abstraction layer.
			Service -> [VMs -> [Resources -> [VmMemory, NumberOfCpus, NumberOfSockets, CoresPerSocket]]]
			Must be implemented in future.

			service name upated by /api/services/{id} endpoint
			all resources updated by /api/vms/{id} endpoint
	*/

	if d.HasChange("name") {
		// TODO: implement service update method
		/*
			POST services/18000000000388
			BODY {
						"action":"edit",
						"resource": {
							"id":"18000000000388",
							"name":"tf-composite#5"
						}
					}
		*/
		var reconfigureRequest ServiceReconfigureRequest
		reconfigureRequest.Action = "edit"
		reconfigureRequest.Resource.ID = d.Id()
		reconfigureRequest.Resource.Name = d.Get("name").(string)

		requestBody, err := json.Marshal(reconfigureRequest)

		if err != nil {
			return err
		}

		body := bytes.NewBuffer(requestBody)

		_, err = requestApi("POST", fmt.Sprintf("services/%s", d.Id()), body)

		if err != nil {
			return err
		}
	}

	if d.HasChange("vms") {
		/*
				ahrechushkin"
			 		CpuCores, MemoryMb updating by create reconfigure request
					Network can update only running automation task
					they are two different requests to update vm
		*/

		if d.HasChange("vms.0.cpu_cores") || d.HasChange("vms.0.memory_mb") {
			var vmReconfigureRequest VmReconfigureRequest
			vmReconfigureRequest.Action = "reconfigure"
			vmReconfigureRequest.Resource.RequestType = "vm_reconfigure"
			vmReconfigureRequest.Resource.VmMemory = d.Get("vms.0.memory_mb").(string)
			vmReconfigureRequest.Resource.NumberOfCpus = d.Get("vms.0.cpu_cores").(string)
			vmReconfigureRequest.Resource.NumberOfSockets = "1"
			vmReconfigureRequest.Resource.CoresPerSocket = d.Get("vms.0.cpu_cores").(string)

			requestBody, err := json.Marshal(vmReconfigureRequest)

			if err != nil {
				return err
			}

			body := bytes.NewBuffer(requestBody)

			vmId := d.Get("vms.0.id").(string)
			_, err = requestApi("POST", fmt.Sprintf("vms/%s", vmId), body)

			if err != nil {
				return err
			}
		}

		if d.HasChange("vms.0.subnet") {
			/*
				POST: api/services/18000000000388/
				BODY: {
					"action":"invoke_custom_button",
					"resource":{
						"task":"call_automation",
						"path":"System/Request/ChangeNetworkType",
						"params":{
							"dialog_subnet_profile":"ycz_icdc_base",
							"new_subnet_name":"Base" !!! Using only for UI.
						}
					}
				}
			*/

			var customButtonRequest ChangeNetworkTypeRequest
			customButtonRequest.Action = "invoke_custom_button"
			customButtonRequest.Resource.Task = "call_automation"
			customButtonRequest.Resource.Path = "System/Request/ChangeNetworkType"
			customButtonRequest.Resource.Params.DialogNetworkProfile = d.Get("vms.0.subnet").(string)

			requestBody, err := json.Marshal(customButtonRequest)

			if err != nil {
				return err
			}

			body := bytes.NewBuffer(requestBody)
			_, err = requestApi("POST", fmt.Sprintf("services/%s", d.Id()), body)

			if err != nil {
				return err
			}

		}
	}
	return nil
}

func resourceServiceDelete(d *schema.ResourceData, m interface{}) error {
	serviceRequest := &ServiceRequest{
		Action: "request_retire",
	}

	requestBody, err := json.Marshal(serviceRequest)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("services/%s", d.Id()), body)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
