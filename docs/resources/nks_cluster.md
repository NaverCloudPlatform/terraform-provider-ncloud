# Resource: ncloud_nks_cluster

Provides a NKS Cluster resource.

## Example Usage

### Basic Usage

```hcl

resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.1.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "subnet-01"
  usage_type     = "GEN" 
}

resource "ncloud_subnet" "subnet_lb" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.100.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  name           = "subnet-lb"
  usage_type     = "LOADB"   
}


data "ncloud_nks_version" "version"{
}

resource "ncloud_login_key" "loginkey" {
  key_name = "sample-login-key"
}


resource "ncloud_nks_cluster" "test" {
  cluster_type                = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
  k8s_version                 = data.ncloud_nks_version.version.versions.0.value
  login_key_name              = ncloud_login_key.loginkey.key_name
  name                        = "sample-cluster"
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [ ncloud_subnet.subnet.id ]
  vpc_no                      = ncloud_vpc.vpc.id
  zone_no                     = "2"

  node_pool {
    is_default     = true
    name           = "default-node-pool"
    node_count     = 1
    product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
    subnet_no_list = [ ncloud_subnet.subnet.id ]
  }

  node_pool {
    is_default     = false
    name           = "add-node-pool"
    node_count     = 1
    product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to create.
* `cluster_type` - (Required) Cluster type. 
  * 10 nodes : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
  * 50 nodes : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
* `login_key_name` - (Required) Key name to configure worker nodes.
* `zone_no` - (Required) Available zone nubmer where the cluster will be placed physically.
* `vpc_no` - (Required) The ID of the VPC where you want to place the cluster.
* `subnet_no_list` - (Required) The ID list of the Subnet where you want to place the cluster.
* `subnet_lb` - (Required) The ID the Subnet where you want to place the Loadbalancer.
* `node_pool` - (Required) Cluster nodepool list.
  * `is_Default` - (Optional) `Boolean` Default node YN. Only one default nodepool is allowed.
  * `name` - (Required) The name of node pool.
  * `node_count` - (Required) Number of worker nodes in nodepool.
  * `product_code` - (Required) Product code of worker nodes in node pool
* `log` - (Optional) Use Log.
    * `audit` - (Required) `Boolean`.
* `k8s_version` - (Optional) NKS Version to create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of Cluster.
* `acg_name` - Name of the ACG, which is configured in the cluster.
* `capacity` - Worker node Capacity.
* `cpu_count` - Worker node CPU cores.
* `created_at` - Cluster creation time.
* `endpoint` - Cluster endpoint.
* `memory_size` - Worker node memory size
* `node_count` - Total number of worker nodes in the cluster.
* `node_max_count` - ??
* `region_code` - Resion code.
* `status` - Cluster status.
* `subnet_lb_name` - Name of Lb Subnet.
* `subnet_name` - Name of Subnet.
* `updated_at` - Cluster updated time.
* `uuid` - The ID of the Cluster. (It is the same result as `id`)
* `vpc_name` - The name of VPC.
* `node_pool` - Cluster node pool list.
  * `autoscale`- Autoscale Config.
    * `enable` - Autoscale enable YN.
    * `max` - Max node count.
    * `min` - Min node count.
  * `instance_no` - Instance number of node pool.
  * `status` - Node pool status.

