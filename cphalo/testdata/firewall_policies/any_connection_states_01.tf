resource "cphalo_firewall_policy" "fw_policy" {
  name = "{{.Prefix}}tf_acc_any_conn_states_fw_policy"

  rule {
    chain = "INPUT"
    action = "ACCEPT"
    connection_states = "ANY"
    position = 1
  }
}
