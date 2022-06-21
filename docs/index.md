---
page_title: "Provider: NAVER Cloud Platform"
---

# Ncloud Provider

The Ncloud provider is used to interact with
[Ncloud](https://www.ncloud.com) (NAVER Cloud Platform) services.
The provider needs to be configured with the proper credentials before it can be used.


## Example Usage

```hcl
terraform {
  required_providers {
    ncloud = {
      source = "NaverCloudPlatform/ncloud"
    }
  }
  required_version = ">= 0.13"
}

// Configure the ncloud provider
provider "ncloud" {
  access_key  = var.access_key
  secret_key  = var.secret_key
  region      = var.region
  site        = var.site
  support_vpc = var.support_vpc
}

// Create a new server instance
resource "ncloud_server" "server" {
  # ...
}
```

## Authentication


The Ncloud provider offers a flexible means of providing credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding an `access_key` `secret_key` `region` and `site` in-line in the
ncloud provider block:

Usage:

```hcl
terraform {
  required_providers {
    ncloud = {
      source = "NaverCloudPlatform/ncloud"
    }
  }
  required_version = ">= 0.13"
}

provider "ncloud" {
  access_key  = var.access_key
  secret_key  = var.secret_key
  region      = var.region
  site        = var.site
  support_vpc = var.support_vpc
}
```


### Environment variables

You can provide your credentials via `NCLOUD_ACCESS_KEY` and `NCLOUD_SECRET_KEY`,
environment variables, representing your Ncloud Access Key and Secret Key, respectively.
`NCLOUD_REGION` is also used, if applicable:

```hcl
terraform {
  required_providers {
    ncloud = {
      source = "NaverCloudPlatform/ncloud"
    }
  }
  required_version = ">= 0.13"
}

provider "ncloud" {}
```

Usage:

```shell
$ export NCLOUD_ACCESS_KEY="accesskey"
$ export NCLOUD_SECRET_KEY="secretkey"
$ export NCLOUD_REGION="KR"
$ terraform plan
```


## Argument Reference

The following arguments are supported:

* `access_key` - (Required) Ncloud access key.
  it can also be sourced from the `NCLOUD_ACCESS_KEY` environment variable.
  Ref to : [Get authentication keys for your account](http://docs.ncloud.com/en/api_new/api_new-1-1.html#preparation)

* `secret_key` - (Required) Ncloud secret key. it can also be sourced from the `NCLOUD_SECRET_KEY` environment variable.
* `region` - (Required) Ncloud region. it can also be sourced from the `NCLOUD_REGION` environment variables. It can be
  obtained through `data.ncloud_regions`
  - [`ncloud_regions` data source](data-sources/regions.md)

* `site` - (Optional) Ncloud site. By default, the value is "public". You can specify only the following value: "public"
  , "gov", "fin". "public" is for `www.ncloud.com`. "gov" is for `www.gov-ncloud.com`. "fin" is for `www.fin-ncloud.com`
  .

~> **Note** `access_key`, `secret_key` : [Get authentication keys for your account](http://docs.ncloud.com/en/api_new/api_new-1-1.html#preparation)

~> **Note** If the `access_key` and `secret_key` are externally exposed, then others may use them to access users' accounts.
Therefore, please carefully manage `access_key` and `secret_key`. Take special care to keep `access_key` and `secret_key` from being uploaded to the public version control system


* `support_vpc` - (Optional) Whether to use VPC. By default, the value is `false` on "public" site. If you want to use VPC environment. Please set this value `true`.  

~> **Note** `support_vpc` is only support if `site` is `public`.

## Testing

Credentials must be provided via the `NCLOUD_ACCESS_KEY`, and `NCLOUD_SECRET_KEY` environment variables in order to run acceptance tests.



