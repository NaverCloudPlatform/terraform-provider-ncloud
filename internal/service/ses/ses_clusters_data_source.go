package ses

import (
	"context"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSESClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSESClustersRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_group_instance_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSESClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	clusters, err := getSESClusters(ctx, config)
	if err != nil {
		LogErrorResponse("GetSESClusters", err, "")
		return diag.FromErr(err)
	}

	if clusters == nil {
		d.SetId("")
		return nil
	}

	resources := []map[string]interface{}{}

	for _, r := range clusters.AllowedClusters {
		instance := map[string]interface{}{
			"id":                        ncloud.StringValue(r.ServiceGroupInstanceNo),
			"service_group_instance_no": ncloud.StringValue(r.ServiceGroupInstanceNo),
			"cluster_name":              ncloud.StringValue(r.ClusterName),
		}

		resources = append(resources, instance)
	}

	for _, r := range clusters.DisallowedClusters {
		instance := map[string]interface{}{
			"id":                        ncloud.StringValue(r.ServiceGroupInstanceNo),
			"service_group_instance_no": ncloud.StringValue(r.ServiceGroupInstanceNo),
			"cluster_name":              ncloud.StringValue(r.ClusterName),
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSESClusters().Schema["clusters"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("clusters", resources); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
