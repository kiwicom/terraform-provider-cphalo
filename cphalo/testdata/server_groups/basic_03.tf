resource "cphalo_server_group" "root_group" {
  name = "{{.Prefix}}changed_name"
  tag = "{{.Prefix}}added_tag"
  description = "and added some interesting description"
}

resource "cphalo_server_group" "child_group_01" {
  name = "{{.Prefix}}child group 01"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_02" {
  name = "{{.Prefix}}child group 02"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_03" {
  name = "{{.Prefix}}child group 03"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_04" {
  name = "{{.Prefix}}child group 04"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_05" {
  name = "{{.Prefix}}child group 05"
  parent_id = "${cphalo_server_group.root_group.id}"
}
