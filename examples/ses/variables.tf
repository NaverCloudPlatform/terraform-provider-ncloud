variable name {
  default = "tf-ses"
}

variable ses_version {
  default = "133"
}

variable os_version {
  default = "SW.VELST.OS.LNX64.CNTOS.0708.B050"
}

variable access_key {
  default = "ACCESS_KEY"
}

variable secret_key {
  default = "SECRET_KEY"
}

variable login_key {
  default = "tf-login-test"
}

variable ses_user_password {
  description = "SES Cluster User Password"
  type = string
  sensitive =  true
}
