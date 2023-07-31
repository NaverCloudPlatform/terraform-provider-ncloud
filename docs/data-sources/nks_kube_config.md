---
subcategory: "Kubernetes Service"
---


# Data Source: ncloud_nks_versions

Provides a kubeconfig from Kubernetes Service cluster.

## Example Usage

```hcl
variable "cluster_uuid" {}

data "ncloud_nks_kube_config" "kube_config"{
  cluster_uuid = var.cluster_uuid
}
```

## Argument Reference

The following arguments are supported:

* `cluster_uuid` - (Required) Cluster uuid.

## Attributes Reference

* `id` - Cluster uuid.
* `host` - Host on kubeconfig.
* `client_certificate` - Client certificate on kubeconfig.
* `client_key` - Client key on kubeconfig.
* `cluster_ca_certificate` - Cluster CA certificate on kubeconfig.
