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
			"service_group_instance_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"search_engine": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"software_product_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager_node": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_dual_manager": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"count": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"data_node": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"storage_size": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"master_node": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_master_only_node_activated": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"count": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Optional: true,
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

	uuid := d.Get("service_group_instance_no").(string)
	cluster, err := getSESCluster(ctx, config, uuid)
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(cluster.ServiceGroupInstanceNo))
	d.Set("uuid", cluster.ServiceGroupInstanceNo)
	d.Set("service_group_instance_no", cluster.ServiceGroupInstanceNo)
	d.Set("cluster_name", cluster.ClusterName)
	d.Set("software_product_code", cluster.SoftwareProductCode)
	d.Set("vpc_no", strconv.Itoa(int(ncloud.Int32Value(cluster.VpcNo))))
	d.Set("login_key_name", cluster.LoginKeyName)

	searchEngineList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["search_engine"].Elem.(*schema.Resource)), []interface{}{})
	searchEngineList.Add(map[string]interface{}{
		"version_code": *cluster.SearchEngineVersionCode,
		"user_name":    *cluster.SearchEngineUserName,
	})
	if err := d.Set("search_engine", searchEngineList.List()); err != nil {
		log.Printf("[WARN] Error setting search_engine set for (%s): %s", d.Id(), err)
	}

	managerNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["manager_node"].Elem.(*schema.Resource)), []interface{}{})
	managerNodeList.Add(map[string]interface{}{
		"is_dual_manager": *cluster.IsDualManager,
		"count":           strconv.Itoa(int(ncloud.Int32Value(cluster.ManagerNodeCount))),
		"subnet_no":       strconv.Itoa(int(ncloud.Int32Value(cluster.ManagerNodeSubnetNo))),
		"product_code":    *cluster.ManagerNodeProductCode,
	})
	if err := d.Set("manager_node", managerNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting manager_node set for (%s): %s", d.Id(), err)
	}

	dataNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["data_node"].Elem.(*schema.Resource)), []interface{}{})
	dataNodeList.Add(map[string]interface{}{
		"count":        strconv.Itoa(int(ncloud.Int32Value(cluster.DataNodeCount))),
		"subnet_no":    strconv.Itoa(int(ncloud.Int32Value(cluster.DataNodeSubnetNo))),
		"product_code": *cluster.DataNodeProductCode,
		"storage_size": *cluster.DataNodeStorageSize,
	})
	if err := d.Set("data_node", dataNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting data_node set for (%s): %s", d.Id(), err)
	}

	masterNodeList := schema.NewSet(schema.HashResource(dataSourceNcloudSESCluster().Schema["master_node"].Elem.(*schema.Resource)), []interface{}{})
	if cluster.MasterNodeCount != nil && cluster.MasterNodeSubnetNo != nil && cluster.MasterNodeProductCode != nil {
		masterNodeList.Add(map[string]interface{}{
			"is_master_only_node_activated": true,
			"count":                         strconv.Itoa(int(ncloud.Int32Value(cluster.MasterNodeCount))),
			"subnet_no":                     strconv.Itoa(int(ncloud.Int32Value(cluster.MasterNodeSubnetNo))),
			"product_code":                  *cluster.MasterNodeProductCode,
		})
	} else {
		masterNodeList.Add(map[string]interface{}{
			"is_master_only_node_activated": false,
		})
	}
	if err := d.Set("master_node", masterNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting master_node set for (%s): %s", d.Id(), err)
	}

	return nil
}
