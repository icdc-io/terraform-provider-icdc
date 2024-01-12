package icdc

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/sethvargo/go-password/password"
)

type InstanceGroupRequest struct {
	Action    string                   `json:"action"`
	Resources []InstanceGroupResources `json:"resources"`
}

type Instance struct {
	id   string `json:"id"`
	name string `json:"name"`
}

type InstanceGroupResources struct {
	ServiceName         string `json:"service_name"`
	ServiceDescription  string `json:"service_description"`
	VmMemory            string `json:"vm_memory"`
	Cpu                 string `json:"cpu"`
	SystemDiskType      string `json:"system_disk_type"`
	SystemDiskSize      string `json:"system_disk_size"`
	AdditionalDiskType  string `json:"additional_disk_type"`
	AdditionalDiskSize  string `json:"additional_disk_size"`
	Vlan                string `json:"vlan"`
	PassAuth            string `json:"pass_auth"`
	Password            string `json:"password"`
	SecurityGroup       string `json:"security_group"`
	SshKey              string `json:"ssh_key"`
	NumberOfVms         string `json:"number_of_vms"`
	UserData            string `json:"user_data"`
	ManagedAccess       string `json:"managed_access"`
	ServiceTemplateHref string `json:"service_template_href"`
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

func generateSecurePassword() string {
	p, err := password.Generate(16, 4, 2, false, true)
	if err != nil {
		log.Fatal(err)
	}

	match, err := regexp.MatchString(`^[^\d]*[A-Z][^\d]*\d[^\d]*.{6,}[^\d]*$`, p)

	if err != nil {
		log.Fatal(err)
	}

	if !match {
		return generateSecurePassword()
	}

	if p[len(p)-1] > 47 && p[len(p)-1] < 58 {
		return generateSecurePassword()
	}

	if p[0] > 47 && p[0] < 58 {
		return generateSecurePassword()
	}

	return p
}

func validatePassword(v interface{}, p cty.Path) diag.Diagnostics {
	value := v.(string)

	if value != "" {
		match, err := regexp.MatchString(`^[^\d]*[A-Z][^\d]*\d[^\d]*.{6,}[^\d]*$`, value)

		if err != nil {
			return diag.FromErr(err)
		}

		if !match {
			return diag.FromErr(fmt.Errorf("invalid password, for security reasons password requires minimum 8 symbols, at least 1 uppercase and 1 number (but not at first or last position)"))
		}
	}

	return diag.Diagnostics{}
}
