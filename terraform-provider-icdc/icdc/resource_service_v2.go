package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sethvargo/go-password/password"
)

func resourceServiceV2() *schema.Resource {
	return &schema.Resource{
		Read:          resourceServiceV2Read,
		CreateContext: resourceServiceV2Create,
		Update:        resourceServiceV2Update,
		Delete:        resourceServiceV2Delete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"memory_mb": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"system_disk_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"system_disk_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_disk_size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"additional_disk_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances_count": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"managed_access": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "password_reset",
			},
			"pass_auth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"security_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"instances": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
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
						"networks": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"nic": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"hostname": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceServiceV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vlan := fmt.Sprintf("%s (%s)", d.Get("subnet").(string), d.Get("subnet").(string))

	password := d.Get("password").(string)

	if password == "" {
		password = generate_secure_password()
	}

	serviceRequest := &ServiceV2Request{
		Action: "add",
		Resources: []ServiceV2Resources{ServiceV2Resources{
			ServiceName:         d.Get("name").(string),
			ServiceDescription:  d.Get("description").(string),
			Cpu:                 d.Get("cpu").(string),
			VmMemory:            d.Get("memory_mb").(string),
			SystemDiskType:      d.Get("system_disk_type").(string),
			SystemDiskSize:      d.Get("system_disk_size").(string),
			Vlan:                vlan,
			PassAuth:            d.Get("pass_auth").(string),
			Password:            password,
			ManagedAccess:       d.Get("managed_access").(string),
			SecurityGroup:       d.Get("security_group").(string),
			NumberOfVms:         d.Get("instances_count").(string),
			ServiceTemplateHref: fmt.Sprintf("/api/service_templates/%s", d.Get("template_id").(string)),
		},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

	// prettystruct for logs
	log.Println(PrettyStruct(serviceRequest))

	responseBody, err := requestApi("POST", "api/compute/v1/service_orders/cart/service_requests/", body)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	var serviceRequestResponse *ServiceRequestResponse
	if err = responseBody.Decode(&serviceRequestResponse); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	log.Println(PrettyStruct(serviceRequestResponse))

	serviceRequestId := serviceRequestResponse.Results[0].ServiceRequestId

	var serviceId string

	log.Println("Service", serviceId, "Created")

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		log.Println("Waiting for instance_group creating")

		serviceId, err = fetchDestinationId(serviceRequestId, "Service")
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if serviceId != "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("error: service is not created"))
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	instances_count, _ := strconv.Atoi(d.Get("instances_count").(string))

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		log.Println("Waiting for instances creating")

		count, _ := instancesCount(serviceId)
		if count == instances_count {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("error: instances was not created"))
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	

	
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		requestUrl := fmt.Sprintf("api/compute/v1/services/%s?expand=resources&attributes=networks", serviceId)
		responseBody, err = requestApi("GET", requestUrl, nil)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("error: cant fetch service"))
		}

		var service *Service

		err = responseBody.Decode(&service)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("error: cant parse service object"))
		}
		log.Println("Waiting for networks config applying")

        //we need to fetch all allocations with type - nic and non-empty ip addresses
		allocationsCount := 0
		allocations, _ := vmsAllocationsList(service.Networks)
		log.Println("DEBUG ALLOCATIONS LIST:", allocations)
		for _, allocation := range allocations {
			if (allocation.Ip != "" && allocation.Type == "nic") {
				allocationsCount += 1
			}
		}

		log.Println("INSTANCES COUNT")
		log.Println("DEBUG ALLOCATIONS COUNT:", allocationsCount)
		if allocationsCount >= instances_count {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("error: instances was not created"))
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	instancesList, err := fetchInstanceList(serviceId)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(serviceId)

	err = d.Set("instances", instancesList)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return nil
}

func instancesCount(serviceId string) (int, diag.Diagnostics) {
	var diags diag.Diagnostics

	requestUrl := fmt.Sprintf("api/compute/v1/services/%s?expand=resources&attributes=vms", serviceId)
	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return 0, append(diags, diag.FromErr(err)...)
	}

	var service *Service

	err = responseBody.Decode(&service)
	if err != nil {
		return 0, append(diags, diag.FromErr(err)...)
	}

	return len(service.Vms), nil
}

func fetchInstanceList(serviceId string) ([]interface{}, error) {

	requestUrl := fmt.Sprintf("api/compute/v1/services/%s?expand=resources&attributes=vms,networks", serviceId)

	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return nil, err
	}

	var service *Service

	err = responseBody.Decode(&service)
	if err != nil {
		return nil, err
	}

	instances := service.Vms
	instancesList := make([]interface{}, len(instances))
	vmsAllocations, _ := vmsAllocationsList(service.Networks)

	for ndx, instance := range instances {
		i := make(map[string]interface{})
		i["id"] = instance.ID
		i["name"] = instance.Name

		var vmAllocations []VmAllocation

		for _, allocation := range vmsAllocations {
			if allocation.VmId == instance.ID {
				vmAllocations = append(vmAllocations, allocation)
			}
		}

		allocationsList := make([]interface{}, len(vmAllocations))

		for ndx, allocation := range vmAllocations {
			a := make(map[string]interface{})
			a["subnet"] = allocation.Subnet
			a["mac"] = allocation.Mac
			a["ip"] = allocation.Ip
			a["hostname"] = allocation.Hostname
			a["nic"] = allocation.NicName
			a["type"] = allocation.Type

			allocationsList[ndx] = a
		}

		i["networks"] = allocationsList
		instancesList[ndx] = i
	}

	return instancesList, nil
}

func vmsAllocationsList(networks []ComputeNetwork) ([]VmAllocation, error) {

	var vmAllocations []VmAllocation

	for _, network := range networks {
		for _, allocation := range network.Allocations {
			vmAllocation := VmAllocation{
				VmId:     strconv.Itoa(allocation.VmId),
				NicName:  allocation.NicName,
				Mac:      allocation.Mac,
				Hostname: allocation.Hostname,
				Ip:       allocation.Ip,
				Type:     allocation.Type,
				Subnet:   network.Name,
				Gateway:  network.Gateway,
				Cidr:     network.Cidr,
			}

			vmAllocations = append(vmAllocations, vmAllocation)
		}
	}

	return vmAllocations, nil

}

func resourceServiceV2Read(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceV2Update(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceV2Delete(d *schema.ResourceData, m interface{}) error {
	serviceRequest := &ServiceRequest{
		Action: "request_retire",
	}

	requestBody, err := json.Marshal(serviceRequest)

	if err != nil {
		return fmt.Errorf("error marhsaling instance_group delete request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	_, err = requestApi("POST", fmt.Sprintf("api/compute/v1/services/%s", d.Id()), body)

	if err != nil {
		return fmt.Errorf("error requesting instance_group delete: %w", err)
	}

	d.SetId("")

	return nil
}

func generate_secure_password() string {
	res, err := password.Generate(16, 4, 2, false, true)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
