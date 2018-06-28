---
layout: "ncloud"
page_title: "Provider: NAVER Cloud Platform"
sidebar_current: "docs-ncloud-index"
description: |-
  The Ncloud provider is used to interact with Ncloud (NAVER Cloud Platform) services. The provider needs to be configured with the proper credentials before it can be used.

---

# Ncloud Provider

The Ncloud provider is used to interact with
[Ncloud](https://www.ncloud.com) (NAVER Cloud Platform) services.
The provider needs to be configured with the proper credentials before it can be used.


## Example Usage

```hcl
// Configure the ncloud provider
provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
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

### Static credentials ###

Static credentials can be provided by adding an `access_key` `secret_key` and `region` in-line in the
ncloud provider block:

Usage:

```hcl
provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```


###Environment variables

You can provide your credentials via `NCLOUD_ACCESS_KEY` and `NCLOUD_SECRET_KEY`,
environment variables, representing your Ncloud Access Key and Secret Key, respectively.
`NCLOUD_REGION` is also used, if applicable:

```hcl
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
  Ref to : (Get authentication keys for your account)[http://docs.ncloud.com/en/api_new/api_new-1-1.html#preparation]

* `secret_key` - (Required) Ncloud secret key.
  it can also be sourced from the `NCLOUD_SECRET_KEY` environment variable.

* `region` - (Optional) Ncloud region. default 'KR'
  it can also be sourced from the `NCLOUD_REGION` environment variables.

~> **Note** `access_key`, `secret_key` : (Get authentication keys for your account)[http://docs.ncloud.com/en/api_new/api_new-1-1.html#preparation]


## Testing

Credentials must be provided via the `NCLOUD_ACCESS_KEY`, and `NCLOUD_SECRET_KEY` environment variables in order to run acceptance tests.



