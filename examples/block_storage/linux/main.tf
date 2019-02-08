provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

resource "ncloud_login_key" "key" {
  "key_name" = "${var.login_key_name}"
}
resource "ncloud_server" "server" {
  "name" = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code" = "${var.server_product_code}"
  "login_key_name" = "${ncloud_login_key.key.key_name}"
  "zone_code" = "KR-2"
}

resource "ncloud_block_storage" "storage" {
  "server_instance_no" = "${ncloud_server.server.id}"
  "name" = "${var.block_storage_name}"
  "size" = "10"
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
    source = "scripts/mount-storage.sh"
    destination = "scripts/mount-storage.sh"
  }

  provisioner "file" {
    source = "scripts/unmount-storage.sh"
    destination = "scripts/unmount-storage.sh"
  }

  provisioner "remote-exec" {
    when = "create"
    # CentOS 5.x: mkfs.ext3 /dev/xvdb1
    # CentOS 6.x: mkfs.ext4 /dev/xvdb1
    # CentOS 7.x: mkfs.xfs /dev/xvdb1
    # Ubuntu Server / Desktop: mkfs.ext4 /dev/xvdb1
    inline = [
      "chmod 755 scripts/mount-storage.sh",
      "sh scripts/mount-storage.sh >> scripts/mount-storage.log",
      "mount"
    ]
  }

  provisioner "remote-exec" {
      when = "destroy"
      inline = [
        "chmod 755 scripts/umount-storage.sh",
        "sh scripts/umount-storage.sh >> scripts/umount-storage.log",
        "mount"
      ]
    }
}
