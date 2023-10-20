variable "access_key" { # export TF_VAR_access_key=...
  description = "access_key, provide through environment variables"
}

variable "secret_key" { # export TF_VAR_secret_key=...
  description = "secret_key, provide through environment variables"
}

variable "region" {
  default = "KR"
}

variable "vpc_name" {
  default = "vpc_name"
}

variable "subnet_name" {
  default = "subnet_name"
}

variable "user_name" {
  default = "user_name"
}

variable "password" {
  description = "password, provide through environment variables"
}

variable "service_name" {
  default = "service-name"
}

variable "name_prefix" {
  default = "name-prefix"
}

variable "database_name" {
  default = "database_name"
}