---
subcategory: "Cloud Data Streaming Service"
---


# Data Source: ncloud_cdss_cluster

## Example Usage
```hcl
data "ncloud_cdss_cluster" "cluster"{
  filter {
    name   = "name"
    values = ["YOUR_CLUSTER_NAME"]
  }
}
```

## Argument Reference
The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Cluster id.
* `service_group_instance_no` - Service Group Instance number.
* `name` - Cluster name.
* `kafka_version_code` - Cloud Data Streaming Service version to be used.
* `config_group_no` - ConfigGroup number to be used.
* `vpc_no` - VPC number to be used.
* `os_image` -  OS type to be used.
* `cmak` - .
  * `user_name` - CMAK access ID. Only lowercase alphanumeric characters and non-consecutive hyphens (-) allowed First character must be a letter, but the last character may be a letter or a number.
  * `user_password` - CMAK access password. Must be at least 8 characters and contain at least one of each: English uppercase letter, lowercase letter, special character, and number.
* `manager_node` - .
  * `node_product_code` - HW specifications of the manager node.
  * `subnet_no` - Subnet number where the manager node is to be located.
* `broker_nodes` - .
  * `node_product_code` - HW specifications of the broker node.
  * `subnet_no` - Subnet number where the broker node is to be located.
  * `node_count` - Number of broker nodes. At least 3 units, up to 10 units allowed. (Can only be increased)
  * `storage_size` - Broker node storage capacity. At least 100 GB, up to 2000 GB. Must be in units of 10 GB.
* `endpoints` - .
  * `plaintext` - List of broker nodes (Port 9092).
  * `tls` - List of broker nodes (Port 9093).
  * `public_endpoint_plaintext` - List of public endpoint of broker nodes.
  * `public_endpoint_plaintext_listener_port` - List of listener port for public endpoint of broker nodes.
  * `public_endpoint_tls` - List of public endpoint of broker nodes (TLS).
  * `public_endpoint_tls_listener_port` - List of listener port for public endpoint of broker nodes (TLS).
  * `hosts_private_endpoint_tls` - Editing details of the hosts file (Private IP hostname format).
  * `hosts_public_endpoint_tls` - Editing details of the hosts file (Public IP hostname format).
  * `zookeeper` - List of ZooKeeper nodes (Port 2181).