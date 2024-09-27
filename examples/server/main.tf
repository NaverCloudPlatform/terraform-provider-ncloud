provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_regions" "regions" {
}

data "ncloud_server_image_numbers" "server_images" {
  filter {
    name = "name"
    values = ["ubuntu-20.04-base"]
  }
}

data "ncloud_server_specs" "spec" {
  filter {
    name   = "server_spec_code"
    values = ["c2-g3"]
  }
}

resource "ncloud_login_key" "loginkey" {
  key_name = "test-key"
}

resource "ncloud_vpc" "test" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
  vpc_no         = ncloud_vpc.test.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.test.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.test.default_network_acl_no
  subnet_type    = "PUBLIC"
  usage_type     = "GEN"
}

resource "ncloud_server" "server" {
  subnet_no           = ncloud_subnet.test.id
  name                = var.server_name
  server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
  server_spec_code    = data.ncloud_server_specs.spec.server_spec_list.0.server_spec_code
  login_key_name      = ncloud_login_key.loginkey.key_name
}
