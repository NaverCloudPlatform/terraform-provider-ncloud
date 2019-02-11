variable "access_key" {} # export TF_VAR_access_key=...
variable "secret_key" {} # export TF_VAR_secret_key=...

variable "region" {
  default = "KR"
}

variable "zone" {
  default = "KR-2"
}

variable "acg_name" {
  default = "winrm-acg"
}

variable "login_key_name" {
  default = "tf-bswin-key4"
}

variable "server_name" {
  default = "tf-bswin-vm4"
}

variable "block_storage_name" {
  default = "tf-bswin-stor4"
}

variable "port_forwarding_external_port" {
  default = "4389"
}
