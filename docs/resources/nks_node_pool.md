---
subcategory: "Kubernetes Service"
---


# Resource: ncloud_nks_node_pool

Provides a Kubernetes Service nodepool resource.

## Example Usage

```hcl
resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.1.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "subnet-01"
  usage_type     = "GEN"
}

resource "ncloud_subnet" "subnet_lb" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.100.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "subnet-lb"
  usage_type     = "LOADB"
}

data "ncloud_nks_versions" "version" {
  filter {
    name = "value"
    values = ["1.23"]
    regex = true
  }
}

resource "ncloud_login_key" "loginkey" {
  key_name = "sample-login-key"
}

resource "ncloud_nks_cluster" "cluster" {
  cluster_type           = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
  k8s_version            = data.ncloud_nks_versions.version.versions.0.value
  login_key_name         = ncloud_login_key.loginkey.key_name
  name                   = "sample-cluster"
  lb_private_subnet_no   = ncloud_subnet.subnet_lb.id
  kube_network_plugin    = "cilium"
  subnet_no_list         = [ ncloud_subnet.subnet.id ]
  vpc_no                 = ncloud_vpc.vpc.id
  zone_no                = "2"

}

data "ncloud_server_image" "image" {
  filter {
    name = "product_name"
    values = ["ubuntu-20.04"]
  }
}

data "ncloud_server_product" "product" {
  server_image_product_code = data.ncloud_server_image.image.product_code

  filter {
    name = "product_type"
    values = [ "STAND" ]
  }

  filter {
    name = "cpu_count"
    values = [ 2 ]
  }

  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }

  filter {
    name = "product_code"
    values = [ "SSD" ]
    regex = true
  }
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "sample-node-pool"
  node_count     = 1
  product_code   = data.ncloud_server_product_code.product.product_code
  subnet_no      = ncloud_subnet.subnet.id
  autoscale {
    enabled = true
    min = 1
    max = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `node_pool_name` - (Required) Nodepool name. 
* `cluster_uuid` - (Required) Cluster uuid.
* `node_count` - (Required) Number of nodes.
* `product_code` - (Required) Product code.
* `software_code` - (Optional) Server image code.
* `autoscale`- (Optional) 
  * `enable` - (Required) Auto scaling availability.
  * `max` - (Required) Maximum number of nodes available for auto scaling.
  * `min` - (Required) Minimum number of nodes available for auto scaling.
* `subnet_no` - (Deprecated) Subnet No.
* `subnet_no_list` - Subnet no list.
* `k8s_version` - (Optional) Kubenretes version. Only upgrade is supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of nodepool.`CusterUuid:NodePoolName` 
* `instance_no` - Instance No.
* `nodes`- Running nodes in nodepool.
  * `name` - The name of Server instance.
  * `instance_no` - The ID of server instance.
  * `spec` - Server spec.
  * `private_ip` - Private IP.
  * `public_ip` - Public IP.
  * `node_status` - Node Status.
  * `container_version` - Container version of node.
  * `kernel_version` - kernel version of node.
## Import

NKS Node Pools can be imported using the cluster_name and node_pool_name separated by a colon (:), e.g.,

$ terraform import ncloud_nks_node_pool.my_node_pool uuid:my_node_pool

