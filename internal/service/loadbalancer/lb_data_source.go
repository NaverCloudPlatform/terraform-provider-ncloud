package loadbalancer

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudLb() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchemaContext(ResourceNcloudLb(), fieldMap, dataSourceNcloudLbRead)
}

func dataSourceNcloudLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("datasource `ncloud_lb`"))
	}

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	lbList, err := getVpcLoadBalancerList(config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	lbListMap := ConvertToArrayMap(lbList)
	if f, ok := d.GetOk("filter"); ok {
		lbListMap = ApplyFilters(f.(*schema.Set), lbListMap, DataSourceNcloudLb().Schema)
	}

	if err := ValidateOneResult(len(lbListMap)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbListMap[0]["load_balancer_no"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudLb(), d, lbListMap[0])
	return nil
}

func getVpcLoadBalancerList(config *conn.ProviderConfig, id string) ([]*LoadBalancerInstance, error) {
	reqParams := &vloadbalancer.GetLoadBalancerInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LoadBalancerInstanceNoList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceList(reqParams)
	if err != nil {
		return nil, err
	}

	lbList := make([]*LoadBalancerInstance, 0)
	for _, lb := range resp.LoadBalancerInstanceList {
		lbList = append(lbList, convertVpcLoadBalancer(lb))
	}

	return lbList, nil
}
