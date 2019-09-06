# cloud\_account\_gcp example

This is an example on how to setup a GCP cloud account in vRealize Automation(vRA) along with regions.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, and the remaining are credentials for connecting to the GCP cloud account.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `client_email` - GCP Client email
* `private_key_id` - GCP Private key ID
* `private_key` - GCP Private key
* `project_id` - GCP Project ID

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
