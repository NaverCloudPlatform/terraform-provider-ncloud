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
	RegisterDataSource("ncloud_ses_cluster", dataSourceNcloudSESCluster())
}

func dataSourceNcloudSESCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSESClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_group_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"search_engine": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dashboard_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_no": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_image_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_dual_manager": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"data_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"master_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"manager_node_instance_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cluster_node_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_instance_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"compute_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudSESClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_ses_cluster`"))
	}

	id := d.Get("id").(string)
	cluster, err := getSESCluster(ctx, config, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(cluster.ServiceGroupInstanceNo))
	d.Set("id", cluster.ServiceGroupInstanceNo)
	d.Set("service_group_instance_no", cluster.ServiceGroupInstanceNo)
	d.Set("cluster_name", cluster.ClusterName)
	d.Set("os_image_code", cluster.SoftwareProductCode)
	d.Set("vpc_no", cluster.VpcNo)
	d.Set("login_key_name", cluster.LoginKeyName)
	d.Set("manager_node_instance_no_list", cluster.ManagerNodeInstanceNoList)

	searchEngineList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["search_engine"].Elem.(*schema.Resource)), []interface{}{})
	searchEngineList.Add(map[string]interface{}{
		"version_code":   *cluster.SearchEngineVersionCode,
		"user_name":      *cluster.SearchEngineUserName,
		"port":           *cluster.SearchEnginePort,
		"dashboard_port": *cluster.SearchEngineDashboardPort,
	})
	if err := d.Set("search_engine", searchEngineList.List()); err != nil {
		log.Printf("[WARN] Error setting search_engine set for (%s): %s", d.Id(), err)
	}

	managerNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["manager_node"].Elem.(*schema.Resource)), []interface{}{})
	managerNodeList.Add(map[string]interface{}{
		"is_dual_manager": *cluster.IsDualManager,
		"count":           *cluster.ManagerNodeCount,
		"subnet_no":       *cluster.ManagerNodeSubnetNo,
		"product_code":    *cluster.ManagerNodeProductCode,
		"acg_id":          *cluster.ManagerNodeAcgId,
		"acg_name":        *cluster.ManagerNodeAcgName,
	})
	if err := d.Set("manager_node", managerNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting manager_node set for (%s): %s", d.Id(), err)
	}

	dataNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["data_node"].Elem.(*schema.Resource)), []interface{}{})
	storageSize, _ := strconv.Atoi(*cluster.DataNodeStorageSize)
	dataNodeList.Add(map[string]interface{}{
		"count":        *cluster.DataNodeCount,
		"subnet_no":    *cluster.DataNodeSubnetNo,
		"product_code": *cluster.DataNodeProductCode,
		"acg_id":       *cluster.DataNodeAcgId,
		"acg_name":     *cluster.DataNodeAcgName,
		"storage_size": storageSize,
	})
	if err := d.Set("data_node", dataNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting data_node set for (%s): %s", d.Id(), err)
	}

	if cluster.IsMasterOnlyNodeActivated != nil && *cluster.IsMasterOnlyNodeActivated {
		masterNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["master_node"].Elem.(*schema.Resource)), []interface{}{})
		masterNodeList.Add(map[string]interface{}{
			"count":        *cluster.MasterNodeCount,
			"subnet_no":    *cluster.MasterNodeSubnetNo,
			"product_code": *cluster.MasterNodeProductCode,
			"acg_id":       *cluster.MasterNodeAcgId,
			"acg_name":     *cluster.MasterNodeAcgName,
		})

		if err := d.Set("master_node", masterNodeList.List()); err != nil {
			log.Printf("[WARN] Error setting master_node set for (%s): %s", d.Id(), err)
		}
	}

	clusterNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["cluster_node_list"].Elem.(*schema.Resource)), []interface{}{})
	if cluster.ClusterNodeList != nil {
		for _, clusterNode := range cluster.ClusterNodeList {
			clusterNodeList.Add(map[string]interface{}{
				"compute_instance_no":   clusterNode.ComputeInstanceNo,
				"compute_instance_name": clusterNode.ComputeInstanceName,
				"private_ip":            clusterNode.PrivateIp,
				"server_status":         clusterNode.ServerStatus,
				"node_type":             clusterNode.NodeType,
				"subnet":                clusterNode.Subnet,
			})
		}
	}
	if err := d.Set("cluster_node_list", clusterNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting cluster node list for (%s): %s", d.Id(), err)
	}

	return nil
}
