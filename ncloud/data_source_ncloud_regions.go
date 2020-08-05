package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
	client := meta.(*ProviderConfig).Client
	d.SetId(time.Now().UTC().String())

	regionList, err := getRegions(client)
	if err != nil {
		return err
	}

	code, codeOk := d.GetOk("code")

	var filteredRegions []*Region
	if codeOk {
		for _, region := range regionList {
			if ncloud.StringValue(region.RegionCode) == code {
				filteredRegions = []*Region{region}
				break
			}
		}
	} else {
		filteredRegions = regionList
	}

	if len(filteredRegions) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return regionsAttributes(d, filteredRegions)
}

func regionsAttributes(d *schema.ResourceData, regions []*Region) error {

	var ids []string
	var s []map[string]interface{}
	for _, region := range regions {
		mapping := flattenRegion(region)
		ids = append(ids, *region.RegionNo)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("regions", s); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), s)
	}

	return nil
}

func getRegions(client *NcloudAPIClient) ([]*Region, error) {
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
