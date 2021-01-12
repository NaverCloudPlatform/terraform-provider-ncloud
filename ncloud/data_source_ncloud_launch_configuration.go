package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_launch_configuration", dataSourceNcloudLaunchConfiguration())
}

func dataSourceNcloudLaunchConfiguration() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"name": {
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

	var launchConfig *LaunchConfiguration
	var err error

	if config.SupportVPC {
		launchConfig, err = getVpcLaunchConfiguration(d, config)
	} else {
		launchConfig, err = getClassicLaunchConfiguration(d, config)
	}
	fmt.Print()
	if err != nil {
		return err
	}

	launchConfigMap := ConvertToMap(launchConfig)
	launchConfigArrMap := []map[string]interface{}{launchConfigMap}

	if f, ok := d.GetOk("filter"); ok {
		launchConfigArrMap = ApplyFilters(f.(*schema.Set), launchConfigArrMap, dataSourceNcloudLaunchConfiguration().Schema)
	}

	if err := validateOneResult(len(launchConfigArrMap)); err != nil {
		return err
	}

	d.SetId(StringOrEmpty(launchConfig.LaunchConfigurationNo))
	SetSingularResourceDataFromMap(d, launchConfigArrMap[0])
	return nil
}
