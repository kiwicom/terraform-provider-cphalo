resource "cphalo_firewall_policy" "fw_logging_policy" {
  name = "{{.Prefix}}tf_acc_fw_logging_policy"

  rule {
    chain             = "INPUT"
    action            = "ACCEPT"
    connection_states = "ANY"
    position          = 1
    log               = false
    comment           = "tf_acc_v2"
  }
}
