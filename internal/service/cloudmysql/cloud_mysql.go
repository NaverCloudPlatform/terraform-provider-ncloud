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
		Update: resourceNcloudMySqlUpdate,
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
					// validation.StringMatch(regexp.MustCompile(`.*[^\\-]$`), "Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input."),
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
					// TODO: 영어, 숫자, 특문 1개 이상 포함 체크
					// validation.StringMatch(regexp.MustCompile(`[a-zA-Z]`)),
					// validation.StringMatch(regexp.MustCompile(`[0-9]`)),
					// validation.StringMatch(regexp.MustCompile("[^`&+\\\"'/' ']")),
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
				Optional: true, //G3인 경우, 필수
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
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "default: true",
			},
			"is_multi_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "default: false",
			},
			"is_storage_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "default: false",
			},
			"is_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "default: true",
			},
			"backup_file_retention_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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
				Computed: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				//todo: 3306 or min:10000 max:20000
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

	resource, err := GetMysqlInstance(config, d.Id())
	if err != nil {
		return nil
	}

	if resource == nil {
		d.SetId("")
		return nil
	}

	if config.SupportVPC {
		//todo: vpc 환경에서 해야할 작업들 여기에
	}

	instance := ConvertToMap(resource)

	SetSingularResourceDataFromMapSchema(ResourceNcloudMySql(), d, instance)

	return nil
}

func resourceNcloudMySqlUpdate(d *schema.ResourceData, meta interface{}) error {
	//config := meta.(*conn.ProviderConfig)

	// todo: 업데이트 추가

	return nil
}

func resourceNcloudMySqlDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if err := deleteMysqlInstacne(config, d.Id()); err != nil {
		return err
	}
	d.SetId("")

	if err := WaitForNcloudMysqlDeletion(config, d.Id()); err != nil {
		return err
	}

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
		LogErrorResponse("stopVpcMysqlInstance", err, reqParams)
		return err
	}
	LogResponse("stopVpcMysqlInstance", resp)

	return nil
}

func createMysqlInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcMysqlInstance(d, config)
	}
	// not support classic
	return nil, NotSupportClassic("resource `ncloud_mysql`")
}

func createVpcMysqlInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	//if _, ok := d.GetOk("subnet_no"); !ok{
	//	return nil, ErrorRequiredArgOnVpc("subnet_no")
	//}

	subnet, err := vpc.GetSubnetInstance(config, d.Get("subnet_no").(string))
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		return nil, fmt.Errorf("no matching subnet(%s) found", d.Get("subnet_no"))
	}

	// todo: reqPrarms 안에 컨디션 조건 체크 후 넣을 같 넣기
	reqParams := &vmysql.CreateCloudMysqlInstanceRequest{
		RegionCode:                 &config.RegionCode,
		VpcNo:                      subnet.VpcNo,
		CloudMysqlImageProductCode: StringPtrOrNil(d.GetOk("image_product_code")),
		CloudMysqlProductCode:      StringPtrOrNil(d.GetOk("product_code")),
		DataStorageTypeCode:        StringPtrOrNil(d.GetOk("data_storage_type_code")),
		// TODO: isha = false -> isMultiZone, standby 사용 X
		IsHa: BoolPtrOrNil(d.GetOk("is_ha")),
		// ToDo: isHa 가 true 일때, 멀티존 여부 선택가능
		//IsMultiZone:					BoolPtrOrNil(d.GetOk("is_multi_zone")),
		// isha = True 일때만 입력 받음
		//IsStorageEncryption:			BoolPtrOrNil(d.GetOk("is_storage_encryption")),
		// todo: isha = True -> isBackup =true
		//IsBackup:						BoolPtrOrNil(d.GetOk("is_backup")),
		//BackupFileRetentionPeriod:		Int32PtrOrNil(d.GetOk("backup_file_retention_period")),
		// todo: - 백업이 수행되는 시간을 설정, 백업 여부(isBackup)가 true이고 자동 백업 여부(isAutomaticBackup)가 false이면 반드시 입력
		//BackupTime:						StringPtrOrNil(d.GetOk("backup_time")),
		// todo: 자동으로 백업 시간을 설정할지에 대한 여부 선택, 자동 백업 여부(isAutomaticBackup)가 true이면 backupTime 입력 불가
		//IsAutomaticBackup:				BoolPtrOrNil(d.GetOk("is_automatic_backup")),
		CloudMysqlServiceName:      StringPtrOrNil(d.GetOk("service_name")),
		CloudMysqlServerNamePrefix: StringPtrOrNil(d.GetOk("name_prefix")),
		CloudMysqlUserName:         StringPtrOrNil(d.GetOk("user_name")),
		CloudMysqlUserPassword:     StringPtrOrNil(d.GetOk("user_password")),
		HostIp:                     StringPtrOrNil(d.GetOk("host_ip")),
		CloudMysqlPort:             Int32PtrOrNil(d.GetOk("port")),
		CloudMysqlDatabaseName:     StringPtrOrNil(d.GetOk("database_name")),
		SubnetNo:                   subnet.SubnetNo,
		//todo: Standby Master 서버의 Subnet 번호
		//- 멀티존 여부(isMultiZone)가 false이면 입력받지 않으며 멀티존 여부(isMultiZone)가 true이면 반드시 입력
		//- standbyMasterSubnetNo는 Master 서버의 Subnet과 Zone이 달라야 하며 같은 Public이거나 Private이어야만 함
		//- getCloudMysqlTargetSubnetList 액션을 통해서 획득 가능
		//StandbyMasterSubnetNo:			StringPtrOrNil(d.GetOk("standby_master_subnet_no")),
	}

	// todo: 아래 가독성 좋게 바꾸기
	if is_ha, ok := d.GetOk("is_ha"); !ok || is_ha.(bool) == true {
		//is_ha 없음 || is_ha == true -> is_ha == true
		reqParams.IsHa = ncloud.Bool(true)

		if isMultiZone, ok := d.GetOk("is_multi_zone"); ok {
			reqParams.IsMultiZone = ncloud.Bool(isMultiZone.(bool))

			// 멀티존 true -> standby sever
			if isMultiZone.(bool) {

			}
		}

		if isStorageEncryption, ok := d.GetOk("is_storage_encryption"); ok {
			reqParams.IsStorageEncryption = ncloud.Bool(isStorageEncryption.(bool))
		}

		reqParams.IsBackup = ncloud.Bool(true)
	} else {

		if isBackup, ok := d.GetOk("is_backup"); ok {
			reqParams.IsBackup = ncloud.Bool(isBackup.(bool))

			if isBackup.(bool) {
				if backupPeriod, ok := d.GetOk("backup_file_retention_period"); ok {
					reqParams.BackupFileRetentionPeriod = ncloud.Int32(backupPeriod.(int32))
				}

				if isAutomaticBackup, ok := d.GetOk("is_automatic_backup"); ok {
					reqParams.IsAutomaticBackup = ncloud.Bool(isAutomaticBackup.(bool))

					if !isAutomaticBackup.(bool) {
						backupTime := d.Get("backup_time")
						reqParams.BackupTime = ncloud.String(backupTime.(string))
					}
				}
			}
		}
	}

	// todo: reqParam 추가할거 체크 후 추가하기
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
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
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
		return fmt.Errorf("error waiting for MysqlInstance state to be \"RUN\": %s", err)
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

func createClassicMysqlInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	return nil, nil
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

// Todo: convert 사용해서 구현하는게 좋은지 여쭤보기
func convertVpcMysqlInstance(r *vmysql.CloudMysqlInstance) *CloudMysqlInstance {
	if r == nil {
		return nil
	}
	instance := &CloudMysqlInstance{
		CloudMysqlInstanceNo:         r.CloudMysqlInstanceNo,
		CloudMysqlServiceName:        r.CloudMysqlServiceName,
		CloudMysqlInstanceStatusName: r.CloudMysqlInstanceStatusName,
		CloudMysqlInstanceStatus:     r.CloudMysqlInstanceStatus.Code,
		CloudMysqlInstanceOperation:  r.CloudMysqlInstanceOperation.Code,
		CloudMysqlImageProductCode:   r.CloudMysqlImageProductCode,
		EngineVersion:                r.EngineVersion,
		License:                      r.License.Code,
		CloudMysqlPort:               r.CloudMysqlPort,
		IsHa:                         r.IsHa,
		IsMultiZone:                  r.IsMultiZone,
		IsBackup:                     r.IsBackup,
		BackupFileRetentionPeriod:    r.BackupFileRetentionPeriod,
		BackupTime:                   r.BackupTime,
		CreateDate:                   r.CreateDate,
		AccessControlGroupNoList:     r.AccessControlGroupNoList,
		CloudMysqlConfigList:         r.CloudMysqlConfigList,
		CloudMysqlServerInstanceList: r.CloudMysqlServerInstanceList,
	}

	return instance
}

type CloudMysqlInstance struct {
	// CloudMysql인스턴스번호
	CloudMysqlInstanceNo *string `json:"cloudMysqlInstanceNo,omitempty"`

	// CloudMysql서비스이름
	CloudMysqlServiceName *string `json:"cloudMysqlServiceName,omitempty"`

	// CloudMysql인스턴스상태이름
	CloudMysqlInstanceStatusName *string `json:"cloudMysqlInstanceStatusName,omitempty"`

	// CloudMysql인스턴스상태
	CloudMysqlInstanceStatus *string `json:"cloudMysqlInstanceStatus,omitempty"`

	// CloudMysql인스턴스OP
	CloudMysqlInstanceOperation *string `json:"cloudMysqlInstanceOperation,omitempty"`

	// CloudMysql이미지상품코드
	CloudMysqlImageProductCode *string `json:"cloudMysqlImageProductCode,omitempty"`

	// CloudMysql엔진버전
	EngineVersion *string `json:"engineVersion,omitempty"`

	// CloudMysql라이선스
	License *string `json:"license,omitempty"`

	// CloudMysql포트
	CloudMysqlPort *int32 `json:"cloudMysqlPort,omitempty"`

	// 고가용성여부
	IsHa *bool `json:"isHa,omitempty"`

	// 멀티존여부
	IsMultiZone *bool `json:"isMultiZone,omitempty"`

	// 백업여부
	IsBackup *bool `json:"isBackup,omitempty"`

	// 백업파일보관기간
	BackupFileRetentionPeriod *int32 `json:"backupFileRetentionPeriod,omitempty"`

	// 백업시간
	BackupTime *string `json:"backupTime,omitempty"`

	// 생성일자
	CreateDate *string `json:"createDate,omitempty"`

	// ACG번호리스트
	AccessControlGroupNoList []*string `json:"accessControlGroupNoList,omitempty"`

	// CloudMysqlConfig리스트
	CloudMysqlConfigList []*string `json:"cloudMysqlConfigList,omitempty"`

	// CloudMysql서버인스턴스리스트
	CloudMysqlServerInstanceList []*vmysql.CloudMysqlServerInstance `json:"cloudMysqlServerInstanceList,omitempty"`
}

type CloudMysqlServerInstance struct {

	// CloudMysql서버인스턴스번호
	CloudMysqlServerInstanceNo *string `json:"cloudMysqlServerInstanceNo,omitempty"`

	// CloudMysql서버이름
	CloudMysqlServerName *string `json:"cloudMysqlServerName,omitempty"`

	// CloudMysql서버역할
	CloudMysqlServerRole *string `json:"cloudMysqlServerRole,omitempty"`

	// CloudMysql인스턴스상태이름
	CloudMysqlServerInstanceStatusName *string `json:"cloudMysqlServerInstanceStatusName,omitempty"`

	// CloudMysql서버인스턴스상태
	CloudMysqlServerInstanceStatus *string `json:"cloudMysqlServerInstanceStatus,omitempty"`

	// CloudMysql서버인스턴스OP
	CloudMysqlServerInstanceOperation *string `json:"cloudMysqlServerInstanceOperation,omitempty"`

	// CloudMysql상품코드
	CloudMysqlProductCode *string `json:"cloudMysqlProductCode,omitempty"`

	// REGION코드
	RegionCode *string `json:"regionCode,omitempty"`

	// ZONE코드
	ZoneCode *string `json:"zoneCode,omitempty"`

	// VPC번호
	VpcNo *string `json:"vpcNo,omitempty"`

	// Subnet번호
	SubnetNo *string `json:"subnetNo,omitempty"`

	// PublicSubnet여부
	IsPublicSubnet *bool `json:"isPublicSubnet,omitempty"`

	// 공인도메인명
	PublicDomain *string `json:"publicDomain,omitempty"`

	// 사설도메인명
	PrivateDomain *string `json:"privateDomain,omitempty"`

	// 데이터스토리지타입
	DataStorageType *string `json:"dataStorageType,omitempty"`

	// 데이터스토리지암호화여부
	IsStorageEncryption *bool `json:"isStorageEncryption,omitempty"`

	// 데이터스토리지사이즈
	DataStorageSize *int64 `json:"dataStorageSize,omitempty"`

	// 사용중인데이터스토리지사이즈
	UsedDataStorageSize *int64 `json:"usedDataStorageSize,omitempty"`

	// virtualCPU개수
	CpuCount *int32 `json:"cpuCount,omitempty"`

	// 메모리사이즈
	MemorySize *int64 `json:"memorySize,omitempty"`

	// 업시간
	Uptime *string `json:"uptime,omitempty"`

	// 생성일자
	CreateDate *string `json:"createDate,omitempty"`
}
