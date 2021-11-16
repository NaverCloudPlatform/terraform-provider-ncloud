# Data Source: ncloud_nks_clusters

Retrieve NKS Clusters list

## Example Usage

```hcl
data "ncloud_nks_clusters" "example"{
}

data "ncloud_nks_cluster" "example"{
  for_each = toset(data.ncloud_nks_clusters.example.cluster_names)
  name     = each.value
}

```

## Attributes Reference

* `id` - Ncloud Region.
* `cluster_names` - Set of NKS Clusters names


