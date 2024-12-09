---
subcategory: "Guide"
---

# [Terraform Remote State](https://developer.hashicorp.com/terraform/language/state/remote)

~> **NOTE** This configuration is applicable from Terraform version 1.10.0 onwards.

By default, Terraform stores state locally in a file named `terraform.tfstate`. When working with Terraform in a team, use of a local file makes Terraform usage complicated because each user must make sure they always have the latest state data before running Terraform and make sure that nobody else runs Terraform at the same time.

With remote state, Terraform writes the state data to a remote data store, which can then be shared between all members of a team. Terraform supports storing state in Terraform Cloud, HashiCorp Consul, Amazon S3, Azure Blob Storage, Google Cloud Storage, etcd, and more.

Remote state is implemented by a [backend](https://developer.hashicorp.com/terraform/language/backend/configuration). Backends are configured with a nested `backend` block within the top-level `terraform` block:

```hcl
terraform {
  backend "s3" {
    ...
  }
}
```

There are some important limitations on backend configuration:

- A configuration can only provide one backend block.
- A backend block cannot refer to **named values** (like input variables, locals, or data source attributes).

Ncloud object storage uses the Amazon S3 Compatible API, which now supports conditional writing. Therefore, built-in remote state and state locking features are available.

## Example Configuration

~> **NOTE** The storage must be created before using this feature.

```hcl
terraform {
  backend "s3" {
    bucket = "tfstate-backend"
    key    = "remote-state/terraform.tfstate"
    # Set the region according to your location
    region = "KR"
    access_key = var.access_key
    secret_key = var.secret_key

    # To skip AWS authentication logic
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_s3_checksum = true

    # For state locking
    use_lockfile = true

    endpoints = {
      # Set the endpoint according to your region
      s3 = "https://kr.object.ncloudstorage.com"
    }
  }
}

provider "ncloud" {
  access_key  = var.access_key
  secret_key  = var.secret_key
  region      = var.region
  support_vpc = true
}
```

This configuration is similar to the [AWS backend state configuration with Amazon S3](https://developer.hashicorp.com/terraform/language/backend/s3), but you need to add the following options to skip AWS authentication logic:

```hcl
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_s3_checksum = true
```

Additionally, ensure to use the backend state with the `use_lockfile = true` option to enable state locking functionality.
