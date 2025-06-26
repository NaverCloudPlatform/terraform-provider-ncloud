package autoscaling

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudLaunchConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLaunchConfigurationCreate,
		Read:   resourceNcloudLaunchConfigurationRead,
		Delete: resourceNcloudLaunchConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"launch_configuration_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 255)),
			},
			"server_image_product_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"member_server_image_no"},
				ForceNew:      true,
			},
			"server_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"member_server_image_no": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"server_image_product_code"},
				ForceNew:      true,
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"init_script_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_encrypted_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudLaunchConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	id, err := createLaunchConfiguration(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	return resourceNcloudLaunchConfigurationRead(d, meta)
}

func createLaunchConfiguration(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &vautoscaling.CreateLaunchConfigurationRequest{
		RegionCode:                  &config.RegionCode,
		ServerImageProductCode:      StringPtrOrNil(d.GetOk("server_image_product_code")),
		MemberServerImageInstanceNo: StringPtrOrNil(d.GetOk("member_server_image_no")),
		ServerProductCode:           StringPtrOrNil(d.GetOk("server_product_code")),
		IsEncryptedVolume:           BoolPtrOrNil(d.GetOk("is_encrypted_volume")),
		InitScriptNo:                StringPtrOrNil(d.GetOk("init_script_no")),
		LaunchConfigurationName:     StringPtrOrNil(d.GetOk("name")),
		LoginKeyName:                StringPtrOrNil(d.GetOk("login_key_name")),
	}

	LogCommonRequest("createVpcLaunchConfiguration", reqParams)
	res, err := config.Client.Vautoscaling.V2Api.CreateLaunchConfiguration(reqParams)
	if err != nil {
		LogErrorResponse("createVpcLaunchConfiguration", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcLaunchConfiguration", res)
	return res.LaunchConfigurationList[0].LaunchConfigurationNo, nil
}

func resourceNcloudLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	launchConfig, err := GetLaunchConfiguration(config, d.Id())
	if err != nil {
		return err
	}

	if launchConfig == nil {
		d.SetId("")
		return nil
	}

	launchConfigMap := ConvertToMap(launchConfig)
	d.Set("server_image_product_code", launchConfig.ServerImageProductCode)
	SetSingularResourceDataFromMapSchema(ResourceNcloudLaunchConfiguration(), d, launchConfigMap)
	return nil
}

func GetLaunchConfiguration(config *conn.ProviderConfig, id string) (*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LaunchConfigurationNoList = []*string{ncloud.String(id)}
	}

	LogCommonRequest("getVpcLaunchConfiguration", reqParams)
	resp, err := config.Client.Vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcLaunchConfiguration", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcLaunchConfiguration", resp)

	if len(resp.LaunchConfigurationList) < 1 {
		return nil, nil
	}

	l := resp.LaunchConfigurationList[0]

	return &LaunchConfiguration{
		LaunchConfigurationName:     l.LaunchConfigurationName,
		ServerImageProductCode:      l.ServerImageProductCode,
		MemberServerImageInstanceNo: l.MemberServerImageInstanceNo,
		ServerProductCode:           l.ServerProductCode,
		LoginKeyName:                l.LoginKeyName,
		InitScriptNo:                l.InitScriptNo,
		IsEncryptedVolume:           l.IsEncryptedVolume,
		LaunchConfigurationNo:       l.LaunchConfigurationNo,
	}, nil
}

func resourceNcloudLaunchConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	err := deleteLaunchConfiguration(config, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func deleteLaunchConfiguration(config *conn.ProviderConfig, id string) error {
	reqParams := &vautoscaling.DeleteLaunchConfigurationRequest{
		LaunchConfigurationNo: ncloud.String(id),
	}

	LogCommonRequest("deleteVpcLaunchConfiguration", reqParams)
	res, err := config.Client.Vautoscaling.V2Api.DeleteLaunchConfiguration(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcLaunchConfiguration", err, reqParams)
		return err
	}
	LogResponse("deleteVpcLaunchConfiguration", res)
	return nil
}

type LaunchConfiguration struct {
	LaunchConfigurationNo       *string `json:"launch_configuration_no,omitempty"`
	LaunchConfigurationName     *string `json:"name,omitempty"`
	ServerImageProductCode      *string `json:"server_image_product_code,omitempty"`
	MemberServerImageInstanceNo *string `json:"member_server_image_no,omitempty"`
	ServerProductCode           *string `json:"server_product_code,omitempty"`
	LoginKeyName                *string `json:"login_key_name,omitempty"`
	InitScriptNo                *string `json:"init_script_no,omitempty"`
	IsEncryptedVolume           *bool   `json:"is_encrypted_volume,omitempty"`
}
