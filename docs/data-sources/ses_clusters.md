# Data Source: ncloud_ses_clusters

Provides list of Search Engine Service cluster uuid.

## Example Usage
``` hcl
data "ncloud_ses_clusters" "clusters"{
}

data "ncloud_ses_clusters" "cluster"{
  filter {
    name = "cluster_name"
    values = ["my_cluster"]
  }
}
```

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `clusters` - .
  * `id` - Cluster Instance No.
  * `service_group_instance_no` - Cluster Instance No(Same as Cluster Id).
  * `cluster_name` - Cluster name.