# cphalo_firewall_zone

Provides details about a specific firewall zone in the CPHalo provider.
For further information on firewall zones, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-zones).

## Example Usage

```terraform
data "cphalo_firewall_zone" "example" {
	name = "any"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the firewall zone.

## Attributes Reference

The following attributes are exported:

* `id` - (string) A unique identifier of the firewall zone.
