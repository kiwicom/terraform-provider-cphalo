resource "cphalo_firewall_service" "svc" {
  name = "custom ssh"
  protocol = "TCP"
  port = "2223"
}
