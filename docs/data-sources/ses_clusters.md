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
* `clusters` - A List of Search Engine Service cluster.

### Search Engine Service Cluster Reference
`clusters` are also exported with the following attributes, when there are relevant: Each element supports the following:

  * `id` - Cluster Instance No.
  * `service_group_instance_no` - Cluster Instance No(Same as Cluster Id).
  * `cluster_name` - Cluster name.