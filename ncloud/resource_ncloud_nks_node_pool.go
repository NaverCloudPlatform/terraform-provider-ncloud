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
	"strings"
	"time"
)

func init() {
	RegisterResource("ncloud_nks_node_pool", resourceNcloudNKSNodePool())
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

func resourceNcloudNKSNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudNKSNodePoolCreate,
		ReadContext:   resourceNcloudNKSNodePoolRead,
		UpdateContext: resourceNcloudNKSNodePoolUpdate,
		DeleteContext: resourceNcloudNKSNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultUpdateTimeout),
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_no": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_pool_name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscale": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNcloudNKSNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterName := d.Get("cluster_name").(string)
	nodePoolName := d.Get("node_pool_name").(string)
	id := NodePoolCreateResourceID(clusterName, nodePoolName)
	cluster, err := getNKSClusterWithName(ctx, config, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	reqParams := &vnks.NodePoolCreationBody{
		Name:        ncloud.String(nodePoolName),
		NodeCount:   Int32PtrOrNil(d.GetOk("node_count")),
		ProductCode: StringPtrOrNil(d.GetOk("product_code")),
		SubnetNo:    getInt32FromString(d, "subnet_no"),
	}

	if _, ok := d.GetOk("autoscale"); ok {
		reqParams.Autoscale = expandAutoScale(d.Get("autoscale").([]interface{}))
	}

	logCommonRequest("resourceNcloudNKSNodePoolCreate", reqParams)
	err = config.Client.vnks.V2Api.ClustersUuidNodePoolPost(ctx, reqParams, cluster.Uuid)
	if err != nil {
		logErrorResponse("resourceNcloudNKSNodePoolCreate", err, reqParams)
		return diag.FromErr(err)
	}

	logResponse("resourceNcloudNKSNodePoolCreate", reqParams)
	if err := waitForNKSNodePoolActive(ctx, d, config, cluster.Uuid, reqParams.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceNcloudNKSNodePoolRead(ctx, d, meta)
}

func resourceNcloudNKSNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterName, nodePoolName, err := NodePoolParseResourceID(d.Id())
	cluster, err := getNKSClusterWithName(ctx, config, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	nodePool, err := getNKSNodePool(ctx, config, cluster.Uuid, &nodePoolName)
	if err != nil {
		return diag.FromErr(err)
	}

	if nodePool == nil {
		d.SetId("")
		return nil
	}

	d.Set("cluster_name", clusterName)
	d.Set("instance_no", nodePool.InstanceNo)
	d.Set("node_pool_name", nodePool.Name)
	d.Set("status", nodePool.Status)
	d.Set("product_code", nodePool.ProductCode)
	d.Set("node_count", nodePool.NodeCount)
	d.Set("k8s_version", nodePool.K8sVersion)
	d.Set("subnet_name", nodePool.SubnetNameList[0])
	d.Set("subnet_no", strconv.Itoa(int(ncloud.Int32Value(nodePool.SubnetNoList[0]))))

	if err := d.Set("autoscale", flattenAutoscale(nodePool.Autoscale)); err != nil {
		log.Printf("[WARN] Error setting Autoscale set for (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceNcloudNKSNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterName, nodePoolName, err := NodePoolParseResourceID(d.Id())
	cluster, err := getNKSClusterWithName(ctx, config, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceNo := Int32PtrOrNil(d.GetOk("instance_no"))
	strInstanceNo := ncloud.String(fmt.Sprintf("%d", *instanceNo))

	if d.HasChanges("node_count", "autoscale") {
		if err := waitForNKSNodePoolActive(ctx, d, config, cluster.Uuid, &nodePoolName); err != nil {
			return diag.FromErr(err)
		}
		reqParams := &vnks.NodePoolUpdateBody{
			NodeCount: Int32PtrOrNil(d.GetOk("node_count")),
		}

		if _, ok := d.GetOk("autoscale"); ok {
			reqParams.Autoscale = expandAutoScale(d.Get("autoscale").([]interface{}))
		}

		err := config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoPatch(ctx, reqParams, cluster.Uuid, strInstanceNo)
		if err != nil {
			logErrorResponse("resourceNcloudNKSNodePoolUpdate", err, reqParams)
			return diag.FromErr(err)
		}

		logResponse("resourceNcloudNKSNodePoolUpdate", reqParams)
		if err := waitForNKSNodePoolActive(ctx, d, config, cluster.Uuid, &nodePoolName); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceNcloudNKSNodePoolRead(ctx, d, config)
}

func resourceNcloudNKSNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterName, nodePoolName, err := NodePoolParseResourceID(d.Id())
	cluster, err := getNKSClusterWithName(ctx, config, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceNo := Int32PtrOrNil(d.GetOk("instance_no"))
	strInstanceNo := ncloud.String(fmt.Sprintf("%d", *instanceNo))

	if err := waitForNKSNodePoolActive(ctx, d, config, cluster.Uuid, &nodePoolName); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudNKSNodePoolDelete", d.Id())
	if err := config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoDelete(ctx, cluster.Uuid, strInstanceNo); err != nil {
		logErrorResponse("resourceNcloudNKSNodePoolDelete", err, strInstanceNo)
		return diag.FromErr(err)
	}

	if err := waitForNKSNodePoolDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNKSNodePoolDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"NODE_SCALE_DOWN", "DELETING"},
		Target:  []string{"NULL"},
		Refresh: func() (result interface{}, state string, err error) {

			clusterName, nodePoolName, err := NodePoolParseResourceID(d.Id())
			cluster, err := getNKSClusterWithName(ctx, config, clusterName)
			if err != nil {
				return nil, "", err
			}

			np, err := getNKSNodePool(ctx, config, cluster.Uuid, &nodePoolName)
			if err != nil {
				return nil, "", err
			}

			if np == nil {
				return nodePoolName, "NULL", nil
			}

			return np, ncloud.StringValue(np.Status), nil

		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for NKS NodePool (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForNKSNodePoolActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, clusterId *string, name *string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATING", "NODE_SCALE_OUT", "NODE_SCALE_DOWN"},
		Target:  []string{"RUN"},
		Refresh: func() (result interface{}, state string, err error) {
			np, err := getNKSNodePool(ctx, config, clusterId, name)
			if err != nil {
				return nil, "", err
			}
			if np == nil {
				return np, "NULL", nil
			}
			return np, ncloud.StringValue(np.Status), nil

		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for NKS NodePool (%s) to become activating: %s", *name, err)
	}
	return nil
}

func getNKSNodePool(ctx context.Context, config *ProviderConfig, uuid *string, nodePoolName *string) (*vnks.NodePoolRes, error) {

	nps, err := getNKSNodePools(ctx, config, uuid)
	if err != nil {
		return nil, err
	}
	for _, np := range nps {
		if ncloud.StringValue(np.Name) == *nodePoolName {
			return np, nil
		}
	}
	return nil, nil
}

func getNKSNodePools(ctx context.Context, config *ProviderConfig, uuid *string) ([]*vnks.NodePoolRes, error) {
	resp, err := config.Client.vnks.V2Api.ClustersUuidNodePoolGet(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return resp.NodePool, nil
}

func flattenAutoscale(ao *vnks.AutoscaleOption) (res []map[string]interface{}) {
	if ao == nil {
		return
	}
	m := map[string]interface{}{
		"enabled": ncloud.BoolValue(ao.Enabled),
		"min":     ncloud.Int32Value(ao.Min),
		"max":     ncloud.Int32Value(ao.Max),
	}
	res = append(res, m)
	return
}

func expandAutoScale(as []interface{}) *vnks.AutoscalerUpdate {
	if len(as) == 0 {
		return nil
	}
	autoScale := as[0].(map[string]interface{})
	return &vnks.AutoscalerUpdate{
		Enabled: ncloud.Bool(autoScale["enabled"].(bool)),
		Min:     ncloud.Int32(int32(autoScale["min"].(int))),
		Max:     ncloud.Int32(int32(autoScale["max"].(int))),
	}
}

const nodePoolResourceIDSeparator = ":"

func NodePoolCreateResourceID(clusterName, nodePoolName string) string {
	parts := []string{clusterName, nodePoolName}
	id := strings.Join(parts, nodePoolResourceIDSeparator)

	return id
}

func NodePoolParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, nodePoolResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-name%[2]snode-pool-name", id, nodePoolResourceIDSeparator)
}
