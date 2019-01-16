---
layout: "ncloud"
page_title: "NCLOUD: ncloud_server"
sidebar_current: "docs-ncloud-resource-server"
description: |-
  Provides a ncloud server instance resource.
---

# ncloud_server

Provides a ncloud server instance resource.

## Example Usage

```hcl
resource "ncloud_server" "server" {
    "name" = "tf-test-vm1"
    "server_image_product_code" = "SPSW0LINUX000032"
    "server_product_code" = "SPSVRSTAND000004"

    "tag_list" = [
        {
            "tag_key"   = "samplekey1"
            "tag_value" = "samplevalue1"
        },
        {
            "tag_key"   = "samplekey2"
            "tag_value" = "samplevalue2"
        },
    ]
}
```

## Argument Reference

The following arguments are supported:

* `server_image_product_code` - (Conditional) Server image product code to determine which server image to create. It can be obtained through `data ncloud_server_images`. You are required to select one among two parameters: server image product code (server_image_product_code) and member server image number(member_server_image_no).
* `server_product_code` - (Optional) Server product code to determine the server specification to create. It can be obtained through the getServerProductList action. Default : Selected as minimum specification. The minimum standards are 1. memory 2. CPU 3. basic block storage size 4. disk type (NET,LOCAL)
* `member_server_image_no` - (Conditional) Required value when creating a server from a manually created server image. It can be obtained through the getMemberServerImageList action.
* `name` - (Optional) Server name to create. default: Assigned by ncloud
* `description` - (Optional) Server description to create
* `login_key_name` - (Optional) The login key name to encrypt with the public key. Default : Uses the most recently created login key name
* `is_protect_server_termination` - (Optional) You can set whether or not to protect return when creating. default : false
* `internet_line_type_code` - (Optional) Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)
* `fee_system_type_code` - (Optional) A rate system identification code. There are time plan(MTRAT) and flat rate (FXSUM). Default : Time plan(MTRAT)
* `zone_code` - (Optional) Zone code. You can determine the ZONE where the server will be created. Default : Assigned by NAVER Cloud Platform.
    Get available values using the data source `ncloud_zones`.
    Conflicts with `zone_no`. Only one of `zone_no` and `zone_code` can be used.
* `zone_no` - (Optional) Zone number. You can determine the ZONE where the server will be created. Default : Assigned by NAVER Cloud Platform.
    Get available values using the data source `ncloud_zones`.
    Conflicts with `zone_code`. Only one of `zone_no` and `zone_code` can be used.
* `access_control_group_configuration_no_list` - (Optional) You can set the ACG created when creating the server. ACG setting number can be obtained through the getAccessControlGroupList action. Default : Default ACG number
* `user_data` - (Optional) The server will execute the user data script set by the user at first boot. To view the column, it is returned only when viewing the server instance.
* `raid_type_name` - (Optional) Raid Type Name.
* `tag_list` - (Optional) Server instance tag list.
  * `tag_key` - (Required) Instance tag key
  * `tag_value` - (Required) Instance tag value

## Attributes Reference

* `id` - The instance ID.
* `instance_no` - Server instance number
* `cpu_count` - number of CPUs
* `memory_size` - The size of the memory in bytes.
* `base_block_storage_size` - The size of base block storage in bytes
* `platform_type` - Platform type
    * `code` - Platform type code
    * `code_name` - Platform type name
* `is_fee_charging_monitoring` - Fee charging monitoring
* `public_ip` - Public IP
* `private_ip` - Private IP
* `server_image_name` - Server image name
* `instance_status` - Server instance status
    * `code` - Server instance status code
    * `code_name` - Server instance status code name
* `instance_status_name` - Server instance status name
* `instance_operation` - Server instance operation
    * `code` - Server instance operation code
    * `code_name` - Server instance operation code name
* `create_date` - Creation date of the server instance
* `uptime` - Server uptime
* `port_forwarding_public_ip` - Port forwarding public ip
* `port_forwarding_external_port` - Port forwarding external port
* `port_forwarding_internal_port` - Port forwarding internal port
* `zone` - Zone info
    * `zone_no` - Zone number
    * `zone_code` - Zone code
    * `zone_name` - Zone name
    * `zone_description` - Zone description
    * `region_no` - Region number
* `region` - Region info
    * `region_no` - Region number
    * `region_code` - Region code
    * `region_name` - Region name
* `base_block_storage_disk_type` - Base block storage disk type
    * `code` - Base block storage disk type code
    * `code_name` - Base block storage disk type code name
* `base_block_storage_disk_detail_type` - Base block storage disk detail type
    * `code` - Base block storage disk detail type code
    * `code_name` - Base block storage disk detail type name
* `internet_line_type` - Internet line type
    * `code` - Internet line type code
    * `code_name` - Internet line type code name
