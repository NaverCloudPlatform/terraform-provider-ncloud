package ncloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func init() {
	RegisterDataSource("ncloud_cdss_cluster", dataSourceNcloudCDSSCluster())
}

func dataSourceNcloudCDSSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudCDSSClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_product_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_group_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cmak": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"manager_node": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_product_code": {
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
			"broker_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_product_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"node_count": {
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
			"endpoints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker_node_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"broker_tls_node_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"public_endpoint_broker_node_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"public_endpoint_broker_node_listener_port_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"public_endpoint_broker_tls_node_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"public_endpoint_broker_tls_node_listener_port_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"local_dns_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"local_dns_tls_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zookeeper_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudCDSSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_vcdss_cluster`"))
	}

	cluster, err := getCDSSCluster(ctx, config, *StringPtrOrNil(d.GetOk("id")))
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*StringPtrOrNil(d.GetOk("id")))
	d.Set("name", cluster.ClusterName)
	d.Set("kafka_version_code", cluster.KafkaVersionCode)
	d.Set("os_product_code", cluster.SoftwareProductCode)
	d.Set("vpc_no", strconv.Itoa(int(cluster.VpcNo)))
	d.Set("config_group_no", strconv.Itoa(int(cluster.ConfigGroupNo)))

	var cList []map[string]interface{}
	var mList []map[string]interface{}
	var bList []map[string]interface{}
	var eList []map[string]interface{}

	cList = append(cList, map[string]interface{}{
		"user_name": cluster.KafkaManagerUserName,
	})
	mList = append(mList, map[string]interface{}{
		"node_product_code": cluster.ManagerNodeProductCode,
		"subnet_no":         strconv.Itoa(int(cluster.ManagerNodeSubnetNo)),
	})
	bList = append(bList, map[string]interface{}{
		"node_product_code": cluster.BrokerNodeProductCode,
		"subnet_no":         strconv.Itoa(int(cluster.BrokerNodeSubnetNo)),
		"node_count":        strconv.Itoa(int(cluster.BrokerNodeCount)),
		"storage_size":      strconv.Itoa(int(cluster.BrokerNodeStorageSize)),
	})

	endpoints, err := getBrokerInfo(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	eList = append(eList, map[string]interface{}{
		"broker_node_list":                                   endpoints.BrokerNodeList,
		"broker_tls_node_list":                               endpoints.BrokerTlsNodeList,
		"public_endpoint_broker_node_list":                   endpoints.PublicEndpointBrokerNodeList,
		"public_endpoint_broker_node_listener_port_list":     endpoints.PublicEndpointBrokerNodeListenerPortList,
		"public_endpoint_broker_tls_node_list":               endpoints.PublicEndpointBrokerTlsNodeList,
		"public_endpoint_broker_tls_node_listener_port_list": endpoints.PublicEndpointBrokerTlsNodeListenerPortList,
		"local_dns_list":                                     endpoints.LocalDnsList,
		"local_dns_tls_list":                                 endpoints.LocalDnsTlsList,
		"zookeeper_list":                                     endpoints.ZookeeperList,
	})

	// Only set data intersection between resource and list
	if err := d.Set("cmak", cList); err != nil {
		log.Printf("[WARN] Error setting cmak set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("manager_node", mList); err != nil {
		log.Printf("[WARN] Error setting manager_node set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("broker_nodes", bList); err != nil {
		log.Printf("[WARN] Error setting broker_nodes set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("endpoints", eList); err != nil {
		log.Printf("[WARN] Error setting endpoints set for (%s): %s", d.Id(), err)
	}

	return nil
}
