package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDatasource("ncloud_vpc", dataSourceNcloudVpc())
}

func dataSourceNcloudVpc() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"vpc_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"status": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"INIT", "CREATING", "RUN", "TERMTING"}, false),
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudVpc(), fieldMap, dataSourceNcloudVpcRead)
}

func dataSourceNcloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	resources, err := getVpcListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vpc.GetVpcListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		reqParams.VpcStatusCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("GetVpcList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetVpcList(reqParams)

	if err != nil {
		logErrorResponse("GetVpcList", err, reqParams)
		return nil, err
	}
	logResponse("GetVpcList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.VpcList {
		instance := map[string]interface{}{
			"id":              *r.VpcNo,
			"vpc_no":          *r.VpcNo,
			"name":            *r.VpcName,
			"ipv4_cidr_block": *r.Ipv4CidrBlock,
			"status":          *r.VpcStatus.Code,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudVpc().Schema)
	}

	return resources, nil
}
