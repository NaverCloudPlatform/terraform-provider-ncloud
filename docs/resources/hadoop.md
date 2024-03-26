---
subcategory: "Hadoop"
---


# Resource: ncloud_hadoop

Provides a Hadoop instance resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_login_key" "login_key" {
  key_name = "hadoop-key"
}

resource "ncloud_vpc" "vpc" {
  name            = "hadoop-vpc"
  ipv4_cidr_block = "10.5.0.0/16"
}

resource "ncloud_subnet" "edge_node_subnet" {
  vpc_no          = ncloud_vpc.vpc.vpc_no
  name            = "edge-node-subnet"
  subnet          = "10.5.0.0/18"
  zone            = "KR-2"
  subnet_type     = "PUBLIC"
  network_acl_no  = ncloud_vpc.vpc.default_network_acl_no
}

resource "ncloud_subnet" "master_node_subnet" {
  vpc_no          = ncloud_vpc.vpc.vpc_no
  name            = "master-node-subnet"
  subnet          = "10.5.64.0/19"
  zone            = "KR-2"
  subnet_type     = "PUBLIC"
  network_acl_no  = ncloud_vpc.vpc.default_network_acl_no
}

resource "ncloud_subnet" "worker_node_subnet" {
  vpc_no          = ncloud_vpc.vpc.vpc_no
  name            = "worker-node-subnet"
  subnet          = "10.5.96.0/20"
  zone            = "KR-2"
  subnet_type     = "PRIVATE"
  network_acl_no  = ncloud_vpc.vpc.default_network_acl_no
}

resource "ncloud_hadoop" "hadoop" {
  vpc_no                        = ncloud_vpc.vpc.vpc_no
  cluster_name                  = "hadoopName"
  cluster_type_code             = "CORE_HADOOP_WITH_SPARK"
  admin_user_name               = "admin-test"
  admin_user_password           = "Admin!2Admin"
  login_key_name                = ncloud_login_key.login_key.key_name
  edge_node_subnet_no           = ncloud_subnet.edge_node_subnet.id
  master_node_subnet_no         = ncloud_subnet.master_node_subnet.id
  worker_node_subnet_no         = ncloud_subnet.worker_node_subnet.id
  bucket_name                   = "bucket_name"
  master_node_data_storage_type = "SSD"
  worker_node_data_storage_type = "SSD"
  master_node_data_storage_size = 100
  worker_node_data_storage_size = 100
}
```


## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the VPC where you want to place the Hadoop Instance.
* `cluster_name` - (Required) Cluster name to create. Can only enter English letters, numbers, and dashes (-), and Korean letters. Must start and end with an English letter (lowercase) or a number. Min: 3, Max: 15
* `cluster_type_code` - (Required) Cluster type code to determin the cluster type to create. Options: CORE_HADOOP_WITH_SPARK 
* `admin_user_name` - (Required) Admin user name of cluster to create. It is the administrator account required to access the Ambari management console. Can only be composed of English letters (lowercase), numbers, and dashes (-).  Must start and end with an English letter (lowercase) or a number.  Min: 3, Max: 15
* `admin_user_password` - (Required) Admin user password of cluster to create. Must include at least 1 alphabetical character (capital letter), special character, and number. Special characters, such as single quotations ('), double quotations ("), the KRW symbol (₩), slashes (/), ampersands (&), back quotes (`), and spaces cannot be included. Min: 8, Max: 20
* `login_key` - (Required) Login key name to set the SSH authentication key required when connecting directly to the node.
* `edge_node_subnet_no` - (Required) The Subnet ID of edge node. Can select a subnet that will locate the edge node. Edge nodes are located in private/public subnets.
* `master_node_subnet_no` - (Required) The Subnet ID of master node. Can select a subnet that will locate the master node.  Master nodes are located in private/public subnets
* `worker_node_subnet_no` - (Required) The Subnet ID of worker node. Must be located in Private Subnet.
* `bucket_name` - (Required) Bucket name to space for storing data in Object Storage.
* `master_node_data_storage_type` - (Required) Data storage type of master node. It does not change atfer installation. There are 2-Options(`SSD`, `HDD`). Default: SSD.
* `worker_node_data_storage_type` - (Required) Data Storage type of worker node. It does not change atfer installation. There are 2-Options(`SSD`, `HDD`). Default: SSD.
* `master_node_data_storage_size` - (Required) Data Storage size of master node. Must be between 100(GBi) and 2000(GBi) in 10(GBi) increaments. 4000(GBi) and 6000(GBi) also available.
* `worker_node_data_storage_size` - (Required) Data Storage size of worker node. Must be between 100(GBi) and 2000(GBi) in 10(GBi) increaments. 4000(GBi) and 6000(GBi) also available.
* `image_product_code` - (Optional) Image product code to determine the Hadoop instance server image specification to create. If not entered, the instance is created for default value. Default: Cloud Hadoop's latest version. It can be obtained through [`ncloud_hadoop_image_products` data source](../data-sources/hadoop_image_products.md)
* `edge_node_product_code` - (Optional, Changeable) Edge node product code to determin the edge node server specification to create. The specification upgrade will be performed after a full service stop, so please stop work in advance. Upgrading to a server with more memory than the current specification is only possible and incurs an additional fee. Default: Selected as minimum specification. The minimum standards are 1. memory 2. CPU. It can be obtained through [`ncloud_hadoop_products` data source](../data-source/hadoop_products.md).
* `master_node_product_code` - (Optional, Changeable) Master node product code to determin the master node server specification to create. The specification upgrade will be performed after a full service stop, so please stop work in advance. Upgrading to a server with more memory than the current specification is only possible and incurs an additional fee. Default: Selected as minimum specification. The minimum standards are 1. memory 2. CPU. It can be obtained through [`ncloud_hadoop_products` data source](../data-sources/hadoop_products.md).
* `worker_node_product_code` - (Optional, Changeable) Worker node product code to determin the worker node server specification to create. The specification upgrade will be performed after a full service stop, so please stop work in advance. Upgrading to a server with more memory than the current specification is only possible and incurs an additional fee. Default: Selected as minimum specification. The minimum standards are 1. memory 2. CPU. It can be obtained through [`ncloud_hadoop_products` data source](../data-sources/hadoop_products.md).
* `add_on_code_list` - (Optional) Hadoop add-on list. This argument can only be used in Cloud Hadoop version 1.5 or higher. Options: PRESTO | HBASE | IMPALA | KUDU | TRINO | NIFI
* `worker_node_count` - (Optional, Changeable) Number of worker server. You can only select between 2 and 8 worker nodes when you first create. The minimum number of worker nodes is 2, and the number of nodes that can be changed at once is 10.
* `use_kdc` - (Optional) Whether to use KDC(Kerberos Distribute Center). Default: false
* `kdc_realm` - (Required if `use_kdc` is provided) KDC's Realm information. Can be entered only if useKdc is true. Only realm-format domain rules are allowed. Only uppercase letters (A-Z) are allowed and up to 15 digits are allowed. Only one dot(.) is allowed (ex. EXAMPLE.COM). 
* `kdc_password` - (Required if `use_kdc` is provided) Password of KDC. Can be entered only if useKdc is true. Must include at least 1 alphabetical character (capital letter), special character, and number. Special characters, such as single quotations ('), double quotations ("), the KRW symbol (₩), slashes (/), ampersands (&), back quotes (`), and spaces cannot be included. Min: 8, Max: 20
* `use_bootstrap_script` - (Optional) Whether to use bootstrap script. Default: false.
* `bootstrap_script` - (Required if `use_kdc` is provided) Bootstrap script. Script can only be performed with buckets linked to Cloud Hadoop. Requires entering folder and file names excluding bucket name. Only English is supported. Cannot use spaces or special characters. Available up to 1024 bytes.
* `use_data_catalog` - (Optional) Whether to use data catalog. Available only `public` site. It is provided by using the Cloud Hadoop Hive Metastore as the catalog for the Data Catalog service. Integration is possible only when the catalog status of the Data Catalog service is normal. Intergration is possible only with Cloud Hadoop version 2.0 or higher. Default: false

## Attributes Reference

In addition to all arguments above, the following attributes are exported

* `id` - The ID of hadoop instance.
* `region_code` - Region code.
* `ambari_server_host` - Ambari server host.
* `cluster_direct_access_account` - Account to access the cluster directly.
* `version` - The version of Hadoop.
* `is_ha` - Whether using high availability of the specific Hadoop.
* `domain` - Domain.
* `access_control_group_no_list` - Access control group number list.
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

## Import

### `terraform import` command

* Hadoop can be imported using the `id`. For example:

```console
$ terraform import ncloud_hadoop.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Hadoop using the `id`. For example:

```terraform
import {
  to = ncloud_hadoop.rsc_name
  id = "12345"
}
```
