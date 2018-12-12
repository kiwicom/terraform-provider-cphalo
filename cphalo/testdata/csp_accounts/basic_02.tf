resource "cphalo_csp_account" "main_aws_account" {
  role_arn = "${aws_iam_role.tf_testacc_cloudpassage_role.arn}"
  external_id = "${var.cphalo_external_id}"
  group_id = "${var.cphalo_root_group}"
  account_display_name = "tf_aws_testacc_basic_02"
}