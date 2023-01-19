provider "ncloud" {
  support_vpc = true
  region      = "KR"
  site        = "public"
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.236.0.0/16"
  name            = "tf-vpc-1"
}

resource "ncloud_subnet" "public-subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.236.0.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  name           = "tf-public-1"
  usage_type     = "GEN"
}

resource "ncloud_subnet" "private-subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.236.100.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "tf-private-1"
  usage_type     = "GEN"
}

data "ncloud_cdss_kafka_version" "kafka_version_sample" {
  filter {
    name   = "name"
    values = ["Kafka 2.8.2"]
  }
}

data "ncloud_cdss_kafka_versions" "kafka_versions_sample" {}

data "ncloud_cdss_os_image" "os_image_sample" {
  filter {
    name   = "image_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}

data "ncloud_cdss_os_images" "os_images_sample" {}

data "ncloud_cdss_node_product" "node_sample" {
  os_image  = data.ncloud_cdss_os_image.os_image_sample.id
  subnet_no = ncloud_subnet.public-subnet.id

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

data "ncloud_cdss_node_products" "nodes_sample" {
  os_image  = data.ncloud_cdss_os_image.os_image_sample.id
  subnet_no = ncloud_subnet.public-subnet.id
}

data "ncloud_cdss_config_group" "config_sample" {
  kafka_version_code = data.ncloud_cdss_kafka_version.kafka_version_sample.id

  filter {
    name   = "name"
    values = ["YOUR_CONFIG_GROUP_NAME"]
  }
}

data "ncloud_cdss_cluster" "cluster_sample" {
  filter {
    name   = "name"
    values = ["YOUR_CLUSTER_NAME"]
  }
}

resource "ncloud_cdss_config_group" "config-group" {
  name               = "tf-config-3"
  kafka_version_code = data.ncloud_cdss_kafka_version.kafka_version_sample.id
  description        = "test"
}

resource "ncloud_cdss_cluster" "cluster-12" {
  name               = "from-tf-cdss"
  kafka_version_code = data.ncloud_cdss_kafka_version.kafka_version_sample.id
  config_group_no    = ncloud_cdss_config_group.config-group.id
  vpc_no             = ncloud_vpc.vpc.id
  os_image           = data.ncloud_cdss_os_image.os_image_sample.id

  cmak {
    user_name     = "terraform"
    user_password = var.cmak_user_password
  }

  manager_node {
    node_product_code = data.ncloud_cdss_node_product.node_sample.id
    subnet_no         = ncloud_subnet.public-subnet.id
  }

  broker_nodes {
    node_product_code = data.ncloud_cdss_node_product.node_sample.id
    subnet_no         = ncloud_subnet.private-subnet.id
    node_count        = 3
    storage_size      = 100
  }
}
