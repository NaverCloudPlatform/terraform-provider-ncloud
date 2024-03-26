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

variable "edge_subnet_name" {
  default = "edge_subnet_name"
}

variable "master_subnet_name" {
  default = "master_subnet_name"
}

variable "worker_subnet_name" {
  default = "worker_subnet_name"
}

variable "hadoop_cluster_name" {
  default = "subnet_name"
}

variable "admin_user_name" {
  default = "admin_user_name"
}

variable "admin_user_password" {
  default = "admin_user_password"
}

variable "bucket_name" {
  default = "bucket_name"
}
