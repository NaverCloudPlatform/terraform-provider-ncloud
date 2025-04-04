# Terraform Provider for Naver Cloud Platform

- Website: https://www.terraform.io
- Documentation: https://www.terraform.io/docs/providers/ncloud/index.html
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.jsdelivr.net/gh/hashicorp/terraform-website@master/public/img/logo-hashicorp.svg" width="600px">

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.1.5 or later.
- [Go](https://golang.org/doc/install) v1.23 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/NaverCloudPlatform/terraform-provider-ncloud`

```sh
$ mkdir -p $GOPATH/src/github.com/NaverCloudPlatform; cd $GOPATH/src/github.com/NaverCloudPlatform
$ git clone git@github.com:NaverCloudPlatform/terraform-provider-ncloud.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/NaverCloudPlatform/terraform-provider-ncloud
$ make build
```

## Using the provider

See the [Naver Cloud Platform Provider documentation](http://www.terraform.io/docs/providers/ncloud/index.html) to get started using the Naver Cloud Platform provider.

## Upgrading the provider

To upgrade to the latest stable version of the Naver Cloud Platform provider run `terraform init -upgrade`. See the [Terraform website](https://www.terraform.io/docs/configuration/providers.html#provider-versions) for more information.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ncloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
