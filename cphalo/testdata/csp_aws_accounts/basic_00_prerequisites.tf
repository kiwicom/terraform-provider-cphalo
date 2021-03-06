variable "aws_access_key" {
  default = "add-me-to-env"
}

variable "aws_secret_key" {
  default = "add-me-to-env"
}

variable "aws_region" {
  default = "eu-west-2"
}

provider "aws" {
  region     = var.aws_region
  access_key = "{{.AWSAccessKey}}"
  secret_key = "{{.AWSSecretKey}}"
}

variable "cphalo_service_id" {
  default = "{{.CPHaloServiceID}}"
}
variable "cphalo_root_group" {
  default = "{{.CPHaloRootGroup}}"
}

variable "cphalo_external_id" {
  type    = "string"
  default = "{{.Prefix}}this-is-some-id-for-tf-cphalo-testacc"
}

resource "aws_iam_role" "tf_testacc_cloudpassage_role" {
  name               = "{{.Prefix}}tf_testacc_cloudpassage_role"
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
  name = "{{.Prefix}}tf_testacc_cloudpassage_service_policy"
  policy = file("testdata/csp_aws_accounts/aws_cphalo_policy.json")
}

resource "aws_iam_policy_attachment" "tf_testacc_cloudpassage_role_attach" {
  name = "{{.Prefix}}tf_testacc_cloudpassage_role_attach"
  roles = [
    aws_iam_role.tf_testacc_cloudpassage_role.name
  ]
  policy_arn = aws_iam_policy.tf_testacc_cloudpassage_service_policy.arn
}
