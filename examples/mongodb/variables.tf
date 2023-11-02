variable "access_key" {
  description = "access_key, provide through environment variables"
}

variable "secret_key" {
  description = "secret_key, provide through environment variables"
}

variable "region" {
  default = "KR"
}

variable "user_name" {
  default = "mongodbuser"
  description = "user name"
}

variable "password" {
  description = "password, provide through environment variables"
}