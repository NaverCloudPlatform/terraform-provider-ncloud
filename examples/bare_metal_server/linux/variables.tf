variable "access_key" {} # export TF_VAR_access_key=...
variable "secret_key" {} # export TF_VAR_secret_key=...

variable "region" {
  default = "KR"
}

variable "login_key_name" {
  default = "tf-bmtest-key"
}

variable "server_name" {
  default = "tf-bmtest-vm"
}

variable "block_storage_name" {
  default = "tf-bmtest-storage"
}

variable "nas_volume_name_prefix" {
  default = "_tfvol"
}

variable "load_balancer_name" {
  default =  "tf-bmtest_lb"
}

variable "certificate_name" {
  default =  "tf-bmtest_ssl_cert"
}

variable "port_forwarding_external_port" {
  default = "5022"
}