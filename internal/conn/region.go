package conn

import (
	"fmt"
	"os"
	"sync"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Region struct {
	RegionNo   *string `json:"regionNo,omitempty"`
	RegionCode *string `json:"regionCode,omitempty"`
	RegionName *string `json:"regionName,omitempty"`
}

var regionCacheByCode = sync.Map{}

func ParseRegionNoParameter(d *schema.ResourceData) (*string, error) {
	if regionCode, regionCodeOk := d.GetOk("region"); regionCodeOk {
		regionNo := GetRegionNoByCode(regionCode.(string))
		if regionNo == nil {
			return nil, fmt.Errorf("no region data for region_code `%s`. please change region_code and try again", regionCode.(string))
		}
		return regionNo, nil
	}

	// provider region
	if regionCode := os.Getenv("NCLOUD_REGION"); regionCode != "" {
		regionNo := GetRegionNoByCode(regionCode)
		if regionNo == nil {
			return nil, fmt.Errorf("no region data for region_code `%s`. please change region_code and try again", regionCode)
		}
		return regionNo, nil
	}

	return nil, nil
}

func GetRegionNoByCode(code string) *string {
	if region, ok := regionCacheByCode.Load(code); ok {
		return region.(Region).RegionNo
	}
	return nil
}

func SetRegionCache(client *NcloudAPIClient) error {
	var regionList []*Region
	var err error

	regionList, err = getVpcRegionList(client)
	if err != nil {
		return err
	}

	for _, r := range regionList {
		region := Region{
			RegionCode: r.RegionCode,
			RegionName: r.RegionName,
		}

		regionCacheByCode.Store(*region.RegionCode, region)
	}

	return nil
}

func getVpcRegionList(client *NcloudAPIClient) ([]*Region, error) {
	resp, err := client.Vserver.V2Api.GetRegionList(&vserver.GetRegionListRequest{})
	if err != nil {
		return nil, err
	}
	var regionList []*Region
	for _, r := range resp.RegionList {
		region := &Region{
			RegionCode: r.RegionCode,
			RegionName: r.RegionName,
		}
		regionList = append(regionList, region)
	}

	return regionList, nil
}

func IsValidRegionCode(code string) bool {
	_, ok := regionCacheByCode.Load(code)
	return ok
}
