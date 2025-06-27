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

func DataSourceNcloudLbListener() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"load_balancer_no": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tls_min_version_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ssl_certificate_no": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchemaContext(ResourceNcloudLbListener(), fieldMap, dataSourceNcloudLbListenerRead)
}

func dataSourceNcloudLbListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	listenerList, err := getVpcLoadBalancerListenerList(config, d.Id(), d.Get("load_balancer_no").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	listenerListMap := ConvertToArrayMap(listenerList)
	if f, ok := d.GetOk("filter"); ok {
		listenerListMap = ApplyFilters(f.(*schema.Set), listenerListMap, DataSourceNcloudLbListener().Schema)
	}

	if err := ValidateOneResult(len(listenerListMap)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(listenerListMap[0]["listener_no"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudLbListener(), d, listenerListMap[0])
	return nil
}

func getVpcLoadBalancerListenerList(config *conn.ProviderConfig, id string, loadBalancerNo string) ([]*LoadBalancerListener, error) {
	reqParams := &vloadbalancer.GetLoadBalancerListenerListRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(loadBalancerNo),
	}

	resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerListenerList(reqParams)
	if err != nil {
		return nil, err
	}

	listenerList := make([]*LoadBalancerListener, 0)
	for _, l := range resp.LoadBalancerListenerList {
		listener := &LoadBalancerListener{
			LoadBalancerListenerNo: l.LoadBalancerListenerNo,
			ProtocolType:           l.ProtocolType.Code,
			Port:                   l.Port,
			UseHttp2:               l.UseHttp2,
			SslCertificateNo:       l.SslCertificateNo,
			TlsMinVersionType:      l.TlsMinVersionType.Code,
			LoadBalancerRuleNoList: l.LoadBalancerRuleNoList,
		}
		if id == *listener.LoadBalancerListenerNo {
			return []*LoadBalancerListener{listener}, nil
		}
		listenerList = append(listenerList, listener)
	}

	return listenerList, nil
}
