resource "cphalo_firewall_zone" "zone" {
  name = "tf_acc_fw_zone"
  ip_address = "2.2.2.2"
  description = "fw zone"
}

