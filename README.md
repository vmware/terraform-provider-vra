# Terraform Provider for VMware vRealize Automation

![License](https://img.shields.io/github/license/vmware/terraform-provider-vra?style=for-the-badge) ![Release](https://img.shields.io/github/release/vmware/terraform-provider-vra?style=for-the-badge)

A self-contained deployable integration between Terraform and VMware vRealize Automation which allows Terraform users to request and provision IaaS resources such as machines, networks, load balancers, along with the configuration of cloud accounts, zones, and projects. This provider supports both vRealize Automation Cloud (SaaS) and vRealize Automation 8 (on-premises). 

> Note: There is a separate provider available for [vRealize Automation 7.x](https://github.com/terraform-providers/terraform-provider-vra7).

## Requirements

![Terraform](https://img.shields.io/badge/Terraform-0.12%2B-blue?style=for-the-badge&logo=terraform) ![Go](https://img.shields.io/github/go-mod/go-version/vmware/terraform-provider-vra?style=for-the-badge&logo=go)

* [Terraform 0.12+](https://www.terraform.io/downloads.html)
* [Go 1.16](https://golang.org/dl/) (to build the provider plugin)

## Using the Provider

The [Terraform Provider for VMware vRealize Automation](https://registry.terraform.io/providers/vmware/vra/latest) is a verified provider. Verified providers are owned and maintained by members of the HashiCorp Technology Partner Program. HashiCorp verifies the authenticity of the publisher and the providers are listed on the Terraform Registry with a verified tier label.

To use a released version of the Terraform provider in your environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider from the Terraform Registry. 

See [Installing the Terraform Provider for VMware vRealize Automation](docs/install_provider.md) for additional instructions on automated and manual installation methods. 

For either installation method, documentation about the provider configuration, resources, and data sources can be found on the [provider page](https://registry.terraform.io/providers/vmware/vra/latest/docs) on the Terraform Registry.

Examples on the use of the provider configuration, resources, and data sources can be found in the project's  `examples` directory.

## Controlling the Provider Version

To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

![Terraform](https://img.shields.io/badge/Terraform-0.13%2B-blue?style=for-the-badge&logo=terraform)

Providers listed on the Terraform Registry can be automatically downloaded when initializing a working directory with `terraform init`. The Terraform configuration block is used to configure some behaviors of Terraform itself, such as the Terraform version and the required providers and versions.

**Example**: A Terraform configuration block.

```hcl
terraform {
  required_providers {
    vra = {
      source  = "vmware/vra"
    }
  }
  required_version = ">= 0.13"
}
```

You can use `version` locking and operators to require specific versions of the provider. 

```hcl
terraform {
  required_providers {
    vra = {
      source  = "vmware/vra"
      version = ">= x.y.z"
    }
  }
  required_version = ">= 0.13"
}
```

[Read more](https://www.terraform.io/docs/configuration/providers.html#provider-versions) on Terraform provider configuration.

![Terraform](https://img.shields.io/badge/Terraform-0.12-blue?style=for-the-badge&logo=terraform)

The version meta-argument specifies a version constraint for a provider, and works the same way as the version argument in a `required_providers` block for the Terraform configuration block. The version constraint in a provider configuration is only used if the `required_providers` is not included for the provider in the Terraform configuration block.

```hcl
provider "vra" {
  version = ">= x.y.z"
  ...
}
```

> Important: The version argument in provider configurations is deprecated. In Terraform 0.13 and later. Version constraints should always be declared in the Terraform block using the `required_providers`. 

## Upgrading the Provider

The provider does not upgrade automatically. After each new release, you can run the following command to upgrade the provider: 

```bash
terraform init -upgrade
```

## Contributing

The project team welcomes contributions from the community. Before you start working with terraform-provider-vra, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Copyright 2019-2021 VMware, Inc.

The Terraform Provider for VMware vRealize Automation is available under the [Mozilla Public License, version 2.0 license](LICENSE).
