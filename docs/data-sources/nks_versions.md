---
subcategory: "Kubernetes Service"
---


# Data Source: ncloud_nks_versions

Provides list of available Kubernetes Service versions.

## Example Usage

```hcl
data "ncloud_nks_versions" "versions" {}

data "ncloud_nks_versions" "v1_22" {
  hypervisor_code = "KVM"
  filter {
    name = "value"
    values = ["1.22"]
    regex = true
  }
}

```

## Argument Reference

~> **Note** The `regional_support` argument (multi-zone / Regional clusters) is only supported on the `public` site. It is not supported on the `fin` and `gov` sites.

The following arguments are supported:

* `hypervisor_code` - (Optional) Hypervisor code. (Default `XEN`)
* `regional_support` - (Optional) Whether to filter versions by multi-zone (Regional) cluster support. Set `true` to return only versions that support Regional clusters, or `false` to return only versions that do not. When omitted, all versions are returned. Multi-zone (Regional) clusters are only supported on the `public` site, and are not supported on the `fin` and `gov` sites.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `versions` - A list of verions
  * `label` - Version label
  * `value` - Version value
