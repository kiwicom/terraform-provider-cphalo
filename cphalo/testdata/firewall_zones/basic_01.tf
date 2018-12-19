resource "cphalo_firewall_zone" "zone" {
  name = "{{.Prefix}}tf_acc_fw_zone"
  ip_address = "1.1.1.1"
}

