variable "access_key" { # export TF_VAR_access_key=...
  description = "access_key, provide through environment variables"
}

variable "secret_key" { # export TF_VAR_secret_key=...
  description = "secret_key, provide through environment variables"
}

variable "region" {
  default = "KR"
}

variable "addon_image_product_code" {
  default = "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
}

variable "addon_cluster_type_code" {
  default = "CORE_HADOOP_WITH_SPARK"
}