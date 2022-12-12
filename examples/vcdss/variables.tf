variable name {
  default = "tf-cdss"
}

variable access_key {
  default = "YOUR_ACCESS_KEY"
}

variable secret_key {
  default = "YOUR_SECRET_KEY"
}

variable login_key {
  default = "YOUR_LOGIN_KEY"
}

variable "cmak_user_password" {
  description = "CDSS cluster CMAK user password"
  type        = string
  sensitive   = true
}