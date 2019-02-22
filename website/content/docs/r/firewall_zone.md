# cphalo_firewall_zone

This resource allows you to create and manage CPHalo firewall zones.  
For further information on firewall zones, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-zones).

## Example Usage

```terraform
resource "cphalo_firewall_zone" "example" {
  name = "databases"
  ip_address = "10.20.30.40,10.20.30.41"
  description = "dev"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A unique name given to the firewall interface.

* `ip_address` - (Required, string) The specified IP address(es) of the firewall zone.

* `description` - (Optional, string) Description of the firewall zone.

## Import

CPHalo firewall zone can be imported using an id, e.g.

```bash
$ terraform import cphalo_firewall_zone.example ff123
```
