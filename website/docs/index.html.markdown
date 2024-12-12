---
layout: "vra"
page_title: "Provider: VMware Aria Automation"
sidebar_current: "docs-vra-index"
description: |-
  A Terraform Provider for VMware Aria Automation.
---

# VMware Aria Automation Provider

Use this Terraform provider to interact with resources supported by [VMware Aria Automation][vmware-aria-automation] services, enabling you to deliver a self-service cloud consumption experience with VMware Cloud Foundation.

[vmware-aria-automation]: https://www.vmware.com/products/aria-automation.html

Please use the navigation to the left to read about available data sources and resources.

## Basic Configuration of the Provider

With Terraform 0.13 and later, the `terraform` configuration block should be used in your configurations.

**Example**: Terraform Configuration

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

provider "vra" {
  // Configuration Options
}
```

In order to use the provider you must configure the provider to communicate with the VMware Aria Automation endpoint. The provider configuration requires the `url` and `refresh_token` or `access_token`.

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
export VRA_URL="https://cloud.example.com"
export VRA_REFRESH_TOKEN="***********************"
```

Documentation about the provider resources and data sources can be found within the sidebar, which has examples specific to their use. Additional examples on the use of the provider configuration, resources, and data sources can be found in the `examples` directory of the [project][project-page].

Note that in all of the examples you will need to update attributes - such as `url`, `refresh_token` or `access_token`, and `insecure` - to match your environment.

## Argument Reference

The following arguments are used to configure the Terraform Provider for VMware Aria Automation:

- `url` - (Required) This is the URL to the VMware Aria Automation endpoint. Can also be specified with the `VRA_URL` environment variable.
- `organization` - (Optional) The name of the organization. Required when using VCF Automation, otherwise, this parameter is ignored. Can also be specified with the `VCFA_ORGANIZATION` environment variable.
- `access_token` - (Optional) This is the access token used to create an API refresh token. Can also be specified with the `VRA_ACCESS_TOKEN` environment variable.
- `refresh_token` - (Optional) This is a refresh token used for API access that has been pre-generated. One of `access_token` or `refresh_token` is required. Can also be specified with the `VRA_REFRESH_TOKEN` environment variable.
- `insecure` - (Optional) This specifies whether if the TLS certificates are validated. Can also be specified with the `VRA_INSECURE` environment variable.
- `reauthorize_timeout` - (Optional) This specifies the timeout for how often to reauthorize the access token. Can also be specified with the `VRA_REAUTHORIZE_TIMEOUT` environment variable.
- `api_timeout` - (Optional) This specifies the timeout in seconds for API operations. Can also be specified with the `VRA_API_TIMEOUT` environment variable.

## Bug Reports and Contributing

For more information how how to submit bug reports, feature requests, or details on how to make your own contributions to the provider, see the Terraform provider for VMware Aria Automation [project][project-page].

[project-page]: https://github.com/vmware/terraform-provider-vra
