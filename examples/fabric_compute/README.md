# fabric compute example

**__NOTE: THIS RESOURCE TYPE MUST BE IMPORTED. ATTEMPTING TO CREATE IT WILL ERROR__ **

These resources (fabric computes) are discovered in vRA(C) as part of creating a Cloud Account and can not therefore be "created" or "destroyed" in a traditional sense.

This is an example on how to import a fabric compute resource into terraform and then manage settings on it.

## Getting Started

There are variables which need to be added to a `terraform.tfvars` file:

* `url` - The base url for API operations
* `refresh_token` - The refresh token for API operations
* `insecure` - Specify whether to validate TLS certificates

To facilitate adding these variables, a sample tfvars file can be copied first:

```shell
cp terraform.tfvars.sample terraform.tfvars
```

This examples assumes a cloud account is already set up.

To import the resource you must find the ID of the fabric compute. There are a couple of way this ID can be aquired:

1. Via API calls
2. Viewing the object in a browser and pulling the ID from the url (ex value: `e2a959f4-7ec5-4941-b532-32b798309feb`)

Once the information is added to `terraform.tfvars` the resource can be imported and the manage via:

```shell
terraform import vra_fabric_compute.this <computeID>
```

If the import is successful output from the command should resemble

```shell
vra_fabric_compute.this: Importing from ID "<computeID>"...
vra_fabric_compute.this: Import prepared!
  Prepared vra_fabric_compute for import
vra_fabric_compute.this: Refreshing state... [id=<computeID>]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform
```

To apply your settings you can now perform an apply comand:

```shell
terraform apply
```
