package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/r3labs/diff/v3"
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
						"system_disk_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"system_disk_size": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"additional_disk": {
							// artemsafonau: i think it will be better to make TypeSet like in AWS
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"additional_disk_type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"additional_disk_size": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"filename": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
			// ToDo: make possibility to choose between ssh and generated-password
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
	
	// ToDo: make update? (add other additional disks)
	additional_disk := "f"
	if (d.Get("vms.0.additional_disk.#") != "0") {
		// ToDo: make different APIs endpoints functions
		var tags *TagsResponse
		responseBody, err := requestApi("GET", "tags?expand=resources&attributes=classification&filter[]=name='/managed/storage_type/*'", nil)
		if err != nil {
			return fmt.Errorf("error getting api tags: %w", err)
		}
		err = responseBody.Decode(&tags)
		if err != nil {
			return fmt.Errorf("error decoding tags: %w", err)
		}

		tfDiskType := d.Get("vms.0.additional_disk.0.additional_disk_type").(string)
		// ToDo: additional_disk_size can't be <= 0
		if containsTag(tags, tfDiskType) && (d.Get("vms.0.additional_disk.0.additional_disk_size") != "") {
			additional_disk = "t"
		} else {
			// return error
			// d.Timeout(schema.TimeoutCreate))
			// o, n := d.GetChange("tags_all")
			// return fmt.Errorf("error creating Backup Plan: %w", err)
			// err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			// return diag.Errorf("error waiting for EC2 Network Insights Analysis (%s) create: %s", d.Id(), err)
			return fmt.Errorf("error: unsupported additional disk type")
		}
	}
	
	log.Println(additional_disk)
	//panic("1111111111")

	service := Service{
		Name:              d.Get("name").(string),
		SshKey:            d.Get("ssh_key").(string),
		ServiceTemplateId: d.Get("service_template_id").(string),
		Vms: []VmParams{VmParams{
			MemoryMb:        		d.Get("vms.0.memory_mb").(string),
			CpuCores:    		 		d.Get("vms.0.cpu_cores").(string),
			SystemDiskType:  		d.Get("vms.0.system_disk_type").(string),
			SystemDiskSize:  		d.Get("vms.0.system_disk_size").(string),
			AdditionalDisk:     additional_disk,
			AdditionalDiskType: d.Get("vms.0.additional_disk.0.additional_disk_type").(string),
			AdditionalDiskSize: d.Get("vms.0.additional_disk.0.additional_disk_size").(string),
			Network:     				vlan,
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
			SystemDiskType:  		 service.Vms[0].SystemDiskType,
			SystemDiskSize:  		 service.Vms[0].SystemDiskSize,
			AdditionalDisk:  		 service.Vms[0].AdditionalDisk,
			AdditionalDiskType:  service.Vms[0].AdditionalDiskType,
			AdditionalDiskSize:  service.Vms[0].AdditionalDiskSize,
			AuthType:            "key", // ToDo: update for generate-password
			SshKey:              service.SshKey,
			ServiceTemplateHref: fmt.Sprintf("/api/service_templates/%s", service.ServiceTemplateId),
			RegionNumber:        "18",
		},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)
	// prettystruct
	log.Println(PrettyStruct(serviceRequest))
	responseBody, err := requestApi("POST", "service_orders/cart/service_requests/", body)
	if err != nil {
		return fmt.Errorf("error requesting service: %w", err)
	}

	var serviceRequestResponse *ServiceRequestResponse
	if err = responseBody.Decode(&serviceRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	log.Println(PrettyStruct(serviceRequestResponse))

	/*
		ahrechushkin: We need to wait for the service request to be completed.
			To know service id we will make requests to /api/services with filter service_request_id int the loop.
			And setup ID only after creating service in Compute DB.
			Monkey patching is not the best way to do this, but anyway it works.
	*/

	serviceRequestId := serviceRequestResponse.Results[0].ServiceRequestId
	var serviceId string
	
	// ToDo: read about timeouts
	// infinite loop in case of error? set time (2 minutes?)
	currentTime := time.Now()
	requiredTime := currentTime.Add(time.Minute * 2)
	for {
		serviceId, err = fetchDestinationId(serviceRequestId, "Service")
		if err != nil {
			return err
		}

		if serviceId != "" {
			d.SetId(serviceId)
			break
		}

		if currentTime.After(requiredTime) {
			return fmt.Errorf("error: service creation time is out")
		}

		currentTime = currentTime.Add(10 * time.Second)
		time.Sleep(10 * time.Second)
	}

	log.Println("Service Created")
	// make loop for checking ~ 10 min for vm full creation
	// read terraform aws documentation

	// checking for vm creating (20 min) - too long
	currentTime = time.Now()
	requiredTime = currentTime.Add(time.Minute * 20)
	for {
		vmsId, err := fetchDestinationVm(serviceId)
		if err != nil {
			return fmt.Errorf("error vm fetch destination: %w", err)
		}

		if vmsId != "" {
			d.Set("vms.0.id", vmsId)
			break
		}

		if currentTime.After(requiredTime) {
			return fmt.Errorf("error: vm creation time out")
		}

		currentTime = currentTime.Add(30 * time.Second)
		time.Sleep(30 * time.Second)
	}
	log.Println("Service Vm Created")

	// ToDo: make read
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

func fetchDestinationVm(serviceRequestId string) (string, error) {

	responseBody, err := requestApi("GET", fmt.Sprintf("services/%s?expand=vms", serviceRequestId), nil)
	if err != nil {
		return "", err
	}

	var response *ServiceVmProvisonResponse
	err = responseBody.Decode(&response)
	if err != nil {
		return "", err
	}

	if response.LifecycleState == "provisioned" {
		return response.Vms[0].Id, nil
	}

	return "", nil
}

func resourceServiceRead(d *schema.ResourceData, m interface{}) error {
	responseBody, err := requestApi("GET", fmt.Sprintf("services/%s?expand=resources&attributes=vms", d.Id()), nil)
	if err != nil {
		return fmt.Errorf("error getting api services: %w", err)
	}

	// can also add dialog_ssh_key
	var service *Service
	err = responseBody.Decode(&service)
	if err != nil {
		return fmt.Errorf("error decoding service response body: %w", err)
	}

	vms := flattenVms(service.Vms, d)
	if err := d.Set("vms", vms); err != nil {
		return fmt.Errorf("error setting vms: %w", err)
	}

	log.Println(PrettyStruct(vms))
	
	d.SetId(d.Id())
	return nil
}

func flattenVms(vmsList []VmParams, d *schema.ResourceData) []interface{} {
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

			// i hope that system disk will be first in all disks
			sort.SliceStable(remoteVm.Disks, func(i, j int) bool {
				return remoteVm.Disks[i].Id < remoteVm.Disks[j].Id
			})

			log.Println(PrettyStruct(remoteVm.Disks))

			vml := make(map[string]interface{})
			vml["id"] = remoteVm.Id
			vml["name"] = remoteVm.Name
			vml["memory_mb"] = strconv.Itoa(remoteVm.Hardware.MemoryMb)
			vml["cpu_cores"] = strconv.Itoa(remoteVm.Hardware.CpuCores)
			vml["subnet"] = remoteVm.Network[0].Name
			vml["system_disk_type"] = "nvme"
			vml["system_disk_size"] = strconv.Itoa(remoteVm.Disks[0].Size / (1 << 30))

			// maybe it will be better to use TypeSet
			if len(remoteVm.Disks) > 1 {
				remoteVm.Disks = remoteVm.Disks[1:]
				disks := make([]map[string]interface{}, 0)
				tf_state_disks := d.Get("vms.0.additional_disk").([]interface{})

				log.Println(PrettyStruct(tf_state_disks))

				for index1 := range tf_state_disks {
					disk1 := tf_state_disks[index1].(map[string]interface{})
					for index2, disk2 := range remoteVm.Disks {
						// make for type
						// convert size to gb and string
						disk2.Size = disk2.Size / (1 << 30)
						strDisk2Size := strconv.Itoa(disk2.Size)
						if strDisk2Size == disk1["additional_disk_size"] {
							included_disk := make(map[string]interface{})
							included_disk["id"] = disk2.Id
							included_disk["additional_disk_type"] = "nvme" // disk.type?
							included_disk["additional_disk_size"] = strconv.Itoa(disk2.Size)
							included_disk["filename"] = disk2.Filename
							disks = append(disks, included_disk)
							// remove element from remoteVmDisks
							remoteVm.Disks = append(remoteVm.Disks[:index2], remoteVm.Disks[index2 + 1:]...)
							break
						}
					}
				}

				// append all other disks
				for _, disk := range remoteVm.Disks {
					included_disk := make(map[string]interface{})
					included_disk["id"] = disk.Id
					included_disk["additional_disk_type"] = "nvme" // disk.type?
					included_disk["additional_disk_size"] = strconv.Itoa(disk.Size / (1 << 30))
					included_disk["filename"] = disk.Filename
					disks = append(disks, included_disk)
				}
				log.Println(PrettyStruct(disks))
				vml["additional_disk"] = disks
			}
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
			return fmt.Errorf("error marshaling vm change name request: %w", err)
		}

		body := bytes.NewBuffer(requestBody)

		_, err = requestApi("POST", fmt.Sprintf("services/%s", d.Id()), body)

		if err != nil {
			return fmt.Errorf("error sending vm change name request: %w", err)
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
				return fmt.Errorf("error marhsaling vm reconfigure request: %w", err)
			}

			body := bytes.NewBuffer(requestBody)

			vmId := d.Get("vms.0.id").(string)
			value, err := requestApi("POST", fmt.Sprintf("vms/%s", vmId), body)
			if err != nil {
				return fmt.Errorf("error sending vms request: %w", err)
			}

			var responseBody ReconfigurationResponse
			if err = value.Decode(&responseBody); err != nil {
				return fmt.Errorf("error decoding vms response body: %w", err)
			}

			log.Println(PrettyStruct(responseBody))
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
				return fmt.Errorf("error marhsaling vm subnet request: %w", err)
			}

			body := bytes.NewBuffer(requestBody)
			value, err := requestApi("POST", fmt.Sprintf("services/%s", d.Id()), body)
			if err != nil {
				return fmt.Errorf("error requesting subnet change: %w", err)
			}

			var responseBody ReconfigurationResponse
			if err := value.Decode(&responseBody); err != nil {
				return fmt.Errorf("error decoding subnet change response: %w", err)
			}

			log.Println(PrettyStruct(responseBody))
		}

		if d.HasChange("vms.0.additional_disk") { 
			// ToDo: think about union of vm resources and disks configuration
		 	/*
			  	artemsafonau" IT MUST BE REFACTORED
					need to use typeset because of order of disks?
					update disks is unstable because of kafka queue
					maybe need to make TypeSet but it will be in FUTURE
			  	change storage type logic in future version
					is it needed to wait for changes applyed?
		 	*/

			// aws check backups for notification

			// set additional_disk_request
			var additionalDiskRequest AdditionalDiskRequest
			additionalDiskRequest.Action = "reconfigure"
			additionalDiskRequest.Resource.CoresPerSocket = d.Get("vms.0.cpu_cores").(string)
			additionalDiskRequest.Resource.NumberOfCpus = d.Get("vms.0.cpu_cores").(string)
			additionalDiskRequest.Resource.NumberOfSockets = "1"
			additionalDiskRequest.Resource.RequestType = "vm_reconfigure"
			additionalDiskRequest.Resource.VmMemory = d.Get("vms.0.memory_mb").(string)

			o, n := d.GetChange("vms.0.additional_disk")
			os := o.([]interface{})
			ns := n.([]interface{})
			changelog, _ := diff.Diff(o, n)

			// make pretty log
			log.Println(PrettyStruct(changelog))

			// make request for tags
			response, err := requestApi("GET", "tags?expand=resources&attributes=classification&filter[]=name='/managed/storage_type/*'", nil)
			if err != nil {
				return fmt.Errorf("error requesting storage types: %w", err)
			}
	
			var tags *TagsResponse
			err = response.Decode(&tags)
			if err != nil {
				return fmt.Errorf("error decoding tags response: %w", err)
			}

			// add vars for getting unique path for update
			var existing_paths = make(map[string]bool)
			var paths = []string{}

			for _, value := range changelog {

				// if create -> create
				// if update -> collect unique paths -> destroy and create
				// if delete -> destroy

				index, err := strconv.Atoi(value.Path[0])
				if err != nil {
					return fmt.Errorf("error converting from string to int: %w", err)
				}
				
				switch value.Type {
					// ToDo: divide into subfunctions
				case "create":
					new := ns[index].(map[string]interface{})

					diskType, ok := new["additional_disk_type"].(string)
					if !ok {
						return fmt.Errorf("can not read additional disk type")
					}
					if !containsTag(tags, diskType) {
						return fmt.Errorf("disk type is not available")
					}

					// ToDo: check for simplest types convertion
					strDiskSize, ok := new["additional_disk_size"].(string)
					if !ok {
						return fmt.Errorf("can not read additional disk size")
					}
					intDiskSize, err := strconv.Atoi(strDiskSize)
					if err != nil {
						return fmt.Errorf("error converting from string to int: %w", err)
					}

					diskAdd := DiskAdd{
						StorageType: diskType,
						Name: "",
						Type: fmt.Sprintf("/managed/storage_type/%s", diskType),
						DiskSizeInMb: intDiskSize * (1 << 10),
					}

					additionalDiskRequest.Resource.DiskAdd = append(additionalDiskRequest.Resource.DiskAdd, diskAdd)
				case "update":
					// collect uniq paths
					if existing_paths[value.Path[0]] {
						continue // Already in the map
					}
					paths = append(paths, value.Path[0])
					existing_paths[value.Path[0]] = true
				case "delete":
					old := os[index].(map[string]interface{}) 
					filename, ok := old["filename"].(string)
					if !ok {
						return fmt.Errorf("can not read filename of disk")
					}

					diskRemove := DiskRemove{
						DiskName: filename,
					}
					additionalDiskRequest.Resource.DiskRemove = append(additionalDiskRequest.Resource.DiskRemove, diskRemove)
				}
			}

			// ToDo: wait for changes applyed or not?
			// ToDo: make functions to add/remove disks
			for _, path := range reverse(paths) {
				index, err := strconv.Atoi(path)
				if err != nil {
					return fmt.Errorf("error converting from string to int: %w", err)
				}

				new := ns[index].(map[string]interface{})

				filename, ok := new["filename"].(string)
				if !ok {
					return fmt.Errorf("can not read filename of disk")
				}

				diskRemove := DiskRemove{
					DiskName: filename,
				}
				additionalDiskRequest.Resource.DiskRemove = append([]DiskRemove{diskRemove}, additionalDiskRequest.Resource.DiskRemove...)

				// ToDo: check for simplest type convertion
				strDiskSize, ok := new["additional_disk_size"].(string)
				if !ok {
					return fmt.Errorf("can not read additional disk size")
				}
				intDiskSize, err := strconv.Atoi(strDiskSize)
				if err != nil {
					return fmt.Errorf("error converting from string to int: %w", err)
				}

				diskType, ok := new["additional_disk_type"].(string)
				if !ok {
					return fmt.Errorf("can not read additional disk type")
				}
				if !containsTag(tags, diskType) {
					return fmt.Errorf("disk type is not available")
				}

				diskAdd := DiskAdd{
					StorageType: diskType,
					Name: "",
					Type: fmt.Sprintf("/managed/storage_type/%s", diskType),
					DiskSizeInMb: intDiskSize * (1 << 10),
				}
				additionalDiskRequest.Resource.DiskAdd = append([]DiskAdd{diskAdd}, additionalDiskRequest.Resource.DiskAdd...)
			}

			log.Println(PrettyStruct(additionalDiskRequest))

			// request and response
			requestBody, err := json.Marshal(additionalDiskRequest)
			if err != nil {
				return fmt.Errorf("error marshaling addititonal disk request: %w", err)
			}

			body := bytes.NewBuffer(requestBody)

			vmId := d.Get("vms.0.id").(string)
			value, err := requestApi("POST", fmt.Sprintf("vms/%s", vmId), body)
			if err != nil {
				return fmt.Errorf("error requesting api vms info: %w", err)
			}

			var responseBody ReconfigurationResponse
			if err = value.Decode(&responseBody); err != nil {
				return fmt.Errorf("error decoding api vms response body: %w", err)
			}

			log.Println(PrettyStruct(responseBody))

			// pretty log output
			log.Println("Disk has been changed")
		}
	}

	// wait ? min for disks applyed?
	// best practise to make read at the end of update
	return nil
}

func resourceServiceDelete(d *schema.ResourceData, m interface{}) error {
	serviceRequest := &ServiceRequest{
		Action: "request_retire",
	}

	requestBody, err := json.Marshal(serviceRequest)

	if err != nil {
		return fmt.Errorf("error marhsaling service retire request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("services/%s", d.Id()), body)

	if err != nil {
		return fmt.Errorf("error requesting service retire: %w", err)
	}

	d.SetId("")

	return nil
}

func reverse(s []string) []string{
	last := len(s) - 1
	for i := 0; i < len(s)/2; i++ {
			s[i], s[last-i] = s[last-i], s[i]
	}
	return s
}

func containsTag(s *TagsResponse, str string) bool {
	for _, v := range s.Resources {
		if v.Name == fmt.Sprintf("/managed/storage_type/%s", str) {
			return true
		}
	}

	return false
}