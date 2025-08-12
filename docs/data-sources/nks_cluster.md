---
subcategory: "Kubernetes Service"
---


# Data Source: ncloud_nks_cluster

Provides a Kubernetes Service cluster data.

## Example Usage

```hcl
variable "cluster_uuid" {}

data "ncloud_nks_cluster" "cluster"{
  uuid = var.cluster_uuid
}

```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) Cluster uuid.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Cluster name.
* `id` - Cluster uuid.
* `endpoint` - Control Plane API address.
* `lb_private_subnet_no` - Subnet No. for private loadbalancer only.
* `lb_public_subnet_no` - Subnet No. for public loadbalancer only. (Supported on `public`, `gov` site)
* `subnet_no_list` - Subnet No. list.
* `public_network` - Public Subnet Network
* `kube_network_plugin` - Kubernetes network plugin.
* `hypervisor_code` - Hypervisor code.
* `cluster_type` - Cluster type. `Maximum number of nodes`
  * `XEN` / `RHV`
    * 10 ea : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
    * 50 ea : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
  * `KVM`
    * 250 ea : `SVR.VNKS.STAND.C004.M016.G003`
* `login_key_name` - Login key name.
* `zone` - zone Code.
* `vpc_no` - VPC No.
* `log` 
  * `audit` - Audit log availability.
* `k8s_version` - Kubenretes version.
* `acg_no` - The ID of cluster ACG.
* `oidc`
  * `issuer_url` - Issuer URL.
  * `client_id` - Client ID.
  * `username_prefix` - Username prefix.
  * `username_claim` - Username claim. 
  * `groups_prefix` - Groups prefix.
  * `groups_claim` - Groups claim. 
  * `required_claim` - Required claim.
* `ip_acl_default_action` - IP ACL default action. (Supported on `public`, `gov` site)
* `ip_acl` (Supported on `public`, `gov` site)
  * `action` - `allow`, `deny`
  * `address` - CIDR
  * `comment` - Comment
* `return_protection` - Return Protection.
* `kms_key_tag` - KMS Key Tag for Cluster Secret Encryption.
* `auth_type` - Authentication type for cluster. Valid values are `API` or `CONFIG_MAP`.
* `access_entries` - Access entries for cluster access management. Only available when `auth_type` is `API`.
  * `entry` - NRN (Ncloud Resource Names) of the user or role.
  * `groups` - List of groups assigned to the access entry.
  * `policies` - List of policies for the access entry.
    * `type` - Policy type. Valid values are `NKSClusterAdminPolicy`, `NKSAdminPolicy`, `NKSEditPolicy`, or `NKSViewPolicy`.
    * `scope` - Policy scope. Valid values are `Cluster` or `Namespace`.
    * `namespaces` - List of namespaces.