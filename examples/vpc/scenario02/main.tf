# VPC > User scenario > Scenario 2. Public and Private Subnet
# https://docs.ncloud.com/ko/networking/vpc/vpc_userscenario2.html

provider "ncloud" {
  support_vpc = true
  region      = "KR"
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_login_key" "key_scn_02" {
  key_name = var.name_scn02
}

# VPC
resource "ncloud_vpc" "vpc_scn_02" {
  name            = var.name_scn02
  ipv4_cidr_block = "10.0.0.0/16"
}

# Subnet
resource "ncloud_subnet" "subnet_scn_02_public" {
  name           = "${var.name_scn02}-public"
  vpc_no         = ncloud_vpc.vpc_scn_02.id
  subnet         = cidrsubnet(ncloud_vpc.vpc_scn_02.ipv4_cidr_block, 8, 0)
  // "10.0.0.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_network_acl.network_acl_02_public.id
  subnet_type    = "PUBLIC"
  // PUBLIC(Public)
}

resource "ncloud_subnet" "subnet_scn_02_private" {
  name           = "${var.name_scn02}-private"
  vpc_no         = ncloud_vpc.vpc_scn_02.id
  subnet         = cidrsubnet(ncloud_vpc.vpc_scn_02.ipv4_cidr_block, 8, 1)
  // "10.0.1.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_network_acl.network_acl_02_private.id
  subnet_type    = "PRIVATE"
  // PRIVATE(Private)
}

# Network ACL
resource "ncloud_network_acl" "network_acl_02_public" {
  vpc_no = ncloud_vpc.vpc_scn_02.id
  name   = "${var.name_scn02}-public"
}

resource "ncloud_network_acl" "network_acl_02_private" {
  vpc_no = ncloud_vpc.vpc_scn_02.id
  name   = "${var.name_scn02}-private"
}


# Server
resource "ncloud_server" "server_scn_02_public" {
  subnet_no                 = ncloud_subnet.subnet_scn_02_public.id
  name                      = "${var.name_scn02}-public"
  server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
  login_key_name            = ncloud_login_key.key_scn_02.key_name
  //server_product_code       = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
}

resource "ncloud_server" "server_scn_02_private" {
  subnet_no                 = ncloud_subnet.subnet_scn_02_private.id
  name                      = "${var.name_scn02}-private"
  server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
  login_key_name            = ncloud_login_key.key_scn_02.key_name
  //server_product_code       = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
}

# Public IP
resource "ncloud_public_ip" "public_ip_scn_02" {
  server_instance_no = ncloud_server.server_scn_02_public.id
  description        = "for ${var.name_scn02}"
}

# NAT Gateway
resource "ncloud_nat_gateway" "nat_gateway_scn_02" {
  vpc_no = ncloud_vpc.vpc_scn_02.id
  zone   = "KR-2"
  name   = var.name_scn02
}

# Route Table
resource "ncloud_route" "route_scn_02_nat" {
  route_table_no         = ncloud_vpc.vpc_scn_02.default_private_route_table_no
  destination_cidr_block = "0.0.0.0/0"
  target_type            = "NATGW"
  // NATGW (NAT Gateway) | VPCPEERING (VPC Peering) | VGW (Virtual Private Gateway).
  target_name            = ncloud_nat_gateway.nat_gateway_scn_02.name
  target_no              = ncloud_nat_gateway.nat_gateway_scn_02.id
}


data "ncloud_root_password" "scn_02_root_password" {
  server_instance_no = ncloud_server.server_scn_02_public.id
  private_key        = ncloud_login_key.key_scn_02.private_key
}

resource "null_resource" "ls-al" {
  connection {
    type     = "ssh"
    host     = ncloud_public_ip.public_ip_scn_02.public_ip
    user     = "root"
    port     = "22"
    password = data.ncloud_root_password.scn_02_root_password.root_password
  }

  provisioner "remote-exec" {
    inline = [
      "ls -al",
    ]
  }

  depends_on = [
    ncloud_public_ip.public_ip_scn_02,
    ncloud_server.server_scn_02_public
  ]
}

# You can add ACG rules remove comment If you want
/*
locals {
  default_acg_rules_inbound = [
    ["TCP", "0.0.0.0/0", "80"],
    ["TCP", "0.0.0.0/0", "443"],
    ["TCP", "${var.client_ip}/32", "22"],
    ["TCP", "${var.client_ip}/32", "3389"],
  ]

  default_acg_rules_outbound = [
    ["TCP", "0.0.0.0/0", "1-65535"],
    ["UDP", "0.0.0.0/0", "1-65534"],
    ["ICMP", "0.0.0.0/0", null]
  ]
}

resource "ncloud_access_control_group" "acg_scn_02" {
  description = "for acc test"
  vpc_no      = ncloud_vpc.vpc_scn_02.id
}

resource "ncloud_access_control_group_rule" "acg_rule_scn_02" {
  access_control_group_no = ncloud_access_control_group.acg_scn_02.id

  dynamic "inbound" {
    for_each = local.default_acg_rules_inbound
    content {
      protocol    = inbound.value[0]
      ip_block    = inbound.value[1]
      port_range  = inbound.value[2]
    }
  }

  dynamic "outbound" {
    for_each = local.default_acg_rules_outbound
    content {
      protocol    = outbound.value[0]
      ip_block    = outbound.value[1]
      port_range  = outbound.value[2]
    }
  }
}
*/
