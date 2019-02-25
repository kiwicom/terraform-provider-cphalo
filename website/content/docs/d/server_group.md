# cphalo_server_group

Provides details about a specific server group in the CPHalo provider.  
For further information on firewall zones, consult the
[CloudPassage Halo documentation](https://library.cloudpassage.com/help/cloudpassage-api-documentation#server-groups).

## Example Usage

```terraform
resource "cphalo_server_group" "foo" {
  name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the server group.

## Attributes Reference

The following attributes are exported:

* `id` - (string) A unique identifier of the server group.

* `parent_id` - (string) The Halo ID of this group's parent group.
