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
  * `id` - Cluster Id.
  * `uuid` - Cluster uuid(Same as Cluster Id)
  * `service_group_instance_no` - Cluster uuid.
  * `cluster_name` - Cluster name.