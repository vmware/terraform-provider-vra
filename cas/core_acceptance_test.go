package cas

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCASWordpressInfrastructure_Basic(t *testing.T) {
	nRint := acctest.RandInt()
	mRInt := acctest.RandInt()
	wRint := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASWordpressInfrastructureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASWordpressInfrastructureConfig(nRint, mRInt, wRint),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASWordpressInfrastructureResourceExists("cas_network.network"),
					testAccCheckCASWordpressInfrastructureResourceExists("cas_machine.mysql"),
					testAccCheckCASWordpressInfrastructureResourceExists("cas_machine.wordpress"),
					resource.TestMatchResourceAttr(
						"cas_network.network", "name", regexp.MustCompile("^terraform_cas_network-"+strconv.Itoa(nRint))),
					resource.TestCheckResourceAttr(
						"cas_network.network", "constraints.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_network.network", "constraints.0.mandatory", "true"),
					resource.TestCheckResourceAttr(
						"cas_network.network", "constraints.0.expression", "pci"),

					resource.TestMatchResourceAttr(
						"cas_machine.mysql", "name", regexp.MustCompile("^terraform_cas_mysql-"+strconv.Itoa(mRInt))),
					resource.TestCheckResourceAttr(
						"cas_machine.mysql", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"cas_machine.mysql", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"cas_machine.mysql", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_machine.mysql", "boot_config.#", "1"),

					resource.TestMatchResourceAttr(
						"cas_machine.wordpress", "name", regexp.MustCompile("^terraform_cas_wordpress-"+strconv.Itoa(wRint))),
					resource.TestCheckResourceAttr(
						"cas_machine.wordpress", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"cas_machine.wordpress", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"cas_machine.wordpress", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_machine.wordpress", "boot_config.#", "1"),
				),
			},
		},
	})
}

func testAccCheckCASWordpressInfrastructureResourceExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASWordpressInfrastructureDestroy(s *terraform.State) error {
	/*
		apiClient := testAccProviderCAS.Meta().(*Client).apiClient

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

func testAccCheckCASWordpressInfrastructureConfig(nRInt, mRInt, wRInt int) string {
	return fmt.Sprintf(`
resource "cas_network" "network" {
	name = "terraform_cas_network-%d"

	constraints {
		mandatory = true
		expression = "pci"
	}
}

resource "cas_machine" "mysql" {
	name = "terraform_cas_mysql-%d"
	
	image = "ubuntu"
	flavor = "small"	

	nics {
        network_id = "${cas_network.network.id}"
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

resource "cas_machine" "wordpress" {
	name = "terraform_cas_wordpress-%d"
	
	image = "ubuntu"
	flavor = "small"	

	nics {
        network_id = "${cas_network.network.id}"
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
- i=0; while [ $i -le 10 ]; do mysql --connect-timeout=3 -h ${cas_machine.mysql.address} -u root -pmysqlpassword -e "SHOW STATUS;" && break || sleep 15; i=$$((i+1)); done
- mysql -u root -pmysqlpassword -h ${cas_machine.mysql.address} -e "create database wordpress_blog;"
- mv /var/www/html/mywordpresssite/wp-config-sample.php /var/www/html/mywordpresssite/wp-config.php
- sed -i -e s/"define('DB_NAME', 'database_name_here');"/"define('DB_NAME', 'wordpress_blog');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_USER', 'username_here');"/"define('DB_USER', 'root');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_PASSWORD', 'password_here');"/"define('DB_PASSWORD', 'mysqlpassword');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_HOST', 'localhost');"/"define('DB_HOST', '${cas_machine.mysql.address}');"/ /var/www/html/mywordpresssite/wp-config.php
- service apache2 reload
EOF
    }
}`, nRInt, mRInt, wRInt)
}
