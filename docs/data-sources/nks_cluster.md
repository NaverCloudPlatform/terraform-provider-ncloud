# Data Source: ncloud_nks_cluster

Provides a Kubernetes Service cluster data.

## Example Usage

```hcl
variable "cluster_name" {}

data "ncloud_nks_cluster" "cluster"{
  name = var.cluster_name
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Cluster name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Cluster name.
* `uuid` - Cluster uuid.
* `endpoint` - Control Plane API address.
* `region_code` - Region code.
* `subnet_lb` - Subnet No. for loadbalancer only.
* `subnet_no_list` - Subnet No. list.
* `cluster_type` - Cluster type. `Maximum number of nodes`
  * 10 ea : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
  * 50 ea : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
* `login_key_name` - Login key name.
* `zone` - zone Code.
* `vpc_no` - VPC No.
* `log` 
  * `audit` - Audit log availability.
* `k8s_version` - Kubenretes version.
