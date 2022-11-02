provider "ncloud" {
  support_vpc = true
  region      = "KR"
  site        = "public"
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.234.0.0/16"
  name            = "from-tf-vpc"
}

resource "ncloud_subnet" "public-subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.234.0.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  name           = "from-tf-public"
  usage_type     = "GEN"
}

resource "ncloud_subnet" "private-subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.234.100.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "from-tf-private"
  usage_type     = "GEN"
}

data "ncloud_cdss_kafka_version" "kafka-version" {
  filter {
    name   = "name"
    values = ["Kafka 2.8.2"]
  }
}

resource "ncloud_cdss_config_group" "config-group" {
  name = "from-tf-config"
  kafka_version_code = data.ncloud_cdss_kafka_version.kafka-version.id
  description = "test"
}

data "ncloud_cdss_os_product" "os_sample" {
  filter {
    name = "product_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}

data "ncloud_cdss_node_product" "node_sample" {
  os_product_code = data.ncloud_cdss_os_product.os_sample.id
  subnet_no       = ncloud_subnet.public-subnet.id

  filter {
    name   = "cpu_count"
    values = ["2"]
  }

  filter {
    name   = "memory_size"
    values = ["8GB"]
  }

  filter {
    name   = "product_type"
    values = ["STAND"]
  }
}

resource "ncloud_cdss_cluster" "cluster-1" {
  name = "from-tf-cdss"
  kafka_version_code = data.ncloud_cdss_kafka_version.kafka-version.id
  config_group_no = ncloud_cdss_config_group.config-group.id
  vpc_no = ncloud_vpc.vpc.id
  os_product_code = data.ncloud_cdss_os_product.os_sample.id

  cmak {
    user_name = "terraform"
    user_password = "test123!@#"
  }

  manager_node {
    node_product_code = data.ncloud_cdss_node_product.node_sample.id
    subnet_no = ncloud_subnet.public-subnet.id
  }

  broker_nodes {
    node_product_code = data.ncloud_cdss_node_product.node_sample.id
    subnet_no = ncloud_subnet.private-subnet.id
    node_count = 3
    storage_size = 100
  }
}

data "ncloud_cdss_cluster" "cluster" {
  id = "2683820"
}

resource "ncloud_cdss_cluster" "cluster-2" {
  name               = "asdfasdf"
  kafka_version_code = data.ncloud_cdss_cluster.cluster.kafka_version_code
  os_product_code    = data.ncloud_cdss_cluster.cluster.os_product_code
  vpc_no             = data.ncloud_cdss_cluster.cluster.vpc_no
  config_group_no    = data.ncloud_cdss_cluster.cluster.config_group_no

  cmak {
    user_name     = [for k in data.ncloud_cdss_cluster.cluster.cmak : k][0]["user_name"]
    user_password = "test123!@#"
  }

  manager_node {
    node_product_code = [for k in data.ncloud_cdss_cluster.cluster.manager_node : k][0]["node_product_code"]
    subnet_no         = [for k in data.ncloud_cdss_cluster.cluster.manager_node : k][0]["subnet_no"]
  }

  broker_nodes {
    node_product_code = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["node_product_code"]
    node_count        = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["node_count"]
    subnet_no         = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["subnet_no"]
    storage_size      = [for k in data.ncloud_cdss_cluster.cluster.broker_nodes : k][0]["storage_size"]
  }
}