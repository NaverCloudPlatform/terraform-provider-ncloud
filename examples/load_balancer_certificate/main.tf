provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "ncloud_load_balancer_ssl_certificate" "cert" {
  "certificate_name"      = "tftest_ssl_cert"
  "privatekey"            = "${file("lbtest.privateKey")}"
  "publickey_certificate" = "${file("lbtest.crt")}"
}

resource "ncloud_load_balancer" "lb" {
  "load_balancer_name"                = "tftest_lb"
  "load_balancer_algorithm_type_code" = "SIPHS"
  "load_balancer_description"         = "tftest_lb description"

  "load_balancer_rule_list" = [
    {
      "protocol_type_code"   = "HTTP"
      "load_balancer_port"   = 80
      "server_port"          = 80
      "l7_health_check_path" = "/monitor/l7check"
    },
    {
      "protocol_type_code"   = "HTTPS"
      "load_balancer_port"   = 443
      "server_port"          = 443
      "l7_health_check_path" = "/monitor/l7check"
      "certificate_name"     = "${ncloud_load_balancer_ssl_certificate.cert.certificate_name}"
    },
  ]

  "region_no" = "1"
}
