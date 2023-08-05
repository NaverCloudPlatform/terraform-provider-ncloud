package cloudmysql

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudMysql() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudMySql(), fieldMap, dataSourceNcloudMysqlRead)
}

func dataSourceNcloudMysqlRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_mysql`")
	}

	r, err := getMysqlList(d, config)

	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(r)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, r[0])

	return nil
}

func getMysqlList(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client

	reqParams := &vmysql.GetCloudMysqlInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.CloudMysqlInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getCloudMysqlList", reqParams)

	resp, err := client.Vmysql.V2Api.GetCloudMysqlInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getCloudMysqlList", err, reqParams)
		return nil, err
	}
	LogResponse("getCloudMysqlList", resp)

	var resourcesList []map[string]interface{}
	for _, r := range resp.CloudMysqlInstanceList {

		instance := ConvertToMap(r)
		resourcesList = append(resourcesList, instance)
	}

	return resourcesList, nil
}
