package autoscaling

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudLaunchConfiguration() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchema(ResourceNcloudLaunchConfiguration(), fieldMap, dataSourceNcloudLaunchConfigurationRead)
}

func dataSourceNcloudLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if v, ok := d.GetOk("id"); ok {
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
		launchConfigListMap = ApplyFilters(f.(*schema.Set), launchConfigListMap, DataSourceNcloudLaunchConfiguration().Schema)
	}

	if err := ValidateOneResult(len(launchConfigListMap)); err != nil {
		return err
	}

	d.SetId(launchConfigListMap[0]["launch_configuration_no"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudLaunchConfiguration(), d, launchConfigListMap[0])
	return nil
}

func getLaunchConfigurationList(config *conn.ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	if config.SupportVPC {
		return getVpcLaunchConfigurationList(config, id)
	} else {
		return getClassicLaunchConfigurationList(config, id)
	}
}

func getVpcLaunchConfigurationList(config *conn.ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LaunchConfigurationNoList = []*string{ncloud.String(id)}
	}

	LogCommonRequest("getVpcLaunchConfigurationList", reqParams)
	resp, err := config.Client.Vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcLaunchConfigurationList", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcLaunchConfigurationList", resp)

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

func getClassicLaunchConfigurationList(config *conn.ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}
	LogCommonRequest("getClassicLaunchConfigurationList", reqParams)
	resp, err := config.Client.Autoscaling.V2Api.GetLaunchConfigurationList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicLaunchConfigurationList", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicLaunchConfigurationList", resp)

	list := make([]*LaunchConfiguration, 0)
	for _, l := range resp.LaunchConfigurationList {
		launchConfiguration := &LaunchConfiguration{
			LaunchConfigurationNo:       l.LaunchConfigurationNo,
			LaunchConfigurationName:     l.LaunchConfigurationName,
			ServerImageProductCode:      l.ServerImageProductCode,
			MemberServerImageInstanceNo: l.MemberServerImageNo,
			ServerProductCode:           l.ServerProductCode,
			LoginKeyName:                l.LoginKeyName,
			UserData:                    l.UserData,
			AccessControlGroupNoList:    flattenAccessControlGroupList(l.AccessControlGroupList),
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
