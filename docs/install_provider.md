# Installing the vRA Terraform provider

The provider is currently considered to be a third-party provider and thus won't be automatically downloaded by ```terraform```, which means you have to either install or build it yourself. The provider is made available in a pre-built binary version via the GitHub repository in the form of releases. This document will walk you through installing a released version of the provider. (The following snippets will use v0.3.5, but you will need to update the version as necessary)

## Downloading the provider

 The most recent version of the provider can be found at https://github.com/vmware/terraform-provider-vra/releases/latest

![example release](images/provider_release_example.png)

You can download the appropriate version of the provider for your OS via either your browser or the commandline using a tool like curl or wget.

### Linux

Create a terraform plugins directory with your hardware platform subdirectory. Typically for 64bit Linux this will be in ```~/.terraform.d/plugins/linux_amd64``` on non-Windows platforms.

```bash
mkdir -p ~/.terraform.d/plugins/linux_amd64 
```

Download the plugin (via a browser or command line)

 ```bash
 RELEASE=0.3.5
 wget -q https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra-linux_amd64-v${RELEASE}.tgz
 ```

Untar/unzip the plugin

```shell
tar xvf terraform-provider-vra-linux_amd64-v${RELEASE}.tgz
```

If you already have an existing version of provider, either remove the existing provider file from the terraform plugins directory or update all configuration files to include latest version

```shell
rm ~/.terraform.d/plugins/terraform-provider-vra*
```

Move the extracted plugin to the terraform plugins directory

```shell
mv terraform-provider-vra_v${RELEASE} ~/.terraform.d/plugins/
```

#### Linux Example

![downloading with wget - Linux ](images/wget_release_linux.png)

### Windows

Create a terraform plugins directory typically this will be in ```%APPDATA%\terraform.d\plugins```.

```powershell
 #powershell
 mkdir $ENV:APPDATA\terraform.d\plugins
```

```cmd
#CMD
mkdir %APPDATA%\terraform.d\plugins
```

Download the plugin (via a browser or command line)

 ```powershell
 $RELEASE="0.3.5"
 wget https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra-windows_amd64-v${RELEASE}.tgz -outfile terraform-provider-vra-windows_amd64-v${RELEASE}.tgz
 ```

Untar/unzip the plugin (Depending on your setup this may require two steps)

```powershell
#using 7zip to unzip
7z x .\terraform-provider-vra-windows_amd64-v${RELEASE}.tgz

# then untar resulting file
tar xvf terraform-provider-vra-windows_amd64-v${RELEASE}.tar
```

Move the extracted plugin to the terraform plugins directory

```powershell
#Powershell
move terraform-provider-vra_v${RELEASE}.exe $ENV:APPDATA\terraform.d\plugins
```

```cmd
#CMD
move terraform-provider-vra_v%RELEASE%.exe %APPDATA%\terraform.d\plugins
```

#### Windows Example

![downloading with wget - Powershell ](images/wget_release_pshell.png)

## Validating the install

To validate the installation you can simply change to the location where your terraform configuration is located and run ```terraform init```. You should see a message indicating that terraform has been successfully initialized.

![init success](images/install_success.png)

## Get Provider version
To find the provider version, you can simply change to the location where your terraform configuration is located and run ```terraform -version```. You should see a message indicating the provider version.
