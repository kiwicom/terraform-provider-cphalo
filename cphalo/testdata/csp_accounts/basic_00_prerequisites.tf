provider "aws" {
  region = "eu-west-2"
}

variable "cphalo_service_id" {
  type = "string"
  default = "856192027328"
}

variable "cphalo_external_id" {
  type = "string"
  default = "this-is-some-id-for-tf-cphalo-testacc"
}

variable "cphalo_root_group" {
  type = "string"
  default = "fff04606e97b11e896d9252f8ed31fc8"
}

resource "aws_iam_role" "tf_testacc_cloudpassage_role" {
  name = "tf_testacc_cloudpassage_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${var.cphalo_service_id}:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${var.cphalo_external_id}"
        }
      }
    }
  ]
}
EOF
}

resource "aws_iam_policy" "tf_testacc_cloudpassage_service_policy" {
  name = "tf_testacc_cloudpassage_service_policy"
  policy = "${file("testdata/csp_accounts/aws_cphalo_policy.json")}"
}

resource "aws_iam_policy_attachment" "tf_testacc_cloudpassage_role_attach" {
  name = "tf_testacc_cloudpassage_role_attach"
  roles = [
    "${aws_iam_role.tf_testacc_cloudpassage_role.name}"
  ]
  policy_arn = "${aws_iam_policy.tf_testacc_cloudpassage_service_policy.arn}"
}
