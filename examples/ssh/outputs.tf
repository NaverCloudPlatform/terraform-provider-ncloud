output "server_name_list" {
  value = "${join(",", ncloud_instance.instance.*.server_name)}"
}