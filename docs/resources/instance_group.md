---
page_title: "icdc_instance_group Resource - terraform-provider-icdc"
subcategory: "compute"
description: |-
  Resource icdc_instance_group details.
---

# icdc_instance_group (Resource)
An **instance group** (old Service) is a collection of virtual machines (VMs) instances that you can manage as a single entity.

## Schema
### Required

- `cpu` (String) - CPU size per instance
- `instances_count` (String) - count of instances
- `memory_mb` (String) - RAM size per instance
- `name` (String) - name of instance group
- `pass_auth` (String) - password usage policy, one of "disable_password", "own_password", "temporary_password"
- `password` (String) - your own password, if you're leave empty string - password will be generated automatically
- `subnet` (String) - the name of vpc subnet
- `system_disk_size` (String) - system disk size per instance - in Gb
- `system_disk_type` (String) - system disk type
- `template_id` (String) - base template_id (related to datastore icdc_template)

### Optional

- `additional_disk_size` (String)
- `additional_disk_type` (String)
- `description` (String)
- `managed_access` (String)
- `security_group` (String)
- `ssh_key` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_data` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `instances` (Block List) (see [below for nested schema](#nestedblock--instances))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)


<a id="nestedblock--instances"></a>
### Nested Schema for `instances`

Read-Only:

- `id` (String) The ID of this resource.
- `name` (String)
- `networks` (List of Object) (see [below for nested schema](#nestedatt--instances--networks))

<a id="nestedatt--instances--networks"></a>
### Nested Schema for `instances.networks`

Read-Only:

- `hostname` (String)
- `ip` (String)
- `mac` (String)
- `nic` (String)
- `subnet` (String)
- `type` (String)
