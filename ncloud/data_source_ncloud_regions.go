package ncloud

import (
	"fmt"
	"os"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
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
	conn := meta.(*NcloudSdk).conn
	d.SetId(time.Now().UTC().String())

	regionList, err := getRegions(conn)
	if err != nil {
		return err
	}

	code, codeOk := d.GetOk("code")

	var filteredRegions []common.Region
	if codeOk {
		if filtered, err := getRegionByCode(conn, code.(string)); err != nil {
			filteredRegions = []common.Region{*filtered}
			return err
		}
	} else {
		filteredRegions = regionList
	}

	if len(filteredRegions) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return regionsAttributes(d, filteredRegions)
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

func getRegions(conn *sdk.Conn) ([]common.Region, error) {
	resp, err := conn.GetRegionList()
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching regions found")
	}

	return resp.RegionList, nil
}

func getRegionByCode(conn *sdk.Conn, code string) (*common.Region, error) {
	regionList, err := getRegions(conn)
	if err != nil {
		return nil, err
	}

	var filteredRegion common.Region
	for _, region := range regionList {
		if code == region.RegionCode {
			filteredRegion = region
			break
		}
	}

	return &filteredRegion, nil
}

func getRegionNoByCode(conn *sdk.Conn, name string) string {
	region, err := getRegionByCode(conn, name)
	if err != nil {
		return ""
	}

	return region.RegionNo
}

var regionCache = make(map[string]string)

func parseRegionNoParameter(conn *sdk.Conn, d *schema.ResourceData) string {
	if paramRegionNo, regionNoOk := d.GetOk("region_no"); regionNoOk {
		return paramRegionNo.(string)
	}

	// provider region
	if regionCode := os.Getenv("NCLOUD_REGION"); regionCode != "" {
		regionNo := regionCache[regionCode]
		if regionNo != "" {
			return regionNo
		}
		regionNo = getRegionNoByCode(conn, regionCode)
		if regionNo != "" {
			regionCache[regionCode] = regionNo
		}
	}

	return ""
}
