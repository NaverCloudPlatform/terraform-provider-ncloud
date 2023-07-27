---
subcategory: "Load Balancer"
---


# Resource: load_balancer_ssl_certificate

Provides a ncloud load balancer ssl certificate resource.

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
