provider "cphalo" {}

data "cphalo_server_group" "root_group" {
  "name" = "kiwi.com"
}

output "server_group_root_id" {
  value = "${data.cphalo_server_group.root_group.id}"
}

resource "cphalo_server_group" "main" {
  name = "tf-main"
}

output "server_group_main_id" {
  value = "${cphalo_server_group.main.id}"
}

resource "cphalo_server_group" "main_parent" {
  name = "tf-main-parent"
  parent_id = "${data.cphalo_server_group.root_group.id}"
}

output "server_group_main_parent_id" {
  value = "${cphalo_server_group.main_parent.id}"
}

resource "cphalo_server_group" "sg1" {
  name = "tf-group"
  parent_id = "${cphalo_server_group.main.id}"
}

output "server_group_sg1_id" {
  value = "${cphalo_server_group.sg1.id}"
}

resource "cphalo_server_group" "sg2" {
  name = "tf-group2"
  parent_id = "${cphalo_server_group.main.id}"
  tag = "new_tag"
}

output "server_group_sg2_id" {
  value = "${cphalo_server_group.sg2.id}"
}
