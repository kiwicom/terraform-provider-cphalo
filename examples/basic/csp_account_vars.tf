variable "aws_access_key" {}
variable "aws_secret_key" {}

variable "aws_region" {
  default = "eu-west-2"
}

variable "cphalo_service_id" {}
variable "cphalo_root_group" {}

variable "cphalo_external_id" {
  type = "string"
  default = "this-is-some-id-for-tf-cphalo-testacc"
}
