// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRAWordpressInfrastructure_Basic(t *testing.T) {
	nRint := acctest.RandInt()
	mRInt := acctest.RandInt()
	wRint := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAWordpressInfrastructureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAWordpressInfrastructureConfig(nRint, mRInt, wRint),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAWordpressInfrastructureResourceExists("vra_network.network"),
					testAccCheckVRAWordpressInfrastructureResourceExists("vra_machine.mysql"),
					testAccCheckVRAWordpressInfrastructureResourceExists("vra_machine.wordpress"),
					resource.TestMatchResourceAttr(
						"vra_network.network", "name", regexp.MustCompile("^terraform_vra_network-"+strconv.Itoa(nRint))),
					resource.TestCheckResourceAttr(
						"vra_network.network", "constraints.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_network.network", "constraints.0.mandatory", "true"),
					resource.TestCheckResourceAttr(
						"vra_network.network", "constraints.0.expression", "pci"),

					resource.TestMatchResourceAttr(
						"vra_machine.mysql", "name", regexp.MustCompile("^terraform_vra_mysql-"+strconv.Itoa(mRInt))),
					resource.TestCheckResourceAttr(
						"vra_machine.mysql", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"vra_machine.mysql", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"vra_machine.mysql", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_machine.mysql", "boot_config.#", "1"),

					resource.TestMatchResourceAttr(
						"vra_machine.wordpress", "name", regexp.MustCompile("^terraform_vra_wordpress-"+strconv.Itoa(wRint))),
					resource.TestCheckResourceAttr(
						"vra_machine.wordpress", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"vra_machine.wordpress", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"vra_machine.wordpress", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_machine.wordpress", "boot_config.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVRAWordpressInfrastructureResourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No %s ID is set", n)
		}

		return nil
	}
}

func testAccCheckVRAWordpressInfrastructureDestroy(_ *terraform.State) error {
	/*
		apiClient := testAccProviderVRA.Meta().(*Client).apiClient

		for _, rs := range s.RootModule().Resources {

			selfKey := ""
			for key, value := range rs.Primary.Attributes {
				if value == "self" {
					selfKey = key
					break
				}
			}

			_, err := client.ReadResource(rs.Primary.Attributes[strings.Replace(selfKey, "rel", "href", 1)])

			if err != nil && !strings.Contains(err.Error(), "404") {
				return fmt.Errorf(
					"Error waiting for (%s) to be destroyed: %s",
					rs.Type+"/"+rs.Primary.Attributes["name"], err)
			}
		}
	*/

	return nil
}

func testAccCheckVRAWordpressInfrastructureConfig(nRInt, mRInt, wRInt int) string {
	return fmt.Sprintf(`
resource "vra_network" "network" {
	name = "terraform_vra_network-%d"

	constraints {
		mandatory = true
		expression = "pci"
	}
}

resource "vra_machine" "mysql" {
	name = "terraform_vra_mysql-%d"

	image = "ubuntu"
	flavor = "small"

	nics {
        network_id = "${vra_network.network.id}"
	}

	boot_config {
        content = <<EOF
#cloud-config
repo_update: true
repo_upgrade: all

packages:
 - mysql-server

runcmd:
 - sed -e '/bind-address/ s/^#*/#/' -i /etc/mysql/mysql.conf.d/mysqld.cnf
 - service mysql restart
 - mysql -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%%' IDENTIFIED BY 'mysqlpassword';"
 - mysql -e "FLUSH PRIVILEGES;"
EOF
    }
}

resource "vra_machine" "wordpress" {
	name = "terraform_vra_wordpress-%d"

	image = "ubuntu"
	flavor = "small"

	nics {
        network_id = "${vra_network.network.id}"
	}

	boot_config {
        content = <<EOF
#cloud-config
repo_update: true
repo_upgrade: all

packages:
- apache2
- php
- php-mysql
- libapache2-mod-php
- php-mcrypt
- mysql-client

runcmd:
- mkdir -p /var/www/html/mywordpresssite && cd /var/www/html && wget https://wordpress.org/latest.tar.gz && tar -xzf /var/www/html/latest.tar.gz -C /var/www/html/mywordpresssite --strip-components 1
- i=0; while [ $i -le 10 ]; do mysql --connect-timeout=3 -h ${vra_machine.mysql.address} -u root -pmysqlpassword -e "SHOW STATUS;" && break || sleep 15; i=$$((i+1)); done
- mysql -u root -pmysqlpassword -h ${vra_machine.mysql.address} -e "create database wordpress_blog;"
- mv /var/www/html/mywordpresssite/wp-config-sample.php /var/www/html/mywordpresssite/wp-config.php
- sed -i -e s/"define('DB_NAME', 'database_name_here');"/"define('DB_NAME', 'wordpress_blog');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_USER', 'username_here');"/"define('DB_USER', 'root');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_PASSWORD', 'password_here');"/"define('DB_PASSWORD', 'mysqlpassword');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_HOST', 'localhost');"/"define('DB_HOST', '${vra_machine.mysql.address}');"/ /var/www/html/mywordpresssite/wp-config.php
- service apache2 reload
EOF
    }
}`, nRInt, mRInt, wRInt)
}
