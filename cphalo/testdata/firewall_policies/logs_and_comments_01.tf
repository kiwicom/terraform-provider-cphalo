resource "cphalo_firewall_policy" "fw_logging_policy" {
  name = "{{.Prefix}}tf_acc_fw_logging_policy"

  rule {
    chain = "INPUT"
    action = "ACCEPT"
    connection_states = "ANY"
    position = 1
    log = true
    log_prefix = "tf_acc_test_"
    comment = "tf_acc"
  }
}
