variable "access_key" { # export TF_VAR_access_key=...
}

variable "secret_key" { # export TF_VAR_secret_key=...
}

variable "region" {
  default = "KR"
}

resource "ncloud_port_forwarding_rule" "test" {
  server_instance_no            = "966669"
  port_forwarding_external_port = "2088"
  port_forwarding_internal_port = "22"
}

