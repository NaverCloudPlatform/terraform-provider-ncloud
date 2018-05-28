output "server_name_list" {
  value = "${join(",", ncloud_instance.terraform-test.*.server_name)}"
}