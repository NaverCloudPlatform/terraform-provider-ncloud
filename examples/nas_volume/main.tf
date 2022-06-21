provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

resource "ncloud_nas_volume" "nas" {
  volume_name_postfix            = "tftest_vol"
  volume_size                    = "500"
  volume_allotment_protocol_type = "NFS"
}

