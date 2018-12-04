resource "cphalo_server_group" "root_group" {
  name = "changed_name"
  tag = "added_tag"
  description = "and added some interesting description"
}

resource "cphalo_server_group" "child_group_01" {
  name = "child group 01"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_02" {
  name = "child group 02"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_03" {
  name = "child group 03"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_04" {
  name = "child group 04"
  parent_id = "${cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "child_group_05" {
  name = "child group 05"
  parent_id = "${cphalo_server_group.root_group.id}"
}
