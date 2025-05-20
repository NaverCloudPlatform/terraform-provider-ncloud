---
subcategory: "Server"
---


# Resource: ncloud_server

Provides a Server instance resource.

## Example Usage

#### Basic (Classic)

```terraform
resource "ncloud_server" "server" {
  name = "tf-test-vm1"
  member_server_image_no = "12345"

  tag_list {
    tag_key = "samplekey1"
    tag_value = "samplevalue1"
  }
  
  tag_list {
    tag_key = "samplekey2"
    tag_value = "samplevalue2"
  }
}
```

#### Basic (VPC)

```terraform
resource "ncloud_login_key" "loginkey" {
  key_name = "test-key"
}

resource "ncloud_vpc" "test" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
  vpc_no         = ncloud_vpc.test.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.test.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.test.default_network_acl_no
  subnet_type    = "PUBLIC"
  usage_type     = "GEN"
}

resource "ncloud_server" "server" {
  subnet_no                 = ncloud_subnet.test.id
  name                      = "my-tf-server"
  server_image_number       = "25495367"
  server_spec_code          = "s2-g3"
  login_key_name            = ncloud_login_key.loginkey.key_name
}
```

#### Create VPC instance reference by data source  (retrieve server_image_number and server_spec_code)

```terraform
resource "ncloud_login_key" "loginkey" {
  key_name = "test-key"
}

resource "ncloud_vpc" "test" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
  vpc_no         = ncloud_vpc.test.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.test.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.test.default_network_acl_no
  subnet_type    = "PUBLIC"
  usage_type     = "GEN"
}

data "ncloud_server_image_numbers" "kvm-image" {
  server_image_name = "rocky-8.10-base"
  filter {
    name = "hypervisor_type"
    values = ["KVM"]
  }
}

data "ncloud_server_specs" "kvm-spec" {
  filter {
    name   = "server_spec_code"
    values = ["s2-g3"]
  }
}

resource "ncloud_server" "kvm-server" {
  subnet_no                 = ncloud_subnet.test.id
  name                      = "tf-kvm-server"
  server_image_number       = data.ncloud_server_image_numbers.kvm-image.image_number_list.0.server_image_number
  server_spec_code          = data.ncloud_server_specs.kvm-spec.server_spec_list.0.server_spec_code
  login_key_name            = ncloud_login_key.loginkey.key_name
}

data "ncloud_server_image_numbers" "xen-image" {
  server_image_name = "rocky-8.10-base"
  filter {
    name = "hypervisor_type"
    values = ["XEN"]
  }
}

data "ncloud_server_specs" "xen-spec" {
  filter {
    name   = "server_spec_code"
    values = ["s2-g2-s50"]
  }
}

resource "ncloud_server" "xen-server" {
  subnet_no                 = ncloud_subnet.test.id
  name                      = "tf-xen-server"
  server_image_number       = data.ncloud_server_image_numbers.xen-image.image_number_list.0.server_image_number
  server_spec_code          = data.ncloud_server_specs.xen-spec.server_spec_list.0.server_spec_code
  login_key_name            = ncloud_login_key.loginkey.key_name
}
```

## Argument Reference

The following arguments are supported:

* `server_image_product_code` - (Optional, Required if `member_server_image_no` or `server_image_number` is not provided) Server image product code to determine which server image to create. It can be obtained through `data.ncloud_server_image(s)`.
  - [`ncloud_server_image` data source](../data-sources/server_image.md)
  - [`ncloud_server_images` data source](../data-sources/server_images.md)

* `server_product_code` - (Optional) Server product code to determine the server specification to create. It can be obtained through the `data.ncloud_server_product(s)` action. Default : Selected as minimum specification. The minimum standards are 1. memory 2. CPU 3. basic block storage size 4. disk type (NET,LOCAL)
  - [`ncloud_server_product` data source](../data-sources/server_product.md)
  - [`ncloud_server_products` data source](../data-sources/server_products.md)

* `member_server_image_no` - (Optional, Required if `server_image_product_code` or `server_image_number` is not provided) Required value when creating a server from a manually created server image. KVM hypervisor type server images are not supported. It can be obtained through the `data.ncloud_member_server_image(s)` action.
  - [`ncloud_member_server_image` data source](../data-sources/member_server_image.md)
  - [`ncloud_member_server_images` data source](../data-sources/member_server_images.md)

* `name` - (Optional) Server name to create. default: Assigned by ncloud
* `description` - (Optional) Server description to create.
* `login_key_name` - (Optional) The login key name to encrypt with the public key. Default : Uses the login key name most recently created.
* `is_protect_server_termination` - (Optional) You can set whether or not to protect return when creating. default :false
* `fee_system_type_code` - (Optional) A rate system identification code. There are time plan(MTRAT) and flat rate (FXSUM). Default : Time plan(MTRAT)
* `zone` - (Optional) Zone code. You can determine the ZONE where the server will be created. Default : Assigned by NAVER Cloud Platform. Get available values using the data source `ncloud_zones`.
* `raid_type_name` - (Optional) Raid Type Name. raidTypeName is required to create BareMetal servers. You must request an increase in BareMetal server creation limits through customer support center. Accepted value example : `1` |  `5`

~> **NOTE:** Below arguments only support Classic environment.

* `access_control_group_configuration_no_list` - (Optional) You can set the ACG created when creating the server. ACG setting number can be obtained through the getAccessControlGroupList action. Default : Default ACG number
* `user_data` - (Optional) The server will execute the user data script set by the user at first boot. To view the column, it is returned only when viewing the server instance.
* `tag_list` - (Optional) Server instance tag list.
  * `tag_key` - (Required) Instance tag key
  * `tag_value` - (Required) Instance tag value

~> **NOTE:** Below arguments only support VPC environment. Please set `support_vpc` of provider to `true`

* `subnet_no` - (Required) The ID of the associated Subnet.
* `server_image_number` - (Optional, Required if `server_image_product_code` or `member_server_image_no` is not provided) Required to create a KVM hypervisor type 3rd generation server. Server image number to determine which server image to create. It can be obtained through `data.ncloud_server_image_numbers`.
  - [`ncloud_server_image_numbers` data source](../data-sources/server_image_numbers.md)
* `server_spec_code` - (Optional, Required if to select the spec) Available only if `server_image_number` is entered. Server spec code to determine the server specification to create. It can be obtained through the `data.ncloud_server_specs` action. Default : Selected as minimum specification. The minimum standards are 1. memory 2. CPU 3. basic block storage size 4. disk type (NET,LOCAL)
  - [`ncloud_server_specs` data source](../data-sources/server_specs.md)
* `init_script_no` - (Optional) Set init script ID, The server can run a user-set initialization script at first boot.
* `placement_group_no` - (Optional) Physical placement group that belongs to the server instance.
* `network_interface` - (Optional) List of Network Interface. You can assign up to three network interfaces.
  * `network_interface_no` - (Required) If you want to add a network interface that you created yourself, set the network interface ID.
  * `order` - (Required) Sets the order of network interfaces to be assigned to the server to create. The unit name (eth0, eth1, etc.) is determined in that order. There must be one primary network interface. If you set `0`, network interface is set by default. You can assign up to three network interfaces.
* `is_encrypted_base_block_storage_volume` - (Optional) you can set whether to encrypt basic block storage if server image is RHV. Default `false`.

## Attributes Reference

* `id` - The ID of server instance.
* `instance_no` - The ID of server instance.
* `cpu_count` - number of CPUs.
* `memory_size` - The size of the memory in bytes.
* `base_block_storage_size` - The size of base block storage in bytes.
* `platform_type` - Platform type code.
* `public_ip` - Public IP.
* `private_ip` - Private IP.
* `server_image_name` - Server image name.
* `port_forwarding_public_ip` - Port forwarding public ip.
* `port_forwarding_external_port` - Port forwarding external port.
* `port_forwarding_internal_port` - Port forwarding internal port.
* `base_block_storage_disk_type` - Base block storage disk type code.
* `base_block_storage_disk_detail_type` - Base block storage disk detail type code.

~> **NOTE:** Below attributes only provide VPC environment.

* `vpc_no` - The ID of the VPC where you want to place the Server Instance.
* `hypervisor_type` - Hypervisor type. (`XEN` or `KVM`)
* `network_interface` - List of Network Interface.
  * `subnet_no` - Subnet ID of the network interface.
  * `private_ip` - IP address of the network interface.

## Import

### `terraform import` command

* Server can be imported using the `id`. For example:

```console
$ terraform import ncloud_server.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Server using the `id`. For example:

```terraform
import {
  to = ncloud_server.rsc_name
  id = "12345"
}
```
