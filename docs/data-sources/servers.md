# Data Source: ncloud_servers

Use this data source to get multiple `ncloud_server` ids 

## Example Usage

#### Basic usage

```hcl
data "ncloud_servers" "servers" {
  ids = [ncloud_server.example1.id, ncloud_server.example2.id]
}
```

#### Usage of using filter

```hcl
data "ncloud_servers" "servers" {
  filter {
    name = "subnet_no"
    values = [ncloud_subnet.example1.id, ncloud_subnet.example2.id]
  }
}
```

#### Usage of `ncloud_servers` data source in `ncloud_nas_volume`

```hcl
data "ncloud_servers" "servers" {
  ids = [ncloud_server.example1.id, ncloud_server.example2.id]
}

resource "ncloud_nas_volume" "vol" {
	volume_name_postfix            = "vol"
	volume_size                    = "500"
	volume_allotment_protocol_type = "NFS"
  server_instance_no_list = data.ncloud_servers.servers.ids
}
```

The following arguments are supported:

* `ids` - (Optional) The set of ID of the Server instances.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression. 
