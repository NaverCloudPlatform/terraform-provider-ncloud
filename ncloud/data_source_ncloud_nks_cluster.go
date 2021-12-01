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
	RegisterDataSource("ncloud_nks_cluster", dataSourceNcloudNKSCluster())
}

func dataSourceNcloudNKSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudNKSClusterRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_lb_no": {
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
		},
	}
}

func dataSourceNcloudNKSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_nks_cluster`"))
	}

	uuid := d.Get("uuid").(string)
	cluster, err := getNKSCluster(ctx, config, uuid)
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(cluster.Uuid))
	d.Set("uuid", cluster.Uuid)
	d.Set("name", cluster.Name)
	d.Set("cluster_type", cluster.ClusterType)
	d.Set("endpoint", cluster.Endpoint)
	d.Set("login_key_name", cluster.LoginKeyName)
	d.Set("k8s_version", cluster.K8sVersion)
	d.Set("zone", cluster.ZoneCode)
	d.Set("vpc_no", strconv.Itoa(int(ncloud.Int32Value(cluster.VpcNo))))
	d.Set("subnet_lb_no", strconv.Itoa(int(ncloud.Int32Value(cluster.SubnetLbNo))))

	if err := d.Set("subnet_no_list", flattenInt32ListToStringList(cluster.SubnetNoList)); err != nil {
		log.Printf("[WARN] Error setting subet no list set for (%s): %s", d.Id(), err)
	}

	return nil
}
