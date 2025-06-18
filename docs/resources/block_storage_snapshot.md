---
subcategory: "Server"
---


# Resource: ncloud_block_storage_snapshot

Provides a ncloud Block Storage Snapshot resource.

## Example Usage

```terraform
resource "ncloud_block_storage_snapshot" "snapshot" {
	block_storage_instance_no = "812345"
	name = "tf-test-snapshot1"
	description = "Terraform test snapshot1"
}
```

## Argument Reference

The following arguments are supported:

* `block_storage_instance_no` - (Required) Block storage instance Number for creating snapshot.
* `name` - (Optional) Block storage snapshot name to create. default : Ncloud assigns default values.
* `description` - (Optional) Descriptions on a snapshot to create.

## Attributes Reference

* `instance_no` - Block Storage Snapshot Instance Number
* `volume_size` - Block Storage Snapshot Volume Size
* `instance_status` - Block Storage Snapshot Instance Status code
* `instance_status_name` - Block Storage Snapshot Instance Status Name
* `instance_operation` - Block Storage Snapshot Instance Operation code
* `hypervisor_type` - Hypervisor type. (`XEN` or `KVM`)


## Import

### `terraform import` command

* Block Storage Snapshot can be imported using the `id`. For example:

```console
$ terraform import ncloud_block_storage_snapshot.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Block Storage Snapshot using the `id`. For example:

```terraform
import {
  to = ncloud_block_storage_snapshot.rsc_name
  id = "12345"
}
```
