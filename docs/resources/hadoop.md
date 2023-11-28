---
subcategory: "Hadoop"
---


# Resource: ncloud_hadoop

Provides a Hadoop instance resource.

## Example Usage

#### Basic (Vpc)
```hcl
resource "ncloud_login_key" "loginKey" {
  key_name = "hadoop-key"
}

resource "ncloud_vpc" "vpc" {
  name = "hadoop-vpc"
  ipv4_cidr_block = "10.5.0.0/16"
}

resource "ncloud_subnet" "master_node_subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = "master-node-subnet"
  subnet             = "10.5.64.0/19"
  zone               = "KR-2"
  subnet_type        = "PUBLIC"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
}

resource "ncloud_subnet" "edge_node_subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = "edge-node-subnet"
  subnet             = "10.5.0.0/18"
  zone               = "KR-2"
  subnet_type        = "PUBLIC"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
}

resource "ncloud_subnet" "worker_node_subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "worker-node-subnet"
	subnet             = "10.5.96.0/20"
	zone               = "KR-2"
	subnet_type        = "PRIVATE"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
}

resource "ncloud_hadoop" "hadoop" {
  vpc_no = "49956"
  cluster_name = "hadoopName"
  cluster_type_code = "CORE_HADOOP_WITH_SPARK"
  admin_user_name = "admin-test"
  admin_user_password = "Admin!2Admin"
  login_key_name = loginKey.key_name
  master_node_subnet_no = ncloud_subnet.master_node_subnet.id
  edge_node_subnet_no = ncloud_subnet.edge_node_subnet.id
  worker_node_subnet_no = ncloud_subnet.worker_node_subnet.id
  bucket_name = "bucket_name"
  master_node_data_storage_type = "SSD"
  worker_node_data_storage_type = "SSD"
  master_node_data_storage_size = 100
  worker_node_data_storage_size = 100
}
```

## Argument Reference

The following arguments are supported:

* `image_product_code` - 
* 
* `vpc_no` - (Required) The ID of the VPC where you want to place the Hadoop Instance.
* TODO: 아래 프로덕트 코드들 -> 문서보고 더 추가할것 
* `master_node_product_code` - (Optional) Master node product code to determin the master node server specification to create. Default: Selected as minimum specification. The minimum standards are 1. memory 2. CPU
* `edge_node_product_code` - (Optional) Edge node product code to determin the edge node server specification to create. Default: Selected as minimum specification. The minimum standards are 1. memory 2. CPU
* `worker_node_product_code` - (Optional) Worker node product code to determin the worker node server specification to create. Default: Selected as minimum specification. The minimum standards are 1. memory 2. CPU
* `cluster_name` - (Required) Cluster name to create.
* `cluster_type_code` - (Required) Cluster type code to determin the cluster type to create.
* `add_on_code_list` - (Optional) Hadoop add-on list. There are 4-options(`PRESTO`, `HBASE`, `IMPALA` and `KUDU`). This argument can only be used in Cloud Hadoop version 1.5 or higher.
* `admin_user_name` - (Required) Admin user name of cluster to create. It is the administrator account required to access the Ambari management console.
* `admin_user_password` - (Required) Admin user password of cluster to create.
* `login_key_name` - (Required) Login key name to set the SSH authentication key required when connecting directly to the node.
* `master_node_subnet_no` - (Required) The Subnet ID of master node. 
* `edge_node_subnet_no` - (Required) The Subnet ID of edge node.
* `worker_node_subnet_no` - (Required) The Subnet ID of worker node. Must be located in Private Subnet.
* `bucket_name` - (Required) Bucket name to space for storing data in Object Storage.
* `master_node_data_storage_type` - (Required) Data storage type of master node. It does not change atfer installation. There are 2-Options(`SSD`, `HDD`). Default: SSD.
* `worker_node_data_storage_type` - (Required) Data Storage type of worker node. It does not change atfer installation. There are 2-Options(`SSD`, `HDD`). Default: SSD.
* `master_node_data_storage_size` - (Required) Data Storage size of master node. Must be between 100(GB) and 2000(GB) in 10(GB) increaments. 4000(GB) and 6000(GB) also available.
* `worker_node_data_storage_size` - (Required) Data Storage size of worker node. Must be between 100(GB) and 2000(GB) in 10(GB) increaments. 4000(GB) and 6000(GB) also available.
* `worker_node_count` - (Optional) Count of worker node. Must be between 2 and 8. Default: 2
* `use_kdc` - (Optional) Whether to use KDC(Kerberos Distribute Center). Default: false
* `kdc_realm` - (Required if `use_kdc` is provided) Only domain rules of type Realm are allowed.
* `kdc_password` - (Required if `use_kdc` is provided) Password of KDC.
* `use_bootstrap_script` - (Optional) Whether to use bootstrap script. Default: false.
* `bootstrap_script` - (Required if `use_kdc` is provided) Bootstrap script. Script can only be performed with buckets linked to Cloud Hadoop. Requires entering folder and file names excluding bucket name.
* `use_data_catalog` - (Optional) Whether to use data catalog. It is provided by using the Cloud Hadoop Hive Metastore as the catalog for the Data Catalog service. Integration is possible only when the catalog status of the Data Catalog service is normal. Intergration is possible only with Cloud Hadoop version 2.0 or higher.

## Attributes Reference

* `ID` - The ID of hadoop instance.
* `login_key` - The login key of hadoop instance.
* `object_storage_bucket` - The Object storage.
* `ambari_server_host` - Ambari server host.
* `cluster_direct_access_account` - Account to access the cluster directly.
* `is_ha` - Whether is High Availability or not.
* `access_control_group_no_list` - Access control group number list.
* `hadoop_server_instance_list` - Server instance list of the hadoop instance. 
  * `hadoop_server_name` - Name of the hadoop server instance.
  * `hadoop_server_role` - Role of the hadoop server instance.
  * `hadoop_product_code` - Product code of the hadoop server instance.
  * `region_code` - Region code of the hadoop server instance.
  * `zone_code` - Zone code of the hadoop server instance.
  * `vpc_no` - Vpc no of the hadoop server instance.
  * `subnet_no` - Subnet no of the hadoop server instance.
  * `is_public_subnet` - Whether is Public Subnet or Private Subnet for the hadoop server instance.
  * `data_storage_size` - Data storage size of the hadoop server instance.
  * `cpu_count` - Cpu count of the hadoop server instance.
  * `memory_size` - Memory size of the hadoop server instance.
  * 