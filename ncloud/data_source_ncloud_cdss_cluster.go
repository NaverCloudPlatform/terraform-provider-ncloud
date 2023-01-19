package ncloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_cdss_cluster", dataSourceNcloudCDSSCluster())
}

func dataSourceNcloudCDSSCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSClusterRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_group_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_image": {
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plaintext": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"tls": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"public_endpoint_plaintext_listener_port": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"public_endpoint_tls_listener_port": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"public_endpoint_plaintext": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"public_endpoint_tls": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"zookeeper": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"hosts_private_endpoint_tls": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"hosts_public_endpoint_tls": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudCDSSClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("dataSource `ncloud_vcdss_cluster`")
	}

	resources, err := getCDSSClusterList(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudCDSSKafkaVersion().Schema)
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	id := resources[0]["id"].(string)
	cluster, err := getCDSSCluster(context.Background(), config, id)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("service_group_instance_no", id)
	d.Set("name", cluster.ClusterName)
	d.Set("kafka_version_code", cluster.KafkaVersionCode)
	d.Set("os_image", cluster.SoftwareProductCode)
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

	endpoints, err := getBrokerInfo(context.Background(), config, d.Id())
	if err != nil {
		return err
	}

	commaSplitFn := func(c rune) bool {
		return c == ','
	}
	newlineSplitFn := func(c rune) bool {
		return c == '\n'
	}
	eList = append(eList, map[string]interface{}{
		"plaintext": strings.FieldsFunc(endpoints.BrokerNodeList, commaSplitFn),
		"tls":       strings.FieldsFunc(endpoints.BrokerTlsNodeList, commaSplitFn),
		"public_endpoint_plaintext_listener_port": strings.FieldsFunc(endpoints.PublicEndpointBrokerNodeListenerPortList, newlineSplitFn),
		"public_endpoint_tls_listener_port":       strings.FieldsFunc(endpoints.PublicEndpointBrokerTlsNodeListenerPortList, newlineSplitFn),
		"public_endpoint_plaintext":               strings.FieldsFunc(endpoints.PublicEndpointBrokerNodeList, newlineSplitFn),
		"public_endpoint_tls":                     strings.FieldsFunc(endpoints.PublicEndpointBrokerTlsNodeList, newlineSplitFn),
		"zookeeper":                               strings.FieldsFunc(endpoints.ZookeeperList, commaSplitFn),
		"hosts_private_endpoint_tls":              strings.FieldsFunc(endpoints.LocalDnsList, newlineSplitFn),
		"hosts_public_endpoint_tls":               strings.FieldsFunc(endpoints.LocalDnsTlsList, newlineSplitFn),
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

func getCDSSClusterList(config *ProviderConfig) ([]map[string]interface{}, error) {
	logCommonRequest("GetCDSSClusterList", "")
	resp, _, err := config.Client.vcdss.V1Api.ClusterGetClusterInfoListPost(context.Background(), vcdss.GetClusterRequest{})

	if err != nil {
		logErrorResponse("GetCDSSClusterList", err, "")
		return nil, err
	}

	logResponse("GetCDSSClusterList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.AllowedClusters {
		instance := map[string]interface{}{
			"id":   ncloud.StringValue(&r.ServiceGroupInstanceNo),
			"name": ncloud.StringValue(&r.ClusterName),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
