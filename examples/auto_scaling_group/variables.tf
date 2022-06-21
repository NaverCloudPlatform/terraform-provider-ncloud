variable "access_key" { # export TF_VAR_access_key=...
}

variable "secret_key" { # export TF_VAR_secret_key=...
}

variable "region" {
  default = "KR"
}

variable "acg_name" {
  default = "ncloud-default-acg"
}