# Data Source: ncloud_nks_version

Gets a kubeconfig from Kubernetes Service cluster.

## Example Usage

```hcl
variable "cluster_name" {}

data "ncloud_nks_kube_config" "kube_config"{
  cluster_name = var.cluster_name
}


// Todo: 아래 k8s provider 사용 부분은 추가할지 논의 필요
// Need Kubenetes Provider Below
provider "kubernetes" {
  host                   = "${data.ncloud_nks_kube_config.kube_config.host}"
  client_certificate     = "${base64decode(data.ncloud_nks_kube_config.kube_config.client_certificate)}"
  client_key             = "${base64decode(data.ncloud_nks_kube_config.kube_config.client_key)}"
  cluster_ca_certificate = "${base64decode(data.ncloud_nks_kube_config.kube_config.cluster_ca_certificate)}"
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

* `cluster_name` - (Required) Cluster Name. (It is the same result as `id`)

## Attributes Reference

* `id` - Cluster Name.
* `host` - Host on kubeconfig.
* `client_certificate` - Client certificate on kubeconfig.
* `client_key` - Client key on kubeconfig.
* `cluster_ca_certificate` - Cluster CA certificate on kubeconfig.
