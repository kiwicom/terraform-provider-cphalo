# Terraform provider for CloudPassage

[![coverage report](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/badges/master/coverage.svg)](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/commits/master)
[![pipeline status](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/badges/master/pipeline.svg)](https://gitlab.skypicker.com/devops/terraform-provider-cphalo/commits/master)

- Website: https://www.terraform.io
- API: https://library.cloudpassage.com/help/cloudpassage-api-documentation

Currently only linux is being supported. Windows support is waiting for your kind PR.

## Installation

### Download example

Download zipped binary from http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_darwin_amd64.zip.  
Download checksum file from http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_darwin_amd64.sha256.

To ensure data integrity, run:
```bash
shasum -c terraform-provider-cphalo_alpha_darwin_amd64.sha256
```

### Supported platforms

- [macos 64-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_darwin_amd64.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_darwin_amd64.sha256))
- [linux 32-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_linux_386.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_linux_386.sha256))
- [linux 64-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_linux_amd64.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_linux_amd64.sha256))
- [linux arm](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_linux_arm.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_linux_arm.sha256))
- [windows 32-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_windows_386.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_windows_386.sha256))
- [windows 64-bit](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_windows_amd64.zip) ([checksum](http://s3.eu-west-1.amazonaws.com/kw-terraform-providers/cphalo/alpha/terraform-provider-cphalo_alpha_windows_amd64.sha256))

## Development

### Client

#### Configuration

Create `.env` file:

```bash
make .env
```

Edit credentials in `.env` file.

#### Run sandbox

```bash
make run-sandbox
```

#### Tests

```bash
make test
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

- [ ] **CSP Accounts**
    - [x] basic implementation
    - [ ] implement all properties / methods

- [ ] **CSP Resources**
    - *read-only*

- [ ] **CSP Findings**
    - *read-only*

- [ ] **CSP Scanner Settings**

- [ ] **Server Groups**
    - [x] basic implementation
    - [ ] implement all properties / methods

- [ ] **Servers**
    - *read-only*
    - [x] basic implementation as data resource
    - [ ] implement all properties / methods

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

- [ ] **Firewall Policies**
    - [x] basic implementation
    - [ ] implement all properties / methods

- [ ] **Firewall Rules**
    - [x] basic implementation
    - [ ] implement all properties / methods

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
    - [x] as data source

- [ ] **Saved Searches**

- [ ] **Global Scanner Settings**

- [ ] **Group Scanner Settings**

- [ ] **System Announcements**
    - *read-only*
