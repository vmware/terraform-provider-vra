---
layout: "vra"
page_title: "Provider: VMware Cloud Automation Services"
sidebar_current: "docs-vra-index"
description: |-
  A Terraform provider to work with VMware Cloud Automation Services.
---

# VMware Cloud Automation Services Provider

The VMware Cloud Automation Services provider gives Terraform the ability to work
with [VMware Cloud Assembly][vmware-cloud-assembly]. This provider can be used to
deploy multi-cloud infrastructure components.

[vmware-cloud-assembly]: https://cloud.vmware.com/cloud-assembly

Use the navigation on the left to read about the various resources and data
sources supported by the provider.

## Example Usage

The following abridged example demonstrates a current basic usage of the
provider.

[tf-vra7-deployment]: /docs/providers/vra7/r/deployment.html

```hcl
provider "vra" {
    url = "${var.url}"
    refresh_token = "${var.refresh_token}"
}
```

See the sidebar for usage information on all the resources, which will have
examples specific to their own use cases.

## Argument Reference

The following arguments are used to configure the VMware vRA7 Provider:

* `url` - (Required) This is the URL to the VMware Cloud Automation
  Services endpoint. Can also  be specified with the `VRA_URL` environment variable.
* `access_token` - (Optional) This is the access token used to create an API
  refresh token. Can also be specified with the `VRA_ACCESS_TOKEN` environment variable.
* `refresh_token` - (Optional) This is a refresh_token used for API access that
  has been pre-generated. One of `access_token` or `refresh_token` is required.
  Can also be specified with the `VRA_REFRESH_TOKEN` environment variable.
* `insecure` - (Optional) This boolean allow you to not check server certificate.

### Debugging options

~> **NOTE:** The following options can leak sensitive data and should only be
enabled when instructed to do so by HashiCorp for the purposes of
troubleshooting issues with the provider, or when attempting to perform your
own troubleshooting. Use them at your own risk and do not leave them enabled!

* ***Add info here on debuggings ***

## Bug Reports and Contributing

For more information how how to submit bug reports, feature requests, or
details on how to make your own contributions to the provider, see the vRA7
provider [project page][tf-vra7-project-page].

[tf-vra7-project-page]: https://github.com/vmware/terraform-provider-vra7


