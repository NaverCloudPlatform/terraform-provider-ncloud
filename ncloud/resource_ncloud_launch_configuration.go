package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterResource("ncloud_launch_configuration", resourceNcloudLaunchConfiguration())
}

func resourceNcloudLaunchConfiguration() *schema.Resource {
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_image_product_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"member_server_image_no"},
			},
			"server_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"member_server_image_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"server_image_product_code"},
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_control_group_configuration_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_encrypted_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"init_script_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNcloudLaunchConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	id, err := createLaunchConfiguration(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	return resourceNcloudLaunchConfigurationRead(d, meta)
}

func createLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcLaunchConfiguration(d, config)
	} else {
		return createClassicLaunchConfiguration(d, config)
	}
}

func createVpcLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
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

	logCommonRequest("createVpcLaunchConfiguration", reqParams)
	res, err := config.Client.vautoscaling.V2Api.CreateLaunchConfiguration(reqParams)
	if err != nil {
		logErrorResponse("createVpcLaunchConfiguration", err, reqParams)
		return nil, err
	}
	logResponse("createVpcLaunchConfiguration", res)
	return res.LaunchConfigurationList[0].LaunchConfigurationNo, nil
}

func createClassicLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	reqParams := &autoscaling.CreateLaunchConfigurationRequest{
		LaunchConfigurationName: StringPtrOrNil(d.GetOk("name")),
		ServerImageProductCode:  StringPtrOrNil(d.GetOk("server_image_product_code")),
		ServerProductCode:       StringPtrOrNil(d.GetOk("server_product_code")),
		MemberServerImageNo:     StringPtrOrNil(d.GetOk("member_server_image_no")),
		LoginKeyName:            StringPtrOrNil(d.GetOk("login_key_name")),
		UserData:                StringPtrOrNil(d.GetOk("user_data")),
		RegionNo:                &config.RegionNo,
	}

	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		reqParams.AccessControlGroupConfigurationNoList = expandStringInterfaceList(param.([]interface{}))
	}

	logCommonRequest("createClassicLaunchConfiguration", reqParams)
	res, err := config.Client.autoscaling.V2Api.CreateLaunchConfiguration(reqParams)
	if err != nil {
		logErrorResponse("createClassicLaunchConfiguration", err, reqParams)
		return nil, err
	}
	logResponse("createClassicLaunchConfiguration", res)
	return res.LaunchConfigurationList[0].LaunchConfigurationNo, nil
}

func resourceNcloudLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	launchConfig, err := getLaunchConfiguration(config, d.Id())
	if err != nil {
		return err
	}

	if launchConfig == nil {
		d.SetId("")
		return nil
	}

	launchConfigMap := ConvertToMap(launchConfig)
	SetSingularResourceDataFromMapSchema(resourceNcloudLaunchConfiguration(), d, launchConfigMap)
	return nil
}

func getLaunchConfiguration(config *ProviderConfig, id string) (*LaunchConfiguration, error) {
	if config.SupportVPC {
		return getVpcLaunchConfiguration(config, id)
	} else {
		return getClassicLaunchConfiguration(config, id)
	}
}

func getVpcLaunchConfiguration(config *ProviderConfig, id string) (*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LaunchConfigurationNoList = []*string{ncloud.String(id)}
	}

	logCommonRequest("getVpcLaunchConfiguration", reqParams)
	resp, err := config.Client.vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getVpcLaunchConfiguration", err, reqParams)
		return nil, err
	}
	logResponse("getVpcLaunchConfiguration", resp)

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

func getClassicLaunchConfiguration(config *ProviderConfig, id string) (*LaunchConfiguration, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}
	logCommonRequest("getClassicLaunchConfiguration", reqParams)
	resp, err := config.Client.autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getClassicLaunchConfiguration", err, reqParams)
		return nil, err
	}
	logResponse("getClassicLaunchConfiguration", resp)

	for _, l := range resp.LaunchConfigurationList {
		if *l.LaunchConfigurationNo == *no {
			return &LaunchConfiguration{
				LaunchConfigurationNo:       l.LaunchConfigurationNo,
				LaunchConfigurationName:     l.LaunchConfigurationName,
				ServerImageProductCode:      l.ServerImageProductCode,
				MemberServerImageInstanceNo: l.MemberServerImageNo,
				ServerProductCode:           l.ServerProductCode,
				LoginKeyName:                l.LoginKeyName,
			}, nil
		}
	}

	return nil, nil
}

func resourceNcloudLaunchConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	err := deleteLaunchConfiguration(config, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func deleteLaunchConfiguration(config *ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcLaunchConfiguration(config, id)
	} else {
		return deleteClassicLaunchConfiguration(config, id)
	}
}

func deleteVpcLaunchConfiguration(config *ProviderConfig, id string) error {
	reqParams := &vautoscaling.DeleteLaunchConfigurationRequest{
		LaunchConfigurationNo: ncloud.String(id),
	}

	logCommonRequest("deleteVpcLaunchConfiguration", reqParams)
	res, err := config.Client.vautoscaling.V2Api.DeleteLaunchConfiguration(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcLaunchConfiguration", err, reqParams)
		return err
	}
	logResponse("deleteVpcLaunchConfiguration", res)
	return nil
}

func deleteClassicLaunchConfiguration(config *ProviderConfig, id string) error {
	launchConfig, err := getClassicLaunchConfiguration(config, id)
	if err != nil {
		return err
	}

	if launchConfig == nil {
		return nil
	}

	reqParams := &autoscaling.DeleteAutoScalingLaunchConfigurationRequest{
		LaunchConfigurationName: launchConfig.LaunchConfigurationName,
	}

	logCommonRequest("deleteClassicLaunchConfiguration", reqParams)
	res, err := config.Client.autoscaling.V2Api.DeleteAutoScalingLaunchConfiguration(reqParams)
	if err != nil {
		logErrorResponse("deleteClassicLaunchConfiguration", err, reqParams)
		return err
	}
	logResponse("deleteClassicLaunchConfiguration", res)

	return nil
}

func getClassicLaunchConfigurationNameByNo(no *string, config *ProviderConfig) (*LaunchConfiguration, error) {
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}
	resp, err := config.Client.autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		return nil, err
	}

	for _, l := range resp.LaunchConfigurationList {
		if *l.LaunchConfigurationNo == *no {
			return &LaunchConfiguration{
				LaunchConfigurationNo:       l.LaunchConfigurationNo,
				LaunchConfigurationName:     l.LaunchConfigurationName,
				ServerImageProductCode:      l.ServerImageProductCode,
				MemberServerImageInstanceNo: l.MemberServerImageNo,
				ServerProductCode:           l.ServerProductCode,
				LoginKeyName:                l.LoginKeyName,
			}, nil
		}
	}
	return nil, fmt.Errorf("Not found LaunchConfiguration(%s)", ncloud.StringValue(no))
}

type LaunchConfiguration struct {
	LaunchConfigurationNo       *string `json:"launch_configuration_no,omitempty,omitempty"`
	LaunchConfigurationName     *string `json:"name,omitempty"`
	ServerImageProductCode      *string `json:"server_image_product_code,omitempty"`
	MemberServerImageInstanceNo *string `json:"member_server_image_no,omitempty"`
	ServerProductCode           *string `json:"server_product_code,omitempty"`
	LoginKeyName                *string `json:"login_key_name,omitempty"`
	InitScriptNo                *string `json:"init_script_no,omitempty"`
	IsEncryptedVolume           *bool   `json:"is_encrypted_volume,omitempty"`
}
