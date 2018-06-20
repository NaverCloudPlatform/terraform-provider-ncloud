provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

resource "ncloud_login_key" "key" {
  "key_name" = "${var.login_key_name}"
}

data "ncloud_access_control_group" "acg" {
  "access_control_group_name" = "${var.acg_name}"
}

data "ncloud_server_image" "image" {
  "product_name_regex" = "^Windows Server 2016 \\(64-bit\\) English Edition$"
}

data "ncloud_server_product" "prod" {
  "server_image_product_code" = "${data.ncloud_server_image.image.id}"
  "product_name_regex" = "^vCPU 2EA(.*)Memory 2GB(.*)"
}

resource "ncloud_instance" "instance" {
  "server_name" = "${var.server_name}"
  "server_image_product_code" = "${data.ncloud_server_image.image.id}"
  "server_product_code" = "${data.ncloud_server_product.prod.id}"
  "login_key_name" = "${ncloud_login_key.key.key_name}"
  "access_control_group_configuration_no_list" = ["${data.ncloud_access_control_group.acg.id}"]
  "user_data" =  <<USER_DATA
CreateObject("WScript.Shell").run("cmd.exe /c powershell Set-ExecutionPolicy RemoteSigned & winrm set winrm/config/service/auth @{Basic="true"} & winrm set winrm/config/service @{AllowUnencrypted="true"} & winrm quickconfig -q & sc config WinRM start= auto & winrm get winrm/config/service")
USER_DATA
}

resource "ncloud_public_ip" "public_ip" {
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "zone_no"            = "${var.zone_no}"
}

resource "ncloud_block_storage" "storage" {
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "block_storage_name" = "${var.block_storage_name}"
  "block_storage_size_gb" = "10"
}

data "ncloud_root_password" "rootpwd" {
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "private_key" = "${ncloud_login_key.key.private_key}"
}

resource "null_resource" "winrm" {
  connection {
    type = "winrm"
    user = "Administrator"
    host = "${ncloud_public_ip.public_ip.public_ip}"
    password = "${data.ncloud_root_password.rootpwd.root_password}"
  }

  provisioner "file" {
    source = "mount-storage.ps1"
    destination = "C:\\scripts\\mount-storage.ps1"
  }

  provisioner "remote-exec" {
    inline = [
      "powershell.exe -File C:\\scripts\\upload.ps1"
    ]
  }
//  provisioner "remote-exec" {
//    # CentOS 5.x: mkfs.ext3 /dev/xvdb1
//    # CentOS 6.x: mkfs.ext4 /dev/xvdb1
//    # CentOS 7.x: mkfs.xfs /dev/xvdb1
//    # Ubuntu Server / Desktop: mkfs.ext4 /dev/xvdb1
//    inline = [
//      "chmod 755 mount-storage.sh",
//      "sh mount-storage.sh >> mount-storage.log"
//    ]
//  }
}