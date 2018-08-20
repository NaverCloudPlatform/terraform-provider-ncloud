provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

resource "ncloud_login_key" "key" {
  "key_name" = "${var.login_key_name}"
}
resource "ncloud_server" "server" {
  "server_name" = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code" = "${var.server_product_code}"
  "login_key_name" = "${ncloud_login_key.key.key_name}"
  "zone_code" = "KR-2"
}

resource "ncloud_block_storage" "storage" {
  "server_instance_no" = "${ncloud_server.server.id}"
  "block_storage_name" = "${var.block_storage_name}"
  "block_storage_size_gb" = "10"
}

data "ncloud_root_password" "rootpwd" {
  "server_instance_no" = "${ncloud_server.server.id}"
  "private_key" = "${ncloud_login_key.key.private_key}"
}

data "ncloud_port_forwarding_rules" "rules" {
  "zone_code" = "${ncloud_server.server.zone_code}"
}

resource "ncloud_port_forwarding_rule" "forwarding" {
  "port_forwarding_configuration_no" = "${data.ncloud_port_forwarding_rules.rules.id}"
  "server_instance_no" = "${ncloud_server.server.id}"
  "port_forwarding_external_port" = "${var.port_forwarding_external_port}"
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

  provisioner "file" {
    source = "mount-storage.sh"
    destination = "mount-storage.sh"
  }

  provisioner "remote-exec" {
    # CentOS 5.x: mkfs.ext3 /dev/xvdb1
    # CentOS 6.x: mkfs.ext4 /dev/xvdb1
    # CentOS 7.x: mkfs.xfs /dev/xvdb1
    # Ubuntu Server / Desktop: mkfs.ext4 /dev/xvdb1
    inline = [
      "chmod 755 mount-storage.sh",
      "sh mount-storage.sh >> mount-storage.log"
    ]
  }
}