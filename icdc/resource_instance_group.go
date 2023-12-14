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
)

func resourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Read:          resourceInstanceGroupRead,
		CreateContext: resourceInstanceGroupCreate,
		Update:        resourceInstanceGroupUpdate,
		Delete:        resourceInstanceGroupDelete,
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
				Optional: true,
				Default:  "",
			},
			"cpu": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memory_mb": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_disk_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_disk_size": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"additional_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"managed_access": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "password_reset",
			},
			"pass_auth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validatePassword,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"security_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"networks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nic": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
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

func resourceInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vlan := fmt.Sprintf("%s (%s)", d.Get("subnet").(string), d.Get("subnet").(string))

	password := d.Get("password").(string)

	if password == "" {
		password = generate_secure_password()
	}

	serviceRequest := &InstanceGroupRequest{
		Action: "add",
		Resources: []InstanceGroupResources{
			{
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
				UserData:            d.Get("user_data").(string),
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
		for _, allocation := range allocations {
			if allocation.Ip != "" && allocation.Type == "nic" {
				allocationsCount += 1
			}
		}

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

func resourceInstanceGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInstanceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInstanceGroupDelete(d *schema.ResourceData, m interface{}) error {
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
