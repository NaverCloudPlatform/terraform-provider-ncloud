variable access_key {
  default = "ACCESS_KEY"
}

variable secret_key {
  default = "SECRET_KEY"
}

variable name {
  default = "tf-ses"
}

variable ses_version {
  default = "1.3.3"
}
variable ses_version_type {
  default = "OpenSearch"
}

variable os_version {
  default = "SW.VELST.OS.LNX64.CNTOS.0708.B050"
}

variable ses_product_code {
  default = "SVR.VELST.STAND.C002.M008.NET.SSD.B050.G002"
}

variable ses_produce_cpu_count {
  default = "2"
}

variable login_key {
  default = "tf-login-test"
}

variable ses_user_password {
  description = "SES Cluster User Password"
  type = string
  sensitive =  true
}
