variable "access_key" { # export TF_VAR_access_key=...
  description = "access_key, provide through environment variables"
}

variable "secret_key" { # export TF_VAR_secret_key=...
  description = "secret_key, provide through environment variables"
}

variable "region" {
  default = "KR"
}

variable "site" {
  description = "Ncloud site. By default, the value is public. You can specify only the following value: public, gov, fin. public is for www.ncloud.com. gov is for www.gov-ncloud.    com. fin is for www.fin-ncloud.com."
  default     = "public"
}

variable "user_name" {
  default = "user_name"
}

variable "user_password" {
  description = "password, provide through environment variables"
}

variable "database_name" {
  default = "database_name"
}