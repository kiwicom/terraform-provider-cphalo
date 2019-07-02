resource "cphalo_server_group" "root_group" {
  name        = "{{.Prefix}}changed_name"
  tag         = "{{.Prefix}}added_tag"
  description = "and added some interesting description"
}
