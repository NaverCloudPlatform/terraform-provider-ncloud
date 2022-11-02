package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"time"
)

func init() {
	RegisterResource("ncloud_cdss_cluster", resourceNcloudCDSSCluster())
}

const (
	StatusCreating = "creating"
	StatusChanging = "changing"
	StatusRunning  = "running"
	StatusDeleting = "deleting"
	StatusError    = "error"
	StatusReturn   = "return"
	StatusNull     = "null"
)

func resourceNcloudCDSSCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudCDSSClusterCreate,
		ReadContext:   resourceNcloudCDSSClusterRead,
		DeleteContext: resourceNcloudCDSSClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultCreateTimeout),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 15)),
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"os_product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_group_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cmak": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 15)),
						},
						"user_password": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},
			"manager_node": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_product_code": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"broker_nodes": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_product_code": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"node_count": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
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

func resourceNcloudCDSSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_cluster`"))
	}

	cSet := d.Get("cmak").(*schema.Set)
	cList := cSet.List()
	cMap := cList[0].(map[string]interface{})

	mSet := d.Get("manager_node").(*schema.Set)
	mList := mSet.List()
	mMap := mList[0].(map[string]interface{})

	bSet := d.Get("broker_nodes").(*schema.Set)
	bList := bSet.List()
	bMap := bList[0].(map[string]interface{})

	reqParams := vcdss.CreateCluster{
		ClusterName:              *StringPtrOrNil(d.GetOk("name")),
		KafkaVersionCode:         *StringPtrOrNil(d.GetOk("kafka_version_code")),
		KafkaManagerUserName:     *StringPtrOrNil(cMap["user_name"], true),
		KafkaManagerUserPassword: *StringPtrOrNil(cMap["user_password"], true),
		SoftwareProductCode:      *StringPtrOrNil(d.GetOk("os_product_code")),
		VpcNo:                    *getInt32FromString(d.GetOk("vpc_no")),
		ManagerNodeProductCode:   *StringPtrOrNil(mMap["node_product_code"], true),
		ManagerNodeSubnetNo:      *getInt32FromString(mMap["subnet_no"], true),
		BrokerNodeProductCode:    *StringPtrOrNil(bMap["node_product_code"], true),
		BrokerNodeCount:          *getInt32FromString(bMap["node_count"], true),
		BrokerNodeSubnetNo:       *getInt32FromString(bMap["subnet_no"], true),
		BrokerNodeStorageSize:    *getInt32FromString(bMap["storage_size"], true),
		ConfigGroupNo:            *getInt32FromString(d.GetOk("config_group_no")),
	}

	logCommonRequest("resourceNcloudCDSSClusterCreate", reqParams)
	resp, _, err := config.Client.vcdss.V1Api.ClusterCreateCDSSClusterReturnServiceGroupInstanceNoPost(ctx, reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudCDSSClusterCreate", err, reqParams)
		return diag.FromErr(err)
	}
	logResponse("resourceNcloudCDSSClusterCreate", resp)

	uuid := strconv.Itoa(int(ncloud.Int32Value(&resp.Result.ServiceGroupInstanceNo)))
	if err := waitForCDSSClusterActive(ctx, d, config, uuid); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)
	return resourceNcloudCDSSClusterRead(ctx, d, meta)
}

func resourceNcloudCDSSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_cluster`"))
	}

	cluster, err := getCDSSCluster(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", cluster.ClusterName)
	d.Set("kafka_version_code", cluster.KafkaVersionCode)
	d.Set("os_product_code", cluster.SoftwareProductCode)
	d.Set("vpc_no", cluster.VpcNo)
	d.Set("config_group_no", cluster.ConfigGroupNo)

	cSet := schema.NewSet(schema.HashResource(resourceNcloudCDSSCluster().Schema["cmak"].Elem.(*schema.Resource)), []interface{}{})
	mSet := schema.NewSet(schema.HashResource(resourceNcloudCDSSCluster().Schema["manager_node"].Elem.(*schema.Resource)), []interface{}{})
	bSet := schema.NewSet(schema.HashResource(resourceNcloudCDSSCluster().Schema["broker_nodes"].Elem.(*schema.Resource)), []interface{}{})

	cSet.Add(map[string]interface{}{
		"user_name": cluster.KafkaManagerUserName,
	})
	mSet.Add(map[string]interface{}{
		"node_product_code": cluster.ManagerNodeProductCode,
		"subnet_no":         cluster.ManagerNodeSubnetNo,
	})
	bSet.Add(map[string]interface{}{
		"node_product_code": cluster.BrokerNodeProductCode,
		"subnet_no":         cluster.BrokerNodeSubnetNo,
		"node_count":        cluster.BrokerNodeCount,
		"storage_size":      cluster.BrokerNodeStorageSize,
	})

	// Only set data intersection between resource and list
	if err := d.Set("cmak", cSet.List()); err != nil {
		log.Printf("[WARN] Error setting cmak set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("manager_node", mSet.List()); err != nil {
		log.Printf("[WARN] Error setting manager_node set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("broker_nodes", bSet.List()); err != nil {
		log.Printf("[WARN] Error setting broker_nodes set for (%s): %s", d.Id(), err)
	}

	endpoints, err := getBrokerInfo(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	eSet := schema.NewSet(schema.HashResource(resourceNcloudCDSSCluster().Schema["endpoints"].Elem.(*schema.Resource)), []interface{}{})
	eSet.Add(map[string]interface{}{
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
	if err := d.Set("broker_nodes", bSet.List()); err != nil {
		log.Printf("[WARN] Error setting endpoints set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudCDSSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	if err := waitForCDSSClusterActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudCDSSClusterDelete", d.Id())
	if _, _, err := config.Client.vcdss.V1Api.ClusterDeleteCDSSClusterServiceGroupInstanceNoDelete(ctx, d.Id()); err != nil {
		logErrorResponse("resourceNcloudCDSSClusterDelete", err, d.Id())
		return diag.FromErr(err)
	}

	if err := waitForCDSSClusterDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCDSSClusterDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{StatusDeleting},
		Target:  []string{StatusReturn},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getCDSSCluster(ctx, config, d.Id())
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return d.Id(), StatusNull, nil
			}
			return cluster, cluster.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for VCDSS Cluster (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForCDSSClusterActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, uuid string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{StatusCreating, StatusChanging},
		Target:  []string{StatusRunning},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getCDSSCluster(ctx, config, uuid)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return uuid, StatusNull, nil
			}
			return cluster, cluster.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for CDSS Cluster (%s) to become activating: %s", uuid, err)
	}
	return nil
}

func getCDSSCluster(ctx context.Context, config *ProviderConfig, uuid string) (*vcdss.OpenApiGetClusterInfoResponseVo, error) {
	resp, _, err := config.Client.vcdss.V1Api.ClusterGetClusterInfoListServiceGroupInstanceNoPost(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func getBrokerInfo(ctx context.Context, config *ProviderConfig, uuid string) (*vcdss.GetBrokerNodeListsResponseVo, error) {
	resp, _, err := config.Client.vcdss.V1Api.ClusterGetBrokerInfoServiceGroupInstanceNoGet(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
