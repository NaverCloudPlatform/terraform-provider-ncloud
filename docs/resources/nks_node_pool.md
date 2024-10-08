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

resource "ncloud_subnet" "subnet_lb_pub" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.101.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  name           = "subnet-lb-pub"
  usage_type     = "LOADB"
}

data "ncloud_nks_versions" "version" {
  hypervisor_code = "KVM"
  filter {
    name = "value"
    values = ["1.27"]
    regex = true
  }
}

resource "ncloud_login_key" "loginkey" {
  key_name = "sample-login-key"
}

resource "ncloud_nks_cluster" "cluster" {
  hypervisor_code        = "KVM"
  cluster_type           = "SVR.VNKS.STAND.C002.M008.G003"
  k8s_version            = data.ncloud_nks_versions.version.versions.0.value
  login_key_name         = ncloud_login_key.loginkey.key_name
  name                   = "sample-cluster"
  lb_private_subnet_no   = ncloud_subnet.subnet_lb.id
  lb_public_subnet_no    = ncloud_subnet.subnet_lb_pub.id
  kube_network_plugin    = "cilium"
  subnet_no_list         = [ ncloud_subnet.subnet.id ]
  vpc_no                 = ncloud_vpc.vpc.id
  public_network         = false
  zone                   = "KR-1"

}

data "ncloud_nks_server_images" "image"{
  hypervisor_code = "KVM"
  filter {
    name = "label"
    values = ["ubuntu-22.04"]
    regex = true
  }
}

data "ncloud_nks_server_products" "product"{
  software_code = data.ncloud_nks_server_images.image.images[0].value
  zone = "KR-1"

  filter {
    name = "product_type"
    values = [ "STAND"]
  }

  filter {
    name = "cpu_count"
    values = [ "2"]
  }

  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid     = ncloud_nks_cluster.cluster.uuid
  node_pool_name   = "sample-node-pool"
  node_count       = 2
  software_code    = data.ncloud_nks_server_images.image.images[0].value
  server_spec_code = data.ncloud_nks_server_products.product.products.0.value
  storage_size = 200
  autoscale {
    enabled = false
    min = 2
    max = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `node_pool_name` - (Required) Nodepool name. 
* `cluster_uuid` - (Required) Cluster uuid.
* `node_count` - (Optioanl) Number of nodes. (Required when autoscale is disabled) 
* `product_code` - (Optional) Product code. Required for `XEN`/`RHV` cluster nodepool.
* `server_spec_code` - (Optional) Server spec code. (Required for `KVM` cluster nodepool)
* `storage_size` - (Optional) Default storage size for `KVM` nodepool. (Default `100GB`)
* `software_code` - (Optional) Server image code.
* `autoscale`- (Optional) 
  * `enable` - (Required) Auto scaling availability.
  * `max` - (Required) Maximum number of nodes available for auto scaling.
  * `min` - (Required) Minimum number of nodes available for auto scaling.
* `subnet_no` - (Deprecated) Subnet No.
* `subnet_no_list` - (Optional) Subnet no list.
* `k8s_version` - (Optional) Kubenretes version. Only upgrade is supported.
* `label` - (Optional) NodePool label.
  * `key` - (Required) Label key.
  * `value` - (Required) Label value.
* `taint` - (Optional) NodePool taint.
  * `key` - (Required) Taint key.
  * `value` - (Required) Taint value.
  * `effect` - (Required) Taint effect.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of nodepool.`cluster_uuid:node_pool_name`
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

### `terraform import` command

* Kubernetes Service Node Pool can be imported using the `id`. For example:

```console
$ terraform import ncloud_nks_node_pool.rsc_name a80d6cbb-fdaa-4fdf-a3d9-063b6ffd5e:my-node 
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Kubernetes Service Node Pool using the `id`. For example:

```terraform
import {
  to = ncloud_nks_node_pool.rsc_name
  id = "a80d6cbb-fdaa-4fdf-a3d9-063b6ffd5e:my-node"
}
```
