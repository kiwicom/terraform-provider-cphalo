resource "cphalo_server_group" "root_group" {
  name = "changed_name"
  tag = "added_tag"
  description = "and added some interesting description"
  alert_profile_ids = [
    "${data.cphalo_alert_profile.test_alert_profile.id}"
  ]
}

// this alert profile has to be created manually on cloudpassage GUI
data "cphalo_alert_profile" "test_alert_profile" {
  name = "test alert"
}

output "alert_id" {
  value = "${data.cphalo_alert_profile.test_alert_profile.id}"
}
