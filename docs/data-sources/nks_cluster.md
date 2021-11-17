# Data Source: ncloud_nks_cluster

Provides a Kubernetes Service cluster data.

## Example Usage

### Basic Usage

```hcl
variable "cluster_name" {}

data "ncloud_nks_cluster" "cluster"{
  cluster_name = var.cluster_name
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of Cluster. (It is the same result as `id`)

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
* `cluster_type` - Cluster type.
  * 10 nodes : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
  * 50 nodes : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
* `login_key_name` - Key name to configure worker nodes.
* `zone_no` - Available zone number where the cluster will be placed physically.
* `vpc_no` - The ID of the VPC where you want to place the cluster.
* `subnet_no_list` - The ID list of the Subnet where you want to place the cluster.
* `subnet_lb` - The ID the Subnet where you want to place the Loadbalancer.
* `log` - Use Log.
  * `audit` - `Boolean`.
* `k8s_version` - Kubenretes version to create.
