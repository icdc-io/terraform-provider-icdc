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

resource "icdc_service" "web-server" {
  name = "tf-web-server"
  memory_mb = 2048
  cpu_cores = 1
  network = "ycz_icdc_base"
  storage_type = "nvme"
  storage_mb = 240
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDU2ixuBrQNH/XeWowxh4CeAjvlkT0Pnz6+GWu1UyJR7/N5TC2DF4cFp4EKWdSxkxKmVg8DUBgsUNJofpDfDJLcwP+kKpYEEiMT4VL4FnPZzDD0hbUbfzaBQCUNtJRHLT91qkysOgm08jaFUlWTI6JhaybVowmpiD0nv1UQW98SKzrVYMXxDv1PSAvESJG8YyQ0zf/RslwaHyiyiqm5uLHoXHEO77ddNkRB5e3meQKiIwEr1f0BjUVgh+kINSlOQLl3euDHaniBAbt6qPOtFHSYXs993rqK3TRN180nigfdSGoJc6FrWF7MiuFC4lUmnk2MzFdGM0TWU/1eniQ0WxfE/lUMI4bIa813+z43cllOvQQitxIgVFRWtJsKm6Lbnw20ioT34rrKKWxHYCI5JvrA7vx39IsgrbFsU952BXOTLVvMPUVGyQYTwIRkmPlJ2GHyicDTBUYv7FGFjVz7gw7ZCIH5HNWSn+57rUdJVzZV+eUM8mrfkPnDOQniRxbnnkk= ahrechushkin@hrechushkin-av"
  service_template_id = "18000000000034"
}

resource "icdc_service" "api" {
  name = "tf-api"
  memory_mb = 4096
  cpu_cores = 1
  network = "ycz_icdc_test111"
  storage_type = "nvme"
  storage_mb = 45
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDU2ixuBrQNH/XeWowxh4CeAjvlkT0Pnz6+GWu1UyJR7/N5TC2DF4cFp4EKWdSxkxKmVg8DUBgsUNJofpDfDJLcwP+kKpYEEiMT4VL4FnPZzDD0hbUbfzaBQCUNtJRHLT91qkysOgm08jaFUlWTI6JhaybVowmpiD0nv1UQW98SKzrVYMXxDv1PSAvESJG8YyQ0zf/RslwaHyiyiqm5uLHoXHEO77ddNkRB5e3meQKiIwEr1f0BjUVgh+kINSlOQLl3euDHaniBAbt6qPOtFHSYXs993rqK3TRN180nigfdSGoJc6FrWF7MiuFC4lUmnk2MzFdGM0TWU/1eniQ0WxfE/lUMI4bIa813+z43cllOvQQitxIgVFRWtJsKm6Lbnw20ioT34rrKKWxHYCI5JvrA7vx39IsgrbFsU952BXOTLVvMPUVGyQYTwIRkmPlJ2GHyicDTBUYv7FGFjVz7gw7ZCIH5HNWSn+57rUdJVzZV+eUM8mrfkPnDOQniRxbnnkk= ahrechushkin@hrechushkin-av"
  service_template_id = "18000000000039"
}

resource "icdc_service" "ui" {
  name = "tf-ui"
  memory_mb = 2048
  cpu_cores = 2
  network = "ycz_icdc_test"
  storage_type = "nvme"
  storage_mb = 30
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDU2ixuBrQNH/XeWowxh4CeAjvlkT0Pnz6+GWu1UyJR7/N5TC2DF4cFp4EKWdSxkxKmVg8DUBgsUNJofpDfDJLcwP+kKpYEEiMT4VL4FnPZzDD0hbUbfzaBQCUNtJRHLT91qkysOgm08jaFUlWTI6JhaybVowmpiD0nv1UQW98SKzrVYMXxDv1PSAvESJG8YyQ0zf/RslwaHyiyiqm5uLHoXHEO77ddNkRB5e3meQKiIwEr1f0BjUVgh+kINSlOQLl3euDHaniBAbt6qPOtFHSYXs993rqK3TRN180nigfdSGoJc6FrWF7MiuFC4lUmnk2MzFdGM0TWU/1eniQ0WxfE/lUMI4bIa813+z43cllOvQQitxIgVFRWtJsKm6Lbnw20ioT34rrKKWxHYCI5JvrA7vx39IsgrbFsU952BXOTLVvMPUVGyQYTwIRkmPlJ2GHyicDTBUYv7FGFjVz7gw7ZCIH5HNWSn+57rUdJVzZV+eUM8mrfkPnDOQniRxbnnkk= ahrechushkin@hrechushkin-av"
  service_template_id = "18000000000026"
}