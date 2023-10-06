package region

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": DataSourceFiltersSchema(),
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

	var regions []*conn.Region
	var err error

	if meta.(*conn.ProviderConfig).SupportVPC {
		regions, err = getVpcRegions(d, meta.(*conn.ProviderConfig))
	} else {
		regions, err = getClassicRegions(d, meta.(*conn.ProviderConfig))
	}

	if err != nil {
		return err
	}

	if code, codeOk := d.GetOk("code"); codeOk {
		for _, region := range regions {
			if ncloud.StringValue(region.RegionCode) == code {
				regions = []*conn.Region{region}
				break
			}
		}
	}

	if len(regions) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := FlattenRegions(regions)

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudRegions().Schema["regions"].Elem.(*schema.Resource).Schema)
	}

	if err := d.Set("regions", resources); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("regions"))
	}

	return nil
}

func getClassicRegions(d *schema.ResourceData, config *conn.ProviderConfig) ([]*conn.Region, error) {
	client := config.Client
	resp, err := client.Server.V2Api.GetRegionList(&server.GetRegionListRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching regions found")
	}

	var regions []*conn.Region

	for _, r := range resp.RegionList {
		regions = append(regions, GetRegion(r))
	}

	return regions, nil
}

func getVpcRegions(d *schema.ResourceData, config *conn.ProviderConfig) ([]*conn.Region, error) {
	client := config.Client
	resp, err := client.Vserver.V2Api.GetRegionList(&vserver.GetRegionListRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching regions found")
	}

	var regions []*conn.Region

	for _, r := range resp.RegionList {
		regions = append(regions, GetRegion(r))
	}

	return regions, nil
}
