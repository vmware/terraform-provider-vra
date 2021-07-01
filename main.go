package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vmware/terraform-provider-vra/vra"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	buildTime string
	version   string
)

func main() {
	versionPtr := flag.Bool("v", false, "show version info")
	flag.Parse()

	if *versionPtr {
		fmt.Printf("version: %s\n", version)
		fmt.Printf("build time: %s\n", buildTime)
		os.Exit(0)
	}

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return vra.Provider()
		},
	})
}
