---
subcategory: "Search Engine Service"
---


# Resource: ncloud_ses_cluster

Provides a Search Engine Service cluster resource.

## Example Usage

``` hcl
variable ses_user_password {
  description = "SES Cluster User Password"
  type = string
  sensitive =  true
}

resource "ncloud_vpc" "vpc" {
	name               = "tf-vpc"
	ipv4_cidr_block    = "172.16.0.0/16"
}

resource "ncloud_subnet" "node_subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "tf-subnet"
	subnet             = "172.16.1.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}
data "ncloud_ses_versions" "ses_versions" {
}

data "ncloud_ses_node_os_images" "os_images" {
}

data "ncloud_ses_node_products" "product_codes" {
  os_image_code = data.ncloud_ses_node_os_images.os_images.images.0.id
  subnet_no = ncloud_subnet.node_subnet.id
}

resource "ncloud_login_key" "loginkey" {
  key_name = "tf-login-key"
}

resource "ncloud_ses_cluster" "cluster" {
  cluster_name                  = "tf-cluster"
  os_image_code                 = data.ncloud_ses_node_os_images.os_images.images.0.id
  vpc_no                        = ncloud_vpc.vpc.id
  search_engine {
	  version_code    			= data.ncloud_ses_versions.ses_versions.versions.0.id
	  user_name       			= "admin"
	  user_password   			= var.ses_user_password
	  dashboard_port            = "5601"
  }
  manager_node {  
	  is_dual_manager           = false
	  product_code     			= data.ncloud_ses_node_products.product_codes.codes.0.id
	  subnet_no        			= ncloud_subnet.node_subnet.id
  }
  data_node {
	  product_code       		= data.ncloud_ses_node_products.product_codes.codes.0.id
	  subnet_no           		= ncloud_subnet.node_subnet.id
	  count            		    = 3
	  storage_size        		= 100
  }
  master_node {
	  product_code       		= data.ncloud_ses_node_products.product_codes.codes.0.id
	  subnet_no           		= ncloud_subnet.node_subnet.id
	  count            		    = 3
  }
  login_key_name                = ncloud_login_key.loginkey.key_name
}
```

## Argument Reference
The following arguments are supported:

* `cluster_name` - Cluster name.
* `os_image_code` -  OS type to be used.
* `vpc_no` - VPC number to be used.
* `search_engine` - .
    * `version_code` - Search Engine Service version to be used.
    * `user_name` - Search Engine UserName. Only lowercase alphanumeric characters and non-consecutive hyphens (-) allowed First character must be a letter, but the last character may be a letter or a number.
    * `user_password` - Search Engine User password. Must be at least 8 characters and contain at least one of each: English uppercase letter, lowercase letter, special character, and number.
    * `dashboard_port` - Search Engine Dashboard port.
* `manager_node` - .
    * `is_dual_manager` - Redundancy of manager node
    * `product_code` - HW specifications of the manager node.
    * `subnet_no` - Subnet number where the manager node is to be located.
* `data_node` - .
    * `product_code` - HW specifications of the data node.
    * `subnet_no` - Subnet number where the data node is to be located.
    * `count` - Number of data nodes. At least 3 units. (Can only be increased)
    * `storage_size` - Data node storage capacity. At least 100 GB, up to 2000 GB. Must be in units of 10 GB.
* `master_node(Optional)` - If declared, creates a master-only node.
    * `product_code` - HW specifications of the master node.
    * `subnet_no` - Subnet number where the master node is to be located.
    * `count` - Number of master nodes. Only 3 or 5 units are available.
* `login_key_name` - Required Login key to access Manager node server

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Cluster Instance No.
* `service_group_instance_no` - Cluster Instance No. (It is the same result as `id`)
* `manager_node` - .
  * `acg_id` - The ID of manager node ACG.
  * `acg_name` - The name of manager node ACG. 
* `data_node` - .
  * `acg_id` - The ID of data node ACG.
  * `acg_name` - The name of data node ACG.
* `master_node` - .
  * `acg_id` - The ID of master node ACG.
  * `acg_name` - The name of master node ACG.
* `manager_node_instance_no_list` - List of Manager node's instance number
* `cluster_node_list` - .
  * `compute_instance_name` - The name of Server instance.
  * `compute_instance_no`   - The ID of Server instance.
  * `node_type`             - Node role code
  * `private_ip`            - Private IP
  * `server_status`         - The status of Server Instance.
  * `subnet`                - The name of Server Instance subnet.
