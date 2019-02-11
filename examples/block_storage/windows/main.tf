provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "ncloud_login_key" "key" {
  "key_name" = "${var.login_key_name}"
}

data "ncloud_access_control_group" "acg" {
  "name" = "${var.acg_name}"
}

data "ncloud_server_image" "image" {
  "product_name_regex" = "^Windows Server 2012(.*)"
}

data "ncloud_server_product" "prod" {
  "server_image_product_code" = "${data.ncloud_server_image.image.id}"
  "product_name_regex"        = "^vCPU 2EA(.*)Memory 4GB(.*)"
}

resource "ncloud_server" "server" {
  "name"                                       = "${var.server_name}"
  "server_image_product_code"                  = "${data.ncloud_server_image.image.id}"
  "server_product_code"                        = "${data.ncloud_server_product.prod.id}"
  "login_key_name"                             = "${ncloud_login_key.key.key_name}"
  "access_control_group_configuration_no_list" = ["${data.ncloud_access_control_group.acg.id}"]
  "user_data"                                  = "CreateObject(\"WScript.Shell\").run(\"cmd.exe /c powershell Set-ExecutionPolicy RemoteSigned & winrm set winrm/config/service/auth @{Basic=\"\"true\"\"} & winrm set winrm/config/service @{AllowUnencrypted=\"\"true\"\"} & winrm quickconfig -q & sc config WinRM start= auto & winrm get winrm/config/service\")"
  "zone"                                       = "KR-2"
}

resource "ncloud_public_ip" "public_ip" {
  "server_instance_no" = "${ncloud_server.server.id}"
  "zone"               = "${var.zone}"
}

resource "ncloud_block_storage" "storage" {
  "server_instance_no" = "${ncloud_server.server.id}"
  "name"               = "${var.block_storage_name}"
  "size"               = "10"
}

data "ncloud_root_password" "rootpwd" {
  "server_instance_no" = "${ncloud_server.server.id}"
  "private_key"        = "${ncloud_login_key.key.private_key}"
}

data "ncloud_port_forwarding_rules" "rules" {
  "zone" = "${ncloud_server.server.zone}"
}

resource "ncloud_port_forwarding_rule" "forwarding" {
  "port_forwarding_configuration_no" = "${data.ncloud_port_forwarding_rules.rules.id}"
  "server_instance_no"               = "${ncloud_server.server.id}"
  "port_forwarding_external_port"    = "${var.port_forwarding_external_port}"
  "port_forwarding_internal_port"    = "3389"
}

resource "null_resource" "winrm" {
  depends_on = ["ncloud_public_ip.public_ip", "ncloud_block_storage.storage"]

  connection {
    type     = "winrm"
    user     = "Administrator"
    host     = "${ncloud_public_ip.public_ip.public_ip}"
    password = "${data.ncloud_root_password.rootpwd.root_password}"
  }

  # mount
  provisioner "file" {
    source      = "scripts/mount-storage.ps1"
    destination = "C:\\scripts\\mount-storage.ps1"
  }

  # unmount
  provisioner "file" {
    source      = "scripts/unmount-storage.ps1"
    destination = "C:\\scripts\\unmount-storage.ps1"
  }

  # get mount points
  provisioner "file" {
    source      = "scripts/get-mountpoints.ps1"
    destination = "C:\\scripts\\get-mountpoints.ps1"
  }

  provisioner "remote-exec" {
    when = "create"

    inline = [
      "powershell.exe -File C:\\scripts\\mount-storage.ps1",
      "powershell.exe -File C:\\scripts\\get-mountpoints.ps1",
    ]
  }

  provisioner "remote-exec" {
    when = "destroy"

    inline = [
      "powershell.exe -File C:\\scripts\\unmount-storage.ps1",
      "powershell.exe -File C:\\scripts\\get-mountpoints.ps1",
    ]
  }
}
