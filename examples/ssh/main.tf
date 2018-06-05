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

  connection {
    type = "ssh"
    user = "root"
    host = "${ncloud_instance.instance.port_forwarding_public_ip}"
    port = "${ncloud_instance.instance.port_forwarding_external_port}"
    password = "${ncloud_root_password.rootpwd.root_password}"
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

data "ncloud_root_password" "rootpwd" {
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "private_key" = "${ncloud_login_key.key.private_key}"
}

