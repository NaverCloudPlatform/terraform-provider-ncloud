package ncloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_nks_node_pools", dataSourceNcloudNKSNodePools())
}

func dataSourceNcloudNKSNodePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSNodePoolsRead,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
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
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_nks_node_pools`"))
	}

	clusterName := d.Get("cluster_name").(string)

	cluster, err := getNKSClusterWithName(ctx, config, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	nodePools, err := getNKSNodePools(ctx, config, cluster.Uuid)
	if err != nil {
		return diag.FromErr(err)
	}

	var npNames []*string
	for _, nodePool := range nodePools {
		npNames = append(npNames, nodePool.Name)
	}

	d.SetId(clusterName)

	d.Set("cluster_name", clusterName)
	d.Set("node_pool_names", npNames)

	return nil
}
