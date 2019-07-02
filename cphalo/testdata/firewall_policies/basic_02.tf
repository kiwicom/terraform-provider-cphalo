resource "cphalo_firewall_policy" "fw" {
  name        = "{{.Prefix}}tf_acc_fw_policy_changed"
  description = "awesome"

  rule {
    chain             = "OUTPUT"
    action            = "DROP"
    connection_states = "NEW"
    position          = 1
  }

  rule {
    chain             = "INPUT"
    action            = "DROP"
    connection_states = "NEW"
    position          = 1
  }

  rule {
    chain             = "INPUT"
    action            = "DROP"
    connection_states = "ESTABLISHED"
    position          = 2
  }

  rule {
    chain             = "INPUT"
    action            = "DROP"
    connection_states = "NEW, ESTABLISHED"
    position          = 3
  }
}
