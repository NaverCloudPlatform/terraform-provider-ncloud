package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"time"
)

func init() {
	RegisterResource("ncloud_nks_cluster", resourceNcloudNKSCluster())
}

const (
	NKSStatusCreatingCode = "CREATING"
	NKSStatusWorkingCode  = "WORKING"
	NKSStatusRunningCode  = "RUNNING"
	NKSStatusDeletingCode = "DELETING"
	NKSStatusNoNodeCode   = "NO_NODE"
	NKSStatusNullCode     = "NULL"
)

func resourceNcloudNKSCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudNKSClusterCreate,
		ReadContext:   resourceNcloudNKSClusterRead,
		DeleteContext: resourceNcloudNKSClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultCreateTimeout),
		},
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 20)),
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_network": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lb_private_subnet_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lb_public_subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kube_network_plugin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"acg_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNKSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	reqParams := &vnks.ClusterInputBody{
		RegionCode: &config.RegionCode,
		//Required
		Name:              StringPtrOrNil(d.GetOk("name")),
		ClusterType:       StringPtrOrNil(d.GetOk("cluster_type")),
		LoginKeyName:      StringPtrOrNil(d.GetOk("login_key_name")),
		K8sVersion:        StringPtrOrNil(d.GetOk("k8s_version")),
		ZoneCode:          StringPtrOrNil(d.GetOk("zone")),
		VpcNo:             getInt32FromString(d.GetOk("vpc_no")),
		SubnetLbNo:        getInt32FromString(d.GetOk("lb_private_subnet_no")),
		LbPublicSubnetNo:  getInt32FromString(d.GetOk("lb_public_subnet_no")),
		KubeNetworkPlugin: StringPtrOrNil(d.GetOk("kube_network_plugin")),
	}

	if publicNetwork, ok := d.GetOk("public_network"); ok {
		reqParams.PublicNetwork = ncloud.Bool(publicNetwork.(bool))
	}

	if list, ok := d.GetOk("subnet_no_list"); ok {
		reqParams.SubnetNoList = expandStringInterfaceListToInt32List(list.([]interface{}))
	}

	if log, ok := d.GetOk("log"); ok {
		reqParams.Log = expandNKSClusterLogInput(log.([]interface{}))
	}

	logCommonRequest("resourceNcloudNKSClusterCreate", reqParams)
	resp, err := config.Client.vnks.V2Api.ClustersPost(ctx, reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudNKSClusterCreate", err, reqParams)
		return diag.FromErr(err)
	}
	uuid := ncloud.StringValue(resp.Uuid)

	logResponse("resourceNcloudNKSClusterCreate", resp)
	if err := waitForNKSClusterActive(ctx, d, config, uuid); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)
	return resourceNcloudNKSClusterRead(ctx, d, meta)
}

func resourceNcloudNKSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	cluster, err := getNKSCluster(ctx, config, d.Id())
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
	d.Set("lb_private_subnet_no", strconv.Itoa(int(ncloud.Int32Value(cluster.SubnetLbNo))))
	d.Set("kube_network_plugin", cluster.KubeNetworkPlugin)
	d.Set("acg_no", strconv.Itoa(int(ncloud.Int32Value(cluster.AcgNo))))
	if cluster.LbPublicSubnetNo != nil {
		d.Set("lb_public_subnet_no", strconv.Itoa(int(ncloud.Int32Value(cluster.LbPublicSubnetNo))))
	}
	if cluster.PublicNetwork != nil {
		d.Set("public_network", cluster.PublicNetwork)
	}

	if err := d.Set("log", flattenNKSClusterLogInput(cluster.Log)); err != nil {
		log.Printf("[WARN] Error setting cluster log for (%s): %s", d.Id(), err)
	}

	if err := d.Set("subnet_no_list", flattenInt32ListToStringList(cluster.SubnetNoList)); err != nil {
		log.Printf("[WARN] Error setting subnet no list set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudNKSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	if err := waitForNKSClusterActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudNKSClusterDelete", d.Id())
	if err := config.Client.vnks.V2Api.ClustersUuidDelete(ctx, ncloud.String(d.Id())); err != nil {
		logErrorResponse("resourceNcloudNKSClusterDelete", err, d.Id())
		return diag.FromErr(err)
	}

	if err := waitForNKSClusterDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNKSClusterDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSStatusDeletingCode},
		Target:  []string{NKSStatusNullCode},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getNKSClusterFromList(ctx, config, d.Id())
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return d.Id(), NKSStatusNullCode, nil
			}
			return cluster, ncloud.StringValue(cluster.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for NKS Cluster (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForNKSClusterActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, uuid string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSStatusCreatingCode, NKSStatusWorkingCode},
		Target:  []string{NKSStatusRunningCode, NKSStatusNoNodeCode},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getNKSCluster(ctx, config, uuid)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return uuid, NKSStatusNullCode, nil
			}
			return cluster, ncloud.StringValue(cluster.Status), nil

		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for NKS Cluster (%s) to become activating: %s", uuid, err)
	}
	return nil
}

func getNKSCluster(ctx context.Context, config *ProviderConfig, uuid string) (*vnks.Cluster, error) {

	resp, err := config.Client.vnks.V2Api.ClustersUuidGet(ctx, &uuid)
	if err != nil {
		return nil, err
	}
	return resp.Cluster, nil
}

func getNKSClusterFromList(ctx context.Context, config *ProviderConfig, uuid string) (*vnks.Cluster, error) {
	clusters, err := getNKSClusters(ctx, config)
	if err != nil {
		return nil, err
	}
	for _, cluster := range clusters {
		if ncloud.StringValue(cluster.Uuid) == uuid {
			return cluster, nil
		}
	}
	return nil, nil
}

func getNKSClusters(ctx context.Context, config *ProviderConfig) ([]*vnks.Cluster, error) {
	resp, err := config.Client.vnks.V2Api.ClustersGet(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Clusters, nil
}
