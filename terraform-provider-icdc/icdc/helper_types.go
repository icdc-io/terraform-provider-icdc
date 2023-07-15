package icdc

// Provider structures
type IcdcToken struct {
	ApiGateway string
	Group      string
	Jwt        string
}

type JwtToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Service structures
type Service struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	SshKey            string
	ServiceTemplateId string     `json:"service_template_id"`
	Vms               []VmParams `json:"vms"`
}

type VmParams struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	MemoryMb       string `json:"memory_mb"`
	CpuCores       string `json:"cpu_cores"`
	SystemDiskType string `json:"system_disk_type"`
	SystemDiskSize string `json:"system_disk_size"`
	// mb change
	AdditionalDisk     string `json:"additional_disk"`
	AdditionalDiskType string `json:"additional_disk_type"`
	AdditionalDiskSize string `json:"additional_disk_size"`
	Network            string `json:"network"`
}

type VmParamsForRead struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	MemoryMb       string `json:"memory_mb"`
	CpuCores       string `json:"cpu_cores"`
	SystemDiskType string `json:"system_disk_type"`
	SystemDiskSize string `json:"system_disk_size"`
	AdditionalDisk []struct {
		AdditionalDiskType string `json:"additional_disk_type"`
		AdditionalDiskSize string `json:"additional_disk_size"`
	} `json:"additional_disk"`
	Subnet string `json:"network"`
}

type ServiceResources struct {
	ServiceName         string `json:"service_name"`
	VmMemory            string `json:"vm_memory"`
	NumberOfSockets     string `json:"number_of_sockets"`
	CoresPerSocket      string `json:"cores_per_socket"`
	Hostname            string `json:"hostname"`
	Vlan                string `json:"vlan"`
	SystemDiskType      string `json:"system_disk_type"`
	SystemDiskSize      string `json:"system_disk_size"`
	AdditionalDisk      string `json:"additional_disk"`
	AdditionalDiskType  string `json:"additional_disk_type"`
	AdditionalDiskSize  string `json:"additional_disk_size"`
	AuthType            string `json:"auth_type"`
	Adminpassword       string `json:"adminpassword"`
	SshKey              string `json:"ssh_key"`
	ServiceTemplateHref string `json:"service_template_href"`
	// RegionNumber        string `json:"region_number"`
}

type ServiceRequest struct {
	Action    string             `json:"action"`
	Resources []ServiceResources `json:"resources"`
}

type ServiceRequestResponse struct {
	Results []struct {
		Success            bool   `json:"success"`
		Message            string `json:"message"`
		ServiceRequestId   string `json:"service_request_id"`
		ServiceRequestHref string `json:"service_request_href"`
		Href               string `json:"href"`
	} `json:"results"`
}

type ServiceMiqRequest struct {
	MiqRequestTasks []struct {
		DestinationId   string `json:"destination_id"`
		DestinationType string `json:"destination_type"`
	} `json:"miq_request_tasks"`
}

type TagsResponse struct {
	Name      string `json:"name"`
	Resources []struct {
		Name string `json:"name"`
	} `json:"resources"`
}

type ServiceVmProvisonResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	LifecycleState string `json:"lifecycle_state"`
	Vms            []struct {
		Href string `json:"href"`
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"vms"`
}

type DataStoreResponse struct {
	Id   string `json:"id"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
}

type Vm struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Hardware struct {
		MemoryMb int `json:"memory_mb"`
		CpuCores int `json:"cpu_total_cores"`
	} `json:"hardware"`
	Disks []struct {
		Id        string `json:"id"`
		Size      int    `json:"size"`
		Filename  string `json:"filename"`
		StorageId string `json:"storage_id"`
	} `json:"disks"`
	Network []struct {
		Name string `json:"name"`
	} `json:"lans"`
	Ipaddresses []string `json:"ipaddresses"`
}

type VmReconfigureRequest struct {
	Action   string `json:"action"`
	Resource struct {
		CoresPerSocket  string       `json:"cores_per_socket"`
		DiskAdd         []DiskAdd    `json:"disk_add,omitempty"`
		DiskRemove      []DiskRemove `json:"disk_remove,omitempty"`
		NumberOfCpus    string       `json:"number_of_cpus"`
		NumberOfSockets string       `json:"number_of_sockets"`
		RequestType     string       `json:"request_type"`
		VmMemory        string       `json:"vm_memory"`
	} `json:"resource"`
}

type DiskAdd struct {
	DiskSizeInMb int    `json:"disk_size_in_mb"`
	Name         string `json:"name"`
	StorageType  string `json:"storage_type"`
	Type         string `json:"type"`
}

type DiskRemove struct {
	DiskName string `json:"disk_name"`
}

type ServiceReconfigureRequest struct {
	Action   string `json:"action"`
	Resource struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"resource"`
}

type ChangeNetworkTypeRequest struct {
	Action   string `json:"action"`
	Resource struct {
		Params struct {
			DialogNetworkProfile string `json:"dialog_network_profile"`
		} `json:"params"`
		Path string `json:"path"`
		Task string `json:"task"`
	} `json:"resource"`
}

// Subnet structures
type Network struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Subnets []Subnet `json:"cloud_subnets"`
}
type Subnet struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	EmsRef          string   `json:"ems_ref"`
	EmsId           string   `json:"ems_id"`
	CloudNetworkId  string   `json:"cloud_network_id"`
	Cidr            string   `json:"cidr"`
	Gateway         string   `json:"gateway"`
	IpVersion       int      `json:"ip_version"`
	NetworkProtocol string   `json:"network_protocol"`
	DnsNameservers  []string `json:"dns_nameservers"`
	NetworkRouterId string   `json:"network_router_id"`
}

type SubnetCreateBody struct {
	Cidr            string   `json:"cidr"`
	IpVersion       int      `json:"ip_version"`
	NetworkProtocol string   `json:"network_protocol"`
	Name            string   `json:"name"`
	DnsNameservers  []string `json:"dns_nameservers"`
}

type NetworkCollection struct {
	Resources []Network `json:"resources"`
}

type CloudNetworkRequest struct {
	Action string           `json:"action"`
	Name   string           `json:"name"`
	Subnet SubnetCreateBody `json:"subnet"`
}

// Security groups structures
type GroupCreateBody struct {
	SecurityGroup SecurityGroupBody `json:"security_group"`
}

type RuleCreateBody struct {
	SecurityGroupRule SecurityGroupRuleBody `json:"security_group_rule"`
}

type SecurityGroupRuleBody struct {
	Direction       string `json:"direction"`
	NetworkProtocol string `json:"network_protocol"`
	Ethertype       string `json:"ethertype"`
	PortRangeMin    string `json:"port_range_min"`
	PortRangeMax    string `json:"port_range_max"`
	Protocol        string `json:"protocol"`
	RemoteGroupId   string `json:"remote_group_id"`
	SecurityGroupId string `json:"security_group_id"`
}

type SecurityGroupBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SecurityGroupRequestResponse struct {
	Id                 string               `json:"id"`
	Name               string               `json:"name"`
	TenantId           string               `json:"tenant_id"`
	RouterId           string               `json:"router_id"`
	Description        string               `json:"description"`
	SecurityGroupRules []SecurityGroupRules `json:"security_group_rules"`
}

type SecurityGroupRuleRequestResponse struct {
}

type SecurityGroupRules []struct {
	Id              string `json:"id"`
	Direction       string `json:"direction"`
	SecurityGroupId string `json:"security_group_id"`
	Description     string `json:"description"`
	Erthertype      string `json:"erthertype"`
	RemoteIpPrefix  string `json:"remote_ip_prefix"`
	PortRangeMax    string `json:"port_range_max"`
	PortRangeMin    string `json:"port_range_min"`
	Protocol        string `json:"protocol"`
	RemoteGroupId   string `json:"remote_group_id"`
}

// Security group rules structures
type AddSecurityGroupRule struct {
	Action          string `json:"action"`
	Direction       string `json:"direction"`
	NetworkProtocol string `json:"network_protocol"`
	PortRangeMin    string `json:"port_range_min"`
	PortRangeMax    string `json:"port_range_max"`
	Protocol        string `json:"protocol"`
	RemoteGroupId   string `json:"remote_group_id"`
	SecurityGroupId string `json:"security_group_id"`
	SourceIpRange   string `json:"source_ip_range"`
}

type SecurityGroupRule struct {
	Id              string `json:"id"`
	EmsRef          string `json:"ems_ref"`
	Direction       string `json:"direction"`
	NetworkProtocol string `json:"network_protocol"`
	PortRangeMin    string `json:"port_range_min"`
	PortRangeMax    string `json:"port_range_max"`
	Protocol        string `json:"host_protocol"`
	RemoteGroupId   string `json:"remote_group_id"`
	SecurityGroupId string `json:"security_group_id"`
	SourceIpRange   string `json:"source_ip_range"`
	ResourceId      string `json:"resource_id"`
	ResourceType    string `json:"resource_type"`
}

type SecurityGroupRulesCollection struct {
	Rules []SecurityGroupRule `json:"firewall_rules"`
}

// General structures
type EmsProvider struct {
	Resources []struct {
		Id string `json:"id"`
	} `json:"resources"`
}

type DeleteRequest struct {
	Action string `json:"action"`
	Id     string `json:"id"`
}

type TaskResponse struct {
	Results []struct {
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		TaskId   string `json:"task_id"`
		TaskHref string `json:"task_href"`
	} `json:"results"`
}

// Response structures
type ReconfigurationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Href    string `json:"href"`
}
