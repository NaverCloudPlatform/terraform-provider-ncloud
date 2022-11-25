# Data Source: ncloud_sourcepipieline_trigger_timezone

This data source is useful for look up the list of Sourcepipeline trigger time zone.

## Example Usage

In the example below, Retrieves all Sourcepipeline schedule trigger time zone list.

```hcl
data "ncloud_sourcepipeline_trigger_timezone" "list_timezone" {
}

output "lookup-timezone-output" {
    value = data.ncloud_sourcepipeline_trigger_timezone.list_timezone.timezone
}
```

## Attributes Reference

The following attributes are exported:

*   `timezone` - The list of Timezone for schedule trigger.
