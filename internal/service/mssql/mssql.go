package cloudmssql

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudMssql() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudMssqlCreate,
		ReadContext:   resourceNcloudMssqlRead,
		DeleteContext: resourceNcloudMssqlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(3 * conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC Number of Cloud DB for MSSQL instance.",
				ForceNew:    true,
			},
			"subnet_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet Number of Cloud DB for MSSQL instance.",
				ForceNew:    true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 15),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\-가-힣]+$`), "Composed of alphabets, numbers, korean, hyphen (-)."),
				)),
				Description: "Name of Cloud DB for MSSQL instance.",
				ForceNew:    true,
			},
			"is_ha": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Choice of High Availability.",
				ForceNew:    true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(4, 16),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]+.+`), "starts with an alphabets."),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_]+$`), "Composed of alphabets, numbers, hyphen (-), (\\), (_)."),
				)),
				Description: "Access username, which will be used for DB admin.",
				ForceNew:    true,
			},
			"user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(8, 20),
					validation.StringMatch(regexp.MustCompile(`[a-zA-Z]+`), "Must have at least one alphabet"),
					validation.StringMatch(regexp.MustCompile(`\d+`), "Must have at least one number"),
					validation.StringMatch(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
					validation.StringMatch(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
				)),
				Description: "Access password for user, which will be used for DB admin.",
				ForceNew:    true,
			},
			"mirror_subnet_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet number of Mirror Server. Required when isMultiZone is true.",
				ForceNew:    true,
			},
			"config_group_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Config Group number of Cloud DB for MSSQL instance.",
				Default:     0,
				ForceNew:    true,
			},
			"image_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image Product Code of Cloud DB for MSSQL instance.",
				ForceNew:    true,
			},
			"product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Product Code of Cloud DB for MSSQL instance.",
				ForceNew:    true,
			},
			"data_storage_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HDD", "SSD"}, false)),
				Description:      "Data Storage Type Code.",
				Default:          "SSD",
				ForceNew:         true,
			},
			"is_multi_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Multi Zone option. Required when isHa is true.",
				Default:     false,
				ForceNew:    true,
			},
			"backup_file_retention_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 30)),
				Description:      "Retention period of back-up files.",
				Default:          1,
				ForceNew:         true,
			},
			"backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringMatch(regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3])([0-5][0-9])$`), "Must be in the format HHMM."),
					validation.StringMatch(regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3])([0-5][0-9])(00|15|30|45)$`), "Must be in 15-minute intervals."),
				)),
				Description: "Back-up time. Required when isAutomaticBackup is false.",
				ForceNew:    true,
			},
			"is_automatic_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Automatic backup time.",
				ForceNew:    true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.Any(
					//validation.
					validation.IntBetween(10000, 20000),
					validation.IntBetween(1433, 1433),
				)),
				Description: "Port of Cloud DB for MSSQL instance.",
				Default:     1433,
				ForceNew:    true,
			},
			"character_set_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DB character set.",
				Default:     "Korean_Wansung_CI_AS",
				ForceNew:    true,
			},
		},
	}
}

func resourceNcloudMssqlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	client := meta.(*conn.ProviderConfig).Client
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cloud_mssql`"))
	}

	reqParams := &vmssql.CreateCloudMssqlInstanceRequest{
		RegionCode:                 &config.RegionCode,
		VpcNo:                      StringPtrOrNil(d.GetOk("vpc_no")),
		SubnetNo:                   StringPtrOrNil(d.GetOk("subnet_no")),
		CloudMssqlServiceName:      StringPtrOrNil(d.GetOk("service_name")),
		ConfigGroupNo:              StringPtrOrNil(d.GetOk("config_group_no")),
		CloudMssqlImageProductCode: StringPtrOrNil(d.GetOk("image_product_code")),
		CloudMssqlProductCode:      StringPtrOrNil(d.GetOk("product_code")),
		DataStorageTypeCode:        StringPtrOrNil(d.GetOk("data_storage_type_code")),
		IsHa:                       BoolPtrOrNil(d.GetOk("is_ha")),
		BackupFileRetentionPeriod:  Int32PtrOrNil(d.GetOk("backup_file_retention_period")),
		BackupTime:                 StringPtrOrNil(d.GetOk("backup_time")),
		IsAutomaticBackup:          BoolPtrOrNil(d.GetOk("is_automatic_backup")),
		CloudMssqlUserName:         StringPtrOrNil(d.GetOk("user_name")),
		CloudMssqlUserPassword:     StringPtrOrNil(d.GetOk("user_password")),
		CloudMssqlPort:             Int32PtrOrNil(d.GetOk("port")),
		CharacterSetName:           StringPtrOrNil(d.GetOk("character_set_name")),
	}

	LogCommonRequest("CreateCloudMssqlInstance", reqParams)
	resp, err := client.Vmssql.V2Api.CreateCloudMssqlInstance(reqParams)
	if err != nil {
		LogErrorResponse("CreateCloudMssqlInstance", err, reqParams)
		return diag.FromErr(err)
	}
	LogCommonResponse("CreateCloudMssqlInstance", GetCommonResponse(resp))

	if err := waitForCloudMssqlActive(ctx, d, config, ncloud.StringValue(resp.CloudMssqlInstanceList[0].CloudMssqlInstanceNo)); err != nil {
		return diag.FromErr(err)
	}

	cloudMssqlInstance := resp.CloudMssqlInstanceList[0]
	d.SetId(*cloudMssqlInstance.CloudMssqlInstanceNo)

	return resourceNcloudMssqlRead(ctx, d, meta)
}

func resourceNcloudMssqlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cloud_mssql`"))
	}

	ms, err := GetCloudMssqlInstance(config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if ms == nil {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found

		return nil
	}

	convertedInstance := ConvertToMap(ms)
	SetSingularResourceDataFromMapSchema(ResourceNcloudMssql(), d, convertedInstance)

	return nil
}

func resourceNcloudMssqlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	client := meta.(*conn.ProviderConfig).Client
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cloud_mssql`"))
	}

	deleteInstanceReqParams := &vmssql.DeleteCloudMssqlInstanceRequest{
		RegionCode:           &config.RegionCode,
		CloudMssqlInstanceNo: ncloud.String(d.Id()),
	}

	if err := waitForCloudMssqlActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	LogCommonRequest("resourceNcloudMssqlDelete", deleteInstanceReqParams)
	resp, err := client.Vmssql.V2Api.DeleteCloudMssqlInstance(deleteInstanceReqParams)
	if err != nil {
		LogErrorResponse("resourceNcloudMssqlDelete", err, deleteInstanceReqParams)
		return diag.FromErr(err)
	}
	LogResponse("resourceNcloudMssqlDelete", resp)

	if err := waitForCloudMssqlDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCloudMssqlDeletion(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DEL"},
		Target:  []string{"NULL"},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vmssql.GetCloudMssqlInstanceListRequest{
				RegionCode:               &config.RegionCode,
				CloudMssqlInstanceNoList: []*string{ncloud.String(d.Id())},
			}
			resp, err := config.Client.Vmssql.V2Api.GetCloudMssqlInstanceList(reqParams)
			if err != nil {
				return nil, "", err
			}

			if len(resp.CloudMssqlInstanceList) < 1 {
				return resp, "NULL", nil
			}

			ms := resp.CloudMssqlInstanceList[0]
			status := ncloud.StringValue(ms.CloudMssqlInstanceOperation.Code)

			return ms, status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for Cloud Mssql instance (%s) to become terminating: %s", d.Id(), err)
	}

	return nil
}

func GetCloudMssqlInstance(config *conn.ProviderConfig, cloudMssqlInstanceNo string) (*vmssql.CloudMssqlInstance, error) {
	reqParams := &vmssql.GetCloudMssqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMssqlInstanceNo: ncloud.String(cloudMssqlInstanceNo),
	}
	LogCommonRequest("GetCloudMssqlInstanceList", reqParams)
	resp, err := config.Client.Vmssql.V2Api.GetCloudMssqlInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("GetCloudMssqlInstanceList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("GetCloudMssqlInstanceList", GetCommonResponse(resp))

	if len(resp.CloudMssqlInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudMssqlInstanceList[0], nil
}

func waitForCloudMssqlActive(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			reqParams := &vmssql.GetCloudMssqlInstanceDetailRequest{
				RegionCode:           &config.RegionCode,
				CloudMssqlInstanceNo: ncloud.String(id),
			}
			resp, err := config.Client.Vmssql.V2Api.GetCloudMssqlInstanceDetail(reqParams)
			if err != nil {
				return nil, "", err
			}

			if len(resp.CloudMssqlInstanceList) < 1 {
				return nil, "", fmt.Errorf("not found cloud mssql instance(%s)", id)
			}

			ms := resp.CloudMssqlInstanceList[0]

			status := ms.CloudMssqlInstanceStatus.Code
			op := ms.CloudMssqlInstanceOperation.Code

			if *status == "INIT" && *op == "CREAT" {
				return ms, "creating", nil
			}

			if *status == "CREAT" && *op == "SETUP" {
				return ms, "settingUp", nil
			}

			if *status == "CREAT" && *op == "NULL" {
				return ms, "running", nil
			}

			return resp, ncloud.StringValue(ms.CloudMssqlInstanceOperation.Code), nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for Cloud Mssql instance (%s) to become activating: %s", id, err)
	}

	return nil
}

func convertCloudMssql(instance *vmssql.CloudMssqlInstance) *CloudMssqlInstance {
	return &CloudMssqlInstance{
		CloudMssqlInstanceNo:         instance.CloudMssqlInstanceNo,
		CloudMssqlServiceName:        instance.CloudMssqlServiceName,
		CloudMssqlInstanceStatusName: instance.CloudMssqlInstanceStatusName,
		CloudMssqlImageProductCode:   instance.CloudMssqlImageProductCode,
		IsHa:                         instance.IsHa,
		CloudMssqlPort:               instance.CloudMssqlPort,
		BackupFileRetentionPeriod:    instance.BackupFileRetentionPeriod,
		BackupTime:                   instance.BackupTime,
		ConfigGroupNo:                instance.ConfigGroupNo,
		EngineVersion:                instance.EngineVersion,
		CreateDate:                   instance.CreateDate,
		DbCollation:                  instance.DbCollation,
		AccessControlGroupNoList:     instance.AccessControlGroupNoList,
	}
}
