---
# generated by https://github.com/hashicorp/terraform-plugin-docs
layout: ""
page_title: "icdc_security_rule Resource - terraform-provider-icdc"
subcategory: "networking"
description: |-
  allow to define instance level firewall rules into cloud's managed by icdc platform
---

# icdc_security_rule (Resource)

## Schema

### Required
- `group_id` (String) - ID of security_group
- `direction` (String) - available values: "egress", "ingress"

### Optional
- `port_range` (String) - port range definition (e.g "2200-2205"), if you need to declare specific port just leave one number (e.g. "2200"), by default `Any (0-65432)`
- `protocol` (String) - available values: "", "icmp", "tcp", "udp". By default: `""`, means `Any` protocol
- `network_protocol` (String) - available values: "ipv4", "ipv6". By default: `ipv4`
- `remote_group_id` (String) - by default allow access from all security groups
- `remote_ip_subnet` (String) - ip subnet in canonical form. By default allow access from any subnet (`"0.0.0.0/0"`)

### Read-Only
- `id` (String) The ID of the resource into ICDC applications
- `ems_ref` (String) ID in underlying systems