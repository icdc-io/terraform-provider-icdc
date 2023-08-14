terraform {
  required_providers {
    icdc = {
      version = "1.0.0"
      source = "local.com/nuspenskaya/icdc"
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


resource "icdc_vpc" "vpc_nina1" {
  name = "vpc-name-1"
  router = "nina_test_1"
}

