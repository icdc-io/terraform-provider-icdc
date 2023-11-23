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
	Networks          []ComputeNetwork `json:"networks"`
}

/*
{
"name": "dcz_devel_mm_test_net",
"cidr": "10.10.10.0/24",
"gateway": "10.10.10.1",
"allocations":[{"hostname": "ahrechushkin-tf-02.devel.cmp.dcz.icdc.io", "ip": "10.10.10.12", "mac": "1c:dc:15:00:00:9d",…],
"parameters":{"display_name": "Mm test net"}
}
*/
type ComputeNetwork struct {
	Name string `json:"name"`
	Cidr string `json:"cidr"`
	Gateway string `json:"gateway"`
	Allocations []NetworkAllocation `json:"allocations"`
}

type NetworkAllocation struct {
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Mac      string `json:"mac"`
	VmId     int `json:"vm_id"`
	NicName  string `json:"nic_name"`
	Type     string `json:"type"`
}

type VmAllocation struct {
	Hostname string
	Ip       string
	Mac      string
	VmId     string
	NicName  string
	Subnet   string
	Gateway  string
	Cidr     string
	Type     string
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
//	RegionNumber        string `json:"region_number"`
}

type ServiceV2Resources struct {
	ServiceName 				string `json:"service_name"`
	ServiceDescription  string `json:"service_description"`
	VmMemory						string `json:"vm_memory"`
	Cpu                 string `json:"cpu"`
	SystemDiskType      string `json:"system_disk_type"`
	SystemDiskSize      string `json:"system_disk_size"`
	AdditionalDiskType  string `json:"additional_disk_type"`
	AdditionalDiskSize  string `json:"additional_disk_size"`
	Vlan                string `json:"vlan"`
	PassAuth						string `json:"pass_auth"`
	Password            string `json:"password"`
	SecurityGroup       string `json:"security_group"`
	SshKey              string `json:"ssh_key"`
	NumberOfVms         string `json:"number_of_vms"`
	UserData            string `json:"user_data"`
	ManagedAccess       string `json:"managed_access"`
	ServiceTemplateHref string `json:"service_template_href"`
}

type Instance struct {
	id string `json:"id"`
	name string `json:"name"`
}

type ServiceRequest struct {
	Action    string             `json:"action"`
	Resources []ServiceResources `json:"resources"`
}

type ServiceV2Request struct {
	Action    string             `json:"action"`
	Resources []ServiceV2Resources `json:"resources"`
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
type SecurityGroupCollection struct {
	Resources []SecurityGroup `json:"resources"`
}

type SecurityGroup struct {
	Id     string `json:"id"`
	EmsRef string `json:"ems_ref"`
	Name   string `json:"name"`
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

type MetadataDns struct {
	Account string `json:"account"`
	Owner		string `json:"owner"`
	Zone    bool `json:"zone"`
	Service bool `json:"service"`
}

type DnsZoneResponse struct {
	Data []struct {
		Name string `json:"name"`
		Metadata MetadataDns `json:"metadata"`
	} `json:"data"`
}

type DnsRecordResponse struct {
	Data []struct {
		Type 		 string `json:"type"`
		Id 	 		 string `json:"id"`
		/*
		Priority int `json:"priority"`
		Weight 	 int `json:"weight"`
		Port		 int `json:"port"`
		*/
		Ttl 		 int `json:"ttl"`
		Group    string `json:"group"`
		Data     string `json:"data"`
		Name		 string `json:"name"`
	} `json:"data"`
}

type AddDnsZone struct {
	Zone struct {
		Name string `json:"name"`
	} `json:"zone"`
}

type AddDnsRecord struct {
	Record struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Data     string `json:"data"`
		Ttl 		 int `json:"ttl"`
	} `json:"record"`
}

type AddMxRecord struct {
	Record struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Data     string `json:"data"`
		Priority int `json:"priority"`
		Ttl 		 int `json:"ttl"`
	} `json:"record"`
}

type AddSrvRecord struct {
	Record struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Data     string `json:"data"`
		Priority int `json:"priority"`
		Weight 	 int `json:"weight"`
		Port 		 int `json:"port"`
		Ttl 		 int `json:"ttl"`
	} `json:"record"`
}


type AddDnsZoneResponse struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

type AddDnsRecordResponse struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}