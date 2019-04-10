resource "cphalo_firewall_service" "svc" {
  name = "{{.Prefix}}tf_acc_custom_ssh"
  protocol = "TCP"
  port = "2222"
}

resource "cphalo_firewall_service" "pingpong" {
  name = "{{.Prefix}}tf_acc_ping"
  protocol = "ICMP"
}

data "cphalo_firewall_service" "tf_acc_http_service" {
  name = "http"
}

output "tf_acc_http_service_protocol" {
  value = "${data.cphalo_firewall_service.tf_acc_http_service.protocol}"
}

output "tf_acc_http_service_port" {
  value = "${data.cphalo_firewall_service.tf_acc_http_service.port}"
}
