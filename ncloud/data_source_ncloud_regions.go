package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterDatasource("ncloud_regions", dataSourceNcloudRegions())
}

func dataSourceNcloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": dataSourceFiltersSchema(),
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     regionSchemaResource,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(time.Now().UTC().String())

	var regions []*Region
	var err error

	if meta.(*ProviderConfig).SupportVPC == true {
		regions, err = getVpcRegions(d, meta.(*ProviderConfig))
	} else {
		regions, err = getClassicRegions(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	if code, codeOk := d.GetOk("code"); codeOk {
		for _, region := range regions {
			if ncloud.StringValue(region.RegionCode) == code {
				regions = []*Region{region}
				break
			}
		}
	}

	if len(regions) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := flattenRegions(regions)

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudRegions().Schema["regions"].Elem.(*schema.Resource).Schema)
	}

	if err := d.Set("regions", resources); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), d.Get("regions"))
	}

	return nil
}

func getClassicRegions(d *schema.ResourceData, config *ProviderConfig) ([]*Region, error) {
	client := config.Client
	resp, err := client.server.V2Api.GetRegionList(&server.GetRegionListRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching regions found")
	}

	var regions []*Region

	for _, r := range resp.RegionList {
		regions = append(regions, GetRegion(r))
	}

	return regions, nil
}

func getVpcRegions(d *schema.ResourceData, config *ProviderConfig) ([]*Region, error) {
	client := config.Client
	resp, err := client.vserver.V2Api.GetRegionList(&vserver.GetRegionListRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching regions found")
	}

	var regions []*Region

	for _, r := range resp.RegionList {
		regions = append(regions, GetRegion(r))
	}

	return regions, nil
}
