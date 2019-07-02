resource "cphalo_server_group" "root_group" {
  name = "{{.Prefix}}root group"
}

resource "cphalo_firewall_policy" "firewall_policy" {
  name   = "{{.Prefix}}tf_acc_fw_policy"
  shared = true

  rule {
    chain             = "INPUT"
    action            = "ACCEPT"
    connection_states = "ANY"
    position          = 1

    firewall_source {
      id   = cphalo_server_group.root_group.id
      kind = "Group"
    }
  }
}
