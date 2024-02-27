---
subcategory: "Guides"
page_title: "Building Networking Infrastructure"
description: |-
   A guide for building network infrastructure in ICDC-powered clouds with Terraform provider.
---

### VPC
-> Note: Since we have not released networking v2, by default, each account has 1 preconfigured VPC router.
VPC includes a virtual router, networks, subnets, load balancer gateway, VPN gateway.

### Create Network and Subnet
-> Note: One subnet per one network.

1. Edit the Terraform configuration to add the `icdc_network` resource
    ```terraform
    resource "icdc_network" "example" {
      name = "example"
      subnet {
        cidr = "10.124.0.0/24"
        gateway = "10.124.0.1"
        dns_nameserver = "8.8.8.8"
      }
    } 
    ```
2. Run `terraform apply`

   This will produce a new network and subnet, with automatically enabled DHCP.

Now you can use the network to create an `instance group` or `service` with instances and VMs with network interfaces in this subnet.
### Instance Level Firewall
-> Note: By default, we have disabled `internal firewall` like `ufw`, `firewalld`, etc.

Instance level firewall includes security_groups and security_rules.

1. Edit the Terraform configuration to add the `icdc_security_group` resource
    ```terraform
    resource "icdc_security_group" "example" {
      name = "example_group"
    }
    ```
2. Run `terraform apply` to create security group

   This will produce a new security group without any rules. Keep in mind: when creating security groups using UI, by default, you will have two default egress rules.
3. Edit the Terraform configuration to add the `icdc_security_rule` resources
    ```terraform
    resource "icdc_security_rule" "example1" {
      group_id = icdc_security_group.example.id
      direction = "egress"
      network_protocol = "ipv4"
    }
   
   resource "icdc_security_rule" "example2" {
      group_id = icdc_security_group.example.id
      direction = "egress"
      network_protocol = "ipv6"
    }
   
   resource "icdc_security_rule" "example3" {
      group_id = icdc_security_group.example.id
      direction = "ingress"
      network_protocol = "ipv4"
      port_range = "443"
      protocol = "tcp"
    }
    ```
4. Run `terraform apply` to create security rules

   This will produce three security rules allowing each outbound traffic and inbound TCP IPv4 traffic on port 443.
### Domains
-> Restriction: we currently support only `A, AAAA, CNAME, TXT, MX, SRV, NS` types of records.

-> Note: By default, each account has two domain zones: `ACC.cmp.LOC.icdc.io` and `ACC.vpn.LOC.icdc.io`.

1. Edit the Terraform configuration to add the `icdc_dns_zone` resource
   ```terraform
   resource "icdc_dns_zone" "example" {
      name = "myzone.example.com"
   }

2. Run `terraform apply`

   This will produce a new domain zone `myzone.example.com` with one default `NS` record.

3. Edit the Terraform Configuration to Add the `icdc_dns_record` Resources

   ```terraform
   resource "icdc_dns_record" "example_a" {
      zone = icdc_dns_zone.example
      type = "a"
      name = "typea"
      data = "10.124.0.4"
      ttl = 600
   }
   
   resource "icdc_dns_record" "example_aaaa" {
      zone = icdc_dns_zone.example
      type = "aaaa"
      name = "typeaaaa"
      data = "3919:ec0b:c16b:25ad:595f:c7e4:93a7:63dc"
      ttl = 600
   }
   
   resource "icdc_dns_record" "example_cname" {
      zone = icdc_dns_zone.example
      type = "cname"
      name = "typecname"
      data = "typea.myzone.example.com"
      ttl = 600
   }
   
   resource "icdc_dns_record" "example_mx" {
      zone = icdc_dns_zone.example
      type = "mx"
      name = "typemx"
      data = "smtp.gmail.com"
      ttl = 600
      priority = 10
   }
   
   resource "icdc_dns_record" "example_txt" {
      zone = icdc_dns_zone.example
      type = "txt"
      name = "typetxt"
      data = "whatever"
      ttl = 600
   }
   
   resource "icdc_dns_record" "example_srv" {
      zone = icdc_dns_zone.example
      type = "srv"
      name = "_sip._tcp"
      data = "typea.myzone.example.com"
      ttl = 600
      priority = 10
      weight = 50
      port = 34567
   }
   ```
4. Run `terraform apply`

This will produce new domain records:

- `typea.myzone.example.com`: The record will resolve into the address `10.124.0.4`
- `typeaaaa.myzone.example.com`: The record will resolve into an address `3919:ec0b:c16b:25ad:595f:c7e4:93a7:63dc`
- `typecname.myzone.example.com`: The record will resolve into the hostname `typea.myzone.example.com`
- `typemx.myzone.example.com`: The record will resolve into `smtp.gmail.com` with priority 10
- `typetxt.myzone.example.com`: The record will return text `whatever`
- `_sip._tcp.myzone.example.com`: The record will pointed to `typea.myzone.example.com` listening on port 34567 for SIP protocol services