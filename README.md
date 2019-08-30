# Terraform provider for VMware vRealize Automation
[![Build Status](https://travis-ci.org/vmware/terraform-provider-vra.svg?branch=master)](https://travis-ci.org/vmware/terraform-provider-vra)


Introduction
------------

A self-contained deployable integration between Terraform and VMware vRealize Automation (vRA) which allows Terraform users to request and provision vRA IaaS resources such as machine, network, load_balancer, along with initial setup of cloud accounts, zones, and projects.

Requirements
------------

To get the vra plugin up and running you need the following things.
* [Terraform 0.12 or above](https://www.terraform.io/downloads.html)
* [Go Language 1.12 or above](https://golang.org/dl/)

Using the provider
----------------------

There are some one time setup needed to begin using the IaaS resources. This may
already be done by administrators or can be done via terraform as well.

### Cloud Accounts

Because vRA works across multiple clouds, a cloud account with credentials must
first be setup. Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](examples/cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](examples/cloud_account_azure/README.md)
* Setup [cloud\_account\_vsphere](examples/cloud_account_vsphere/README.md)

### Zones

While the cloud account examples included setting up zones, here is an example
to setup a zone:

* Setup [zone](examples/zone/README.md)


### Projects

While the cloud account examples included setting up a project, here is an example
to setup a project:

* Setup [project](examples/project/README.md)

### Flavor mappings

### Image mappings

### Machine


Upgrading the provider
----------------------

The vra provider doesn't upgrade automatically once you've started using it. After a new release you can run 

```bash
terraform init -upgrade
```

## Execution
These are the Terraform commands that can be used for the vRA plugin:
* `terraform init` - The init command is used to initialize a working directory containing Terraform configuration files.
* `terraform plan` - Plan command shows plan for resources like how many resources will be provisioned and how many will be destroyed.
* `terraform apply` - apply is responsible to execute actual calls to provision resources.
* `terraform refresh` - By using the refresh command you can check the status of the request.
* `terraform show` - show will set a console output for resource configuration and request status.
* `terraform destroy` - destroy command will destroy all the  resources present in terraform configuration file.

Navigate to the location where `main.tf` and binary are placed and use the above commands as needed.

## Contributing

The terraform-provider-vra project team welcomes contributions from the community. Before you start working with terraform-provider-vra, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

terraform-provider-vra is available under the [Mozilla Public License, version 2.0 license](LICENSE).
