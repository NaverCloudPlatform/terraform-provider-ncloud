provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_regions" "regions" {
  filter {
		name   = "region_code"
		values = [".*N$"]
    regex = true
	}
  output_file = "regions.json"
}

