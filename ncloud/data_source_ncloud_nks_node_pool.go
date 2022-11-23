package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func init() {
	RegisterDataSource("ncloud_nks_node_pool", dataSourceNcloudNKSNodePool())
}

func dataSourceNcloudNKSNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSNodePoolRead,
		Schema: map[string]*schema.Schema{
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_no": {
				Type:     schema.TypeString,
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
			"subnet_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_code": {
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
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kernel_version": {
							Type:     schema.TypeString,
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

	clusterUuid := d.Get("cluster_uuid").(string)
	nodePoolName := d.Get("node_pool_name").(string)
	id := NodePoolCreateResourceID(clusterUuid, nodePoolName)

	nodePool, err := getNKSNodePool(ctx, config, clusterUuid, nodePoolName)
	if err != nil {
		return diag.FromErr(err)
	}

	if nodePool == nil {
		d.SetId("")
		return nil
	}

	d.SetId(id)

	d.Set("cluster_uuid", clusterUuid)
	d.Set("instance_no", strconv.Itoa(int(ncloud.Int32Value(nodePool.InstanceNo))))
	d.Set("node_pool_name", nodePool.Name)
	d.Set("product_code", nodePool.ProductCode)
	d.Set("node_count", nodePool.NodeCount)
	d.Set("k8s_version", nodePool.K8sVersion)
	d.Set("subnet_no", strconv.Itoa(int(ncloud.Int32Value(nodePool.SubnetNoList[0]))))

	if err := d.Set("autoscale", flattenNKSNodePoolAutoScale(nodePool.Autoscale)); err != nil {
		log.Printf("[WARN] Error setting Autoscale set for (%s): %s", d.Id(), err)
	}

	nodes, err := getNKSNodePoolWorkerNodes(ctx, config, clusterUuid, nodePoolName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nodes", flattenNKSWorkerNodes(nodes)); err != nil {
		log.Printf("[WARN] Error setting workerNodes set for (%s): %s", d.Id(), err)
	}
	return nil
}
