# cphalo_firewall_service

This resource allows you to create and manage CPHalo firewall services.  
For further information on firewall services, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-services).

## Example Usage

```terraform
resource "cphalo_firewall_service" "example" {
  name = "custom_ssh"
  protocol = "TCP"
  port = "1022"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A unique name given to the firewall interface.

* `protocol` - (Required, string) The specified protocol of the firewall service. TCP, UDP, and ICMP are allowed.

* `port` - (Optional, string) The specified port(s) of the firewall service.

## Import

CPHalo firewall service can be imported using an id, e.g.

```bash
$ terraform import cphalo_firewall_service.example ff123
```
