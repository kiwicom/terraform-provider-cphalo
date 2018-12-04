resource "cphalo_firewall_service" "svc" {
  name = "custom ssh"
  protocol = "TCP"
  port = "2222"
}

resource "cphalo_firewall_service" "pingpong" {
  name = "ping"
  protocol = "ICMP"
}
