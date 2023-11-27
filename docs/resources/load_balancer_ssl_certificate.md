---
subcategory: "Classic Load Balancer"
---


# Resource: load_balancer_ssl_certificate

Provides a ncloud load balancer ssl certificate resource.

~> **NOTE:** This resource only supports Classic environment.

## Example Usage

```hcl
resource "ncloud_load_balancer_ssl_certificate" "cert" {
  certificate_name      = "tftest_ssl_cert"
  privatekey            = "${file("lbtest.privateKey")}"
  publickey_certificate = "${file("lbtest.crt")}"
  certificate_chain     = "${file("lbtest.chain")}"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_name` - (Required) Name of a certificate
* `privatekey` - (Required) Private key for a certificate
* `publickey_certificate` - (Required) Public key for a certificate
* `certificate_chain` - (Optional) Chainca certificate (Required if the certificate is issued with a chainca)

## Import

### `terraform import` command

* Load Balancer SSL Certificate can be imported using the `id`. For example:

```console
$ terraform import ncloud_load_balancer_ssl_certificate.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Load Balancer SSL Certificate using the `id`. For example:

```terraform
import {
  to = ncloud_load_balancer_ssl_certificate.rsc_name
  id = "12345"
}
```
