# Terraform provider for CloudPassage

[![coverage report](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/badges/master/coverage.svg)](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/commits/master)
[![pipeline status](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/badges/master/pipeline.svg)](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/commits/master)

- Website: https://www.terraform.io
- API: https://library.cloudpassage.com/help/cloudpassage-api-documentation

Currently only linux is being supported. Windows support is waiting for your kind PR.

## Installation

**For MacOS (`darwin`)**.
For other platforms, replace `darwin_amd64` in the file names with appropriate platform. *Check links below*.

```bash
# download binary and checksum file
wget http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip
wget http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256

# verify integrity
shasum -c terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256

# install plugin
mkdir -p ~/.terraform.d/plugins
unzip terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip
mv terraform-provider-cphalo_v0.0.0-master ~/.terraform.d/plugins

# cleanup
rm terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip readme.md terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256
```

A `cphalo_client` command line utility is distributed in the `zip`, to allow easier listing of resources currently supported by this provider.

### Supported platforms

- [macos 64-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_darwin_amd64.sha256))
- [linux 32-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_386.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_386.sha256))
- [linux 64-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_amd64.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_amd64.sha256))
- [linux arm](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_arm.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_linux_arm.sha256))
- [windows 32-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_386.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_386.sha256))
- [windows 64-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_amd64.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/v0.0.0-master/terraform-provider-cphalo_v0.0.0-master_windows_amd64.sha256))

## Example code

Examples can be found in [examples](examples/basic) directory.

### Authorization

#### Via Terraform files

```hcl-terraform
provider "cphalo" {
  application_key = "your-application-key"
  application_secret = "your-application-secret"
}
```

#### Via env

```bash
export CP_APPLICATION_KEY="your-application-key"
export CP_APPLICATION_SECRET="your-application-secret"
```

## Endpoint implementation status:

- [ ] **Agent Upgrades**

- [ ] **Halo Connectors**
    - *read-only*

- [ ] **Containers**
    - *read-only*

- [ ] **Container Batch**

- [ ] **Container Events**
    - *read-only*

- [ ] **Container Images**
    - *read-only*

- [ ] **Container Image Issues**
    - *read-only*

- [ ] **Container Image Registries**

- [ ] **Container Image Repository**
    - *read-only*

- [ ] **Container Image Summaries**
    - *read-only*

- [ ] **Container Processes**
    - *read-only*

- [ ] **Container Software Package**
    - *read-only*

- [x] **CSP Accounts**

- [ ] **CSP Resources**
    - *read-only*

- [ ] **CSP Findings**
    - *read-only*

- [ ] **CSP Scanner Settings**

- [x] **Server Groups**

- [ ] **Servers**
    - *read-only*
    - [ ] as data resource

- [ ] **Server Accounts**

- [ ] **Server Commands**
    - *read-only*

- [ ] **Server Connections**
    - *read-only*

- [ ] **Server Processes**
    - *read-only*

- [ ] **Server Scans**
    - *read-only*

- [ ] **Server Local Firewalls**
    - *read-only*

- [ ] **Server Firewall Logs**
    - *read-only*

- [ ] **Local User Accounts**
    - *read-only*

- [ ] **Local User Groups**
    - *read-only*

- [ ] **Scan History**
    - *read-only*

- [ ] **Issues**
    - *read-only*

- [ ] **Configuration Policies**

- [ ] **File Integrity Policy**

- [ ] **File Integrity Baselines**

- [ ] **CVE Details**
    - *read-only*

- [ ] **CVE Exceptions**

- [x] **Firewall Policies**

- [x] **Firewall Rules**

- [x] **Firewall Interfaces**

- [x] **Firewall Services**

- [x] **Firewall Zones**

- [ ] **Log-Based Intrusion Detection Policies**

- [ ] **Special Events Policies**
    - *read-only*

- [ ] **Events**
    - *read-only*

- [x] **Alert Proflies**
    - *read-only*

- [ ] **Saved Searches**

- [ ] **Global Scanner Settings**

- [ ] **Group Scanner Settings**

- [ ] **System Announcements**
    - *read-only*
