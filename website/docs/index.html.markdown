---
layout: "vra"
page_title: "Provider: VMware vRealize Automation Services"
sidebar_current: "docs-vra-index"
description: |-
  A Terraform provider to work with VMware vRealize Automation Services.
---

# VMware vRealize Automation Provider

Use the VMware vRealize Automation provider to interact with resources supported by
[VMware vRealize Automation][vmware-vrealize-automation] services. This provider can
be used to configure multi-cloud infrastructure components to enable multi-cloud automation services.

[vmware-vrealize-automation]: https://www.vmware.com/products/vrealize-automation.html

Use the navigation on the left to read about the various resources and data
sources supported by the provider.

## Example Usage

The following abridged example demonstrates a current basic usage of the
provider.

```hcl
provider "vra" {
    url = "${var.url}"
    refresh_token = "${var.refresh_token}"
}
```

See the sidebar for usage information on all the resources, which will have
examples specific to their own use cases.

## Argument Reference

The following arguments are used to configure the VMware vRealize Automation Provider:

* `url` - (Required) This is the URL to the VMware vRealize Automation
  Services endpoint. Can also  be specified with the `vRA_URL` environment variable.
* `access_token` - (Optional) This is the access token used to create an API
  refresh token. Can also be specified with the `vRA_ACCESS_TOKEN` environment variable.
* `refresh_token` - (Optional) This is a refresh_token used for API access that
  has been pre-generated. One of `access_token` or `refresh_token` is required.
  Can also be specified with the `vRA_REFRESH_TOKEN` environment variable.

## Bug Reports and Contributing

For more information how how to submit bug reports, feature requests, or
details on how to make your own contributions to the provider, see the Terraform provider for VMware vRealize Automation [project page][tf-vra-project-page].

[tf-vra-project-page]: https://github.com/vmware/terraform-provider-vra


