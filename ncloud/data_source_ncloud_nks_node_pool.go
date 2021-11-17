package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func init() {
	RegisterDataSource("ncloud_nks_node_pool", dataSourceNcloudNKSNodePool())
}

func dataSourceNcloudNKSNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSNodePoolRead,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_no": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_name_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscale": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudNKSNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_nks_node_pool`"))
	}

	clusterName := d.Get("cluster_name").(string)
	nodePoolName := d.Get("node_pool_name").(string)
	id := NodePoolCreateResourceID(clusterName, nodePoolName)

	cluster, err := getNKSClusterWithName(ctx, config, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}
	if cluster == nil {
		return diag.FromErr(fmt.Errorf("cluster \"%s\"  not found ", clusterName))
	}

	nodePool, err := getNKSNodePool(ctx, config, cluster.Uuid, &nodePoolName)
	if err != nil {
		return diag.FromErr(err)
	}

	if nodePool == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(&id))

	d.Set("instance_no", nodePool.InstanceNo)
	d.Set("node_pool_name", nodePool.Name)
	d.Set("status", nodePool.Status)
	d.Set("product_code", nodePool.ProductCode)

	if err := d.Set("subnet_no_list", flattenSubnetNoList(nodePool.SubnetNoList)); err != nil {
		log.Printf("[WARN] Error setting subet no list set for (%s): %s", d.Id(), err)
	}
	if err := d.Set("autoscale", flattenAutoscale(nodePool.Autoscale)); err != nil {
		log.Printf("[WARN] Error setting Autoscale set for (%s): %s", d.Id(), err)
	}
	return nil
}
