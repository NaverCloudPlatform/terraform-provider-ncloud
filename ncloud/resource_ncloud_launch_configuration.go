package ncloud

import (
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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_image_product_code": {
				Type:          schema.TypeString,
				Optional:      true,
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
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
			"launch_configuration_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudLaunchConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var err error
	var id *string
	if config.SupportVPC {
		id, err = createVpcLaunchConfiguration(d, config)
	} else {
		id, err = createClassicLaunchConfiguration(d, config)
	}
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	d.Set("name", d.Id())
	return resourceNcloudLaunchConfigurationRead(d, meta)
}

func resourceNcloudLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var err error

	var launchConfig *LaunchConfiguration
	if config.SupportVPC {
		launchConfig, err = getVpcLaunchConfiguration(d, config)
	} else {
		launchConfig, err = getClassicLaunchConfiguration(d, config)
	}
	launchConfigMap := ConvertToMap(launchConfig)
	if err != nil {
		return err
	}

	SetSingularResourceDataFromMapSchema(resourceNcloudLaunchConfiguration(), d, launchConfigMap)
	return nil
}

func resourceNcloudLaunchConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var err error
	if config.SupportVPC {
		var launchConfig *LaunchConfiguration
		if launchConfig, err = getVpcLaunchConfiguration(d, config); err != nil {
			return err
		}

		err = deleteVpcLaunchConfiguration(config, launchConfig.LaunchConfigurationNo)
	} else {
		err = deleteClassicLaunchConfiguration(d, config)
	}
	if err != nil {
		return err
	}

	return nil
}

func deleteClassicLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &autoscaling.DeleteAutoScalingLaunchConfigurationRequest{
		LaunchConfigurationName: ncloud.String(d.Id()),
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

func deleteVpcLaunchConfiguration(config *ProviderConfig, launchConfigNo *string) error {
	reqParams := &vautoscaling.DeleteLaunchConfigurationRequest{
		LaunchConfigurationNo: launchConfigNo,
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
	return res.LaunchConfigurationList[0].LaunchConfigurationName, nil
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
	return res.LaunchConfigurationList[0].LaunchConfigurationName, nil
}

func getClassicLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) (*LaunchConfiguration, error) {
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{}
	if v, ok := d.GetOk("name"); ok {
		reqParams.LaunchConfigurationNameList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getClassicLaunchConfiguration", reqParams)
	res, err := config.Client.autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getClassicLaunchConfiguration", err, reqParams)
		return nil, err
	}
	logResponse("getClassicLaunchConfiguration", res)

	if err := validateOneResult(len(res.LaunchConfigurationList)); err != nil {
		return nil, err
	}
	configuration := res.LaunchConfigurationList[0]
	return &LaunchConfiguration{
		LaunchConfigurationName:     configuration.LaunchConfigurationName,
		ServerImageProductCode:      configuration.ServerImageProductCode,
		MemberServerImageInstanceNo: configuration.MemberServerImageNo,
		ServerProductCode:           configuration.ServerProductCode,
		LoginKeyName:                configuration.LoginKeyName,
		Region:                      &config.RegionCode,
	}, nil
}

func getVpcLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) (*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode:                  &config.RegionCode,
		LaunchConfigurationNameList: []*string{StringPtrOrNil(d.GetOk("name"))},
	}

	logCommonRequest("getVpcLaunchConfiguration", reqParams)
	res, err := config.Client.vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getVpcLaunchConfiguration", err, reqParams)
		return nil, err
	}
	logResponse("getVpcLaunchConfiguration", res)

	if err := validateOneResult(len(res.LaunchConfigurationList)); err != nil {
		return nil, err
	}
	configuration := res.LaunchConfigurationList[0]
	return &LaunchConfiguration{
		LaunchConfigurationName:     configuration.LaunchConfigurationName,
		ServerImageProductCode:      configuration.ServerImageProductCode,
		MemberServerImageInstanceNo: configuration.MemberServerImageInstanceNo,
		ServerProductCode:           configuration.ServerProductCode,
		LoginKeyName:                configuration.LoginKeyName,
		Region:                      configuration.RegionCode,
		Status:                      configuration.LaunchConfigurationStatus.Code,
		InitScriptNo:                configuration.InitScriptNo,
		IsEncryptedVolume:           configuration.IsEncryptedVolume,
		LaunchConfigurationNo:       configuration.LaunchConfigurationNo,
	}, nil
}

type LaunchConfiguration struct {
	LaunchConfigurationName     *string `json:"name,omitempty"`
	ServerImageProductCode      *string `json:"server_image_product_code,omitempty"`
	MemberServerImageInstanceNo *string `json:"member_server_image_no,omitempty"`
	ServerProductCode           *string `json:"server_product_code,omitempty"`
	LoginKeyName                *string `json:"login_key_name,omitempty"`
	Region                      *string `json:"region,omitempty"`
	Status                      *string `json:"status,omitempty"`
	InitScriptNo                *string `json:"init_script_no,omitempty"`
	IsEncryptedVolume           *bool   `json:"is_encrypted_volume,omitempty"`
	LaunchConfigurationNo       *string `json:"launch_configuration_no,omitempty,omitempty"`
}
