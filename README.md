# Terraform Provider for VMware vRealize Automation

A self-contained deployable integration between Terraform and VMware vRealize Automation which allows Terraform users to request and provision IaaS resources such as machines, networks, load balancers, along with the configuration of cloud accounts, zones, and projects. This provider supports both vRealize Automation Cloud (SaaS) and vRealize Automation 8 (on-premises). 

> Note: There is a separate provider available for [vRealize Automation 7.x](https://github.com/terraform-providers/terraform-provider-vra7).

## Requirements

* [Terraform 0.12+](https://www.terraform.io/downloads.html)
* [Go 1.16](https://golang.org/dl/) (to build the provider plugin)

## Using the Provider

The [Terraform Provider for VMware vRealize Automation](https://registry.terraform.io/providers/vmware/vra/latest) is a verified provider. Verified providers are owned and maintained by members of the HashiCorp Technology Partner Program. HashiCorp verifies the authenticity of the publisher and the providers are listed on the Terraform Registry with a verified tier label.

To use a released version of the Terraform provider in your environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider from the Terraform Registry. 

See [Installing the Terraform Provider for VMware vRealize Automation](docs/install_provider.md) for additional instructions on automated and manual installation methods and how to control the provider version.

For either installation method, documentation about the provider configuration, resources, and data sources can be found on the [provider page](https://registry.terraform.io/providers/vmware/vra/latest/docs) on the Terraform Registry.

The provider accepts either a `refresh_token` or an `access_token` to interact with the vRealize Automation API, but not both at the same time. 

* For more information on obtaining a `refresh_token` for the provider configuration the provider, see [Get Your Refresh Token for the vRealize Automation API](docs/refresh_token.md).

* For more information on obtaining an `access_token` for the provider configuration, see [Get Your Access Token for the vRealize Automation API](https://code.vmware.com/docs/14701/vrealize-automation-8-6-api-programming-guide/GUID-AC1E4407-6139-412A-B4AA-1F102942EA94.html) on VMware {code}.

Examples on the use of the provider configuration, resources, and data sources can be found in the project's  `examples` directory.

## Upgrading the Provider

The provider does not upgrade automatically. After each new release, you can run the following command to upgrade the provider: 

```bash
terraform init -upgrade
```

## Contributing

The project team welcomes contributions from the community. Before you start working with terraform-provider-vra, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

The Terraform Provider for VMware vRealize Automation is available under the [Mozilla Public License, version 2.0 license](LICENSE).
