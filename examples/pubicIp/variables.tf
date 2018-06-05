variable "access_key" {} # export TF_VAR_access_key=...
variable "secret_key" {} # export TF_VAR_secret_key=...

variable "region" {
  default = "KR"
}

variable "server_instance_no" {
  default = "805853"
}

variable "region_no" {
  default = "1"
}

variable "zone_no" {
  default = "3"
}
