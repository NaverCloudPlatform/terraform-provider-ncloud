---
subcategory: "Guide"
---

# Remote State Backend

~> **NOTE** This configuration is applicable from Terraform version 1.10.0 onwards.

Ncloud object storage uses the Amazon S3 Compatible API, which now supports conditional writing. Therefore, built-in remote state and state locking features are available.

## Example Configuration

~> **NOTE** The storage must be created before using the remote state locking feature.

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