package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterResource("ncloud_launch_configuration", resourceNcloudLaunchConfiguration())
}

func resourceNcloudLaunchConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLaunchConfigurationCreate,
		Read:   resourceNcloudLaunchConfigurationRead,
		Update: nil,
		Delete: resourceNcloudLaunchConfigurationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"server_image_product_code": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"member_server_image_no"},
			},
			"server_product_code": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"member_server_image_no": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"server_image_product_code"},
			},
			"login_key_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"access_control_group_configuration_no_list": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"is_encrypted_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"init_script_no": {
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
	return resourceNcloudLaunchConfigurationRead(d, meta)
}

func resourceNcloudLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var err error

	if config.SupportVPC {
		err = getVpcLaunchConfiguration(d, config)
	} else {
		err = getClassicLaunchConfiguration(d, config)
	}

	if err != nil {
		return err
	}

	return nil
}

// TODO : Implementation
func resourceNcloudLaunchConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var err error
	if config.SupportVPC {
		err = deleteVpcLaunchConfiguration(d, config)
	} else {
		err = deleteClassicLaunchConfiguration(d, config)
	}
	if err != nil {
		return err
	}

	return nil
}

// TODO : Implementation
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

// TODO : Implementation
func deleteVpcLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) error {
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

func getClassicLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		LaunchConfigurationNameList: []*string{ncloud.String(d.Id())},
	}

	logCommonRequest("getClassicLaunchConfiguration", reqParams)
	res, err := config.Client.autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getClassicLaunchConfiguration", err, reqParams)
		return err
	}
	logResponse("getClassicLaunchConfiguration", res)

	configuration := res.LaunchConfigurationList[0]
	instance := map[string]interface{}{
		"name":                      *configuration.LaunchConfigurationName,
		"server_image_product_code": *configuration.ServerImageProductCode,
		"server_product_code":       *configuration.ServerProductCode,
		"member_server_image_no":    *configuration.MemberServerImageNo,
		"login_key_name":            *configuration.LoginKeyName,
		"user_data":                 *configuration.UserData,
	}
	d.Set("region", &config.RegionCode)
	SetSingularResourceDataFromMapSchema(resourceNcloudLaunchConfiguration(), d, instance)
	return nil
}

// TODO : Implementation
func getVpcLaunchConfiguration(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode:                  &config.RegionCode,
		LaunchConfigurationNameList: []*string{ncloud.String(d.Id())},
	}

	logCommonRequest("getVpcLaunchConfiguration", reqParams)
	res, err := config.Client.vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getVpcLaunchConfiguration", err, reqParams)
		return err
	}
	logResponse("getVpcLaunchConfiguration", res)

	configuration := res.LaunchConfiguration[0]
	instance := map[string]interface{}{
		"region":                    *configuration.RegionCode,
		"name":                      *configuration.LaunchConfigurationName,
		"server_image_product_code": *configuration.ServerImageProductCode,
		"member_server_image_no":    *configuration.MemberServerImageInstanceNo,
		"server_product_code":       *configuration.ServerProductCode,
		"login_key_name":            *configuration.LoginKeyName,
		"status":                    *configuration.LaunchConfigurationStatus.Code,
		"init_script_no":            *configuration.InitScriptNo,
		"is_encrypted_volume":       *configuration.IsEncryptedVolume,
	}

	SetSingularResourceDataFromMapSchema(resourceNcloudLaunchConfiguration(), d, instance)
	return nil
}
