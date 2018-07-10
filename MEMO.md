# Terraform provider ncloud

[Terraform](https://www.terraform.io/)

https://oss.navercorp.com/ncloud-paas/terraform-provider-ncloud
˜

## Install

#### Terraform

https://www.terraform.io/intro/getting-started/install.html 을 참고하여 terraform 설치

#### Terrform provider ncloud

아직 terraform provider ncloud 은 개발 단계로 공식 terraform provider 사이트에 등록전입니다.
전달받은 파일은 $HOME/.terraform.d/plugins 에 옮겨 놓습니다.

참고 : https://www.terraform.io/docs/extend/how-terraform-works.html#discovery

## 시작하기

인프라 구성을 하기 위한 코드를 작성합니다.

```
$ cat main.tf
provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "ncloud_server" "server" {
  "server_name"               = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code"       = "${var.server_product_code}"
}

$ cat variables.tf
variable "access_key" {} # export TF_VAR_access_key=...
variable "secret_key" {} # export TF_VAR_secret_key=...

variable "region" {
  default = "KR"
}

variable "server_name" {
  default = "tf-test"
}

variable "server_image_product_code" {
  default = "SPSW0LINUX000032"
}

variable "server_product_code" {
  default = "SPSVRSTAND000004" #SPSVRSTAND000056
}
```

## Terraform init

테라폼 사용 초기화 단계

```
$ terraform init

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

## Terraform plan

Generate and show an execution plan

```
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + ncloud_server.server
      id:                                    <computed>
      base_block_storage_disk_detail_type.%: <computed>
      base_block_storage_disk_type.%:        <computed>
      base_block_storage_size:               <computed>
      cpu_count:                             <computed>
      create_date:                           <computed>
      internet_line_type.%:                  <computed>
      is_fee_charging_monitoring:            <computed>
      memory_size:                           <computed>
      platform_type.%:                       <computed>
      port_forwarding_external_port:         <computed>
      port_forwarding_internal_port:         <computed>
      port_forwarding_public_ip:             <computed>
      private_ip:                            <computed>
      public_ip:                             <computed>
      region.%:                              <computed>
      server_create_count:                   "1"
      server_image_name:                     <computed>
      server_image_product_code:             "SPSW0LINUX000032"
      server_instance_no:                    <computed>
      server_instance_operation.%:           <computed>
      server_instance_status.%:              <computed>
      server_instance_status_name:           <computed>
      server_name:                           "terraform-test"
      server_product_code:                   "SPSVRSTAND000004"
      uptime:                                <computed>
      zone.%:                                <computed>


Plan: 1 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

## Terraform graph

Create a visual graph of Terraform resources

```
$ terraform graph | dot -Tsvg > graph.svg
```

## Terraform apply

```
$ terraform apply

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + ncloud_server.server
      id:                                    <computed>
      base_block_storage_disk_detail_type.%: <computed>
      base_block_storage_disk_type.%:        <computed>
      base_block_storage_size:               <computed>
      cpu_count:                             <computed>
      create_date:                           <computed>
      internet_line_type.%:                  <computed>
      is_fee_charging_monitoring:            <computed>
      memory_size:                           <computed>
      platform_type.%:                       <computed>
      port_forwarding_external_port:         <computed>
      port_forwarding_internal_port:         <computed>
      port_forwarding_public_ip:             <computed>
      private_ip:                            <computed>
      public_ip:                             <computed>
      region.%:                              <computed>
      server_create_count:                   "1"
      server_image_name:                     <computed>
      server_image_product_code:             "SPSW0LINUX000032"
      server_instance_no:                    <computed>
      server_instance_operation.%:           <computed>
      server_instance_status.%:              <computed>
      server_instance_status_name:           <computed>
      server_name:                           "tf-test"
      server_product_code:                   "SPSVRSTAND000004"
      uptime:                                <computed>
      zone.%:                                <computed>


Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

ncloud_server.server: Creating...
  base_block_storage_disk_detail_type.%: "" => "<computed>"
  base_block_storage_disk_type.%:        "" => "<computed>"
  base_block_storage_size:               "" => "<computed>"
  cpu_count:                             "" => "<computed>"
  create_date:                           "" => "<computed>"
  internet_line_type.%:                  "" => "<computed>"
  is_fee_charging_monitoring:            "" => "<computed>"
  memory_size:                           "" => "<computed>"
  platform_type.%:                       "" => "<computed>"
  port_forwarding_external_port:         "" => "<computed>"
  port_forwarding_internal_port:         "" => "<computed>"
  port_forwarding_public_ip:             "" => "<computed>"
  private_ip:                            "" => "<computed>"
  public_ip:                             "" => "<computed>"
  region.%:                              "" => "<computed>"
  server_create_count:                   "" => "1"
  server_image_name:                     "" => "<computed>"
  server_image_product_code:             "" => "SPSW0LINUX000032"
  server_instance_no:                    "" => "<computed>"
  server_instance_operation.%:           "" => "<computed>"
  server_instance_status.%:              "" => "<computed>"
  server_instance_status_name:           "" => "<computed>"
  server_name:                           "" => "tf-test"
  server_product_code:                   "" => "SPSVRSTAND000004"
  uptime:                                "" => "<computed>"
  zone.%:                                "" => "<computed>"
ncloud_server.server: Still creating... (10s elapsed)
ncloud_server.server: Still creating... (20s elapsed)
ncloud_server.server: Still creating... (30s elapsed)
ncloud_server.server: Still creating... (40s elapsed)
ncloud_server.server: Still creating... (50s elapsed)
ncloud_server.server: Still creating... (1m0s elapsed)
ncloud_server.server: Still creating... (1m10s elapsed)
ncloud_server.server: Still creating... (1m20s elapsed)
ncloud_server.server: Still creating... (1m30s elapsed)
ncloud_server.server: Still creating... (1m40s elapsed)
ncloud_server.server: Still creating... (1m50s elapsed)
ncloud_server.server: Still creating... (2m0s elapsed)
ncloud_server.server: Still creating... (2m10s elapsed)
ncloud_server.server: Still creating... (2m20s elapsed)
ncloud_server.server: Still creating... (2m30s elapsed)
ncloud_server.server: Still creating... (2m40s elapsed)
ncloud_server.server: Still creating... (2m50s elapsed)
ncloud_server.server: Still creating... (3m0s elapsed)
ncloud_server.server: Still creating... (3m10s elapsed)
ncloud_server.server: Still creating... (3m20s elapsed)
ncloud_server.server: Still creating... (3m30s elapsed)
ncloud_server.server: Still creating... (3m40s elapsed)
ncloud_server.server: Still creating... (3m50s elapsed)
ncloud_server.server: Creation complete after 3m51s (ID: 841982)

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## Terraform show

Inspect Terraform state or plan

```
$ terraform show

ncloud_server.server:
  id = 841982
  base_block_storage_disk_detail_type.% = 2
  base_block_storage_disk_detail_type.code =
  base_block_storage_disk_detail_type.code_name =
  base_block_storage_disk_type.% = 2
  base_block_storage_disk_type.code = NET
  base_block_storage_disk_type.code_name = Network Storage
  base_block_storage_size = 53687091200
  cpu_count = 2
  create_date = 2018-07-01T13:01:53+0900
  internet_line_type.% = 2
  internet_line_type.code = PUBLC
  internet_line_type.code_name = PUBLC
  is_fee_charging_monitoring = false
  memory_size = 4294967296
  platform_type.% = 2
  platform_type.code = LNX32
  platform_type.code_name = Linux 32 Bit
  port_forwarding_public_ip = 106.10.41.173
  private_ip = 10.41.5.103
  public_ip =
  region.% = 3
  region.region_code = KR
  region.region_name = Korea
  region.region_no = 1
  server_create_count = 1
  server_image_name = centos-6.3-32
  server_image_product_code = SPSW0LINUX000032
  server_instance_no = 841982
  server_instance_operation.% = 2
  server_instance_operation.code = NULL
  server_instance_operation.code_name = Server NULL OP
  server_instance_status.% = 2
  server_instance_status.code = RUN
  server_instance_status.code_name = Server run state
  server_instance_status_name = running
  server_name = tf-test
  server_product_code = SPSVRSTAND000004
  uptime = 2018-07-01T13:05:39+0900
  user_data =
  zone.% = 3
  zone.zone_description = 평촌 zone
  zone.zone_name = KR-2
  zone.zone_no = 3
```

## Terraform destroy

Destroy Terraform-managed infrastructure

```
$ terraform destroy
ncloud_server.server: Refreshing state... (ID: 841982)

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  - ncloud_server.server


Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

ncloud_server.server: Destroying... (ID: 841982)
ncloud_server.server: Still destroying... (ID: 841982, 10s elapsed)
ncloud_server.server: Destruction complete after 14s

Destroy complete! Resources: 1 destroyed.
```

- 서버 생성 두대
- block storage 추가
- nginx 설치
- 공공 IP 추가
- lb 추가
