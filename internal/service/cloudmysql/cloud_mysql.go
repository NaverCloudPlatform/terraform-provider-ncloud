package cloudmysql

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	"log"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func ResourceNcloudMySql() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudMySqlCreate,
		Read:   resourceNcloudMySqlRead,
		Delete: resourceNcloudMySqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[ㄱ-ㅣ가-힣A-Za-z0-9-]+$`), "Composed of alphabets, numbers, hyphen (-)."),
				)),
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[a-z]+[a-z0-9-]+$`), "Composed of lowercase alphabets, numbers, hyphen (-)."),
				)),
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(4, 16),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]+.+`), "starts with an alphabets."),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_]+$`), "Composed of alphabets, numbers, hyphen (-), (\\), (_)."),
				)),
				Description: "mysql user id",
			},
			"user_password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(8, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9~!@#$%^*()\-_=\[\]\{\};:,.<>?]{8,20}$`), "Must Combine at least one each of alphabets, numbers, special characters except ` & + \\ \" ' / and white space."),
					validation.StringMatch(regexp.MustCompile(`.+[a-zA-Z]{1,}.+|.+[a-zA-Z]{1,}|[a-zA-Z]{1,}.+`), "Must have at least 1 alphabet."),
					validation.StringMatch(regexp.MustCompile(`.+[0-9]{1,}.+|.+[0-9]{1,}|[0-9]{1,}.+`), "Must have at least 1 Number."),
					validation.StringMatch(regexp.MustCompile(`.+[~!@#$%^*()\-_=\[\]\{\};:,.<>?].+|.+[~!@#$%^*()\-_=\[\]\{\};:,.<>?]|[~!@#$%^*()\-_=\[\]\{\};:,.<>?].+`), "Must have at least 1 special characters except ` & + \\ \" ' / and white space."),
				)),
				Description: "mysql user password",
			},
			"host_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 30),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_]+$`), ""),
				)),
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_version_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"data_storage_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "default: SSD",
			},
			"is_ha": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
				Description: "default: true",
			},
			"is_multi_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "default: false",
			},
			"is_storage_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "default: false",
			},
			"is_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "default: true",
			},
			"backup_file_retention_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: "default: 1(1 day)",
			},
			"backup_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ex) 01:15",
			},
			"is_automatic_backup": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "default: 3306",
			},
			"standby_master_subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudMySqlCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	id, err := createMysqlInstance(d, config)

	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	log.Printf("[INFO] Mysql instance ID: %s", d.Id())

	return resourceNcloudMySqlRead(d, meta)
}

func resourceNcloudMySqlRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	rs, err := GetMysqlInstance(config, d.Id())
	if err != nil {
		return nil
	}

	if rs == nil {
		d.SetId("")
		return nil
	}

	instance := ConvertToMap(rs)

	SetSingularResourceDataFromMapSchema(ResourceNcloudMySql(), d, instance)

	return nil
}

func resourceNcloudMySqlDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if err := deleteMysqlInstacne(config, d.Id()); err != nil {
		return err
	}

	//When an API is developed, uncomment below. and erase `time.Sleep(5 * time.Minute)`
	//if err := WaitForNcloudMysqlDeletion(config, d.Id()); err != nil {
	//	return err
	//}

	time.Sleep(3 * time.Minute)
	d.SetId("")
	return nil
}

func deleteMysqlInstacne(config *conn.ProviderConfig, id string) error {
	var err error

	if config.SupportVPC {
		err = deleteVpcMysqlInstance(config, id)
	} else {
		return NotSupportClassic("resource `ncloud_mysql`")
	}

	if err != nil {
		return err
	}

	return nil
}

func deleteVpcMysqlInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &vmysql.DeleteCloudMysqlInstanceRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(id),
	}
	LogCommonRequest("deleteVpcMysqlInstance", reqParams)

	resp, err := config.Client.Vmysql.V2Api.DeleteCloudMysqlInstance(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcMysqlInstance", err, reqParams)
		return err
	}
	LogResponse("deleteVpcMysqlInstance", resp)

	return nil
}

func createMysqlInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcMysqlInstance(d, config)
	}
	return nil, NotSupportClassic("resource `ncloud_mysql`")
}

func createVpcMysqlInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {

	subnet, err := vpc.GetSubnetInstance(config, d.Get("subnet_no").(string))
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		return nil, fmt.Errorf("no matching subnet(%s) found", d.Get("subnet_no"))
	}
	reqParams := &vmysql.CreateCloudMysqlInstanceRequest{
		RegionCode:                 &config.RegionCode,
		VpcNo:                      subnet.VpcNo,
		CloudMysqlImageProductCode: StringPtrOrNil(d.GetOk("image_product_code")),
		CloudMysqlProductCode:      StringPtrOrNil(d.GetOk("product_code")),
		DataStorageTypeCode:        StringPtrOrNil(d.GetOk("data_storage_type_code")),
		CloudMysqlServiceName:      StringPtrOrNil(d.GetOk("service_name")),
		CloudMysqlServerNamePrefix: StringPtrOrNil(d.GetOk("name_prefix")),
		CloudMysqlUserName:         StringPtrOrNil(d.GetOk("user_name")),
		CloudMysqlUserPassword:     StringPtrOrNil(d.GetOk("user_password")),
		HostIp:                     StringPtrOrNil(d.GetOk("host_ip")),
		CloudMysqlPort:             Int32PtrOrNil(d.GetOk("port")),
		CloudMysqlDatabaseName:     StringPtrOrNil(d.GetOk("database_name")),
		SubnetNo:                   subnet.SubnetNo,
	}

	if isHa:= d.Get("is_backup"); isHa == nil || isHa.(bool){
		reqParams.IsHa = ncloud.Bool(true)
	} else {
		reqParams.IsHa = ncloud.Bool(false)
	}

	if *reqParams.IsHa {
		if isMultiZone := d.Get("is_multi_zone"); isMultiZone == nil || !isMultiZone.(bool){
			reqParams.IsMultiZone = ncloud.Bool(false)

		} else{
			reqParams.IsMultiZone = ncloud.Bool(true)
			reqParams.StandbyMasterSubnetNo = ncloud.String(d.Get("standby_master_subnet_no").(string))
		}

		if isStorageEncryption := d.Get("is_storage_encryption"); isStorageEncryption==nil || !isStorageEncryption.(bool) {
			reqParams.IsStorageEncryption = ncloud.Bool(false)
		} else {
			reqParams.IsStorageEncryption = ncloud.Bool(true)
		}

		if isBackup := d.Get("is_backup"); isBackup != nil && !isBackup.(bool){
			return nil, fmt.Errorf("when is_ha is true, is_backup must be true")
		}

		reqParams.IsBackup = ncloud.Bool(true)
	}else {
		if isBackup:= d.Get("is_backup"); isBackup == nil {
			reqParams.IsBackup = ncloud.Bool(true)
		} else if isBackup.(bool) {
			reqParams.IsBackup = ncloud.Bool(true)
		} else {
			reqParams.IsBackup = ncloud.Bool(false)
		}
	}

	if *reqParams.IsBackup {
		if backupPeriod, ok := d.GetOk("backup_file_retention_period"); ok {
			reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(backupPeriod.(int)))
		}

		if isAutomaticBackup:= d.Get("is_automatic_backup"); isAutomaticBackup == nil || isAutomaticBackup.(bool){
			reqParams.IsAutomaticBackup = ncloud.Bool(true)
		} else {
			reqParams.IsAutomaticBackup = ncloud.Bool(false)
		}
	}

	if !(*reqParams.IsAutomaticBackup){
		if backupTime, ok := d.GetOk("backup_time"); ok {
			reqParams.BackupTime = ncloud.String(backupTime.(string))
		} else {
			return nil, fmt.Errorf("when is_automatic_backup is false, must input backup_time")
		}
	}

	LogCommonRequest("createVpcServerInstance", reqParams)
	resp, err := config.Client.Vmysql.V2Api.CreateCloudMysqlInstance(reqParams)
	if err != nil {
		LogErrorResponse("createVpcMysqlInstance", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcMysqlInstance", resp)
	mysqlInstance := resp.CloudMysqlInstanceList[0]

	if err := waitStateNcloudMysqlForCreation(config, *mysqlInstance.CloudMysqlInstanceNo); err != nil {
		return nil, err
	}

	return mysqlInstance.CloudMysqlInstanceNo, nil
}

func waitStateNcloudMysqlForCreation(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlInstance(config, id)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.CloudMysqlInstanceStatusName), nil
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for MysqlInstance state to be \"CREAT\": %s", err)
	}

	return nil
}

func WaitForNcloudMysqlDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"DEL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlInstance(config, id)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.CloudMysqlInstanceStatus.Code), nil
		},

		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for MysqlInstance state to be \"DEL\": %s", err)
	}

	return nil
}

func GetMysqlInstance(config *conn.ProviderConfig, id string) (*vmysql.CloudMysqlInstance, error) {
	if config.SupportVPC {
		return getVpcMysqlInstance(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_mysql`")
}

func getVpcMysqlInstance(config *conn.ProviderConfig, id string) (*vmysql.CloudMysqlInstance, error) {

	reqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("getVpcMysqlInstance", reqParams)
	resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getVpcMysqlInstance", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcMysqlInstance", resp)

	if len(resp.CloudMysqlInstanceList) == 0 {
		return nil, nil
	}
	if err := ValidateOneResult(len(resp.CloudMysqlInstanceList)); err != nil {
		return nil, err
	}

	return resp.CloudMysqlInstanceList[0], nil
}