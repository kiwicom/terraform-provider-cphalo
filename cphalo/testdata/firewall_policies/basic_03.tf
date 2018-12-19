resource "cphalo_firewall_policy" "fw" {
  name = "{{.Prefix}}tf_acc_fw_policy_changed"
  description = "awesome"
  ignore_forwarding_rules = true
  shared = true

  rule {
    chain = "OUTPUT"
    action = "DROP"
    connection_states = "NEW"
    position = 1
  }

  rule {
    chain = "INPUT"
    action = "DROP"
    connection_states = "NEW"
    position = 1
  }
}
