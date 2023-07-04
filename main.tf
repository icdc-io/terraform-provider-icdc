terraform {
  required_providers {
    icdc = {
      version = "1.0.0"
      source = "local.com/nuspenskaya/icdc"
    }
  }
}

provider "icdc" {
  username = "nuspenskaya@ibagroup.eu"
  password = "Ilive4me12345"
  location = "ycz"
  location_number = 18
  account = "icdc"
  role = "member"
  platform = "icdc"
}

resource "icdc_network" "net-nina" {
  vpc_id = "3be5b80f-61de-4bc9-9fdc-1ff1b123bc11"
  name = "tf-net-nina"
  mtu = "1500"
  ip_version = "6"
  dns_nameservers = ["194.213.212.130"]
  network_id = "a6ed4a28-edff-4379-ace6-6c167be8a578"
  enable_dhcp = "true"
  cidr = "192.168.1.0/22"
  gateway_ip = "192.168.1.1"
  ipv6_address_mode = "dhcpv6-stateful"
}
