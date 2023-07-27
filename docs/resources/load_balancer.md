---
subcategory: "Load Balancer"
---


# Resource: ncloud_load_balancer
Provides a ncloud load balancer instance resource.

## Example Usage

```hcl
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
    certificate_name     = "cert"
  }
  
  server_instance_no_list = ["812345", "812346"]
  network_usage_type      = "PBLIP"

  region = "KR"
}
```

## Argument Reference

The following arguments are supported:

* `rule_list` - (Required) Load balancer rules.
  * `protocol_type` - (Required) Protocol type code of load balancer rules. The following codes are available. [HTTP | HTTPS | TCP | SSL]
  * `load_balancer_port` - (Required) Load balancer port of load balancer rules
  * `server_port` - (Required) Server port of load balancer rules
  * `l7_health_check_path` - Health check path of load balancer rules. Required when the `protocol_type` is HTTP/HTTPS.
  * `certificate_name` - Load balancer SSL certificate name. Required when the `protocol_type` value is SSL/HTTPS.
  * `proxy_protocol_use_yn` - (Optional) Use 'Y' if you want to check client IP addresses by enabling the proxy protocol while you select TCP or SSL.
* `name` - (Optional) Name of a load balancer instance. Default: Automatically specified by Ncloud.
* `algorithm_type` - (Optional) Load balancer algorithm type code. The available algorithms are as follows: [ROUND ROBIN (RR) | LEAST_CONNECTION (LC)]. Default: ROUND ROBIN (RR)
* `description` - (Optional) Description of a load balancer instance.
* `server_instance_no_list` - (Optional) List of server instance numbers to be bound to the load balancer
* `network_usage_type` - (Optional) Network usage identification code. PBLIP(PublicIP), PRVT(PrivateIP). default : PBLIP(PublicIP)
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `zone` - (Optional) Zone code. Zone in which you want to create a NAS volume. Default: The first zone of the region.
    Get available values using the data source `ncloud_zones`.

## Attributes Reference

* `instance_no` - Load balancer instance No
* `virtual_ip` - Virtual IP address
* `create_date` - Creation date of the load balancer instance
* `domain_name` - Domain name
* `instance_status_name` - Load balancer instance status name
* `instance_status` - Load balancer instance status code
* `instance_operation` - Load balancer instance operation code
* `is_http_keep_alive` - Http keep alive value [true | false]
* `connection_timeout` - Connection timeout
* `load_balanced_server_instance_list` - Load balanced server instance list
