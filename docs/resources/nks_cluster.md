---
subcategory: "Kubernetes Service"
---


# Resource: ncloud_nks_cluster

Provides a Kubernetes Service Cluster resource.

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
  return_protection      = false
  log {
    audit = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Cluster name.
* `hypervisor_code` - (Optional) Hypervisor code. `XEN`(Default), `KVM`, `RHV` ( `KVM` supported on `public(site) / KR (region)`)
* `cluster_type` -(Required) Cluster type. `Maximum number of nodes`
  * `XEN` / `RHV`
    * 10 ea : `SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002`
    * 50 ea : `SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002`
  * `KVM`
    * 250 ea : `SVR.VNKS.STAND.C004.M016.G003`
* `login_key_name` - (Required) Login key name.
* `zone` - (Required) zone Code.
* `vpc_no` - (Required) VPC No.
* `subnet_no_list` - (Required) Subnet No. list.
* `public_network` - (Optional) Public Subnet Network (`boolean`)
* `lb_private_subnet_no` - (Required) Subnet No. for private loadbalancer only.
* `lb_public_subnet_no` - (Optional) Subnet No. for public loadbalancer only. (Required in `public` and `gov` site)
* `kube_network_plugin` - (Optional) Specifies the network plugin. Only Cilium is supported.
* `log` - (Optional)
  * `audit` - (Required) Audit log availability. (`boolean`)
* `k8s_version` - (Optional) Kubenretes version. Only upgrade is supported.
* `oidc` - (Optional)
  * `issuer_url` - (Required) Issuer URL.
  * `client_id` - (Required) Client ID.
  * `username_prefix` - (Optional) Username prefix.
  * `username_claim` - (Optional) Username claim.
  * `groups_prefix` - (Optional) Groups prefix.
  * `groups_claim` - (Optional) Groups claim.
  * `required_claim` - (Optional) Required claim.
* `ip_acl_default_action` - (Optional) IP ACL default action. `allow`, `deny` (Supported on `public`, `gov` site)
* `ip_acl` (Optional) (Supported on `public`, `gov` site)
  * `action` - (Required) `allow`, `deny`
  * `address` - (Required) CIDR
  * `comment` - (Optional) Comment
* `return_protection` - (Optional) Return Protection.
* `kms_key_tag` - (Optional) KMS Key Tag for Cluster Secret Encryption.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Cluster uuid.
* `uuid` - Cluster uuid.  (It is the same result as `id`)
* `endpoint` - Control Plane API address.
* `acg_no` - The ID of cluster ACG.

## Import

### `terraform import` command

* Kubernetes Service Cluster can be imported using the `id`. For example:

```console
$ terraform import ncloud_nks_cluster.rsc_name a80d6cbb-fdaa-4fdf-a3d9-063b6ffd5e
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Kubernetes Service Cluster using the `id`. For example:

```terraform
import {
  to = ncloud_nks_cluster.rsc_name
  id = "a80d6cbb-fdaa-4fdf-a3d9-063b6ffd5e"
}
```
