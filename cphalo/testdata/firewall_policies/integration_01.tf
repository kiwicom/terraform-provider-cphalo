resource "cphalo_firewall_policy" "fw_policy" {
  name = "{{.Prefix}}tf_acc_fw_policy"
  description = "awesome"
  ignore_forwarding_rules = true
  shared = true

  rule {
    chain = "INPUT"
    action = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position = 1

    firewall_interface = "${cphalo_firewall_interface.fw_interface.id}"
    firewall_service = "${cphalo_firewall_service.fw_service.id}"

    firewall_source = "${cphalo_firewall_zone.fw_in_zone.id}"
  }

  rule {
    chain = "OUTPUT"
    action = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position = 1

    firewall_interface = "${cphalo_firewall_interface.fw_interface.id}"
    firewall_service = "${cphalo_firewall_service.fw_service.id}"

    firewall_target = "${cphalo_firewall_zone.fw_out_zone.id}"
  }
}

resource "cphalo_firewall_zone" "fw_in_zone" {
  name = "{{.Prefix}}tf_acc_fw_in_zone"
  ip_address = "1.1.1.1"
}

resource "cphalo_firewall_zone" "fw_out_zone" {
  name = "{{.Prefix}}tf_acc_fw_out_zone"
  ip_address = "10.10.10.10"
}

resource "cphalo_firewall_service" "fw_service" {
  name = "{{.Prefix}}tf_acc_fw_svc"
  protocol = "TCP"
  port = "2222"
}

resource "cphalo_firewall_interface" "fw_interface" {
  name = "{{.Prefix}}eth42"
}
