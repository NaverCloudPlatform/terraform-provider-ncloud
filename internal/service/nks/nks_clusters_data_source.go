package nks

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudNKSClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_uuids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNcloudNKSClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	clusters, err := GetNKSClusters(ctx, config)
	if err != nil {
		return diag.FromErr(err)
	}

	var cUuids []*string
	for _, cluster := range clusters {
		cUuids = append(cUuids, cluster.Uuid)
	}

	d.SetId(config.RegionCode)
	d.Set("cluster_uuids", cUuids)

	return nil
}
