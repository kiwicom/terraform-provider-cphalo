# cphalo_alert_profile

Provides details about a specific alert profile in the CPHalo provider.  
For further information on alert profiles, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#alert-profiles).

## Example Usage

```terraform
data "cphalo_alert_profile" "foo" {
	name = "alert_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the alert profile.

## Attributes Reference

The following attributes are exported:

* `id` - (string) The Halo ID (a unique identifier) of the alert profile.
