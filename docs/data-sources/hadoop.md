---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop

This module can be useful for getting detail of Hadoop instance created before.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take Hadoop instance ID and obtain the data.

```terraform
data "ncloud_hadoop" "by_id" {
  id = 1234567
}

data "ncloud_hadoop" "by_name" {
  cluster_name = "example"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) Hadoop instance number. Either `id` or `cluster_name` must be provided.
* `cluster_name` - (Required) Hadoop service name. Either `id` or `cluster_name` must be provided.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `region_code` - Region code.
* `vpc_no` - The ID of the associated VPC.
* `edge_node_subnet_no` - The subnet ID of the associated edge node.
* `master_node_subnet_no` - The subnet ID of the associated master node.
* `worker_node_subnet_no` - The subnet ID of the associated worker node.
* `master_node_data_storage_type` - Data storage type of master node. There are 2-Options(`SSD`, `HDD`).
* `worker_node_data_storage_type` - Data storage type of master node. There are 2-Options(`SSD`, `HDD`).
* `master_node_data_storage_size` - Data Storage size of master node. Must be between 100(GBi) and 2000(GBi). 4000(GBi) and 6000(GBi) also available.
* `worker_node_data_storage_size` - Data Storage size of master node. Must be between 100(GBi) and 2000(GBi). 4000(GBi) and 6000(GBi) also available.
* `image_product_code` - The image product code of the Hadoop instance.
* `edge_node_product_code` - Edge server product code.
* `master_node_product_code` - Master server product code.
* `worker_node_product_code` - Worker server product code.
* `worker_node_count` - Number of worker server.
* `cluster_type_code` - The cluster type code.
* `version` - The version of Hadoop.
* `ambari_server_host` - The name of ambari host.
* `cluster_direct_access_account` - Account name with direct access to the cluster.
* `login_key` - The login key name.
* `bucket_name` - The name of object storage bucket.
* `use_kdc` - Whether to use Kerberos authentication configuration.
* `kdc_realm` - Realm information of kerberos authentication.
* `domain` - Domain.
* `is_ha` - Whether using high availability of the specific Hadoop.
* `add_on_code_list` - The list of Hadoop Add-On.
* `access_control_group_no_list` - The list of access control group number.
* `hadoop_server_list` The list of Hadoop server instance.
  * `server_instance_no` - Server instance number.
  * `server_name` - Name of the server.
  * `server_role` - Server role code. ex) M(Master), H(Standby Master)
  * `zone_code` - Zone code.
  * `subnet_no` - The ID of the associated Subnet.
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `cpu_count` - the number of the virtual CPU.
  * `memory_size` - Memory size.
  * `data_storage_type` - The type of data storage.
  * `data_storage_size` - Data storage size.
  * `uptime` - Running start time.
  * `create_date` - Server create date.
