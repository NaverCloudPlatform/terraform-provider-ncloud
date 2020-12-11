# Network ACL Rule
locals {
  public_subnet_inbound = [
    [1, "TCP", "0.0.0.0/0", "80", "ALLOW"],
    [2, "TCP", "0.0.0.0/0", "443", "ALLOW"],
    [3, "TCP", "${var.client_ip}/32", "22", "ALLOW"],
    [4, "TCP", "${var.client_ip}/32", "3389", "ALLOW"],
    [5, "TCP", "0.0.0.0/0", "9001-65535", "ALLOW"],
    [197, "TCP", "0.0.0.0/0", "1-65535", "DROP"],
    [198, "UDP", "0.0.0.0/0", "1-65535", "DROP"],
    [199, "ICMP", "0.0.0.0/0", null, "DROP"],
  ]

  public_subnet_outbound = [
    [1, "TCP", "0.0.0.0/0", "80", "ALLOW"],
    [2, "TCP", "0.0.0.0/0", "443", "ALLOW"],
    [3, "TCP", "0.0.0.0/0", "32768-65535", "ALLOW"],
    [4, "TCP", "${ncloud_server.server_scn_02_private.network_interface[0].private_ip}/32", "8080", "ALLOW"],
    // Allow 8080 port to private server
    [197, "TCP", "0.0.0.0/0", "1-65535", "DROP"],
    [198, "UDP", "0.0.0.0/0", "1-65535", "DROP"],
    [199, "ICMP", "0.0.0.0/0", null, "DROP"]
  ]
}

resource "ncloud_network_acl_rule" "network_acl_02_rule_public" {
  network_acl_no = ncloud_network_acl.network_acl_02_public.id
  dynamic "inbound" {
    for_each = local.public_subnet_inbound
    content {
      priority    = inbound.value[0]
      protocol    = inbound.value[1]
      ip_block    = inbound.value[2]
      port_range  = inbound.value[3]
      rule_action = inbound.value[4]
    }
  }

  dynamic "outbound" {
    for_each = local.public_subnet_outbound
    content {
      priority    = outbound.value[0]
      protocol    = outbound.value[1]
      ip_block    = outbound.value[2]
      port_range  = outbound.value[3]
      rule_action = outbound.value[4]
    }
  }
}

locals {
  private_subnet_inbound = [
    [1, "TCP", "${ncloud_server.server_scn_02_public.network_interface[0].private_ip}/32", "8080", "ALLOW"],
    // Allow 8080 port from public server
    [2, "TCP", "0.0.0.0/0", "32768-65535", "ALLOW"],
    [197, "TCP", "0.0.0.0/0", "1-65535", "DROP"],
    [198, "UDP", "0.0.0.0/0", "1-65535", "DROP"],
    [199, "ICMP", "0.0.0.0/0", null, "DROP"],
  ]

  private_subnet_outbound = [
    [1, "TCP", "${ncloud_server.server_scn_02_public.network_interface[0].private_ip}/32", "32768-65535", "ALLOW"],
    // Allow 32768-65535 port to public server
    [197, "TCP", "0.0.0.0/0", "1-65535", "DROP"],
    [198, "UDP", "0.0.0.0/0", "1-65535", "DROP"],
    [199, "ICMP", "0.0.0.0/0", null, "DROP"]
  ]
}

resource "ncloud_network_acl_rule" "network_acl_02_private" {
  network_acl_no = ncloud_network_acl.network_acl_02_private.id

  dynamic "inbound" {
    for_each = local.private_subnet_inbound
    content {
      priority    = inbound.value[0]
      protocol    = inbound.value[1]
      ip_block    = inbound.value[2]
      port_range  = inbound.value[3]
      rule_action = inbound.value[4]
    }
  }

  dynamic "outbound" {
    for_each = local.private_subnet_outbound
    content {
      priority    = outbound.value[0]
      protocol    = outbound.value[1]
      ip_block    = outbound.value[2]
      port_range  = outbound.value[3]
      rule_action = outbound.value[4]
    }
  }
}