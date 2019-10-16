output "server_name_list" {
  value = join(",", ncloud_server.server.*.name)
}

