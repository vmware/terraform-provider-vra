---
layout: "vra"
page_title: "Provider: VMware vRealize Automation"
sidebar_current: "docs-vra-index"
description: |-
  A Terraform provider for VMware vRealize Automation.
---

# VMware vRealize Automation Provider

Use the VMware vRealize Automation provider to interact with resources supported by
[VMware vRealize Automation][vmware-vrealize-automation] services. This provider can be used to configure multi-cloud infrastructure components to enable multi-cloud automation services.

[vmware-vrealize-automation]: https://www.vmware.com/products/vrealize-automation.html

Please use the navigation to the left to read about available data sources and resources.

## Basic Configuration of the Provider

With Terraform 0.13 and later, the `terraform` configuration block should be used in your configurations.

**Example**: Terraform Configuration

```hcl
terraform {
  required_providers {
    vra = {
      source  = "vmware/vra"
    }
  }
  required_version = ">= 0.13"
}

provider "vra" {
  // Configuration Options
}
```
In order to use the provider you must configure the provider to communicate with the vRealize Automation endpoint. The provider configuration requires the `url` and `refresh_token` or `access_token`.

The provider also can accept both signed and self-signed server certificates. It is recommended that in production environments you only use certificates signed by a certificate authority. Setting the `insecure` parameter to `true` will direct the Terraform to skip certificate verification. This is **not recommended** in production deployments. It is recommended that you use a trusted connection using certificates signed by a certificate authority.

**Example**: Configuration with Credentials

```hcl
provider "vra" {
  url           = var.vra_url
  refresh_token = var.vra_refresh_token
  insecure      = false
}
```
**Example**: Setting Environment Variables

```shell
export VRA_URL="https://api.mgmt.cloud.vmware.com"
export VRA_REFRESH_TOKEN="***********************"
```

Documentation about the provider resources and data sources can be found within the sidebar, which has examples specific to their use. Additional examples on the use of the provider configuration, resources, and data sources can be found in the `examples` directory of the [project][tf-vra-project-page].

Note that in all of the examples you will need to update attributes - such as `url`, `refresh_token` or `access_token`, and `insecure` - to match your environment.

## Argument Reference

The following arguments are used to configure the Terraform Provider for VMware vRealize Automation:

* `url` - (Required) This is the URL to the VMware vRealize Automation endpoint. Can also be specified with the `VRA_URL` environment variable.
* `access_token` - (Optional) This is the access token used to create an API refresh token. Can also be specified with the `VRA_ACCESS_TOKEN` environment variable.
* `refresh_token` - (Optional) This is a refresh_token used for API access that has been pre-generated. One of `access_token` or `refresh_token` is required. Can also be specified with the `VRA_REFRESH_TOKEN` environment variable.
* `insecure` - (Optional) This specifies whether if the TLS certificates are validated. Can also be specified with the `VRA7_INSECURE` environment variable.

## Bug Reports and Contributing

For more information how how to submit bug reports, feature requests, or details on how to make your own contributions to the provider, see the Terraform provider for VMware vRealize Automation [project][tf-vra-project-page].

[tf-vra-project-page]: https://github.com/vmware/terraform-provider-vra
