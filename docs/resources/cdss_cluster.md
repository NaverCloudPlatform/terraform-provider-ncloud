# Resource: ncloud_cdss_cluster

## Example Usage

``` hcl
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

resource "ncloud_cdss_cluster" "cluster" {
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
```

## Argument Reference
The following arguments are supported:

* `name` - Cluster name.
* `kafka_version_code` - Cloud Data Streaming Service version to be used.
* `os_product_code` -  OS type to be used.
* `vpc_no` - VPC number to be used.
* `config_group_no` - ConfigGroup number to be used.
* `cmak` - .
    * `user_name` - CMAK access ID. Only lowercase alphanumeric characters and non-consecutive hyphens (-) allowed First character must be a letter, but the last character may be a letter or a number.
    * `user_password` - CMAK access password. Must be at least 8 characters and contain at least one of each: English uppercase letter, lowercase letter, special character, and number.
* `manager_node` - .
    * `node_product_code` - HW specifications of the manager node.
    * `subnet_no` - Subnet number where the manager node is to be located.
* `broker_nodes` - .
    * `node_product_code` - HW specifications of the broker node.
    * `subnet_no` - Subnet number where the broker node is to be located.
    * `node_count` - Number of broker nodes. At least 3 units, up to 10 units allowed.
    * `storage_size` - Broker node storage capacity. At least 100 GB, up to 2000 GB. Must be in units of 10 GB.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Cluster uuid.
* `endpoints` - .
    * `broker_node_list` - List of broker nodes (Port 9092).
    * `broker_tls_node_list` - List of broker nodes (Port 9093).
    * `public_endpoint_broker_node_list` - List of public endpoint of broker nodes.
    * `public_entpoint_broker_node_listener_port_list` - List of listener port for public endpoint of broker nodes.
    * `public_endpoint_broker_tls_node_list` - List of public endpoint of broker nodes (TLS).
    * `public_endpoint_broker_tls_node_listener_port_list` - List of listener port for public endpoint of broker nodes (TLS).
    * `local_dns_list` - Editing details of the hosts file (Private IP hostname format).
    * `local_dns_tls_list` - Editing details of the hosts file (Private IP hostname format) (TLS).
    * `zookeeper_list` - List of ZooKeeper nodes (Port 2181).