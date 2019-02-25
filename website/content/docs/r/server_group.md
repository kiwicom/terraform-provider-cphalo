# cphalo_server_group

This resource allows you to create and manage CPHalo server groups.  
For further information on server groups, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#server-groups).

## Example Usage

```terraform
resource "cphalo_server_group" "foo" {
  name                     = "foo"
  description              = "some description"
  parent_id                = "123"
  tag                      = "a_tag"
  linux_firewall_policy_id = "123"
  alert_profile_ids        = ["1", "2"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A unique name for this group.

* `description` - (Optional, string) Additional descriptive information.

* `parent_id` - (Optional, string) The Halo ID of this group's parent group. If not provided, this group is the root group.

* `tag` - (Optional, string) A unique tag assigned to this group.

* `linux_firewall_policy_id` - (Optional, string) Halo ID of the Linux firewall policy assigned to this group.

* `alert_profile_ids` - (Optional, list) Alert profiles assigned to this group.

## Import

CPHalo server group can be imported using an id, e.g.

```bash
$ terraform import cphalo_server_group.foo ff123
```
