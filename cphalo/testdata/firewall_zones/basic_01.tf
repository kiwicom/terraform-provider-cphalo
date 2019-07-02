resource "cphalo_firewall_zone" "zone" {
  name       = "{{.Prefix}}tf_acc_fw_zone"
  ip_address = [
    "1.1.1.1",
    "2.2.2.2",
  ]
}

data "cphalo_firewall_zone" "tf_acc_any_zone" {
  name = "any"
}

output "tf_acc_any_zone_ip_address" {
  value = data.cphalo_firewall_zone.tf_acc_any_zone.ip_address
}

output "tf_acc_any_zone_description" {
  value = data.cphalo_firewall_zone.tf_acc_any_zone.description
}
