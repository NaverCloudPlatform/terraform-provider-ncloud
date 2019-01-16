provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "ncloud_server" "server" {
  "name"               = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code"       = "${var.server_product_code}"
}

resource "ncloud_load_balancer" "lb" {
  "name"                = "tftest_lb"
  "algorithm_type_code" = "SIPHS"
  "description"         = "tftest_lb description"

  "rule_list" = [
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
      "certificate_name"     = "aaa"
    },
  ]

  "server_instance_no_list" = ["${ncloud_server.server.id}"]
  "internet_line_type_code" = "PUBLC"
  "network_usage_type_code" = "PBLIP"
  "region_no"               = "1"
}
