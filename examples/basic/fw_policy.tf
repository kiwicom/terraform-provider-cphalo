resource "cphalo_firewall_policy" "tf_examples_basic_fw_policy" {
  name = "tf_examples_basic_fw_policy"
  description = "tf examples basic description"
  ignore_forwarding_rules = true

  rule {
    chain = "INPUT"
    action = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position = 1

    firewall_interface = "${cphalo_firewall_interface.tf_examples_basic_fw_interface.id}"
    firewall_service = "${cphalo_firewall_service.tf_examples_basic_fw_service.id}"

    firewall_source = "${cphalo_firewall_zone.tf_examples_basic_fw_in_zone.id}"
  }

  rule {
    chain = "OUTPUT"
    action = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position = 1

    firewall_interface = "${cphalo_firewall_interface.tf_examples_basic_fw_interface.id}"
    firewall_service = "${cphalo_firewall_service.tf_examples_basic_fw_service.id}"

    firewall_target = "${cphalo_firewall_zone.tf_examples_basic_fw_out_zone.id}"
  }
}

resource "cphalo_firewall_zone" "tf_examples_basic_fw_in_zone" {
  name = "tf_examples_basic_fw_in_zone"
  ip_address = "1.1.1.1"
}

resource "cphalo_firewall_zone" "tf_examples_basic_fw_out_zone" {
  name = "tf_examples_basic_fw_out_zone"
  ip_address = "10.10.10.10"
}

resource "cphalo_firewall_service" "tf_examples_basic_fw_service" {
  name = "tf_examples_basic_fw_svc"
  protocol = "TCP"
  port = "2222"
}

resource "cphalo_firewall_interface" "tf_examples_basic_fw_interface" {
  name = "eth42"
}

resource "cphalo_firewall_policy" "tf_examples_basic_fw_subpolicy" {
  name = "tf_examples_basic_fw_subpolicy"
  description = "tf examples basic description"
  ignore_forwarding_rules = true

  rule {
    chain = "INPUT"
    action = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position = 1

    firewall_interface = "${cphalo_firewall_interface.tf_examples_basic_fw_sub_interface.id}"
    firewall_service = "${cphalo_firewall_service.tf_examples_basic_fw_sub_service.id}"

    firewall_source = "${cphalo_firewall_zone.tf_examples_basic_fw_sub_zone.id}"
  }
}

resource "cphalo_firewall_zone" "tf_examples_basic_fw_sub_zone" {
  name = "tf_examples_basic_fw_sub_zone"
  ip_address = "2.2.2.2"
}

resource "cphalo_firewall_service" "tf_examples_basic_fw_sub_service" {
  name = "tf_examples_basic_fw_sub_service"
  protocol = "UDP"
  port = "53"
}

resource "cphalo_firewall_interface" "tf_examples_basic_fw_sub_interface" {
  name = "eth43"
}
