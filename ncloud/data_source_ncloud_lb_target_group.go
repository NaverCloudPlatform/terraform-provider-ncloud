package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_lb_target_group", dataSourceNcloudLbTargetGroup())
}

func dataSourceNcloudLbTargetGroup() *schema.Resource {
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
		"health_check": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cycle": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"down_threshold": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"up_threshold": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"http_method": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"protocol": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"url_path": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"filter": dataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchemaContext(resourceNcloudLbTargetGroup(), fieldMap, dataSourceNcloudLbTargetGroupRead)
}

func dataSourceNcloudLbTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("datasource `ncloud_lb_target_group`"))
	}

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	targetGroupList, err := getVpcLoadBalancerTargetGroupList(config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	targetGroupListMap := ConvertToArrayMap(targetGroupList)
	if f, ok := d.GetOk("filter"); ok {
		targetGroupListMap = ApplyFilters(f.(*schema.Set), targetGroupListMap, dataSourceNcloudLbTargetGroup().Schema)
	}

	if err := validateOneResult(len(targetGroupListMap)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetGroupListMap[0]["target_group_no"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudLbTargetGroup(), d, targetGroupListMap[0])
	return nil
}

func getVpcLoadBalancerTargetGroupList(config *ProviderConfig, id string) ([]*TargetGroup, error) {
	reqParams := &vloadbalancer.GetTargetGroupListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.TargetGroupNoList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.vloadbalancer.V2Api.GetTargetGroupList(reqParams)
	if err != nil {
		return nil, err
	}

	targetGroupList := make([]*TargetGroup, 0)
	for _, tg := range resp.TargetGroupList {
		targetGroupList = append(targetGroupList, convertVpcTargetGroup(tg))
	}

	return targetGroupList, nil
}
