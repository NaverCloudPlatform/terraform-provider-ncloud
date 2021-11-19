# Data Source: ncloud_nks_node_pool

Provides a Kubernetes Service nodepool data.

## Example Usage

### Basic Usage

```hcl
variable "cluster_name" {}
variable "node_pool_name" {}

data "ncloud_nks_node_pool" "node_pool"{
  node_pool_name = var.node_pool_name
  cluster_name = var.cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `node_pool_name` - (Required) The name of nodepool.
* `cluster_name` - (Required) The name of Cluster.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of nodepool.`CusterName:NodePoolName`
* `node_count` - Number of nodes.
* `product_code` - Product code.
* `autoscale`
  * `enable` - Auto scaling availability.
  * `max` - Maximum number of nodes available for auto scaling.
  * `min` - Minimum number of nodes available for auto scaling.
* `subnet_no_list` - Subnet No. list.
* `instance_no` - Instance No.
* `subnet_name_list` - Subnet name list.
* `status` - Nodepool status.
