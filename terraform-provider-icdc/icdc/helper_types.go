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
	Network string `json:"network"`
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

// network
type NetworkRequestResponse struct {
	Network struct {
		Id                  string `json:"id"`
		Name                string `json:"name"`
		TenantId            string `json:"tenant_id"`
		Status              string `json:"status"`
		PortSecurityEnabled bool   `json:"port_security_enabled"`
		Mtu                 int    `json:"mtu"`
		Subnet              struct {
			Id              string `json:"id"`
			Cidr            string `json:"cidr"`
			NetworkId       string `json:"network_id"`
			IpVersion       int    `json:"ip_version"`
			TenantId        string `json:"tenant_id"`
			RouterId        string `json:"router_id"`
			EnableDhcp      bool   `json:"enable_dhcp"`
			AllocationPools []struct {
				Start string `json:"start"`
				Stop  string `json:"stop"`
			} `json:"allocation_pools"`
			DnsNameservers []string `json:"dns_nameservers"`
			Name           string   `json:"name"`
			GatewayIp      string   `json:"gateway_ip"`
		} `json:"subnet"`
		Ports []struct {
			Id           string `json:"id"`
			Name         string `json:"name"`
			NetworkId    string `json:"network_id"`
			MacAddress   string `json:"mac_address"`
			AdminStateUp bool   `json:"admin_state_up"`
			DeviceId     string `json:"device_id"`
			DeviceOwner  string `json:"device_owner"`
			FixedIps     []struct {
				IpAddress string `json:"ip_address"`
				SubnetId  string `json:"subnet_id"`
			} `json:"fixed_ips"`
			SecuriyGroups []string `json:"security_groups"`
			Type          string   `json:"type"`
		} `json:"ports"`
	} `json:"network"`
}

type NetworkCreateBody struct {
	Name   string       `json:"name"`
	Mtu    int          `json:"mtu"`
	Subnet SubnetParams `json:"subnet"`
}

type Network struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	TenantId            string `json:"tenant_id"`
	Status              string `json:"status"`
	PortSecurityEnabled bool   `json:"port_security_enabled"`
	Mtu                 int    `json:"mtu"`
	Subnet              struct {
		Id              string `json:"id"`
		Cidr            string `json:"cidr"`
		NetworkId       string `json:"network_id"`
		IpVersion       int    `json:"ip_version"`
		TenantId        string `json:"tenant_id"`
		RouterId        string `json:"router_id"`
		EnableDhcp      bool   `json:"enable_dhcp"`
		AllocationPools []struct {
			Start string `json:"start"`
			Stop  string `json:"stop"`
		} `json:"allocation_pools"`
		DnsNameservers []string `json:"dns_nameservers"`
		Name           string   `json:"name"`
		GatewayIp      string   `json:"gateway_ip"`
	} `json:"subnet"`
	Ports []struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		NetworkId    string `json:"network_id"`
		MacAddress   string `json:"mac_address"`
		AdminStateUp bool   `json:"admin_state_up"`
		DeviceId     string `json:"device_id"`
		DeviceOwner  string `json:"device_owner"`
		FixedIps     []struct {
			IpAddress string `json:"ip_address"`
			SubnetId  string `json:"subnet_id"`
		} `json:"fixed_ips"`
		SecuriyGroups []string `json:"security_groups"`
		Type          string   `json:"type"`
	} `json:"ports"`
	Metadata int `json:"metadata"`
}

type SubnetParams struct {
	Name      string `json:"name"`
	IpVersion int    `json:"ip_version"`
	Cidr      string `json:"cidr"`
	GatewayIp string `json:"gateway_ip"`
	//Ipv6AddressMode string `json:"ipv6_address_mode"`
	EnableDhcp     bool   `json:"enable_dhcp"`
	DnsNameservers string `json:"dns_nameservers"`
}

type CloudNetworkRequest struct {
	Network NetworkCreateBody `json:"network"`
}

// security group/rule
type SecurityGroupRequestResponse struct {
	SecurityGroup SecurityGroup `json:"security_group"`
}

type SecurityGroup struct {
	Id                 string                    `json:"id"`
	TenantId           string                    `json:"tenant_id"`
	RouterId           string                    `json:"router_id"`
	Description        string                    `json:"description"`
	Name               string                    `json:"name"`
	SecurityGroupRules []SecurityGroupRulesInVpc `json:"security_group_rules"`
}

type SecurityGroupInVpc struct {
	Id                 string                    `json:"id"`
	TenantId           string                    `json:"tenant_id"`
	RouterId           string                    `json:"router_id"`
	Description        string                    `json:"description"`
	Name               string                    `json:"name"`
	SecurityGroupRules []SecurityGroupRulesInVpc `json:"security_group_rules"`
}

type SecurityGroupRuleRequestResponse struct {
	SecurityGroupRule struct {
		Id              string `json:"id"`
		Direction       string `json:"direction"`
		SecurityGroupId string `json:"security_group_id"`
		Description     int    `json:"description"`
		Erthertype      string `json:"erthertype"`
		RemoteIpPrefix  int    `json:"remote_ip_prefix"`
		PortRangeMax    string `json:"port_range_max"`
		PortRangeMin    string `json:"port_range_min"`
		Protocol        string `json:"protocol"`
		RemoteGroupId   string `json:"remote_group_id"`
	} `json:"security_group_rule"`
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

type SecurityGroupRulesInVpc struct {
	Id              string `json:"id"`
	Direction       string `json:"direction"`
	SecurityGroupId string `json:"security_group_id"`
	Description     int    `json:"description"`
	Erthertype      string `json:"erthertype"`
	RemoteIpPrefix  int    `json:"remote_ip_prefix"`
	PortRangeMax    int    `json:"port_range_max"`
	PortRangeMin    int    `json:"port_range_min"`
	Protocol        int    `json:"protocol"`
	RemoteGroupId   int    `json:"remote_group_id"`
}

type GroupCreateBody struct {
	SecurityGroup SecurityGroupBody `json:"security_group"`
}

type SecurityGroupBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RuleCreateBody struct {
	SecurityGroupRule SecurityGroupRuleBody `json:"security_group_rule"`
}

type SecurityGroupRuleBody struct {
	Direction       string `json:"direction"`
	Ethertype       string `json:"ethertype"`
	PortRangeMin    string `json:"port_range_min"`
	PortRangeMax    string `json:"port_range_max"`
	Protocol        string `json:"protocol"`
	RemoteGroupId   string `json:"remote_group_id"`
	SecurityGroupId string `json:"security_group_id"`
}

// vpc
type VpcGetResponse []struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	TenantId string `json:"tenant_id"`
}

type VpcRequestResponse struct {
	Vpc struct {
		Id             string               `json:"id"`
		Name           string               `json:"name"`
		TenantId       string               `json:"tenant_id"`
		Router         RouterResponse       `json:"router"`
		Networks       []Network            `json:"networks"`
		SecurityGroups []SecurityGroupInVpc `json:"security_groups"`
	} `json:"vpc"`
}
type RouterResponse struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	AdminStateUp        bool     `json:"admin_state_up"`
	Status              string   `json:"status"`
	TenantId            string   `json:"tenant_id"`
	ExternalGatewayInfo int      `json:"external_gateway_info"`
	Routes              []string `json:"routes"`
}

type VpcCreateBody struct {
	Vpc VpcStructBody `json:"vpc"`
}

type VpcStructBody struct {
	Name   string           `json:"name"`
	Router RouterCreateBody `json:"router"`
}

type RouterCreateBody struct {
	Name string `json:"name"`
}

type CloudVpcRequest struct {
	Vpc VpcCreateBody `json:"vpc"`
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
