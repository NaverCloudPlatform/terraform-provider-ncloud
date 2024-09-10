---
subcategory: "Server"
---

# Data Source: ncloud_login_key

Get a list of Server login keys.

## Example Usage

```terraform
data "ncloud_login_key" "all" {
    filter {
        name = "key_name"
        values = ["test-key1"]
    }
    output_file = "keys.json"
}
output "key_list" {
    value = {
        for key in data.ncloud_login_key.all.login_key_list:
            key.key_name => key.fingerprint
    }
}
```

Outputs:
```terraform
key_list = {
    "test-key1": "2e:fa:e7:f8:fb:4c:18:0e:cd:f2:5b:20:79:c6:77:bd"
}
```

## Argument Reference

The following arguments are supported:
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `vlaues` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the argumments above:

* `login_key_list` - List of login keys.
  * `key_name` - The name of login key.
  * `fingerprint` -  Fingerprint of the login key
  * `create_date` - Login key create date.