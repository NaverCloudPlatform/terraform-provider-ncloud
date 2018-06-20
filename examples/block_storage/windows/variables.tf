variable "access_key" {} # export TF_VAR_access_key=...
variable "secret_key" {} # export TF_VAR_secret_key=...

variable "region" {
  default = "KR"
}

variable "zone_no" {
  default = "3"
}

variable "acg_name" {
  default = "winrm-acg"
}

variable "login_key_name" {
  default = "tf-bswin-key3"
}

variable "server_name" {
  default = "tf-bswin-vm3"
}

variable "block_storage_name" {
  default = "tf-bswin-sto3"
}

