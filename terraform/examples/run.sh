export TCA_USERNAME=administrator@vsphere.local
export TCA_PASSWORD=VMware1!
rm -rf .terraform.lock.hcl
rm -rf .terraform; terraform init && terraform apply --auto-approve


