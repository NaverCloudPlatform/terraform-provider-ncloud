package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_lb", dataSourceNcloudLb())
}

func dataSourceNcloudLb() *schema.Resource {
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
		"filter": dataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchemaContext(resourceNcloudLb(), fieldMap, dataSourceNcloudLbRead)
}

func dataSourceNcloudLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("datasource `ncloud_lb`"))
	}

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	reqParams := &vloadbalancer.GetLoadBalancerInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if d.Id() != "" {
		reqParams.LoadBalancerInstanceNoList = []*string{ncloud.String(d.Id())}
	}

	resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceList(reqParams)
	if err != nil {
		return diag.FromErr(err)
	}

	lbList := make([]*LoadBalancerInstance, 0)
	for _, lb := range resp.LoadBalancerInstanceList {
		lbList = append(lbList, &LoadBalancerInstance{
			LoadBalancerInstanceNo:   lb.LoadBalancerInstanceNo,
			LoadBalancerDescription:  lb.LoadBalancerDescription,
			LoadBalancerName:         lb.LoadBalancerName,
			LoadBalancerDomain:       lb.LoadBalancerDomain,
			LoadBalancerIpList:       lb.LoadBalancerIpList,
			LoadBalancerType:         lb.LoadBalancerType.Code,
			LoadBalancerNetworkType:  lb.LoadBalancerNetworkType.Code,
			ThroughputType:           lb.ThroughputType.Code,
			IdleTimeout:              lb.IdleTimeout,
			VpcNo:                    lb.VpcNo,
			SubnetNoList:             lb.SubnetNoList,
			LoadBalancerListenerList: lb.LoadBalancerListenerNoList,
		})
	}

	lbListMap := ConvertToArrayMap(lbList)
	if f, ok := d.GetOk("filter"); ok {
		lbListMap = ApplyFilters(f.(*schema.Set), lbListMap, dataSourceNcloudLb().Schema)
	}

	if err := validateOneResult(len(lbListMap)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbListMap[0]["load_balancer_no"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudLb(), d, lbListMap[0])
	return nil
}
