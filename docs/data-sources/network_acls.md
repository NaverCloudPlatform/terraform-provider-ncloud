---
subcategory: "VPC"
---


# Data Source: ncloud_network_acls

This resource is useful for look up the list of Network ACL in the region.

## Example Usage

The example below shows how to make multiple network acl rule using an existing Network ACLs.

```hcl
data "ncloud_network_acls" "nacl" {}

resource "ncloud_network_acl_rule" "nacl_rule" {
  count             = length(data.ncloud_network_acls.nacl.network_acls)
  network_acl_no    = data.ncloud_network_acls.nacl.network_acls[count.index].id

  inbound {
    priority          = 100
    protocol          = "TCP"
    rule_action       = "ALLOW"
    ip_block          = "0.0.0.0/0"
    port_range        = "22"
  }

  outbound {
    priority          = 100
    protocol          = "TCP"
    rule_action       = "ALLOW"
    ip_block          = "0.0.0.0/0"
    port_range        = "1-65535"
  }
}
```

## Argument Reference

The following arguments are supported:  

* `network_acl_no_list` - (Optional) List of Network ACL ID to retrieve.
* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `name` - (Optional) name of the specific Network ACLs to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

The following attributes are exported:

* `network_acls` - The list of Network ACL

### Network ACL Reference

`network_acls` are also exported with the following attributes, where are relevant: Each element supports the following:

* `id` - The ID of Network ACL.
* `network_acl_no` - The ID of Network ACL. (It is the same result as `id`)
* `vpc_no` - The ID of the associated VPC.
* `is_default` - Whether default or not by VPC creation.
* `name` - The name of Network ACL.
* `description` - Description of Network ACL.
