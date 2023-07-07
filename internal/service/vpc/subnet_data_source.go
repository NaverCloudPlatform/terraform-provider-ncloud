package vpc

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudSubnet() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"zone": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"network_acl_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false)),
		},
		"usage_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"GEN", "LOADB", "BM", "NATGW"}, false)),
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudSubnet(), fieldMap, dataSourceNcloudSubnetRead)
}

func dataSourceNcloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_subnet`")
	}

	resources, err := getSubnetListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getSubnetListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vpc.GetSubnetListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.SubnetNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet"); ok {
		reqParams.Subnet = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		reqParams.ZoneCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("network_acl_no"); ok {
		reqParams.NetworkAclNo = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_type_code"); ok {
		reqParams.SubnetTypeCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("usage_type_code"); ok {
		reqParams.UsageTypeCode = ncloud.String(v.(string))
	}

	LogCommonRequest("GetSubnetList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetSubnetList(reqParams)

	if err != nil {
		LogErrorResponse("GetSubnetList", err, reqParams)
		return nil, err
	}
	LogResponse("GetSubnetList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.SubnetList {
		instance := map[string]interface{}{
			"id":             *r.SubnetNo,
			"subnet_no":      *r.SubnetNo,
			"vpc_no":         *r.VpcNo,
			"zone":           *r.ZoneCode,
			"name":           *r.SubnetName,
			"subnet":         *r.Subnet,
			"subnet_type":    *r.SubnetType.Code,
			"usage_type":     *r.UsageType.Code,
			"network_acl_no": *r.NetworkAclNo,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudSubnet().Schema)
	}

	return resources, nil
}
