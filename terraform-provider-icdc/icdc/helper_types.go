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
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemoryMb    string `json:"memory_mb"`
	CpuCores    string `json:"cpu_cores"`
	StorageType string `json:"storage_type"`
	StorageMb   string `json:"storage_mb"`
	Network     string `json:"network"`
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
	AuthType            string `json:"auth_type"`
	Adminpassword       string `json:"adminpassword"`
	SshKey              string `json:"ssh_key"`
	ServiceTemplateHref string `json:"service_template_href"`
	RegionNumber        string `json:"region_number"`
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

type Vm struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Hardware struct {
		MemoryMb int `json:"memory_mb"`
		CpuCores int `json:"cpu_total_cores"`
	} `json:"hardware"`
	Disks []struct {
		Id   string `json:"id"`
		Size int    `json:"size"`
	}
	Network []struct {
		Name string `json:"name"`
	} `json:"lans"`
	Ipaddresses []string `json:"ipaddresses"`
}

type VmReconfigureRequest struct {
	Action   string `json:"action"`
	Resource struct {
		RequestType     string `json:"request_type"`
		VmMemory        string `json:"vm_memory"`
		NumberOfCpus    string `json:"number_of_cpus"`
		NumberOfSockets string `json:"number_of_sockets"`
		CoresPerSocket  string `json:"cores_per_socket"`
	} `json:"resource"`
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
type SecurityGroupCollection struct {
	Resources []SecurityGroup `json:"resources"`
}

type SecurityGroup struct {
	Id            string `json:"id"`
	EmsRef        string `json:"ems_ref"`
	Name          string `json:"name"`
	FirewallRules []struct {
		Id                    string `json:"id"`
		EmsRef                string `json:"ems_ref"`
		Direction             string `json:"direction"`
		NetworkProtocol       string `json:"network_protocol"`
		Port                  int    `json:"port"`
		SourceIpRange         string `json:"source_ip_range"`
		SourceSecurityGroupId string `json:"source_security_group_id"`
	} `json:"firewall_rules"`
}

type SecurityGroupCreateRequest struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

type SecurityGroupTaskResult struct {
	TaskResults struct {
		SecurityGroups SecurityGroupResource `json:"security_group"`
	} `json:"task_results"`
}

type SecurityGroupResource struct {
	EmsRef string `json:"id"`
	Name   string `json:"name"`
}

type SecurityGroupDeleteRequest struct {
	Action string `json:"action"`
	Id     string `json:"id"`
	Name   string `json:"name"`
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
