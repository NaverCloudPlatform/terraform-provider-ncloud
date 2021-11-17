# Data Source: ncloud_nks_cluster

Provides a Kubernetes Service cluster resource.

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
  subnet_type    = "PUBLIC"
  name           = "subnet-lb"
  usage_type     = "LOADB"
}


data "ncloud_nks_version" "version"{
}

resource "ncloud_login_key" "loginkey" {
  key_name = "sample-login-key"
}


resource "ncloud_nks_cluster" "cluster" {
  cluster_type                = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
  k8s_version                 = data.ncloud_nks_version.version.versions.0.value
  login_key_name              = ncloud_login_key.loginkey.key_name
  name                        = "sample-cluster"
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [ ncloud_subnet.subnet.id ]
  vpc_no                      = ncloud_vpc.vpc.id
  zone_no                     = "2"
}


```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of Cluster. (It is the same result as `id`)
* `cluster_type` - (Required) Cluster type.
  * 10 nodes : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
  * 50 nodes : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
* `login_key_name` - (Required) Key name to configure worker nodes.
* `zone_no` - (Required) Available zone number where the cluster will be placed physically.
* `vpc_no` - (Required) The ID of the VPC where you want to place the cluster.
* `subnet_no_list` - (Required) The ID list of the Subnet where you want to place the cluster.
* `subnet_lb` - (Required) The ID the Subnet where you want to place the Loadbalancer.
* `log` - (Optional) Use Log.
  * `audit` - (Required) `Boolean`.
* `k8s_version` - (Optional) Kubenretes version to create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Name of Cluster.
* `uuid` - The UUID of Cluster.
* `acg_name` - Name of the ACG, which is configured in the cluster.
* `created_at` - Cluster creation time.
* `endpoint` - Cluster endpoint.
* `region_code` - Resion code.
* `status` - Cluster status.
* `subnet_lb_name` - Name of Lb Subnet.
* `subnet_name` - Name of Subnet.
* `vpc_name` - The name of VPC.

## Import

Kubernetes Service Cluster can be imported using the name, e.g.,

$ terraform import ncloud_nks_cluster.my_cluster my_cluster