# cphalo_firewall_policy

This resource allows you to create and manage CPHalo firewall policies.  
For further information on firewall policies, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#firewall-policies).

## Example Usage

```terraform
resource "cphalo_firewall_policy" "example" {
  name                    = "example_policy"
  description             = "example"
  ignore_forwarding_rules = false
  shared                  = true

  rule {
    chain             = "INPUT"
    action            = "ACCEPT"
    active            = false
    connection_states = "NEW, ESTABLISHED"
    position          = 1

    firewall_interface = "${cphalo_firewall_interface.example.id}"
    firewall_service   = "${cphalo_firewall_service.example.id}"

    firewall_source {
      id   = "${cphalo_firewall_zone.example.id}"
      kind = "FirewallZone"
    }
  }

  rule {
    chain             = "OUTPUT"
    action            = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position          = 1

    firewall_interface = "${cphalo_firewall_interface.example.id}"
    firewall_service   = "${cphalo_firewall_service.example.id}"

    firewall_target {
      id   = "All Active Servers"
      kind = "Group"
    }
  }
}

resource "cphalo_firewall_zone" "example" {
  name        = "databases"
  ip_address  = "10.20.30.40,10.20.30.41"
  description = "dev"
}

resource "cphalo_firewall_service" "example" {
  name     = "custom_ssh"
  protocol = "TCP"
  port     = "1022"
}

resource "cphalo_firewall_interface" "example" {
  name = "eth2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A unique name given to the firewall policy.

* `platform` - (Optional, string) The OS platform of the firewall policy. Either `windows` or `linux`.

* `description` - (Optional, string) A description of the firewall policy.

* `ignore_forwarding_rules` - (Optional, string) Linux-only. `true` if the iptables default forwarding rules should be ignored; otherwise `false`. Default is `false`. For servers running in a Docker environment, the value should be `true`.

* `rule` - (Optional, object) A firewall rule. Firewall policies can contain multiple rules.

    * `chain` - (Optional, string) Whether the firewall rule covers `INPUT` or `OUTPUT` connections. Allowed values are `INPUT` and `OUTPUT`.
    
    * `action` - (Optional, string) The specified action to take if this rule is matched. Allowed values are `ACCEPT`, `DROP`, and `REJECT` (`REJECT` is Linux-only).
    
    * `active` - (Optional, bool) Whether the firewall rule is active or not.
    
    * `connection_states` - (Optional, string) Linux-only. The specified firewall connection state(s) for this rule. `NEW`, `RELATED`, and `ESTABLISHED` are allowed.
    
    * `position` - (Optional, int) The position order of the rule in the chain.
    
    * `firewall_interface` - (Optional, string) Linux-only. The specified firewall interface for this rule. Specify the ID of the interface you wish to use.
    
    * `firewall_service` - (Optional, string) The specified firewall service for this rule. Specify the ID of the service you wish to use.
    
    * `firewall_source` - (Optional, object) The specified source/zone for an `INPUT` rule. 
    
        * `id` - (Optional, string) ID of the source/zone you wish to use.
        
        * `kind` - (Optional, string) Type of source you wish to use. Allowed values are `FirewallZone`, `Group`, `User`, or `UserGroup`. 
            
            > **Note**: When using `UserGroup` you must specify the name, and not the ID of the source. Currently, only `All GhostPorts users` is a valid `UserGroup`. `All Active Servers` is a special `Group` that has no ID, so you must specify it by name.
    
    * `firewall_target` - (Optional, object) The specified destination/zone for an `OUTPUT` rule.
        
        * `id` - (Optional, string) ID of the destination/zone you wish to use.
        
        * `kind` - (Optional, string) Type of destination you wish to use. Allowed values are `FirewallZone` or `Group`.
        
            > **Note**: `All Active Servers` is a special `Group` that has no ID, so you must specify it by name.

## Import

CPHalo firewall policy can be imported using an id, e.g.

```bash
$ terraform import cphalo_firewall_policy.example ff123
```
