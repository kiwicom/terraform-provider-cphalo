# cphalo_firewall_interface

Provides details about a specific firewall interface in the CPHalo provider.
For further information on firewall interfaces, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-interfaces).

## Example Usage

```terraform
data "cphalo_firewall_interface" "example" {
	name = "eth0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the firewall interface.

## Attributes Reference

The following attributes are exported:

* `id` - (string) A unique identifier of the firewall interface.
