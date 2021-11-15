# Data Source: ncloud_nks_version

Gets a KubeConfig from nks cluster.

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

* `id` - (Required) NKS Cluster UUID.

## Attributes Reference

* `host` -Host on KubeConfig.
* `client_certificate` - ClientCertificate on KubeConfig.
* `client_key` - ClientKey on KubeConfig.
* `cluster_ca_certificate` - Cluster CA Certificate on KubeConfig.