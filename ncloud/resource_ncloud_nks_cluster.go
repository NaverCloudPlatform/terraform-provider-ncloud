package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-ncloud/sdk/vnks"
	"log"
	"strconv"
	"time"
)

func init() {
	RegisterResource("ncloud_nks_cluster", resourceNcloudNKSCluster())
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
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultUpdateTimeout),
		},
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
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
			},
			"cluster_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002", "SVR.VNKS.STAND.C004.M016.NET.SSD.B050.G002"}, false)),
				Required:         true,
				ForceNew:         true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
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
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func getInt32(d *schema.ResourceData, key string) *int32 {

	if v, ok := d.GetOk(key); ok {
		intV, err := strconv.Atoi(v.(string))
		if err == nil {
			return ncloud.Int32(int32(intV))
		}
	}
	return nil
}

func getInt32List(d *schema.ResourceData, key string) (int32List []*int32) {
	if list, ok := d.GetOk(key); ok {
		int32List = stringListToInt32List(list.([]interface{}))
	}
	return
}

func stringListToInt32List(list []interface{}) (int32List []*int32) {
	for _, v := range list {
		intV, err := strconv.Atoi(v.(string))
		if err == nil {
			int32List = append(int32List, ncloud.Int32(int32(intV)))
		}
	}
	return
}

func resourceNcloudNKSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	reqParams := &vnks.ClusterInputBody{
		RegionCode: &config.RegionCode,
		//Required
		Name:         StringPtrOrNil(d.GetOk("name")),
		ClusterType:  StringPtrOrNil(d.GetOk("cluster_type")),
		LoginKeyName: StringPtrOrNil(d.GetOk("login_key_name")),
		K8sVersion:   StringPtrOrNil(d.GetOk("k8s_version")),
		ZoneCode:     StringPtrOrNil(d.GetOk("zone")),
		VpcNo:        getInt32(d, "vpc_no"),
		SubnetNoList: getInt32List(d, "subnet_no_list"),
		SubnetLbNo:   getInt32(d, "subnet_lb_no"),
	}

	if _, ok := d.GetOk("log"); ok {
		reqParams.Log = expandLogInput(d.Get("log").([]interface{}))
	}

	logCommonRequest("resourceNcloudNKSClusterCreate", reqParams)
	resp, err := config.Client.vnks.V2Api.ClustersPost(ctx, reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudNKSClusterCreate", err, reqParams)
		return diag.FromErr(err)
	}

	logResponse("resourceNcloudNKSClusterCreate", resp)
	if err := waitForNKSClusterActive(ctx, d, config, ncloud.StringValue(reqParams.Name)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ncloud.StringValue(reqParams.Name))
	return resourceNcloudNKSClusterRead(ctx, d, meta)
}

func resourceNcloudNKSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	cluster, err := getNKSClusterWithName(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(cluster.Name))
	d.Set("uuid", cluster.Uuid)
	d.Set("acg_name", cluster.AcgName)
	d.Set("name", cluster.Name)
	d.Set("cluster_type", cluster.ClusterType)
	d.Set("created_at", cluster.CreatedAt)
	d.Set("endpoint", cluster.Endpoint)
	d.Set("region_code", cluster.RegionCode)
	d.Set("status", cluster.Status)
	d.Set("subnet_name", cluster.SubnetName)
	d.Set("login_key_name", cluster.LoginKeyName)
	d.Set("k8s_version", cluster.K8sVersion)
	d.Set("zone", cluster.ZoneCode)
	d.Set("vpc_no", fmt.Sprintf("%d", *cluster.VpcNo))
	d.Set("vpc_name", cluster.VpcName)
	d.Set("subnet_lb_no", fmt.Sprintf("%d", *cluster.SubnetLbNo))
	d.Set("subnet_lb_name", cluster.SubnetLbName)

	if err := d.Set("subnet_no_list", flattenSubnetNoList(cluster.SubnetNoList)); err != nil {
		log.Printf("[WARN] Error setting subet no list set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudNKSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	uuid := d.Get("uuid").(string)

	if err := waitForNKSClusterActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudNKSClusterDelete", d.Id())
	if err := config.Client.vnks.V2Api.ClustersUuidDelete(ctx, &uuid); err != nil {
		logErrorResponse("resourceNcloudNKSClusterDelete", err, uuid)
		return diag.FromErr(err)
	}

	if err := waitForNKSClusterDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNKSClusterDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"NULL"},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getNKSClusterWithName(ctx, config, d.Id())
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return d.Id(), "NULL", nil
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

func waitForNKSClusterActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, name string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATING", "WORKING"},
		Target:  []string{"RUNNING"},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getNKSClusterWithName(ctx, config, name)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return name, "NULL", nil
			}
			return cluster, ncloud.StringValue(cluster.Status), nil

		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for NKS Cluster (%s) to become activating: %s", name, err)
	}
	return nil
}

func getNKSClusterWithName(ctx context.Context, config *ProviderConfig, name string) (*vnks.Cluster, error) {

	clusters, err := getNKSClusters(ctx, config)
	if err != nil {
		return nil, err
	}
	for _, cluster := range clusters {
		if ncloud.StringValue(cluster.Name) == name {
			return cluster, nil
		}
	}
	return nil, nil
}

func getNKSCluster(ctx context.Context, config *ProviderConfig, uuid string) (*vnks.Cluster, error) {

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

func expandLogInput(logList []interface{}) *vnks.ClusterLogInput {
	if len(logList) == 0 {
		return nil
	}
	log := logList[0].(map[string]interface{})
	return &vnks.ClusterLogInput{
		Audit: ncloud.Bool(log["audit"].(bool)),
	}
}
