provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

data "ncloud_server_image" "image" {
  "infra_resource_detail_type_code" = "BM" // Bare Metal Image
  "platform_type_code_list" = ["LNX64"] // Linux 64bit
}

data "ncloud_server_product" "prod" {
  "server_image_product_code" = "${data.ncloud_server_image.image.id}"
  "product_name_regex" = "^(.*)2\\.6 GHz(.*)8 cores(.*)"
}

resource "ncloud_login_key" "key" {
  "key_name" = "${var.login_key_name}"
}

resource "ncloud_server" "bm" {
  "name" = "${var.server_name}"
  "server_image_product_code" = "${data.ncloud_server_image.image.id}"
  "server_product_code" = "${data.ncloud_server_product.prod.id}"
  "login_key_name" = "${ncloud_login_key.key.key_name}"
  "raid_type_name" = "5"
  "zone_no" = "2"
}

resource "ncloud_block_storage" "storage" {
  "server_instance_no" = "${ncloud_server.bm.id}"
  "name" = "${var.block_storage_name}"
  "size" = "10"
}

resource "ncloud_nas_volume" "nas" {
  "volume_name_postfix" = "${var.nas_volume_name_prefix}"
  "volume_size_gb" = "500"
  "volume_allotment_protocol_type_code" = "NFS"
  "server_instance_no_list" = ["${ncloud_server.bm.id}"]
}

resource "ncloud_public_ip" "public_ip" {
  "server_instance_no" = "${ncloud_server.bm.id}"
}

resource "ncloud_load_balancer_ssl_certificate" "cert" {
  "certificate_name" = "${var.certificate_name}"
  "privatekey" = "${file("certs/lbtest.privateKey")}"
  "publickey_certificate" = "${file("certs/lbtest.crt")}"
}

resource "ncloud_load_balancer" "lb" {
  "name" = "${var.load_balancer_name}"
  "algorithm_type_code" = "SIPHS"
  "description" = "${var.load_balancer_name} description"

  "rule_list" = [
    {
      "protocol_type_code" = "HTTP"
      "load_balancer_port" = 80
      "server_port" = 80
      "l7_health_check_path" = "/monitor/l7check"
    },
    {
      "protocol_type_code" = "HTTPS"
      "load_balancer_port" = 443
      "server_port" = 443
      "l7_health_check_path" = "/monitor/l7check"
      "certificate_name"     = "${ncloud_load_balancer_ssl_certificate.cert.certificate_name}"
    },
  ]
  "server_instance_no_list" = ["${ncloud_server.bm.id}"]
  "internet_line_type_code" = "PUBLC"
  "network_usage_type_code" = "PBLIP"
}


data "ncloud_root_password" "pwd" {
  "server_instance_no" = "${ncloud_server.bm.id}"
  "private_key" = "${ncloud_login_key.key.private_key}"
}

data "ncloud_port_forwarding_rules" "rules" {
  "zone_no" = "${ncloud_server.bm.zone_no}"
}

resource "ncloud_port_forwarding_rule" "rule" {
  "port_forwarding_configuration_no" = "${data.ncloud_port_forwarding_rules.rules.id}"
  "server_instance_no" = "${ncloud_server.bm.id}"
  "port_forwarding_external_port" = "${var.port_forwarding_external_port}"
  "port_forwarding_internal_port" = "22"
}

resource "null_resource" "ssh" {
  connection {
    type = "ssh"
    user = "root"
    host = "${ncloud_port_forwarding_rule.rule.port_forwarding_public_ip}"
    port = "${ncloud_port_forwarding_rule.rule.port_forwarding_external_port}"
    password = "${data.ncloud_root_password.pwd.root_password}"
  }

  provisioner "file" {
    source = "scripts/mount-storage.sh"
    destination = "scripts/mount-storage.sh"
  }

  provisioner "remote-exec" {
    when = "create"
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
