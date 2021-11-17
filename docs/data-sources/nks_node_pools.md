# Data Source: ncloud_nks_clusters

Retrieve the NKS Node Pools associated with a named NKS cluster. This will allow you to pass a list of Node Pool names to other resources.

## Example Usage

```hcl
var cluster_name {}

data "ncloud_nks_node_pools" "node_pools"{
  cluster_name = var.cluster_name
}

data "ncloud_nks_node_pool" "example"{
  for_each = data.ncloud_nks_node_pools.node_pools.node_pool_names

  cluster_name    = var.cluster_name
  node_pool_name = each.value
}

```
## Argument Reference

* `cluster_name` - (Required) The name of the cluster.

## Attributes Reference

* `id` - Cluster name.
* `node_pool_names` - Set of all node pool naems in NKS Clusters.


