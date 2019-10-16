provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

resource "ncloud_load_balancer_ssl_certificate" "cert" {
  certificate_name      = "tftest_ssl_cert"
  privatekey            = file("lbtest.privateKey")
  publickey_certificate = file("lbtest.crt")
}

resource "ncloud_load_balancer" "lb" {
  name           = "tftest_lb"
  algorithm_type = "SIPHS"
  description    = "tftest_lb description"

  rule_list {
    protocol_type        = "HTTP"
    load_balancer_port   = 80
    server_port          = 80
    l7_health_check_path = "/monitor/l7check"
  }
  rule_list {
    protocol_type        = "HTTPS"
    load_balancer_port   = 443
    server_port          = 443
    l7_health_check_path = "/monitor/l7check"
    certificate_name     = ncloud_load_balancer_ssl_certificate.cert.certificate_name
  }

  region = "KR"
}

