# cphalo_server_group_firewall_policy

This resource allows you to assign and remove firewall policies to a CPHalo server group.  
For further information on server groups and firewall policies, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#server-groups).

## Example Usage

```terraform
resource "cphalo_server_group_firewall_policy" "foo" {
  group_id                 = "123"
  linux_firewall_policy_id = "456"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, string) ID of the server group

* `linux_firewall_policy_id` - (Required, string) Halo ID of the Linux firewall policy assigned to this group.

## Import

CPHalo server group firewall policy can be imported using an id, e.g.

```bash
$ terraform import cphalo_server_group_firewall_policy.foo ff123
```

Where `ff123` is the Halo server group ID.
