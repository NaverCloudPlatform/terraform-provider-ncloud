variable "access_key" { # export TF_VAR_access_key=...
  description = "access_key, provide through environment variables"
}

variable "secret_key" { # export TF_VAR_secret_key=...
  description = "secret_key, provide through environment variables"
}

variable "region" {
  default = "KR"
}