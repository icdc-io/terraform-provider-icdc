---
subcategory: "Guides"
page_title: "Publishing websites"
description: |-
   A guide for publishing websites in ICDC-powered clouds with Terraform provider.
---

### Description
To deploy a website we are need to create an instance group with enabled webservice, 
load balancer route pointed to an instance and domain record with type CNAME which will point 
to an account load balancer public hostname

### Create an instance group
->Note read [article](networking_infrastructure.md#create-network-and-subnet) to know how create a new subnet

For creating instance group need to know `template_id`. We can copy it from UI of service of fetch using datasource `icdc_template`
Also we need to automate installation of our webservice, as the example we will use `user_data` field to inject part of cloud-init script.
Injected script will be ignored if contains some lint errors.



1. Edit the Terraform configuration to add the `icdc_instance_group` resource
    ```terraform
    data "icdc_template" "centos" {
      name = "CentOS Stream"
      version = "9-230914"
    }
   
   locals {
     httpd_installation = yamlencode({
     runcmd: [
       "dnf install -y httpd",
       "systemctl enable httpd --now"
     ]})
   }

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
     user_data = locals.httpd_installation
   }
    ```

2. Run `terraform apply`

This will produce a new instance group with two instances 1 CPU, 4 Gb RAM and 30 Gb disk size into network named `example`
and OS CentOS. Also, will be automatically installed and ran httpd web-server

### Create load balancer route
For creating load balancer route we need instance group and dns record with cname type

->Note read [article](networking_infrastructure.md#domains) to know how create domain record

1. Edit the Terraform configuration to add the `icdc_alb_route` resource

    ```terraform
    resource "icdc_dns_zone" "z1" {
      name = "my.awesome.zone"
    }

    resource "icdc_dns_record" "r1" {
      zone = icdc_dns_zone.z1.name
      name = "website"
      type = "cname"
      ttl = 600
      data = "alb.public.hostname"
    }

    resource "icdc_alb_route" "website" {
      name = "mywebsite"
      hostname = "${icdc_dns_record.r1.name}.${icdc_dns_zone.z1.name}"
      services = [
        icdc_instance_group.instance-group1.id
      ]
      tls_termination = "edge"
      insecure = "redirect"

      depends_on = [ icdc_dns_record.r1 ]
    }
    ```
2. Run `terraform apply`

Now your website will be available at website.my.awesome.zone

