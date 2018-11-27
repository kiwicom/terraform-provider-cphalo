resource "cphalo_firewall_policy" "fw" {
  name = "complicated firewall"
  description = "awesome"
  ignore_forwarding_rules = true
  shared = true
}
