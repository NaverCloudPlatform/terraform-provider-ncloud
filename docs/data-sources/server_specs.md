---
subcategory: "Server"
---


# Data Source: ncloud_server_specs

To create a server instance (VM), you should select a server spec. This data source gets a list of server specs. You can look up the spec code of servers.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take a list of Server spec.

```terraform
data "ncloud_server_specs" "example" {
  output_file = "spec.json" 
}

output "spec_list" {
  value = {
    for spec in data.ncloud_server_specs.example.server_spec_list:
    spec.server_spec_code => [spec.description, spec.generation_code]
  }
}
```

```terraform
data "ncloud_server_specs" "example" {
  filter {
    name = "server_spec_code"
    values = ["c2-g3"]
  }
}
```

Outputs: 
```terraform
spec_list = {
  "c2-g3" = [
    "vCPU 2EA, Memory 4GB",
    "G3",
  ]
  "m2-g3" = [
    "vCPU 2EA, Memory 16GB",
    "G3",
  ]
  "c2-g2-s50" = [
    "vCPU 2EA, Memory 4GB, [SSD]Disk 50GB",
    "G2",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `server_spec_list` - List of SEVER spec.
  * `server_spec_code` - Server spec code.
  * `hypervisor_type` - Server hypervisor type. (`XEN` or `KVM`)
  * `generation_code` - Server generation code. (`G2` or `G3`)
  * `cpu_architecture_type` - Server cpu type.
  * `cpu_count` - Server cpu count.
  * `memory_size` - Server memory size(Byte).
  * `block_storage_max_count` - Maximum number of BlockStorage that can be allocated. 
  * `block_storage_max_iops` - BlockStorage max IOPS.
  * `block_storage_max_throughput` - BlockStorage max throughput(Mbps).
  * `network_performance` - Network performance(bps).
  * `network_interface_max_count` - Maximum number of network interfaces that can be allocated.
  * `gpu_count` - GPU count.
  * `description` - Server sepc description.
  * `product_code` - The code of server product.
