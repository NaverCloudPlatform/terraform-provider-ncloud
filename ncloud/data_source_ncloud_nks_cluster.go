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
	RegisterDataSource("ncloud_nks_cluster", dataSourceNcloudNKSCluster())
}

const (
//NKSOperationChangeCode             = "CHANG"
//NKSOperationCreateCode             = "CREAT"
//NKSOperationDisUseCode             = "DISUS"
//NKSOperationNullCode               = "NULL"
//NKSOperationPendingTerminationCode = "PTERM"
//NKSOperationTerminateCode          = "TERMT"
//NKSOperationUseCode                = "USE"
)

func dataSourceNcloudNKSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSClusterRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acg_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"node_max_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_lb_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_lb_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"init_script_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"init_script_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pod_security_policy_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_pool": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceNcloudNKSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_nks_cluster`"))
	}

	cluster, err := getNKSClusterCluster(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(cluster.Uuid))
	d.Set("acg_name", cluster.AcgName)
	d.Set("name", cluster.Name)
	d.Set("cluster_type", cluster.ClusterType)
	d.Set("capacity", cluster.Capacity)
	d.Set("node_count", cluster.NodeCount)
	d.Set("node_max_count", cluster.NodeMaxCount)
	d.Set("created_at", cluster.CreatedAt)
	d.Set("updated_at", cluster.UpdatedAt)
	d.Set("cpu_count", cluster.CpuCount)
	d.Set("endpoint", cluster.Endpoint)
	d.Set("memory_size", cluster.MemorySize)
	d.Set("region_code", cluster.RegionCode)
	d.Set("status", cluster.Status)
	d.Set("subnet_name", cluster.SubnetName)
	d.Set("login_key_name", cluster.LoginKeyName)
	d.Set("k8s_version", fmt.Sprintf("%s-nks.1", cluster.K8sVersion))
	d.Set("zone_no", fmt.Sprintf("%d", *cluster.ZoneNo))
	d.Set("vpc_no", fmt.Sprintf("%d", *cluster.VpcNo))
	d.Set("vpc_name", cluster.VpcName)
	d.Set("subnet_lb_no", fmt.Sprintf("%d", *cluster.SubnetLbNo))
	d.Set("subnet_lb_name", cluster.SubnetLbName)

	if err := d.Set("subnet_no_list", flattenSubnetNoList(cluster.SubnetNoList)); err != nil {
		log.Printf("[WARN] Error setting subet no list set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("node_pool", flattenNksNodePoolList(cluster.NodePool)); err != nil {
		log.Printf("[WARN] Error setting node pool set for (%s): %s", d.Id(), err)
	}

	return nil
}
