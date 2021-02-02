provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

resource "ncloud_launch_configuration" "lc" {
  name = "my-lc"
  server_image_product_code = "SPSW0LINUX000046"
  server_product_code = "SPSVRSSD00000003"
}

resource "ncloud_auto_scaling_group" "asg" {
  name = "my-auto"
  launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
  min_size = 1
  max_size = 1
  zone_no_list = ["2"]
  wait_for_capacity_timeout = "0"
}

resource "ncloud_auto_scaling_policy" "policy" {
  name = "my-policy"
  adjustment_type_code = "CHANG"
  scaling_adjustment = 2
  auto_scaling_group_no = ncloud_auto_scaling_group.asg.auto_scaling_group_no
}

resource "ncloud_auto_scaling_schedule" "schedule" {
  name = "my-schedule"
  min_size = 1
  max_size = 1
  desired_capacity = 1
  start_time = "" # 2021-02-02T15:00:00+0900
  end_time = "" # 2021-02-02T17:00:00+0900
  auto_scaling_group_no = ncloud_auto_scaling_group.asg.auto_scaling_group_no
}