package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceInstanceGroupRead,
		CreateContext: resourceInstanceGroupCreate,
		UpdateContext: resourceInstanceGroupUpdate,
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
		password = generateSecurePassword()
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
				AdditionalDiskType:  d.Get("additional_disk_type").(string),
				AdditionalDiskSize:  d.Get("additional_disk_size").(string),
				Vlan:                vlan,
				PassAuth:            d.Get("pass_auth").(string),
				Password:            password,
				ManagedAccess:       d.Get("managed_access").(string),
				SecurityGroup:       d.Get("security_group").(string),
				NumberOfVms:         d.Get("instances_count").(string),
				ServiceTemplateHref: fmt.Sprintf("/api/service_templates/%s", d.Get("template_id").(string)),
				UserData:            d.Get("user_data").(string),
				SshKey:              d.Get("ssh_key").(string),
			},
		},
	}

	requestBody, err := json.Marshal(serviceRequest)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	body := bytes.NewBuffer(requestBody)

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

	if serviceRequestResponse == nil || len(serviceRequestResponse.Results) == 0 {
		return append(diags, diag.FromErr(fmt.Errorf("empty service request response"))...)
	}
	if serviceRequestResponse.Results[0].Success != true {
		err = fmt.Errorf(serviceRequestResponse.Results[0].Message)
		return append(diags, diag.FromErr(err)...)
	}

	// Create is asynchronous: this ID is for the request/task, not the final service object.
	serviceRequestId := serviceRequestResponse.Results[0].ServiceRequestId

	var serviceId string

	log.Println("Service", serviceId, "Created")

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		log.Println("Waiting for instance_group creating")

		serviceId, err = fetchDestinationId(serviceRequestId, "Service")
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if serviceId == "" {
			return resource.RetryableError(fmt.Errorf("service destination id is not ready yet"))
		}

		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	log.Printf("[DEBUG] Service %s found for request %s", serviceId, serviceRequestId)

	// Fail fast on provisioning errors instead of waiting for create timeout.
	// The lifecycle state is driven by backend request task status.
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, reqErr := requestApi("GET", fmt.Sprintf("api/compute/v1/services/%s?expand=resources", serviceId), nil)
		if reqErr != nil {
			return resource.RetryableError(fmt.Errorf("failed to fetch service %s lifecycle: %w", serviceId, reqErr))
		}

		var serviceLifecycle *ServiceVmProvisonResponse
		if decodeErr := resp.Decode(&serviceLifecycle); decodeErr != nil {
			return resource.RetryableError(fmt.Errorf("failed to decode service %s lifecycle response: %w", serviceId, decodeErr))
		}
		if serviceLifecycle == nil {
			return resource.RetryableError(fmt.Errorf("empty lifecycle response for service %s", serviceId))
		}

		log.Printf("[DEBUG] Service %s lifecycle_state=%s", serviceId, serviceLifecycle.LifecycleState)
		switch serviceLifecycle.LifecycleState {
		case "error_in_provisioning":
			return resource.NonRetryableError(fmt.Errorf("service %s provisioning failed (lifecycle_state=%s)", serviceId, serviceLifecycle.LifecycleState))
		case "provisioned":
			return nil
		default:
			return resource.RetryableError(fmt.Errorf("waiting for service %s provisioning, current lifecycle_state=%s", serviceId, serviceLifecycle.LifecycleState))
		}
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	iCount, _ := strconv.Atoi(d.Get("instances_count").(string))

	// Wait until the expected VM count appears under the service.
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		log.Println("Waiting for instances creating")

		count, countDiags := instancesCount(serviceId)
		if countDiags.HasError() {
			return resource.RetryableError(fmt.Errorf("failed to count instances for service %s", serviceId))
		}
		if count == iCount {
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
		if networkBody, prettyErr := PrettyStruct(service.Networks); prettyErr == nil {
			log.Printf("[DEBUG] Network response for service %s:\n%s", serviceId, networkBody)
		} else {
			log.Printf("[DEBUG] Failed to pretty print network response for service %s: %v", serviceId, prettyErr)
		}
		log.Println("Waiting for networks config applying")

		// Network readiness means each VM has at least one NIC allocation with a non-empty IP.
		// This avoids returning from create before networking is actually usable.
		allocationsCount := 0
		allocations, _ := vmsAllocationsList(service.Networks)
		for _, allocation := range allocations {
			if allocation.Ip != "" && allocation.Type == "nic" {
				allocationsCount += 1
			}
		}

		if allocationsCount >= iCount {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("error: instances was not created"))
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(serviceId)

	// Final read normalizes all computed/derived fields in state from live API data.
	return resourceInstanceGroupRead(ctx, d, m)
}


func resourceInstanceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	serviceId := d.Id()
	log.Printf("[DEBUG] === START READ for Service ID: %s ===", serviceId)

	_ = ctx
	_ = m

	resp, err := requestApi("GET", fmt.Sprintf("api/compute/v1/services/%s?expand=resources&attributes=vms,networks", serviceId), nil)
	if err != nil {
		log.Printf("[DEBUG] Service %s not found in API. Marking as deleted.", serviceId)
		d.SetId("")
		return nil
	}

	var service *Service

	if err := resp.Decode(&service); err != nil {
		return diag.FromErr(fmt.Errorf("error decoding service data: %w", err))
	}

	sServ, _ := PrettyStruct(service)
	log.Printf("[DEBUG] Service API Response:\n%s", sServ)

	instancesList := buildInstanceListFromService(service)

	var groupCpu, groupMem, groupSysDiskType, groupSysDiskSize, groupAdditionalDiskType, groupAdditionalDiskSize, groupSubnet string

	// Prefer subnet from service networks allocations (instances list),
	// because VM lans can be empty in some API responses.
	if len(instancesList) > 0 {
		firstInstance, ok := instancesList[0].(map[string]interface{})
		if ok {
			networks, ok := firstInstance["networks"].([]interface{})
			if ok && len(networks) > 0 {
				firstNetwork, ok := networks[0].(map[string]interface{})
				if ok {
					if subnet, ok := firstNetwork["subnet"].(string); ok {
						groupSubnet = subnet
					}
				}
			}
		}
	}

	if len(service.Vms) > 0 {
		// Instance group is treated as a homogeneous set for shared fields.
		// We read VM-level hardware/disks from the first VM as the representative
		// source of group CPU/memory/system-disk/additional-disk values.
		firstVmID := service.Vms[0].ID
		// Read hardware/disks from a VM endpoint because service-level VMs payload does not
		// reliably include complete hardware/disk details.
		vmResp, err := requestApi("GET", fmt.Sprintf("api/compute/v1/vms/%s?attributes=hardware,disks,lans", firstVmID), nil)
		if err != nil {
			log.Printf("[WARN] Failed to fetch details for first VM %s: %v", firstVmID, err)
		} else {
			var remoteVm Vm
			if err := vmResp.Decode(&remoteVm); err != nil {
				return diag.FromErr(fmt.Errorf("error decoding first vm data: %w", err))
			}

			sort.SliceStable(remoteVm.Disks, func(i, j int) bool {
				return remoteVm.Disks[i].Id < remoteVm.Disks[j].Id
			})

			// Fallback subnet source: if network allocations were missing in service payload,
			// use VM LAN data as best effort.
			if groupSubnet == "" && len(remoteVm.Network) > 0 {
				groupSubnet = remoteVm.Network[0].Name
			}
			groupCpu = strconv.Itoa(remoteVm.Hardware.CpuCores)
			groupMem = strconv.Itoa(remoteVm.Hardware.MemoryMb)
			// Disk convention used by provider:
			// - index 0 is system disk
			// - index 1 (if present) is the managed additional disk
			if len(remoteVm.Disks) > 0 {
				groupSysDiskSize = strconv.Itoa(remoteVm.Disks[0].Size / (1 << 30))
				groupSysDiskType = diskType(remoteVm.Disks[0].StorageId)
			}
			if len(remoteVm.Disks) > 1 {
				groupAdditionalDiskSize = strconv.Itoa(remoteVm.Disks[1].Size / (1 << 30))
				groupAdditionalDiskType = diskType(remoteVm.Disks[1].StorageId)
			}
		}
	}

	preparedState := map[string]interface{}{
		"id":                   serviceId,
		"name":                 service.Name,
		"description":          service.Description,
		"cpu":                  groupCpu,
		"memory_mb":            groupMem,
		"system_disk_size":     groupSysDiskSize,
		"system_disk_type":     groupSysDiskType,
		"subnet":               groupSubnet,
		"additional_disk_size": groupAdditionalDiskSize,
		"additional_disk_type": groupAdditionalDiskType,
		"instances_count":      strconv.Itoa(len(instancesList)),
		"instances":            instancesList,
	}
	if prepared, err := PrettyStruct(preparedState); err == nil {
		log.Printf("[DEBUG] Prepared instance group state before d.Set:\n%s", prepared)
	} else {
		log.Printf("[DEBUG] Failed to pretty print prepared instance group state: %v", err)
	}

	// Apply all derived values to Terraform state in one read pass.
	if err := d.Set("name", service.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("description", service.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %w", err))
	}
	if err := d.Set("cpu", groupCpu); err != nil {
		return diag.FromErr(fmt.Errorf("error setting cpu: %w", err))
	}
	if err := d.Set("memory_mb", groupMem); err != nil {
		return diag.FromErr(fmt.Errorf("error setting memory_mb: %w", err))
	}
	if err := d.Set("system_disk_size", groupSysDiskSize); err != nil {
		return diag.FromErr(fmt.Errorf("error setting system_disk_size: %w", err))
	}
	if err := d.Set("system_disk_type", groupSysDiskType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting system_disk_type: %w", err))
	}
	if err := d.Set("additional_disk_size", groupAdditionalDiskSize); err != nil {
		return diag.FromErr(fmt.Errorf("error setting additional_disk_size: %w", err))
	}
	if err := d.Set("additional_disk_type", groupAdditionalDiskType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting additional_disk_type: %w", err))
	}
	if err := d.Set("subnet", groupSubnet); err != nil {
		return diag.FromErr(fmt.Errorf("error setting subnet: %w", err))
	}
	if err := d.Set("instances_count", strconv.Itoa(len(instancesList))); err != nil {
		return diag.FromErr(fmt.Errorf("error setting instances_count: %w", err))
	}

	if err := d.Set("instances", instancesList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting instances list: %w", err))
	}

	log.Printf("[DEBUG] === READ FINISHED. Instances found: %d ===", len(instancesList))
	return diags
}

func resourceInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serviceId := d.Id()
	log.Printf("[DEBUG] === START UPDATE (ASYNC SUBMIT MODE) for Service ID: %s ===", serviceId)
	_ = m
	diags := diag.Diagnostics{}

	// Group-level metadata update (service object), independent from per-VM hardware updates.
	if d.HasChanges("name", "description") {
		var req ServiceReconfigureRequest
		req.Action = "edit"
		req.Resource.ID = serviceId
		req.Resource.Name = d.Get("name").(string)
		req.Resource.Description = d.Get("description").(string)

		if payload, err := PrettyStruct(req); err == nil {
			log.Printf("[DEBUG] Sending Service Edit Payload:\n%s", payload)
		}

		body, err := json.Marshal(req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling service edit request for service %s: %w", serviceId, err))
		}

		resp, err := requestApi("POST", fmt.Sprintf("api/compute/v1/services/%s", serviceId), bytes.NewBuffer(body))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error requesting service edit for service %s: %w", serviceId, err))
		}

		var editResp map[string]interface{}
		if err := resp.Decode(&editResp); err != nil {
			return diag.FromErr(fmt.Errorf("error decoding service edit response for service %s: %w", serviceId, err))
		}

		if payload, prettyErr := PrettyStruct(editResp); prettyErr == nil {
			log.Printf("[DEBUG] Service Edit Response for %s:\n%s", serviceId, payload)
		}

		successVal, hasSuccess := editResp["success"]
		message, _ := editResp["message"].(string)
		if hasSuccess {
			if success, ok := successVal.(bool); ok && !success && strings.TrimSpace(message) != "" {
				return diag.FromErr(fmt.Errorf("service edit rejected for service %s: %s", serviceId, strings.TrimSpace(message)))
			}
		}
	}

	// Network profile change is applied through service automation flow (not VM reconfigure endpoint).
	if d.HasChange("subnet") {
		targetSubnet := d.Get("subnet").(string)

		newNetworkName := convertName(targetSubnet)
		if newNetworkName == targetSubnet || newNetworkName == "" {
			parts := strings.Split(targetSubnet, "_")
			newNetworkName = parts[len(parts)-1]
		}
		if newNetworkName != "" {
			newNetworkName = strings.ToUpper(newNetworkName[:1]) + newNetworkName[1:]
		}

		var networkReq ChangeNetworkTypeRequest
		networkReq.Action = "change_network_type"
		networkReq.Resource.Task = "call_automation"
		networkReq.Resource.Path = "Service/Network/StateMachines/ChangeNetworkType/Default"
		networkReq.Resource.Params.DialogNetworkProfile = targetSubnet
		networkReq.Resource.Params.DialogVmID = "all"
		networkReq.Resource.Params.NewNetworkName = newNetworkName

		sReq, _ := PrettyStruct(networkReq)
		log.Printf("[DEBUG] Sending Network Update Payload:\n%s", sReq)

		body, err := json.Marshal(networkReq)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling subnet change request for service %s: %w", serviceId, err))
		}
		resp, err := requestApi("POST", fmt.Sprintf("api/compute/v1/services/%s", serviceId), bytes.NewBuffer(body))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error requesting subnet change for service %s: %w", serviceId, err))
		}

		var networkResp ReconfigurationResponse
		if err := resp.Decode(&networkResp); err != nil {
			return diag.FromErr(fmt.Errorf("error decoding subnet change response for service %s: %w", serviceId, err))
		}
		sResp, _ := PrettyStruct(networkResp)
		log.Printf("[DEBUG] Network Change API Response:\n%s", sResp)
		if !networkResp.Success {
			return diag.FromErr(fmt.Errorf("subnet change rejected for service %s: %s", serviceId, networkResp.Message))
		}
	}

	hasCpuOrMemChange := d.HasChanges("cpu", "memory_mb")
	hasAdditionalDiskChange := d.HasChanges("additional_disk_type", "additional_disk_size")
	oldCpuRaw, newCpuRaw := d.GetChange("cpu")
	oldMemRaw, newMemRaw := d.GetChange("memory_mb")

	oldCpu := getStringValue(oldCpuRaw)
	newCpu := getStringValue(newCpuRaw)
	oldMem := getStringValue(oldMemRaw)
	newMem := getStringValue(newMemRaw)

	isCPUChange := oldCpu != newCpu
	isMemoryChange := oldMem != newMem

	isMemoryIncreaseHotplug := false
	isMemoryDecreaseNonHotplug := false
	if isMemoryChange {
		oldMemInt, oldErr := strconv.Atoi(oldMem)
		newMemInt, newErr := strconv.Atoi(newMem)
		if oldErr == nil && newErr == nil {
			if newMemInt > oldMemInt {
				isMemoryIncreaseHotplug = true
			}
			if newMemInt < oldMemInt {
				isMemoryDecreaseNonHotplug = true
			}
		}
	}

	// Hotplug classification used by this resource (based on current platform behavior):
	// - additional disk add/remove: hotplug
	// - memory increase: hotplug
	// - CPU core count change: non-hotplug
	// - memory decrease: non-hotplug
	//
	// Note: this resource currently manages CPU cores + memory + additional disk only.
	// Socket-level hotplug controls are not exposed in schema yet.
	hasHotplugChange := hasAdditionalDiskChange || isMemoryIncreaseHotplug
	hasNonHotplugChange := isCPUChange || isMemoryDecreaseNonHotplug

	// VM hardware/disk updates are executed per VM in the group.
	// The group is treated as homogeneous: same target values are sent to each VM.
	if hasCpuOrMemChange || hasAdditionalDiskChange {
		instances := d.Get("instances").([]interface{})
		if len(instances) == 0 {
			log.Printf("[WARN] No instances found in state for service %s, fetching instances from API for reconfigure.", serviceId)
			liveInstances, err := fetchInstanceList(serviceId)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch instances for reconfigure in service %s: %w", serviceId, err))
			}
			instances = liveInstances
		}

		if len(instances) == 0 {
			return diag.FromErr(fmt.Errorf("no instances found for service %s, cannot apply vm reconfigure request", serviceId))
		}

		targetCpu := d.Get("cpu").(string)
		targetMem := d.Get("memory_mb").(string)
		targetCpuInt, err := strconv.Atoi(targetCpu)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid cpu value %q: %w", targetCpu, err))
		}
		targetMemInt, err := strconv.Atoi(targetMem)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid memory_mb value %q: %w", targetMem, err))
		}

		targetAdditionalDiskType := getOptionalString(d, "additional_disk_type")
		targetAdditionalDiskSize := getOptionalString(d, "additional_disk_size")
		hasTargetAdditionalDisk := targetAdditionalDiskType != "" || targetAdditionalDiskSize != ""
		// Additional disk settings are managed as a pair to avoid ambiguous requests.
		if hasAdditionalDiskChange && (targetAdditionalDiskType == "") != (targetAdditionalDiskSize == "") {
			return diag.FromErr(fmt.Errorf("`additional_disk_type` and `additional_disk_size` must be set together"))
		}

		targetAdditionalDiskSizeGb := 0
		if hasTargetAdditionalDisk {
			targetAdditionalDiskSizeGb, err = strconv.Atoi(targetAdditionalDiskSize)
			if err != nil {
				return diag.FromErr(fmt.Errorf("invalid additional_disk_size value %q: %w", targetAdditionalDiskSize, err))
			}
		}

		var tags *TagsResponse
		if hasAdditionalDiskChange && hasTargetAdditionalDisk {
			// Validate requested disk storage type against provider tags before submitting updates.
			tagsResp, err := requestApi("GET", "api/compute/v1/tags?expand=resources&attributes=classification&filter[]=name='/managed/storage_type/*'", nil)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error requesting storage types: %w", err))
			}
			if err := tagsResp.Decode(&tags); err != nil {
				return diag.FromErr(fmt.Errorf("error decoding tags response: %w", err))
			}
			if !containsTag(tags, targetAdditionalDiskType) {
				return diag.FromErr(fmt.Errorf("disk type %q is not available", targetAdditionalDiskType))
			}
		}

		appliedVMs := make([]string, 0, len(instances))

		for i, inst := range instances {
			vmData, ok := inst.(map[string]interface{})
			if !ok {
				return diag.FromErr(fmt.Errorf("invalid instance format at index %d for service %s", i, serviceId))
			}
			vmIDRaw, exists := vmData["id"]
			if !exists {
				return diag.FromErr(fmt.Errorf("instance id missing at index %d for service %s", i, serviceId))
			}
			vmId, ok := vmIDRaw.(string)
			if !ok || vmId == "" {
				return diag.FromErr(fmt.Errorf("invalid instance id at index %d for service %s", i, serviceId))
			}

			log.Printf("[DEBUG] ---> Processing VM %d/%d (ID: %s) <---", i+1, len(instances), vmId)
			log.Printf("[DEBUG] Target CPU: %s, Target RAM: %s", targetCpu, targetMem)

			var vmReq VmReconfigureRequest
			vmReq.Action = "reconfigure"
			vmReq.Resource.RequestType = "vm_reconfigure"
			vmReq.Resource.VmMemory = targetMemInt
			vmReq.Resource.NumberOfCpus = targetCpuInt
			vmReq.Resource.CoresPerSocket = targetCpuInt
			vmReq.Resource.NumberOfSockets = 1
			vmReq.Resource.DiskAdd = []DiskAdd{}
			vmReq.Resource.DiskRemove = []DiskRemove{}

			if hasAdditionalDiskChange {
				vmResp, err := requestApi("GET", fmt.Sprintf("api/compute/v1/vms/%s?attributes=disks", vmId), nil)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to fetch VM %s disks: %w", vmId, err))
				}

				var liveVm Vm
				if err := vmResp.Decode(&liveVm); err != nil {
					return diag.FromErr(fmt.Errorf("failed to decode VM %s disks: %w", vmId, err))
				}

				sort.SliceStable(liveVm.Disks, func(i, j int) bool {
					return liveVm.Disks[i].Id < liveVm.Disks[j].Id
				})

				// Provider disk convention:
				// - index 0 is system disk
				// - remaining disks are treated as additional disks managed by this resource
				additionalDisks := make([]struct {
					Filename string
					SizeGb   int
					DiskType string
				}, 0)
				if len(liveVm.Disks) > 1 {
					for _, disk := range liveVm.Disks[1:] {
						additionalDisks = append(additionalDisks, struct {
							Filename string
							SizeGb   int
							DiskType string
						}{
							Filename: disk.Filename,
							SizeGb:   disk.Size / (1 << 30),
							DiskType: diskType(disk.StorageId),
						})
					}
				}

				if hasTargetAdditionalDisk {
					// Keep existing additional disk when it already matches requested type/size.
					alreadyMatches := len(additionalDisks) == 1 &&
						additionalDisks[0].SizeGb == targetAdditionalDiskSizeGb &&
						additionalDisks[0].DiskType == targetAdditionalDiskType

					if !alreadyMatches {
						// Reconcile by removing current additional disks and adding requested one.
						for _, disk := range additionalDisks {
							vmReq.Resource.DiskRemove = append(vmReq.Resource.DiskRemove, DiskRemove{DiskName: disk.Filename})
						}
						vmReq.Resource.DiskAdd = append(vmReq.Resource.DiskAdd, DiskAdd{
							StorageType:  targetAdditionalDiskType,
							Name:         "",
							Type:         fmt.Sprintf("/managed/storage_type/%s", targetAdditionalDiskType),
							DiskSizeInMb: targetAdditionalDiskSizeGb * (1 << 10),
						})
					}
				} else {
					// Empty additional disk target means remove all additional disks from VM.
					for _, disk := range additionalDisks {
						vmReq.Resource.DiskRemove = append(vmReq.Resource.DiskRemove, DiskRemove{DiskName: disk.Filename})
					}
				}
			}

			if !hasCpuOrMemChange && len(vmReq.Resource.DiskAdd) == 0 && len(vmReq.Resource.DiskRemove) == 0 {
				log.Printf("[DEBUG] No disk changes needed for VM %s, skipping reconfigure.", vmId)
				continue
			}

			sVmReq, _ := PrettyStruct(vmReq)
			log.Printf("[DEBUG] Sending Reconfigure Payload for VM %s:\n%s", vmId, sVmReq)

			reqBody, err := json.Marshal(vmReq)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to marshal reconfigure request for VM %s: %w", vmId, err))
			}
			resp, err := requestApi("POST", fmt.Sprintf("api/compute/v1/vms/%s", vmId), bytes.NewBuffer(reqBody))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to reconfigure VM %s: %w", vmId, err))
			}

			var reconfResp ReconfigurationResponse
			if err := resp.Decode(&reconfResp); err == nil {
				sResp, _ := PrettyStruct(reconfResp)
				log.Printf("[DEBUG] Reconfigure API Response for VM %s:\n%s", vmId, sResp)
				if !reconfResp.Success {
					return diag.FromErr(fmt.Errorf("reconfigure rejected for VM %s: %s", vmId, reconfResp.Message))
				}
			} else {
				return diag.FromErr(fmt.Errorf("failed to decode reconfigure response for VM %s: %w", vmId, err))
			}

			appliedVMs = append(appliedVMs, vmId)
		}

		if hasHotplugChange && len(appliedVMs) > 0 {
			// Hotplug changes should be observable quickly; poll live VM state to confirm apply.
			waitErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
				for _, vmID := range appliedVMs {
					resp, err := requestApi("GET", fmt.Sprintf("api/compute/v1/vms/%s?attributes=hardware,disks", vmID), nil)
					if err != nil {
						return resource.RetryableError(fmt.Errorf("failed to read VM %s while waiting hotplug apply: %w", vmID, err))
					}

					var liveVm Vm
					if err := resp.Decode(&liveVm); err != nil {
						return resource.RetryableError(fmt.Errorf("failed to decode VM %s while waiting hotplug apply: %w", vmID, err))
					}

					if isMemoryIncreaseHotplug {
						if strconv.Itoa(liveVm.Hardware.MemoryMb) != newMem {
							return resource.RetryableError(fmt.Errorf("waiting memory hotplug apply for VM %s", vmID))
						}
					}

					if hasAdditionalDiskChange {
						sort.SliceStable(liveVm.Disks, func(i, j int) bool {
							return liveVm.Disks[i].Id < liveVm.Disks[j].Id
						})

						additionalDisks := make([]struct {
							Id        string
							Size      int
							Filename  string
							StorageId string
						}, 0)
						if len(liveVm.Disks) > 1 {
							for _, disk := range liveVm.Disks[1:] {
								additionalDisks = append(additionalDisks, struct {
									Id        string
									Size      int
									Filename  string
									StorageId string
								}{
									Id:        disk.Id,
									Size:      disk.Size,
									Filename:  disk.Filename,
									StorageId: disk.StorageId,
								})
							}
						}
						targetAdditionalDiskType := getOptionalString(d, "additional_disk_type")
						targetAdditionalDiskSize := getOptionalString(d, "additional_disk_size")
						hasTargetAdditionalDisk := targetAdditionalDiskType != "" && targetAdditionalDiskSize != ""

						if !hasTargetAdditionalDisk && len(additionalDisks) != 0 {
							return resource.RetryableError(fmt.Errorf("waiting disk remove hotplug apply for VM %s", vmID))
						}

						if hasTargetAdditionalDisk {
							targetAdditionalDiskSizeGb, err := strconv.Atoi(targetAdditionalDiskSize)
							if err != nil {
								return resource.NonRetryableError(fmt.Errorf("invalid additional_disk_size value %q: %w", targetAdditionalDiskSize, err))
							}
							if len(additionalDisks) != 1 {
								return resource.RetryableError(fmt.Errorf("waiting disk add hotplug apply for VM %s", vmID))
							}
							currentDiskType := diskType(additionalDisks[0].StorageId)
							currentDiskSizeGb := additionalDisks[0].Size / (1 << 30)
							if currentDiskType != targetAdditionalDiskType || currentDiskSizeGb != targetAdditionalDiskSizeGb {
								return resource.RetryableError(fmt.Errorf("waiting disk reconfigure hotplug apply for VM %s", vmID))
							}
						}
					}
				}
				return nil
			})
			if waitErr != nil {
				return diag.FromErr(fmt.Errorf("hotplug changes were submitted but did not apply in time: %w", waitErr))
			}
		}

		if hasNonHotplugChange {
			// Non-hotplug changes may be accepted but applied only after stop/start cycle.
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Non-hotplug changes may require stop/start",
				Detail:   "Some configuration changes require the instance to be stopped and started again before they take effect.",
			})
		}

		if hasHotplugChange {
			// Final read loop ensures Terraform state reflects hotplug changes after async apply.
			readSyncErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
				readDiags := resourceInstanceGroupRead(ctx, d, m)
				if readDiags.HasError() {
					return resource.RetryableError(fmt.Errorf("read after hotplug update returned error diagnostics"))
				}

				if isMemoryIncreaseHotplug {
					currentMem := getOptionalString(d, "memory_mb")
					if currentMem != newMem {
						return resource.RetryableError(fmt.Errorf("waiting state sync for memory_mb: expected=%s got=%s", newMem, currentMem))
					}
				}

				if hasAdditionalDiskChange {
					currentDiskType := getOptionalString(d, "additional_disk_type")
					currentDiskSize := getOptionalString(d, "additional_disk_size")

					if hasTargetAdditionalDisk {
						if currentDiskType != targetAdditionalDiskType || currentDiskSize != targetAdditionalDiskSize {
							return resource.RetryableError(fmt.Errorf("waiting state sync for additional disk: expected=%s/%s got=%s/%s", targetAdditionalDiskType, targetAdditionalDiskSize, currentDiskType, currentDiskSize))
						}
					} else if currentDiskType != "" || currentDiskSize != "" {
						return resource.RetryableError(fmt.Errorf("waiting state sync for additional disk removal: got=%s/%s", currentDiskType, currentDiskSize))
					}
				}

				return nil
			})
			if readSyncErr != nil {
				return diag.FromErr(fmt.Errorf("hotplug changes applied but state did not synchronize in time: %w", readSyncErr))
			}
		}
	}

	log.Printf("[DEBUG] === UPDATE FINISHED (requests submitted, async apply) for Service ID: %s ===", serviceId)
	return diags
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

// getOptionalString safely reads optional string fields from Terraform state.
// It returns empty string for nil/non-string values.
func getOptionalString(d *schema.ResourceData, key string) string {
	value := d.Get(key)
	if value == nil {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}

// getStringValue converts generic values (for example from GetChange) into string
// for comparison/logging, with nil mapped to empty string.
func getStringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
