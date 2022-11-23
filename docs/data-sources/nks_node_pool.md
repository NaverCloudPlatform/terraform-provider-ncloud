# Data Source: ncloud_nks_node_pool

Provides a Kubernetes Service nodepool data.

## Example Usage

### Basic Usage

```hcl
variable "cluster_uuid" {}
variable "node_pool_name" {}

data "ncloud_nks_node_pool" "node_pool"{
  node_pool_name = var.node_pool_name
  cluster_uuid   = var.cluster_uuid
}
```

## Argument Reference

The following arguments are supported:

* `node_pool_name` - (Required) The name of nodepool.
* `cluster_uuid` - (Required) Cluster uuid.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of nodepool.`CusterUuid:NodePoolName`
* `node_count` - Number of nodes.
* `product_code` - Product code.
* `autoscale`
  * `enable` - Auto scaling availability.
  * `max` - Maximum number of nodes available for auto scaling.
  * `min` - Minimum number of nodes available for auto scaling.
* `subnet_no` - Subnet No.
* `instance_no` - Nodepool instance No.
* `nodes`- Running nodes in nodepool.
  * `name` - The name of Server instance.
  * `instance_no` - The ID of server instance.
  * `spec` - Server spec.
  * `private_ip` - Private IP.
  * `public_ip` - Public IP.
  * `node_status` - Node Status.
  * `container_version` - Container version of node.
  * `kernel_version` - kernel version of node.