# Data Source: ncloud_nks_version

Gets a kubeconfig from Kubernetes Service cluster.

## Example Usage

```hcl
variable "cluster_id" {}

data "ncloud_nks_kube_config" "kc"{
  id = var.cluster_id
}

// Need Kubenetes Provider Below
provider "kubernetes" {
  host                   = "${data.ncloud_nks_kube_config.kc.host}"
  client_certificate     = "${base64decode(data.ncloud_nks_kube_config.kc.client_certificate)}"
  client_key             = "${base64decode(data.ncloud_nks_kube_config.kc.client_key)}"
  cluster_ca_certificate = "${base64decode(data.ncloud_nks_kube_config.kc.cluster_ca_certificate)}"
}

data "kubernetes_all_namespaces" "allns" {}

output "all-ns" {
  value = data.kubernetes_all_namespaces.allns.namespaces
}

output "ns-present" {
  value = contains(data.kubernetes_all_namespaces.allns.namespaces, "kube-system")
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required) Cluster UUID.

## Attributes Reference

- `host` - Host on kubeconfig.
- `client_certificate` - Client certificate on kubeconfig.
- `client_key` - Client key on kubeconfig.
- `cluster_ca_certificate` - Cluster CA certificate on kubeconfig.
