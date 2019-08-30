# zone example

This is an example on how to create a compute placement zone within a region of an existing cloud account. 

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, there are credentials for connecting to the vCenter instance, and the name of the data collector already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```
To create a zone, a cloud account must be setup in vRA and the region must be identified within which the zone will be created.

Once the information is added to `terraform.tfvars`, the zones can be created via:

```shell
terraform init
terraform apply
```
