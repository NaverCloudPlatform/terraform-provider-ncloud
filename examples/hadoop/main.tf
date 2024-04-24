provider "ncloud" {
  support_vpc = true
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

resource "ncloud_login_key" "login_key" {
  key_name = "hadoop-key"
}

resource "ncloud_vpc" "vpc" {
  name               = var.vpc_name
  ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "edge_subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = var.edge_subnet_name
  subnet             = "10.5.0.0/18"
  zone               = "KR-2"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
  subnet_type        = "PUBLIC"
  usage_type         = "GEN"
}

resource "ncloud_subnet" "master_subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = var.master_subnet_name
  subnet             = "10.5.64.0/19"
  zone               = "KR-2"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
  subnet_type        = "PUBLIC"
  usage_type         = "GEN"
}

resource "ncloud_subnet" "worker_subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = var.worker_subnet_name
  subnet             = "10.5.96.0/20"
  zone               = "KR-2"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
  subnet_type        = "PRIVATE"
  usage_type         = "GEN"
}

resource "ncloud_hadoop" "hadoop" {
  vpc_no = ncloud_vpc.vpc.vpc_no
  cluster_name = var.hadoop_cluster_name
  cluster_type_code = "CORE_HADOOP_WITH_SPARK"
  admin_user_name = var.admin_user_name
  admin_user_password = var.admin_user_password
  login_key_name = ncloud_login_key.login_key.key_name
  edge_node_subnet_no = ncloud_subnet.edge_subnet.subnet_no
  master_node_subnet_no = ncloud_subnet.master_subnet.subnet_no
  worker_node_subnet_no = ncloud_subnet.worker_subnet.subnet_no
  bucket_name = var.bucket_name
  master_node_data_storage_type = "SSD"
  worker_node_data_storage_type = "SSD"
  master_node_data_storage_size = 100
  worker_node_data_storage_size = 100
}

data "ncloud_hadoop" "by_cluster_name" {
  cluster_name = ncloud_hadoop.hadoop.cluster_name
}
