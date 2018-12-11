resource "cphalo_csp_account" "tf_examples_basic_aws_account" {
  role_arn = "${aws_iam_role.tf_examples_basic_cloudpassage_role.arn}"
  external_id = "${var.cphalo_external_id}"
  group_id = "${var.cphalo_root_group}"
  account_display_name = "tf_examples_basic_aws_account"
}
