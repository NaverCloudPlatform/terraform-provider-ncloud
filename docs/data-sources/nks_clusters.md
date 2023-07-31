---
subcategory: "Kubernetes Service"
---


# Data Source: ncloud_nks_clusters

Provides list of Kubernetes Service cluster uuid.

## Example Usage

```hcl
data "ncloud_nks_clusters" "clusters"{
}

data "ncloud_nks_cluster" "cluster"{
  for_each = toset(data.ncloud_nks_clusters.example.cluster_uuids)
  uuid     = each.value
}

```

## Attributes Reference

* `id` - Ncloud Region.
* `cluster_uuids` - Set of NKS Clusters uuids


