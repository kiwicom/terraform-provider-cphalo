data "cphalo_server_group" "base_server_group" {
  name = "kiwi.com"
}

output "name" {
  value = "${data.cphalo_server_group.base_server_group.name}"
}

output "id" {
  value = "${data.cphalo_server_group.base_server_group.id}"
}

resource "cphalo_server_group" "tf_examples_basic_root_group" {
  name = "tf_examples_basic_root_group"
  tag = "tf_examples_basic_tag"
  description = "tf examples basic description"
  linux_firewall_policy_id = "${cphalo_firewall_policy.tf_examples_basic_fw_policy.id}"
}

resource "cphalo_server_group" "tf_examples_basic_child_group_01" {
  name = "tf_examples_basic_child_group_01"
  parent_id = "${cphalo_server_group.tf_examples_basic_root_group.id}"
  linux_firewall_policy_id = "${cphalo_firewall_policy.tf_examples_basic_fw_subpolicy.id}"
}

resource "cphalo_server_group" "tf_examples_basic_child_group_02" {
  name = "tf_examples_basic_child_group_02"
  parent_id = "${cphalo_server_group.tf_examples_basic_root_group.id}"
}
