resource "cphalo_firewall_policy" "tf_acc_data_source_policy" {
  name = "tf_acc_data_source_policy"

  rule {
    chain = "INPUT"
    action = "ACCEPT"
    connection_states = "NEW, ESTABLISHED"
    position = 1

    firewall_interface = "${data.cphalo_firewall_interface.tf_acc_eth0_interface.id}"
    firewall_service = "${data.cphalo_firewall_service.tf_acc_http_service.id}"

    firewall_source = "${data.cphalo_firewall_zone.tf_acc_any_zone.id}"
  }
}

data "cphalo_firewall_zone" "tf_acc_any_zone" {
  name = "any"
}

data "cphalo_firewall_service" "tf_acc_http_service" {
  name = "http"
}

data "cphalo_firewall_interface" "tf_acc_eth0_interface" {
  name = "eth0"
}
