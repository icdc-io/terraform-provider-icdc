# Terraform ICDC Provider

The [ICDC Provider](https://registry.terraform.io/providers/icdc-io/icdc/latest/docs) allows [Terraform](https://terraform.io) to manage [ICDC](https://icdc.io) resources.

## Usage example

```hcl
# 1. Specify the version of the ICDC Provider to use
terraform {
  required_providers {
    icdc = {
      source = "icdc-io/icdc"
      version = "=1.0.0"
    }
  }
}

# 2. Configure the ICDC Provider
provider "icdc" {
    username = "username"
    location = "user-region"
    auth_group = "user-authgroup"
}

# 3. Create a virtual network
resource icdc_network example {
  name = "example-network"
  subnet {
    cidr = "11.0.0.0/26"
    gateway = "11.0.0.1"
    dns_nameserver = "8.8.8.8"
  }
}

# 4. Use data-source for select needed version of OS
data icdc_template centos-stream{
  name = "CentOS Stream"
  version = "9-230519"
}

# 5. Deploy instance group with selected OS into created network
resource icdc_instance_group instance-group1 {
  name = "instance-group-1"
  template_id = data.icdc_template.centos-stream.id
  subnet = icdc_network.example.name
  cpu = "1"
  memory_mb = "4096"
  system_disk_type = "nvme"
  system_disk_size = "30"
  pass_auth = "temporary_password"
  instances_count = "2"
  user_data = <<-EOT
            runcmd:
            - dnf install -y httpd
            - systemctl enable httpd --now
        EOT
}

# 6. You can find out the list of supported resources below.
```

## Supported services and resources

- [ICDC Networking](https://icdc.io/networking)
  - [VPC Networks](./docs/resources/network.md)
  - Security [Groups](./docs/resources/security_group.md) and [Rules](./docs/resources/security_group_rule.md)
  - [Load Balancer routes](./docs/resources/alb_route.md)
  - DNS [zones](./docs/resources/dns_zone.md) and [records](./docs/resources/dns_record.md)
- [ICDC Compute](https://icdc.io/compute)
  - [Services](./docs/resources/service.md)
  - [Instance Groups](./docs/resources/instance_group.md)
  - Datasource - [templates](./docs/data-sources/template.md)

You can find out more examples of using ICDC provider [here](./docs/guides/)