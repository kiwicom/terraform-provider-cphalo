resource "cphalo_firewall_service" "svc" {
  name = "tf_acc_custom_ssh"
  protocol = "TCP"
  port = "2222"
}

resource "cphalo_firewall_service" "pingpong" {
  name = "tf_acc_ping"
  protocol = "ICMP"
}
