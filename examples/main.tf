terraform {
  required_providers {
    icdc = {
      version = "1.0.0"
    }
  }
}

provider "icdc" {
    username = "YOU_EMAIL.TEST"
    location = "LOC"
    auth_group = "ACCOUNT.ROLE"
}

/*
// CREATING VPC NETWORK RESOURCE

resource icdc_network tfdemo {
  name = ""
  mtu = 1500
  subnet {
    cidr = "10.97.4.0/26"
    gateway = "10.97.4.1"
    dns_nameserver = "194.213.212.2"
  }
}

// CREATING INSTANCE_GROUP RESOURCE (service_v2)
resource icdc_instance_group tfdemoig {
  name = "tfdemo1"
  description = "Terraform Demo"
  cpu = "1"
  memory_mb = "4096"
  system_disk_type = "nvme"
  system_disk_size = "30"
  subnet = icdc_network.tfdemo.subnet.0.name
  pass_auth = "temporary_password"
  password = ""
  template_id = "27000000000001"
  instances_count = 2
  user_data = "runcmd:\n- dnf install -y httpd\n- systemctl enable httpd --now"

  depends_on = [ icdc_network.tfdemo ]
}

//CREATING DNS_ZONE RESOURCE
resource icdc_dns_zone tfdemoz {
  name = "tfdemo"
}

//CREATING DNS_RECORD RESOURCE
resource icdc_dns_record tfdemo {
  type = "cname"
  name = "webserver"
  zone = icdc_dns_zone.tfdemoz.name
  data = "MYHOSTNAME"
  ttl = 600

  depends_on = [ icdc_dns_zone.tfdemoz ]

}

//CREATING LB_ROUTE RESOURCE
resource icdc_alb_route tfdemor {
  name = "tfdemo"
  hostname = "${icdc_dns_record.tfdemo.name}.${icdc_dns_record.tfdemo.zone}"
  insecure = "redirect"
  tls_termination = "edge"
  ip_version = 4

  services = [
    icdc_instance_group.tfdemoig.id
  ]

  depends_on = [ 
    icdc_instance_group.tfdemoig, 
    icdc_dns_record.tfdemo
  ]
}
*/