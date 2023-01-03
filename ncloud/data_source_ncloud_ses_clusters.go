package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func init() {
	RegisterDataSource("ncloud_ses_clusters", dataSourceNcloudSESClusters())
}

func dataSourceNcloudSESClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSESClustersRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_ses_clusters`"))
	}

	clusters, err := getSESClusters(ctx, config)
	if err != nil {
		logErrorResponse("GetSESClusters", err, "")
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
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSESClusters().Schema["clusters"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("clusters", resources); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
