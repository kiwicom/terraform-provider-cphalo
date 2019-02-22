# cphalo_firewall_interface

This resource allows you to create and manage CPHalo firewall interfaces.  
For further information on firewall interfaces, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-interfaces).

## Example Usage

```terraform
resource "cphalo_firewall_interface" "example" {
  name = "eth2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A unique name given to the firewall interface.

## Import

CPHalo firewall interface can be imported using an id, e.g.

```bash
$ terraform import cphalo_firewall_interface.example ff123
```
