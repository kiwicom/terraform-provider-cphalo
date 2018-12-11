provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region = "${var.aws_region}"
}

resource "aws_iam_role" "tf_examples_basic_cloudpassage_role" {
  name = "tf_examples_basic_cloudpassage_role"
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
  // after role is created we have to wait a little bit, so cloudpassage can access it
  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "aws_iam_policy" "tf_examples_basic_cloudpassage_service_policy" {
  name = "tf_examples_basic_cloudpassage_service_policy"
  policy = "${file("aws_cphalo_policy.json")}"
}

resource "aws_iam_policy_attachment" "tf_examples_basic_cloudpassage_role_attach" {
  name = "tf_examples_basic_cloudpassage_role_attach"
  roles = [
    "${aws_iam_role.tf_examples_basic_cloudpassage_role.name}"
  ]
  policy_arn = "${aws_iam_policy.tf_examples_basic_cloudpassage_service_policy.arn}"
}
