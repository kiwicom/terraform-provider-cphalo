# cphalo_firewall_policy

Provides details about a specific firewall policy in the CPHalo provider.  
For further information on firewall policies, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-policies).

## Example Usage

```terraform
data "cphalo_firewall_policy" "example" {
	name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the firewall policy.

## Attributes Reference

The following attributes are exported:

* `id` - (string) A unique identifier of the firewall policy.
