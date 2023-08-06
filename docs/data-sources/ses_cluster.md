---
subcategory: "Search Engine Service"
---


# Data Source: ncloud_ses_cluster

Provides a Search Engine Service cluster data.

## Example Usage
``` hcl
variable "ses_cluster_name" {}

data "ncloud_ses_clusters" "clusters"{
  filter {
    name   = "cluster_name"
    values = [var.ses_cluster_name]
  }
}

data "ncloud_ses_cluster" "my_cluster"{
  id = data.ncloud_ses_clusters.clusters.clusters.0.id
}
```

## Argument Reference
The following arguments are supported

* `id` - (Required) Cluster Instance No.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `cluster_name` - Cluster name.
* `service_group_instance_no` - Cluster Instance No.
* `id` - Cluster Instance No.
* `os_image_code` -  OS type to be used.
* `vpc_no` - VPC number to be used.
* `search_engine` - .
    * `version_code` - Search Engine Service version to be used.
    * `user_name` - Search Engine UserName. Only lowercase alphanumeric characters and non-consecutive hyphens (-) allowed First character must be a letter, but the last character may be a letter or a number.
    * `port` - Search Engine Port.
    * `dashboard_port` - Search Engine Dashboard Port.
* `manager_node` - .
    * `is_dual_manager` - Redundancy of manager node
    * `product_code` - HW specifications of the manager node.
    * `subnet_no` - Subnet number where the manager node is to be located.
    * `acg_id` - The ID of manager node ACG.
    * `acg_name` - The name of manager node ACG.
* `data_node` - .
    * `product_code` - HW specifications of the data node.
    * `subnet_no` - Subnet number where the data node is to be located.
    * `node_count` - Number of data nodes. At least 3 units, up to 10 units allowed.
    * `storage_size` - Data node storage capacity.
    * `acg_id` - The ID of data node ACG.
    * `acg_name` - The name of data node ACG.
* `master_node(Optional)` - .
  * `product_code` - HW specifications of the master node.
  * `subnet_no` - Subnet number where the master node is to be located.
  * `node_count` - Number of master nodes.
  * `acg_id` - The ID of master node ACG.
  * `acg_name` - The name of master node ACG.
* `manager_node_instance_no_list` - List of Manager node's instance number 
* `cluster_node_list` - .
  * `compute_instance_name` - The name of Server instance.
  * `compute_instance_no`   - The ID of Server instance.
  * `node_type`             - Node role code
  * `private_ip`            - Private IP
  * `server_status`         - The status of Server Instance.
  * `subnet`                - The name of Server Instance subnet.
* `login_key_name` - Required Login key to access Manager node server