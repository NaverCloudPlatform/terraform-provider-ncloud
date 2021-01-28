package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_launch_configuration", dataSourceNcloudLaunchConfiguration())
}

func dataSourceNcloudLaunchConfiguration() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"launch_configuration_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchema(resourceNcloudLaunchConfiguration(), fieldMap, dataSourceNcloudLaunchConfigurationRead)
}

func dataSourceNcloudLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if v, ok := d.GetOk("launch_configuration_no"); ok {
		d.SetId(v.(string))
	}

	launchConfigList, err := getLaunchConfigurationList(config, d.Id())
	if err != nil {
		return err
	}

	if launchConfigList == nil {
		return nil
	}

	launchConfigListMap := ConvertToArrayMap(launchConfigList)
	if f, ok := d.GetOk("filter"); ok {
		launchConfigListMap = ApplyFilters(f.(*schema.Set), launchConfigListMap, dataSourceNcloudLaunchConfiguration().Schema)
	}

	if err := validateOneResult(len(launchConfigListMap)); err != nil {
		return err
	}

	d.SetId(launchConfigListMap[0]["launch_configuration_no"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudLaunchConfiguration(), d, launchConfigListMap[0])
	return nil
}

func getLaunchConfigurationList(config *ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	if config.SupportVPC {
		return getVpcLaunchConfigurationList(config, id)
	} else {
		return getClassicLaunchConfigurationList(config, id)
	}
}

func getVpcLaunchConfigurationList(config *ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LaunchConfigurationNoList = []*string{ncloud.String(id)}
	}

	logCommonRequest("getVpcLaunchConfigurationList", reqParams)
	resp, err := config.Client.vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getVpcLaunchConfigurationList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcLaunchConfigurationList", resp)

	if len(resp.LaunchConfigurationList) < 1 {
		return nil, nil
	}

	list := make([]*LaunchConfiguration, 0)
	for _, l := range resp.LaunchConfigurationList {
		list = append(list, &LaunchConfiguration{
			LaunchConfigurationName:     l.LaunchConfigurationName,
			ServerImageProductCode:      l.ServerImageProductCode,
			MemberServerImageInstanceNo: l.MemberServerImageInstanceNo,
			ServerProductCode:           l.ServerProductCode,
			LoginKeyName:                l.LoginKeyName,
			InitScriptNo:                l.InitScriptNo,
			IsEncryptedVolume:           l.IsEncryptedVolume,
			LaunchConfigurationNo:       l.LaunchConfigurationNo,
		})
	}

	return list, nil
}

func getClassicLaunchConfigurationList(config *ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}
	logCommonRequest("getClassicLaunchConfigurationList", reqParams)
	resp, err := config.Client.autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		logErrorResponse("getClassicLaunchConfigurationList", err, reqParams)
		return nil, err
	}
	logResponse("getClassicLaunchConfigurationList", resp)

	list := make([]*LaunchConfiguration, 0)
	for _, l := range resp.LaunchConfigurationList {
		launchConfiguration := &LaunchConfiguration{
			LaunchConfigurationNo:       l.LaunchConfigurationNo,
			LaunchConfigurationName:     l.LaunchConfigurationName,
			ServerImageProductCode:      l.ServerImageProductCode,
			MemberServerImageInstanceNo: l.MemberServerImageNo,
			ServerProductCode:           l.ServerProductCode,
			LoginKeyName:                l.LoginKeyName,
		}

		if *l.LaunchConfigurationNo == *no {
			return []*LaunchConfiguration{launchConfiguration}, nil
		}

		list = append(list, launchConfiguration)
	}

	if *no != "" {
		return nil, nil
	}

	return list, nil
}
