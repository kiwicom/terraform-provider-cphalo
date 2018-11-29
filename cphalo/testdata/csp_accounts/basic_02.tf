resource "cphalo_csp_account" "main_aws_account" {
  role_arn = "arn:aws:iam::782106534067:role/CloudPassage-Service-Role" // TODO: move this outside
  external_id = "this-is-external-id-1" // TODO: move this outside
  group_id = "fff04606e97b11e896d9252f8ed31fc8"
  account_display_name = "tf_aws_testacc_basic_02"
}
