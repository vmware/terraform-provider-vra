package main

import (
	"github.com/vmware/terraform-provider-vra/vra"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return vra.Provider()
		},
	})
}
