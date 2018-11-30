resource "cphalo_firewall_policy" "fw" {
  name = "awesome firewall policy"
}

resource "cphalo_firewall_rule" "first" {
  parent_id = "${cphalo_firewall_policy.fw.id}"
  chain = "INPUT"
  action = "DROP"
  active = true
  connection_states = "NEW, ESTABLISHED"
  position = 1
}

resource "cphalo_firewall_rule" "third" {
  parent_id = "${cphalo_firewall_policy.fw.id}"
  chain = "OUTPUT"
  action = "DROP"
  active = true
  connection_states = "NEW"
  position = 1
}
