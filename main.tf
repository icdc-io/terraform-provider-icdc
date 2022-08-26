terraform {
  required_providers {
    icdc = {
      version = "0.1"
      source = "local.com/ahrechushkin/icdc"
    }
  }
}

provider "icdc" {
  username = "ahrechushkin@ibagroup.eu"
  password = ""
  location = "ycz"  // ycz, xby, zby, sbg
  group    = "icdc.admin" // account.role
  platform = "icdc" // or "scdc"
}

resource "icdc_service" "composite-resource" {
  name = "tf-resource"
  service_template_id = "18000000000035"
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDU2ixuBrQNH/XeWowxh4CeAjvlkT0Pnz6+GWu1UyJR7/N5TC2DF4cFp4EKWdSxkxKmVg8DUBgsUNJofpDfDJLcwP+kKpYEEiMT4VL4FnPZzDD0hbUbfzaBQCUNtJRHLT91qkysOgm08jaFUlWTI6JhaybVowmpiD0nv1UQW98SKzrVYMXxDv1PSAvESJG8YyQ0zf/RslwaHyiyiqm5uLHoXHEO77ddNkRB5e3meQKiIwEr1f0BjUVgh+kINSlOQLl3euDHaniBAbt6qPOtFHSYXs993rqK3TRN180nigfdSGoJc6FrWF7MiuFC4lUmnk2MzFdGM0TWU/1eniQ0WxfE/lUMI4bIa813+z43cllOvQQitxIgVFRWtJsKm6Lbnw20ioT34rrKKWxHYCI5JvrA7vx39IsgrbFsU952BXOTLVvMPUVGyQYTwIRkmPlJ2GHyicDTBUYv7FGFjVz7gw7ZCIH5HNWSn+57rUdJVzZV+eUM8mrfkPnDOQniRxbnnkk= ahrechushkin@hrechushkin-av"
  vms {
    cpu = "2"
    memory_mb = "2048"
    storage_type = "nvme"
    storage_gb = "30"
    network = "ycz_icdc_base"
  }
}