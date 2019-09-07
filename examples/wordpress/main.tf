provider "vra" {
  url          = var.url
  access_token = var.access_token
}

resource "vra_machine" "database" {
  name   = "terraform-vra-mysql"
  image  = "ubuntu"
  flavor = "small"

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
 - mysql -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'mysqlpassword';"
 - mysql -e "FLUSH PRIVILEGES;"
EOF
  }
}

resource "vra_machine" "wordpress" {
  name   = "terraform_vra_wordpress"
  image  = "ubuntu"
  flavor = "small"

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
- i=0; while [ $i -le 10 ]; do mysql --connect-timeout=3 -h ${vra_machine.database.address} -u root -pmysqlpassword -e "SHOW STATUS;" && break || sleep 15; i=$$((i+1)); done
- mysql -u root -pmysqlpassword -h ${vra_machine.database.address} -e "create database wordpress_blog;"
- mv /var/www/html/mywordpresssite/wp-config-sample.php /var/www/html/mywordpresssite/wp-config.php
- sed -i -e s/"define('DB_NAME', 'database_name_here');"/"define('DB_NAME', 'wordpress_blog');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_USER', 'username_here');"/"define('DB_USER', 'root');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_PASSWORD', 'password_here');"/"define('DB_PASSWORD', 'mysqlpassword');"/ /var/www/html/mywordpresssite/wp-config.php && sed -i -e s/"define('DB_HOST', 'localhost');"/"define('DB_HOST', '${vra_machine.database.address}');"/ /var/www/html/mywordpresssite/wp-config.php
- service apache2 reload
EOF
  }
}
