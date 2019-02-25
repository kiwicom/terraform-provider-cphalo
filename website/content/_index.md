---
title: Provider
type: docs
---

# Terraform CloudPassage Halo Provider `v0.0.0-master`

The CloudPassage Halo provider is used to interact with CloudPassage Halo resources.  
To read more about CloudPassage, consult the [official docs](https://library.cloudpassage.com/help/) (*account is needed for access*).

Since it's a non official Terraform provider, you have to install it manually.

## Installation

### Download binary for you system

**Available systems:**

- [macos 64-bit](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip) ([checksum](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256))
- [linux 32-bit](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_386.zip) ([checksum](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_386.sha256))
- [linux 64-bit](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_amd64.zip) ([checksum](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_amd64.sha256))
- [linux arm](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_arm.zip) ([checksum](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_arm.sha256))
- [windows 32-bit](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_386.zip) ([checksum](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_386.sha256))
- [windows 64-bit](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_amd64.zip) ([checksum](https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_amd64.sha256))

**Download binary**

```bash
wget https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip
```

**Download checksum**

```bash
wget https://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/Releases/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256
```

**Verify integrity**

```bash
shasum -c terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256
```

### Install

**Create plugins directory**

```bash
mkdir -p ~/.terraform.d/plugins
```

**Uncompress**

```bash
unzip terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip
```

**Move plugin**

```bash
mv terraform-provider-cphalo_v0.0.0-master ~/.terraform.d/plugins
```

## Authorization

Provider needs to be configured with the proper credentials before it can be used.

```terraform
provider "cphalo" {
  application_key    = "your-application-key"
  application_secret = "your-application-secret"
}
```

Alternatively, you can provide credentials via environmental variables.

```bash
export CPHALO_APPLICATION_KEY="your-application-key"
export CPHALO_APPLICATION_SECRET="your-application-secret"
```

```terraform
provider "cphalo" {}
```

## Example usage

```terraform
provider "cphalo" {
  application_key    = "your-application-key"
  application_secret = "your-application-secret"
}

resource "cphalo_server_group" "docs_group" {
  name = "docs_group"
}
```

You can check for more examples in the [source code](https://gitlab.com/kiwicom/terraform-provider-cphalo/tree/master/examples).
