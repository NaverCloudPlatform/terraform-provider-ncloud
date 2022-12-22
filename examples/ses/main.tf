provider "ncloud" {
  support_vpc = true
  region      = "KR"        # FKR
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_vpc" "vpc" {
  name            = "tf-ses-vpc"
  ipv4_cidr_block = "172.16.0.0/16"
}

resource "ncloud_subnet" "node_subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "172.16.1.0/24"
  zone           = "KR-1"           # FKR-1
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "tf-ses-private-subnet"
  usage_type     = "GEN"
}

data "ncloud_ses_versions" "ses_versions" {
  filter {
    name = "id"
    values = [var.ses_version]
    regex = true
  }
}

data "ncloud_ses_node_os_images" "os_versions" {
  filter {
    name = "id"
    values = [var.os_version]
    regex = true
  }
}

data "ncloud_ses_node_products" "product_codes" {
  os_image_code = data.ncloud_ses_node_os_images.os_versions.codes.0.id
  subnet_no = ncloud_subnet.node_subnet.id
  filter {
    name = "id"
    values = [var.ses_product_code]
  }
  filter {
    name   = "cpu_count"
    values = [var.ses_produce_cpu_count]
  }
}

resource "ncloud_login_key" "loginkey" {
  key_name = var.login_key
}

resource "ncloud_ses_cluster" "cluster" {
  cluster_name                  = "tf-ses"
  os_image_code                 = data.ncloud_ses_node_os_images.os_versions.codes.0.id
  vpc_no                        = ncloud_vpc.vpc.id
  search_engine {
    version_code    			= data.ncloud_ses_versions.ses_versions.codes.0.id
    user_name       			= "admin"
    user_password   			= var.ses_user_password
    dashboard_port              = "5601"
  }
  manager_node {
    is_dual_manager             = false
    product_code     			= data.ncloud_ses_node_products.product_codes.codes.0.id
    subnet_no        			= ncloud_subnet.node_subnet.id
  }
  data_node {
    product_code       		    = data.ncloud_ses_node_products.product_codes.codes.0.id
    subnet_no           		= ncloud_subnet.node_subnet.id
    count            		    = 4
    storage_size        		= 100
  }
  master_node {
    is_master_only_node_activated = true
    product_code       		      = data.ncloud_ses_node_products.product_codes.codes.0.id
    subnet_no           		  = ncloud_subnet.node_subnet.id
    count            		      = 3
  }
  login_key_name                = ncloud_login_key.loginkey.key_name
}

data "ncloud_ses_cluster" "cluster" {
  uuid = ncloud_ses_cluster.cluster.uuid
}

output "cluster" {
  value = data.ncloud_ses_cluster.cluster
}

data "ncloud_ses_clusters" "clusters" {
    filter {
      name   = "cluster_name"
      values = ["oss"]
    }
}

output "clusters" {
  value = data.ncloud_ses_clusters.clusters
}

data "ncloud_ses_cluster" "target_cluster" {
  uuid = data.ncloud_ses_clusters.clusters.clusters.0.id
}

output "target_cluster" {
  value = data.ncloud_ses_cluster.target_cluster
}