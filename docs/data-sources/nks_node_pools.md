# Data Source: ncloud_nks_clusters

Retrieve the NKS Node Pools associated with a named NKS cluster. This will allow you to pass a list of Node Pool names to other resources.

## Example Usage

```hcl
var cluster_uuid {}

data "ncloud_nks_node_pools" "node_pools"{
  cluster_uuid = var.cluster_uuid
}

data "ncloud_nks_node_pool" "example"{
  for_each = data.ncloud_nks_node_pools.node_pools.node_pool_names

  cluster_uuid    = var.cluster_uuid
  node_pool_name = each.value
}

```
## Argument Reference

* `cluster_uuid` - (Required) Cluster uuid.

## Attributes Reference

* `id` - Cluster uuid.
* `node_pool_names` - Set of all node pool names in NKS Clusters.


