provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = "KR"
  support_vpc = true
}

resource "ncloud_login_key" "loginkey" {
  key_name = "tete1"
}
