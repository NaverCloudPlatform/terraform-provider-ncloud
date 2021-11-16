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
* `node_count` - Number of worker nodes in nodepool.
* `product_code` - Product code of worker nodes in node pool
* `autoscale`- Autoscale config.
  * `enable` - Autoscale enable YN.
  * `max` - Max node count.
  * `min` - Min node count.
* `subnet_no_list` - The ID list of the Subnet where you want to place the nodepool.
* `instance_no` - Instance number of nodepool.
* `subnet_name_list` - The name list of the Subnet where you want to place the nodepool.
* `status` - Nodepool status.
