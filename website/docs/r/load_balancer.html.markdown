---
layout: "ncloud"
page_title: "NCLOUD: ncloud_load_balancer"
sidebar_current: "docs-ncloud-resource-load-balancer"
description: |-
  Provides a ncloud load balancer instance resource.
---

# ncloud_load_balancer
Provides a ncloud load balancer instance resource.

## Example Usage

```hcl
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
      "certificate_name"     = "cert"
    },
  ]

  "server_instance_no_list" = ["812345", "812346"]
  "internet_line_type_code" = "PUBLC"
  "network_usage_type_code" = "PBLIP"

  "region_no"               = "1"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` - (Optional) Name of a load balancer instance. Default: Automatically specified by Ncloud.
* `load_balancer_algorithm_type_code` - (Optional) Load balancer algorithm type code. The available algorithms are as follows: [ROUND ROBIN (RR) | LEAST_CONNECTION (LC)]. Default: ROUND ROBIN (RR)
* `load_balancer_description` - (Optional) Description of a load balancer instance.
* `load_balancer_rule_list` - (Required) Load balancer rules.
  * `protocol_type_code` - (Required) Protocol type code of load balancer rules. The following codes are available. [HTTP | HTTPS | TCP | SSL]
  * `load_balancer_port` - (Required) Load balancer port of load balancer rules
  * `server_port` - (Required) Server port of load balancer rules
  * `l7_health_check_path` - Health check path of load balancer rules. Required when the `protocol_type_code` is HTTP/HTTPS.
  * `certificate_name` - Load balancer SSL certificate name. Required when the `protocol_type_code` value is SSL/HTTPS.
  * `proxy_protocol_use_yn` - (Optional) Use 'Y' if you want to check client IP addresses by enabling the proxy protocol while you select TCP or SSL.
* `server_instance_no_list` - (Optional) List of server instance numbers to be bound to the load balancer
* `internet_line_type_code` - (Optional) Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)
* `network_usage_type_code` - (Optional) Network usage identification code. PBLIP(PublicIP), PRVT(PrivateIP). default : PBLIP(PublicIP)
* `region_no` - (Optional) Region number. Get available values using the `data ncloud_regions`. Default: KR region
* `zone_no` - (Optional) Zone number. Zone in which you want to create a NAS volume. Get available values using the `data ncloud_zones`.
    Default: The first zone of the region.

## Attributes Reference

* `load_balancer_instance_no` - Load balancer instance No
* `virtual_ip` - Virtual IP address
* `load_balancer_algorithm_type` - Load balancer algorithm type
    * `code` - Load balancer algorithm type code
    * `code_name` - Load balancer algorithm type code name
* `create_date` - Creation date of the load balancer instance
* `domain_name` - Domain name
* `internet_line_type` - Internet line identification type
    * `code` - Internet line identification type code
    * `code_name` - Internet line identification type code name
* `load_balancer_instance_status_name` - Load balancer instance status name
* `load_balancer_instance_status` - Load balancer instance status
    * `code` - Load balancer instance status code
    * `code_name` - Load balancer instance status code name
* `load_balancer_instance_operation` - Load balancer instance operation
    * `code` - Load balancer instance operation code
    * `code_name` - Load balancer instance operation code name
* `network_usage_type` - Network usage type
    * `code` - Network usage type code
    * `code_name` - Network usage type code name
* `is_http_keep_alive` - Http keep alive value [true | false]
* `connection_timeout` - Connection timeout
* `load_balanced_server_instance_list` - Load balanced server instance list
