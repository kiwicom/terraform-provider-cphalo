# Terraform Provider for CloudPassage Halo

[![coverage report](https://gitlab.com/kiwicom/terraform-provider-cphalo/badges/master/pipeline.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/pipelines)
[![pipeline status](https://gitlab.com/kiwicom/terraform-provider-cphalo/badges/master/coverage.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/commits/master)
[![mit license](https://img.shields.io/badge/license-MIT-green.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/blob/master/LICENSE)
[![go report](https://goreportcard.com/badge/gitlab.com/kiwicom/terraform-provider-cphalo)](https://goreportcard.com/report/gitlab.com/kiwicom/terraform-provider-cphalo)
[![go doc](https://godoc.org/gitlab.com/kiwicom/terraform-provider-cphalo?status.svg)](https://godoc.org/gitlab.com/kiwicom/terraform-provider-cphalo)
[![contribute](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/forks/new)

The Cloudpassage Halo Terraform Provider plugin allows you to configure and manage the security of your virtual servers. It provides control over file integrity monitoring (FIM), firewall automation, vulnerability monitoring, network access control, security event alerting, and assessment.

This GitLab repository runs an acceptance test every night to ensure that this Terraform Provider is working correctly. The `pipeline` status badge, which you can find at the top of the project details page of this repository, indicates whether the current release of the Provider passed this acceptance test or not.

Please find the documentation of this plugin by clicking on the links below:

- Reference Guide: https://kiwicom.gitlab.io/terraform-provider-cphalo/
- API Documentation: https://library.cloudpassage.com/help/cloudpassage-api-documentation

## System Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

**Note**: If you're using Terraform version **0.11.xx or older**, then use [version 0.1.3](https://gitlab.com/kiwicom/terraform-provider-cphalo/-/tags/v0.1.3) of this Provider.

## Building the Provider

Please follow the instructions below if you want to build the Provider from source code.

1. Clone the repository

```sh
$ git clone git@gitlab.com:kiwicom/terraform-provider-cphalo.git
```

2. Change directory into `terraform-provider-cphalo` and build the provider

```sh
$ cd terraform-provider-cphalo
$ make build
```

3. Find the resulting binary in `bin/current_system`

## Running Tests

Our Terraform Provider contains acceptance tests, which will exercise its code by executing real `plan`, `apply`, `refresh`, and `destroy` life-cycles using your Cloudpassage Halo account, and evaluate your `CSP AWS account` integration.

Running this acceptance test, therefore, requires you to have AWS and Cloudpassage Halo access credentials. Please consult Cloudpassage Halo's [documentation](https://library.cloudpassage.com/help/) on how to obtain these credentials.

#### How to run acceptance tests

First, you will need to enter your CloudPassage Halo and AWS credentials in the environment variables of this Provider.

Fill out these credentials by running `make .env` in your terminal, opening the resulting file, and entering them in the corresponding fields.

Next, you can begin testing by running the `make testacc` command in your terminal.

## Contributing

Contributions are always welcome. Our repository will perform the following checks once your submit your pull request:

- lint `make lint`
- tests `make test`
- website `make website-build`
- acceptance tests `make testacc`
- build `make release`

For more information on existing tools use `make help`.
