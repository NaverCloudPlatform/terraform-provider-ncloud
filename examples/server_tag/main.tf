provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_zones" "korea" {
  "region" = "KR"
}

resource "ncloud_server" "server" {
  "server_image_product_code" = "SPSW0LINUX000046" #Conditional "server_image_product_code" OR "member_server_image_no"
  "server_product_code"       = "SPSVRSTAND000004" #Optional

  "tag_list" = [
    {
      "tag_key"   = "samplekey1"
      "tag_value" = "samplevalue1"
    },
    {
      "tag_key"   = "samplekey2"
      "tag_value" = "samplevalue2"
    },
  ]
}
