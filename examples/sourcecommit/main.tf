provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

resource "ncloud_sourcecommit_repository" "test-repository" {
  name = "tf-test-repository"
  desciptype = "test repository"
  filesafer = false
}