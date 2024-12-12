<!--
© Broadcom. All Rights Reserved.
The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
SPDX-License-Identifier: BSD-2
-->

<!-- markdownlint-disable first-line-h1 no-inline-html -->

# Terraform Provider for VMware Aria Automation

A self-contained deployable integration between Terraform and VMware Aria Automation which allows Terraform users to request and provision IaaS resources such as machines, networks, load balancers, along with the configuration of cloud accounts, zones, and projects.

This provider supports VMware Aria Automation 8.

## Requirements

- [Terraform 0.13+][terraform-install]

  For general information about Terraform, visit [developer.hashicorp.com][terraform-install] and [the project][terraform-github] on GitHub.

- [Go 1.23.2][golang-install]

  Required if building the provider.

## Using the Provider

The [Terraform Provider for VMware Aria Automation](https://registry.terraform.io/providers/vmware/vra/latest) is a partner provider. Partner providers are owned and maintained by members of the HashiCorp Technology Partner Program. HashiCorp verifies the authenticity of the publisher and the providers are listed on the Terraform Registry with a `Partner` label.

To use a released version of the Terraform provider in your environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider from the Terraform Registry.

Refer to [Installing the Terraform Provider for VMware Aria Automation](docs/install_provider.md) for additional instructions on automated and manual installation methods and how to control the provider version.

For either installation method, documentation about the provider configuration, resources, and data sources can be found on the [provider page](https://registry.terraform.io/providers/vmware/vra/latest/docs) on the Terraform Registry.

The provider accepts either a `refresh_token` or an `access_token` to interact with the product API, but not both at the same time.

For more information on obtaining a `refresh_token` for the provider configuration the provider, refer [Get Your Refresh Token for the VMware Aria Automation API](docs/refresh_token.md).

Examples on the use of the provider configuration, resources, and data sources can be found in the project's `examples` directory.

## Upgrading the Provider

The provider does not upgrade automatically. After each new release, you can run the following command to upgrade the provider:

```bash
terraform init -upgrade
```

## Contributing

The Terraform Provider for VMware Aria Automation is the work of many contributors and the project team appreciates your help!

If you discover a bug or would like to suggest an enhancement, submit [an issue][provider-issues].

If you would like to submit a pull request, please read the [contribution guidelines][provider-contributing] to get started. In case of enhancement or feature contribution, we kindly ask you to open an issue to discuss it beforehand.

## License

© Broadcom. All Rights Reserved.
The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.

The Terraform Provider for VMware Aria Automation is available under the [Mozilla Public License, version 2.0][provider-license] license.

[golang-install]: https://golang.org/doc/install
[provider-contributing]: CONTRIBUTING.md
[provider-issues]: https://github.com/vmware/terraform-provider-vra/issues/new/choose
[provider-license]: LICENSE
[terraform-install]: https://developer.hashicorp.com/terraform/install
[terraform-github]: https://github.com/hashicorp/terraform
