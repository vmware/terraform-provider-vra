# fabric network vsphere example

**__NOTE: THIS RESOURCE TYPE MUST BE IMPORTED. ATTEMPTING TO CREATE IT WILL ERROR__ **

These resources (fabric networks) are discovered in vRA(C) as part of creating a vSphere Cloud Account and can not therefore be "created" or "destroyed" in a traditional sense.

This is an example on how to import a vsphere fabric network resource into terraform and then manage settings on it.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. There are names of cloud_account, region, zone, project, already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `cidr` - The CIDR specification for the network (ex: 10.0.0.0/24)
* `gateway` - The default gateway for the network
* `domain` - The DNS domain machines will be customized with

To facilitate adding these variables, a sample tfvars file can be copied first:

```shell
cp terraform.tfvars.sample terraform.tfvars
```

This examples assumes a vsphere cloud account is already set up, if not please refer to vsphere cloud account example

* Setup [cloud\_account\_vsphere](../cloud_account_vsphere/README.md)

To import the resource you must find the ID of the fabric network. There are a couple of way this ID can be aquired:

1. Use of the vra_fabric_network data source. ex:  
```
  data "vra_fabric_network" "subnet" {
  filter = "name eq '${var.subnet_name}' and cloudAccountId eq '${var.cloud_account_vsphere_id}' and externalRegionId eq '${var.region_id_vsphere}'"
}
```
2. via API calls 
3. Viewing the object in a browser and pulling the ID from the url. (ex value: bea2cb32-ba08-4876-bb6c-9ce8af6cd90c)

Once the information is added to `terraform.tfvars` the resource can be imported and the manage via:

```shell
terraform import vra_fabric_network_vsphere.simple  <networkID>
```

If the import is successful output from the command should resemble

```shell

vra_fabric_network_vsphere.simple: importing from ID "<networkID>"
vra_fabric_network_vsphere.simple: Import Prepared!
  Prepared vra_fabric_network_vsphere for import
  vra_fabric_network_vsphere.simple: Refreshing state... [id=<networkID>]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

To apply your settings you can now perform an apply comand:

```shell
terraform apply
```
