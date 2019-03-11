resource "cphalo_firewall_zone" "zone" {
  name        = "{{.Prefix}}tf_acc_fw_zone"
  description = "fw zone"
  ip_address  = [
    "3.3.3.3",
    "4.4.4.4",
  ]
}

