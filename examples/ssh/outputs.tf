output "server_name_list" {
  value = "${join(",", ncloud_instance.instance.*.server_name)}"
}

output "port_forward_info" {
  value = "host = ${ncloud_instance.instance.port_forwarding_public_ip}, port = ${ncloud_instance.instance.port_forwarding_external_port}"
}