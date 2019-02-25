# Terraform provider for CloudPassage Halo

[![coverage report](https://gitlab.com/kiwicom/terraform-provider-cphalo/badges/master/pipeline.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/pipelines)
[![pipeline status](https://gitlab.com/kiwicom/terraform-provider-cphalo/badges/master/coverage.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/commits/master)
[![mit license](https://img.shields.io/badge/license-MIT-green.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/blob/master/LICENSE)
[![go report](https://goreportcard.com/badge/gitlab.com/kiwicom/terraform-provider-cphalo)](https://goreportcard.com/report/gitlab.com/kiwicom/terraform-provider-cphalo)
[![go doc](https://godoc.org/gitlab.com/kiwicom/terraform-provider-cphalo?status.svg)](https://godoc.org/gitlab.com/kiwicom/terraform-provider-cphalo)
[![contribute](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](https://gitlab.com/kiwicom/terraform-provider-cphalo/forks/new)

- Website: TODO
- API Docs: https://library.cloudpassage.com/help/cloudpassage-api-documentation

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Building the Provider

Clone repository

```sh
$ git clone git@gitlab.com:kiwicom/terraform-provider-cphalo.git
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-cphalo
$ make build
```

*You can then find binary in `bin/current_system`.*

## Running tests

The Terraform Provider has acceptance tests, these can run against CloudPassage Halo service. Credentials are required.  
For more information on getting the credentials, consult the [official docs](https://library.cloudpassage.com/help/) of CloudPassage Halo.

AWS credentials are also needed, since tests need access to AWS to properly test `CSP AWS account` integration.

**Create `.env` file**

```bash
make .env
```

**Fill `.env` with your credentials**

**And run tests**

```bash
make testacc
```

## Contributing

Contributions are always welcome. Pull requests have to pass the following checks:

- lint `make lint`
- tests `make test`
- website `make website-build`
- acceptance tests `make testacc`
- build `make release`

For more information on existing tools use `make help`.
