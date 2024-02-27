---
page_title: "icdc_alb_route Resource - terraform-provider-icdc"
subcategory: "networking"
description: |-
    Resource icdc_alb_route details.
---
# icdc_alb_route (Resource)
Application Load Balancer (ALB) is a network service that distributes incoming public web traffic between virtual servers to provide fault tolerance for websites and applications.

How it works:
The Load Balancer (ALB) allows to easily configure web or TLS traffic to virtual machines, belonging to the same Compute service based on a domain name in the HTTP protocol or the Server Name Indication (SNI) extension of TLS.
When the volume of incoming traffic changes dramatically, the balancer evenly distributes the entire volume of requests between resources according to the Round-robin algorithm. ALB provides automatic renewal of LetsEncrypt certificates and the ability to upload your own certificates. Forwarding through ALB does not require opening virtual machine destination ports in the Network Security Group in the Firewall tab. ALB provides access on the main public address of the account only on ports 80 and 443. These ports can be reassigned in the Port-Forwarding application.

The Health Check feature allows to exclude unhealthy servers from the load balancing rotation. ALB monitors the status of the virtual machines responsible for handling the web traffic route, and will consider VMs healthy as long as they return status codes between 2XX and 3XX to the health check requests (carried out every interval). If a VM is found to be unavailable (powered off or broken), the HealthCheck ensures that web traffic is not directed to that particular VM.

ALB management is only available to users with the Account Admin role.

[Read official documentation for more details.](https://docs.icdc.io/networking/load_balancer/load_balancer/#route-creation)

## Schema

### Required

- `hostname` (String)
- `name` (String)
- `services` (List of String)

### Optional

- `cloudgw_name` (String)
- `healthcheck` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--healthcheck))
- `insecure` (String)
- `ip_version` (Number)
- `path` (String)
- `target_port` (Number)
- `tls_termination` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--healthcheck"></a>
### Nested Schema for `healthcheck`

Optional:

- `follow_redirects` (Boolean)
- `hostname` (String)
- `interval` (Number)
- `method` (String)
- `path` (String)
- `port` (Number)
- `scheme` (String)
- `timeout` (Number)
