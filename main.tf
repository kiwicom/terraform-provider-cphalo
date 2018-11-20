provider "cphalo" {}

data "cphalo_server_group" "root_group" {
  "name" = "main"
}

resource "cphalo_server_group" "sg1" {
  name = "tf-group"
  parent_id = "${data.cphalo_server_group.root_group.id}"
}

output "main_group_id" {
  value = "${data.cphalo_server_group.root_group.id}"
}

output "created_group_id" {
  value = "${cphalo_server_group.sg1.id}"
}

resource "cphalo_server_group" "sg2" {
  name = "tf-group44"
  //  parent_id = "${cphalo_server_group.sg1.id}"
  parent_id = "${data.cphalo_server_group.root_group.id}"
  tag = "new_tag"
}

output "created_sub_group_id" {
  value = "${cphalo_server_group.sg2.id}"
}
