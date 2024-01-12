---
layout: ""
page_title: "Provider: ICDC"
description: |-
  The ICDC provider is used to configure your infrastructure under ICDC platform control
---

# ICDC Provider
The ICDC provider is used to configure your infrastructure under [ICDC Platform](https://icdc.io/) control.

## Example Usage
A typical provider configuration will look something like:

```hcl
provider "icdc" {
    username = "username"
    location = "user-region"
    auth_group = "user-authgroup"
}
```

See the [provider reference](#icdc-provider-schema) page for details on authentication and configuring the provider.


## Resources and available features
- [icdc_service](resources/service.md) - produce and operate [Compute](https://docs.at.icdc.io/compute/overview/) services
  - create
  - update
  - delete
- [icdc_instance_group](resources/instance_group.md) - produce and operate [Compute](https://docs.at.icdc.io/compute/overview/) services, but use a new feature "provisioning v2".
  - create
  - delete
- [icdc_network](resources/network.md) - produce and operate [VPC networks with subnets](https://docs.at.icdc.io/networking/vpc_networks/vpc_networks)
  - create
  - delete
- [icdc_dns_zone](resources/dns_zone.md) - produce and operate [domain zones](https://docs.at.icdc.io/networking/dns_domains/dns_domains/)
  - create
  - delete
- [icdc_dns_record](resources/dns_record.md) - produce and operate [domain records](https://docs.at.icdc.io/networking/dns_domains/dns_domains/)
  - create
  - delete
- [icdc_alb_route](resources/alb_route.md) - produce and operate [load balancer routes](https://docs.at.icdc.io/networking/load_balancer/load_balancer/)
  - create
  - update
  - delete
- [icdc_security_group](resources/security_group.md) - produce and operate [security groups](https://docs.at.icdc.io/networking/firewall/)
  - create
  - delete
- [icdc_security_rule](resources/security_rule.md) - produce and operate [security group rules](https://docs.at.icdc.io/networking/firewall/)
  - create
  - delete

## ICDC Provider Schema
### Required
- `auth_group` (String, Sensitive) - User active group, contains needed account and role
- `username` (String, Sensitive)
- `location` (String, Sensitive) - operated location

### Optional
- `password` (String, Sensitive) - user password, also user can declare it using env variable `ICDC_PASSWORD`
- `sso_client_id` (String, Sensitive)
- `sso_realm` (String, Sensitive) - basically operator name
- `sso_url` (String, Sensitive)


