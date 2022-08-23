terraform {
  required_providers {
    icdc = {
      version = "0.1"
      source = "local.com/ahrechushkin/icdc"
    }
  }
}

provider "icdc" {
  //username = ""
  //password = ""
  //location = "ycz"
  //group    = "icdc.admin"
}

resource "icdc_service_request" "new" {
  service_name = "terraform-test"
  description = "terraform-plugin"
  vm_memory = "4096"
  number_of_sockets = "1"
  cores_per_socket = "1"
  hostname = "generated-hostname"
  vlan = "ycz_icdc_base (ycz_icdc_base)"
  system_disk_type = "nvme"
  system_disk_size = "30"
  auth_type = "key"
  adminpassword = "generate_password"
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDU2ixuBrQNH/XeWowxh4CeAjvlkT0Pnz6+GWu1UyJR7/N5TC2DF4cFp4EKWdSxkxKmVg8DUBgsUNJofpDfDJLcwP+kKpYEEiMT4VL4FnPZzDD0hbUbfzaBQCUNtJRHLT91qkysOgm08jaFUlWTI6JhaybVowmpiD0nv1UQW98SKzrVYMXxDv1PSAvESJG8YyQ0zf/RslwaHyiyiqm5uLHoXHEO77ddNkRB5e3meQKiIwEr1f0BjUVgh+kINSlOQLl3euDHaniBAbt6qPOtFHSYXs993rqK3TRN180nigfdSGoJc6FrWF7MiuFC4lUmnk2MzFdGM0TWU/1eniQ0WxfE/lUMI4bIa813+z43cllOvQQitxIgVFRWtJsKm6Lbnw20ioT34rrKKWxHYCI5JvrA7vx39IsgrbFsU952BXOTLVvMPUVGyQYTwIRkmPlJ2GHyicDTBUYv7FGFjVz7gw7ZCIH5HNWSn+57rUdJVzZV+eUM8mrfkPnDOQniRxbnnkk= ahrechushkin@hrechushkin-av"
  service_template_href = "/api/service_templates/18000000000034"
  region_number = "18"
}