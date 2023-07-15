terraform {
  required_providers {
    icdc = {
      version = "1.0.0"
      source = "local.com/ahrechushkin/icdc"
    }
  }
}

provider "icdc" {
    username = "ahrechushkin@ibagroup.eu"
    #password = ""
    location = "zby"
    auth_group = "icdc.admin"
    # auth_server - optional parameter, needed for development goals
    #auth_server = "login.icdc.io"
}

data "icdc_template" "debian" {
  name = "Debian:10.7"
}

data "icdc_template" "centos8" {
  name = "CentOS:8.3"
}


resource icdc_service tf_srv3 {
  name = "tf_srv3"
  service_template_id = data.icdc_template.centos8.id
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDYW/MTYTUGvoFrvPZNJwDGMx6i81VPHVLAb28HSVLG1zQVCZCKk80lqUqU0lRNSyCxQYoCTHl0e1IUnNR0tVxyYzDU88VpZMDaGvxND2Gcpv+UwpE6vJseiScSrRkW2VSEbBYD9joscysSsm6BAM3gb8oR6WBbzRb5C8X2Hz5jmlXqVMEK2qJU565OJa7BkzcvIcD/0swjcG6cjOMFoiwWpP/j9qELFxrdU5lbM82ucmv8YnZ3MzS2RrwHpV+TqhDuVP6+TjkCW1gswUU6HQK6d91O63nJZT2cQmQzjumGRfJ3U08zowSS6dJWv3e+/7zKI/Ylcy06qnpqrnYI7gkgQWdNpfLX5mfx33aYIyN0GYIytahDDhXOnVCdF+nHg+02mNmglB28KwTlK1LYRuBiAtxesTU2C33pOV3GS16Z+EmhgqtYiI0W+ryvl6pmpqyzrQ13fHOQuaKvYZpCQd9GtDZwkyB0zdqQd6n++b1K1Fq9Y2CDOOnD/4PrEoprTnU= ahrechushkin@workstation"
  vms {
    cpu_cores = "1"
    memory_mb = "2024"
    system_disk_type = "nvme"
    system_disk_size = "30"
    subnet = "zby_icdc_base"
  }
}


resource icdc_service tf_srv4 {
  name = "tf_srv4"
  service_template_id = data.icdc_template.debian.id
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDYW/MTYTUGvoFrvPZNJwDGMx6i81VPHVLAb28HSVLG1zQVCZCKk80lqUqU0lRNSyCxQYoCTHl0e1IUnNR0tVxyYzDU88VpZMDaGvxND2Gcpv+UwpE6vJseiScSrRkW2VSEbBYD9joscysSsm6BAM3gb8oR6WBbzRb5C8X2Hz5jmlXqVMEK2qJU565OJa7BkzcvIcD/0swjcG6cjOMFoiwWpP/j9qELFxrdU5lbM82ucmv8YnZ3MzS2RrwHpV+TqhDuVP6+TjkCW1gswUU6HQK6d91O63nJZT2cQmQzjumGRfJ3U08zowSS6dJWv3e+/7zKI/Ylcy06qnpqrnYI7gkgQWdNpfLX5mfx33aYIyN0GYIytahDDhXOnVCdF+nHg+02mNmglB28KwTlK1LYRuBiAtxesTU2C33pOV3GS16Z+EmhgqtYiI0W+ryvl6pmpqyzrQ13fHOQuaKvYZpCQd9GtDZwkyB0zdqQd6n++b1K1Fq9Y2CDOOnD/4PrEoprTnU= ahrechushkin@workstation"
  vms {
    cpu_cores = "1"
    memory_mb = "2024"
    system_disk_type = "nvme"
    system_disk_size = "30"
    subnet = "zby_icdc_base"
  }
}

resource icdc_subnet tf_sbnt {
  name = "tf_sbnt"
  cidr = "9.110.0.0/26"
  gateway = "9.110.0.1"
  dns_nameserver = "178.172.238.130"
  network_protocol = "ipv4"
}

resource "icdc_security_group" "security_group_nina1" {
  name = "RDP_SSH"
  vpc_id = "3be5b80f-61de-4bc9-9fdc-1ff1b123bc11"
  description = "Allow incoming 22 and 3389 tcp"
  direction = "egress"
  ethertype = "IPv4"
  remote_ip_prefix = "string"
  port_range_max = "2203"
  port_range_min = "2200"
  protocol = "icmp"
  remote_group_id = "3fa85f64-5717-4562-b3fc-2c963f66afa6"
}