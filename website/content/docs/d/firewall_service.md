# cphalo_firewall_service

Provides details about a specific firewall service in the CPHalo provider.
For further information on firewall services, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-services).

## Example Usage

```terraform
data "cphalo_firewall_service" "example" {
	name = "http"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the firewall service.

## Attributes Reference

The following attributes are exported:

* `id` - (string) A unique identifier of the firewall service.
