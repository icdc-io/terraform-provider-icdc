# Terraform Provider ICDC

## Available resources
||||||
|---|---|---|---|---|
||**Name**|**Internal name**|**Description**|**Documentation link**|
|1.|Service|icdc_service|Service - abstract entity which includes inside - virtual machine, network port, backups, schedules, accesses|https://help.icdc.io/compute/en/Use_of_the_Services/index.html|
|2.|Subnet|icdc_subnet| Fully isolated subnet which allow to assign virtual machines network ports|https://help.icdc.io/networking/en/VPC_Networks.html|
|3.|Security group|icdc_security_group|OVN security group, allow to enable/disable incoming/upcoming trafic|
|4.|Security group rule|icdc_security_group_rule| A security group rules|

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
|icdc_subnet||||||
||**parameter**|**description**|**type**|||
||name|subnet **external name(internal name will be applied after creating)**|string|||
||cidr|cidr of subnet|string|||
||network_protocol|IP protocol|string|||
||ip_version|IP version *(will be removed in future)*|int|||
||gateway|Address of subnet gateway|string|||
||dns_nameservers|List of dns nameservers|list of strings|||
|icdc_security_group||||||
||**parameter**|**description**|**type**|||
||name|Security group name|string|||
|icdc_security_group_rule||||||
||**parameter**|**description**|**type**|||
||resource_id|security group id (database)|string|||
||security_group_id|security group ems_ref|string|||
||direction|trafik direction (ingress/egress)|string|||
||protocol|TCP/UDP/Any|string|||
||network_protocol|network protocol IPV4/IPV6|string|||
||port|start_port value|string|||
||end_port|end_port_value|string|||
||source_ip_range|source ip range|string|||
||

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
