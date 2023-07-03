package nks

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudNKSNodePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSNodePoolsRead,
		Schema: map[string]*schema.Schema{
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_pool_names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNcloudNKSNodePoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_nks_node_pools`"))
	}

	clusterUuid := d.Get("cluster_uuid").(string)

	nodePools, err := getNKSNodePools(ctx, config, clusterUuid)
	if err != nil {
		return diag.FromErr(err)
	}

	var npNames []*string
	for _, nodePool := range nodePools {
		npNames = append(npNames, nodePool.Name)
	}

	d.SetId(clusterUuid)

	d.Set("cluster_uuid", clusterUuid)
	d.Set("node_pool_names", npNames)

	return nil
}
