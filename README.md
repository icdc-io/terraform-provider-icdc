# Terraform Provider ICDC

## Available resources
||||||
|---|---|---|---|---|
||**Name**|**Internal name**|**Description**|**Documentation link**|
|1.|Service|icdc_service|Service - abstract entity which includes inside - virtual machine, network port, backups, schedules, accesses|https://help.icdc.io/compute/en/Use_of_the_Services/index.html|
|2.|Subnet|icdc_network| Fully isolated subnet which allow to assign virtual machines network ports|https://help.icdc.io/networking/en/VPC_Networks.html|
|3.|Security group|icdc_security_group|OVN security group, allow to enable/disable incoming/upcoming trafic|
|4.|Security group rule|icdc_security_group_rule| A security group rules|
|5.|Vpc|icdc_vpc| A virtual private cloud|

## Resources description
|||||||
|---|---|---|---|---|---|
|icdc_service||||||
||**parameter**|**description**|**type**|||
||name|Service Name|string|||
||service_temaplte_id|Template ID (define OS)|string|||
||ssh_key|User ssh key for access to VMs|string|||
||vms|A list of **vms(currently available only 1 per service)**|list|||
|||**parameter**|**description**|**type**||
|||cpu_cores|Count of CPU cores per vm|string||
|||memory_mb|VM RAM size|string||
|||system_disk_type|OS disk type|string||
|||system_disk_size|OS disk size|string||
|||subnet|Name of VM network|string||
|||additional_disk|A list of additional_disks **(optional)** |list||
||||**parameter**|**description**|**type**|
||||additional_disk_type|Type of additional disk |string|
||||additional_disk_size|Size of additional disk (in gb)|string|
|---|---|---|---|---|---|
|icdc_network||||||
||**parameter**|**description**|**type**|||
||name|subnet **external name(internal name will be applied after creating)**|string|||
||cidr|cidr of subnet|string|||
||network_protocol|IP protocol|string|||
||ip_version|IP version *(will be removed in future)*|int|||
||gateway|Address of subnet gateway|string|||
||dns_nameservers|dns nameserver(yes, only one)|string|||
||mtu|Maximum transmission unit|string|||
||vpc_id|Vpc id|internal parameter, should point to previous created vpc|||
|icdc_security_group||||||
||**parameter**|**description**|**type**|||
||name|Security group name|string|||
||description|Security group description|string|||
||vpc_id|Vpc id|internal parameter, should point to previous created vpc|||
|icdc_security_group_rule||||||
||**parameter**|**description**|**type**|||
||ethertype|The layer 3 protocol type, valid values are IPv4 or IPv6|string|||
||security_group_id|internal parameter, should point to previous created security group|string|||
||direction|The direction of the rule, valid values are ingress or egress|string|||
||protocol|TCP/UDP/Any|string|||
||remote_ip_prefix|The remote CIDR, the value needs to be a valid CIDR (i.e. 192.168.0.0/16).|string|||
||port_range_max|The higher part of the allowed port range, valid integer value needs to be between 1 and 65535|string|||
||port_range_min|The lower part of the allowed port range, valid integer value needs to be between 1 and 65535|string|||
||remote_group_id|The remote group id, the value needs to be an Openstack ID of a security group in the same tenant|string|||


## Provider parameters

|||
|---|---|
|**name**|**description**|
|username|userid/user email which using to access to cloud|
|password|user password|
|location|managed location *(mb will be moved to resources)*|
|location_number|the number of location *(will be removed in future)*| 
|account|user account|
|role|user role|
|platform|platform name *(icdc/scdc)*|



## Usage

1. Clone respository and build provider for specified OS arch.
```bash
> git clone git@github.com:icdc-io/terraform-provider-icdc.git
> cd terraform-provider-icdc/terraform-provider-icdc
> make
```

2. Define TF Plan (see [example](examples/main.tf))

3. Control your cloud resources with [terraform](https://www.terraform.io/docs)
