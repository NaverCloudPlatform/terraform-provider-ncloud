---
subcategory: "Apigw"
---


# Resource: ncloud_apigw_product

Provides APIGW product resource.

## Example Usage

```terraform
resource "ncloud_apigw_product" "example" {
  product_name = "tf-prod"
  subscription_code = "PUBLIC"
  description = "test"
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required) Name of the APIGW product to create. Min: 1, Max: 100
* `subscription_code` - (Required) Subscription method of the APIGW product to create. `PUBLIC`: The API can be used by anyone without approval. `PROTECTED`: Publisher approval is required to use the API.
* `description` - (Optional) Description of the product to create. Max: 300

## Attribute Reference

* `id` - Unique ID for product. ID is same as `invoke_id`.
* `tenant_id` - tenant id.
* `published` - published or not (true/false).
* `modifier` - modifier id.
* `domain_code` - NCLOUD domain code (PUB/FIN/GOV).
* `deleted` - deleted or not (true/false).
* `mod_time` - modify time.
* `zone_code` - zone code.

## Import

### `terraform import` command

* APIGW product can be imported using the `id`. For example:

```console
$ terraform import ncloud_apigw_product.rsc_name example-id
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import APIGW product using the `id`. For example:

```terraform
import {
    to = ncloud_apigw_product.rsc_name
    id = "example-id"
}
```
