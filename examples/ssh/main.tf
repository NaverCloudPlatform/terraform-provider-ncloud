provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

resource "ncloud_login_key" "key" {
  "key_name" = "${var.login_key_name}"
}

resource "ncloud_instance" "instance" {
  "server_name" = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code" = "${var.server_product_code}"
  "login_key_name" = "${ncloud_login_key.key.key_name}"
}

data "ncloud_root_password" "rootpwd" {
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "private_key" = "${ncloud_login_key.key.private_key}"
}

data "ncloud_port_forwarding_rules" "rules" {}

resource "ncloud_port_forwarding_rule" "forwarding" {
  "port_forwarding_configuration_no" = "${data.ncloud_port_forwarding_rules.rules.id}"
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "port_forwarding_external_port" = "6088"
  "port_forwarding_internal_port" = "22"
}

resource "null_resource" "ssh" {
  connection {
    type = "ssh"
    user = "root"
    host = "${ncloud_port_forwarding_rule.forwarding.port_forwarding_public_ip}"
    port = "${ncloud_port_forwarding_rule.forwarding.port_forwarding_external_port}"
    password = "${data.ncloud_root_password.rootpwd.root_password}"
  }

  # Copies the file as the root user using SSH
  provisioner "file" {
    source = "myapp.conf"
    destination = "/etc/myapp.conf"
  }

  provisioner "remote-exec" {
    inline = [
      "echo 'hello'"
    ]
  }
}