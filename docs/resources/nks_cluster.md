# Resource: ncloud_nks_cluster

Provides a Kubernetes Service cluster resource.

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
  subnet_type    = "PUBLIC"
  name           = "subnet-lb"
  usage_type     = "LOADB"
}


data "ncloud_nks_versions" "version" {
  filter {
    name = "value"
    values = ["1.20"]
    regex = true
  }
}

resource "ncloud_login_key" "loginkey" {
  key_name = "sample-login-key"
}


resource "ncloud_nks_cluster" "cluster" {
  cluster_type                = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
  k8s_version                 = data.ncloud_nks_versions.version.versions.0.value
  login_key_name              = ncloud_login_key.loginkey.key_name
  name                        = "sample-cluster"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  kube_network_plugin         = "cilium"
  subnet_no_list              = [ ncloud_subnet.subnet.id ]
  vpc_no                      = ncloud_vpc.vpc.id
  zone                        = "KR-1"
  log {
    audit = true
  }
}


```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Cluster name.
* `cluster_type` -(Required) Cluster type. `Maximum number of nodes`
  * 10 ea : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
  * 50 ea : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
* `login_key_name` - (Required) Login key name.
* `zone` - (Required) zone Code.
* `vpc_no` - (Required) VPC No.
* `subnet_no_list` - (Required) Subnet No. list.
* `public_network` - (Optional) Public Subnet Network (`boolean`)
* `lb_private_subnet_no` - (Required) Subnet No. for private loadbalancer only.
* `lb_public_subnet_no` - (Optional) Subnet No. for public loadbalancer only. (Available only `SGN` region)
* `log` - (Optional)
  * `audit` - (Required) Audit log availability. (`boolean`)
* `k8s_version` - (Optional) Kubenretes version .

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Cluster uuid.
* `uuid` - Cluster uuid.  (It is the same result as `id`)
* `endpoint` - Control Plane API address.
* `acg_no` - The ID of cluster ACG.

## Import

Kubernetes Service Cluster can be imported using the name, e.g.,

$ terraform import ncloud_nks_cluster.my_cluster uuid

