# Data Source: ncloud_cdss_cluster

## Example Usage
``` hcl
variable "cdss_cluster_uuid" {}

data "ncloud_cdss_cluster" "cluster"{
  id = var.cdss_cluster_uuid
}

resource "ncloud_cdss_cluster" "cluster-2" {
  name               = "test-cluster"
  kafka_version_code = data.ncloud_cdss_cluster.cluster.kafka_version_code
  os_product_code    = data.ncloud_cdss_cluster.cluster.os_product_code
  vpc_no             = data.ncloud_cdss_cluster.cluster.vpc_no
  config_group_no    = data.ncloud_cdss_cluster.cluster.config_group_no

  cmak {
    user_name     = [for k in data.ncloud_cdss_cluster.cluster.cmak : k][0]["user_name"]
    user_password = "test123!@#"
  }

  manager_node {
    node_product_code = [for k in data.ncloud_cdss_cluster.cluster.manager_node : k][0]["node_product_code"]
    subnet_no         = [for k in data.ncloud_cdss_cluster.cluster.manager_node : k][0]["subnet_no"]
  }

  broker_nodes {
    node_product_code = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["node_product_code"]
    node_count        = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["node_count"]
    subnet_no         = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["subnet_no"]
    storage_size      = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["storage_size"]
  }
}
```

## Argument Reference
The following arguments are supported

* `id` - (Required) Cluster uuid.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `name` - Cluster name.
* `kafka_version_code` - Cloud Data Streaming Service version to be used.
* `os_product_code` -  OS type to be used.
* `vpc_no` - VPC number to be used.
* `config_group_no` - ConfigGroup number to be used.
* `cmak` - .
    * `user_name` - CMAK access ID. Only lowercase alphanumeric characters and non-consecutive hyphens (-) allowed First character must be a letter, but the last character may be a letter or a number.
    * `user_password` - CMAK access password. Must be at least 8 characters and contain at least one of each: English uppercase letter, lowercase letter, special character, and number.
* `manager_node` - .
    * `node_product_code` - HW specifications of the manager node.
    * `subnet_no` - Subnet number where the manager node is to be located.
* `broker_nodes` - .
    * `node_product_code` - HW specifications of the broker node.
    * `subnet_no` - Subnet number where the broker node is to be located.
    * `node_count` - Number of broker nodes. At least 3 units, up to 10 units allowed.
    * `storage_size` - Broker node storage capacity. At least 100 GB, up to 2000 GB. Must be in units of 10 GB.
* `endpoints` - .
    * `broker_node_list` - List of broker nodes (Port 9092).
    * `broker_tls_node_list` - List of broker nodes (Port 9093).
    * `public_endpoint_broker_node_list` - List of public endpoint of broker nodes.
    * `public_entpoint_broker_node_listener_port_list` - List of listener port for public endpoint of broker nodes.
    * `public_endpoint_broker_tls_node_list` - List of public endpoint of broker nodes (TLS).
    * `public_endpoint_broker_tls_node_listener_port_list` - List of listener port for public endpoint of broker nodes (TLS).
    * `local_dns_list` - Editing details of the hosts file (Private IP hostname format).
    * `local_dns_tls_list` - Editing details of the hosts file (Private IP hostname format) (TLS).
    * `zookeeper_list` - List of ZooKeeper nodes (Port 2181).