package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	log.Printf("[DEBUG] Get Region List")
	d.SetId(time.Now().UTC().String())

	resp, err := conn.GetRegionList()
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("no matching regions found")
	}

	code, codeOk := d.GetOk("code")

	var filterRegions []common.Region

	for _, region := range resp.RegionList {
		if codeOk {
			if code == region.RegionCode {
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		filterRegions = append(filterRegions, region)
	}

	if len(filterRegions) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return regionsAttributes(d, filterRegions)
}

func regionsAttributes(d *schema.ResourceData, regions []common.Region) error {

	var ids []string
	var s []map[string]interface{}
	for _, region := range regions {
		mapping := map[string]interface{}{
			"region_no":   region.RegionNo,
			"region_code": region.RegionCode,
			"region_name": region.RegionName,
		}

		log.Printf("[DEBUG] ncloud_regions - adding region mapping: %v", mapping)
		ids = append(ids, string(region.RegionNo))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("regions", s); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
