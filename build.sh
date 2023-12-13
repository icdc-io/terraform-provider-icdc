#/bin/sh
cd terraform-provider-icdc && make && cd .. && rm -f .terraform.lock.hcl && terraform init
