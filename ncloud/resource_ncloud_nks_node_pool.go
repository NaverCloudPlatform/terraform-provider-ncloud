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
	"strings"
	"time"
)

func init() {
	RegisterResource("ncloud_nks_node_pool", resourceNcloudNKSNodePool())
}

const (
	NKSNodePoolStatusRunCode       = "RUN"
	NKSNodePoolStatusNodeScaleDown = "NODE_SCALE_DOWN"
	NKSNodePoolStatusNodeScaleOut  = "NODE_SCALE_OUT"
	NKSNodePoolIDSeparator         = ":"
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
			Update: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultCreateTimeout),
		},
		Schema: map[string]*schema.Schema{
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_no": {
				Type:     schema.TypeString,
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
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
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
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kernel_version": {
							Type:     schema.TypeString,
							Computed: true,
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

	clusterUuid := d.Get("cluster_uuid").(string)
	nodePoolName := d.Get("node_pool_name").(string)
	id := NodePoolCreateResourceID(clusterUuid, nodePoolName)

	reqParams := &vnks.NodePoolCreationBody{
		Name:        ncloud.String(nodePoolName),
		NodeCount:   Int32PtrOrNil(d.GetOk("node_count")),
		ProductCode: StringPtrOrNil(d.GetOk("product_code")),
		SubnetNo:    getInt32FromString(d.GetOk("subnet_no")),
	}

	if _, ok := d.GetOk("autoscale"); ok {
		reqParams.Autoscale = expandNKSNodePoolAutoScale(d.Get("autoscale").([]interface{}))
	}

	logCommonRequest("resourceNcloudNKSNodePoolCreate", reqParams)
	err := config.Client.vnks.V2Api.ClustersUuidNodePoolPost(ctx, reqParams, ncloud.String(clusterUuid))
	if err != nil {
		logErrorResponse("resourceNcloudNKSNodePoolCreate", err, reqParams)
		return diag.FromErr(err)
	}

	logResponse("resourceNcloudNKSNodePoolCreate", reqParams)
	if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, ncloud.StringValue(reqParams.Name)); err != nil {
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

	clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
	nodePool, err := getNKSNodePool(ctx, config, clusterUuid, nodePoolName)
	if err != nil {
		return diag.FromErr(err)
	}

	if nodePool == nil {
		d.SetId("")
		return nil
	}

	d.Set("cluster_uuid", clusterUuid)
	d.Set("instance_no", strconv.Itoa(int(ncloud.Int32Value(nodePool.InstanceNo))))
	d.Set("node_pool_name", nodePool.Name)
	d.Set("product_code", nodePool.ProductCode)
	d.Set("node_count", nodePool.NodeCount)
	d.Set("k8s_version", nodePool.K8sVersion)
	d.Set("subnet_no", strconv.Itoa(int(ncloud.Int32Value(nodePool.SubnetNoList[0]))))

	if err := d.Set("autoscale", flattenNKSNodePoolAutoScale(nodePool.Autoscale)); err != nil {
		log.Printf("[WARN] Error setting Autoscale set for (%s): %s", d.Id(), err)
	}

	nodes, err := getNKSNodePoolWorkerNodes(ctx, config, clusterUuid, nodePoolName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nodes", flattenNKSWorkerNodes(nodes)); err != nil {
		log.Printf("[WARN] Error setting workerNodes set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudNKSNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceNo := StringPtrOrNil(d.GetOk("instance_no"))

	if d.HasChanges("node_count", "autoscale") {
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}
		reqParams := &vnks.NodePoolUpdateBody{
			NodeCount: Int32PtrOrNil(d.GetOk("node_count")),
		}

		if _, ok := d.GetOk("autoscale"); ok {
			reqParams.Autoscale = expandNKSNodePoolAutoScale(d.Get("autoscale").([]interface{}))
		}

		err := config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoPatch(ctx, reqParams, ncloud.String(clusterUuid), instanceNo)
		if err != nil {
			logErrorResponse("resourceNcloudNKSNodePoolUpdate", err, reqParams)
			return diag.FromErr(err)
		}

		logResponse("resourceNcloudNKSNodePoolUpdate", reqParams)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
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

	clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceNo := StringPtrOrNil(d.GetOk("instance_no"))
	if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudNKSNodePoolDelete", d.Id())
	if err := config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoDelete(ctx, ncloud.String(clusterUuid), instanceNo); err != nil {
		logErrorResponse("resourceNcloudNKSNodePoolDelete", err, instanceNo)
		return diag.FromErr(err)
	}

	if err := waitForNKSNodePoolDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNKSNodePoolDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSNodePoolStatusNodeScaleDown, NKSStatusDeletingCode},
		Target:  []string{NKSStatusNullCode},
		Refresh: func() (result interface{}, state string, err error) {

			clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
			if err != nil {
				return nil, "", err
			}

			np, err := getNKSNodePool(ctx, config, clusterUuid, nodePoolName)
			if err != nil {
				return nil, "", err
			}

			if np == nil {
				return nodePoolName, NKSStatusNullCode, nil
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

func waitForNKSNodePoolActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, clusterUuid string, nodePoolName string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSStatusCreatingCode, NKSNodePoolStatusNodeScaleOut, NKSNodePoolStatusNodeScaleDown},
		Target:  []string{NKSNodePoolStatusRunCode},
		Refresh: func() (result interface{}, state string, err error) {
			np, err := getNKSNodePool(ctx, config, clusterUuid, nodePoolName)
			if err != nil {
				return nil, "", err
			}
			if np == nil {
				return np, NKSStatusNullCode, nil
			}
			return np, ncloud.StringValue(np.Status), nil

		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for NKS NodePool (%s) to become activating: %s", nodePoolName, err)
	}
	return nil
}

func getNKSNodePool(ctx context.Context, config *ProviderConfig, uuid string, nodePoolName string) (*vnks.NodePoolRes, error) {
	nps, err := getNKSNodePools(ctx, config, uuid)
	if err != nil {
		return nil, err
	}
	for _, np := range nps {
		if ncloud.StringValue(np.Name) == nodePoolName {
			return np, nil
		}
	}
	return nil, nil
}

func getNKSNodePools(ctx context.Context, config *ProviderConfig, uuid string) ([]*vnks.NodePoolRes, error) {
	resp, err := config.Client.vnks.V2Api.ClustersUuidNodePoolGet(ctx, ncloud.String(uuid))
	if err != nil {
		return nil, err
	}
	return resp.NodePool, nil
}

func NodePoolCreateResourceID(clusterName, nodePoolName string) string {
	parts := []string{clusterName, nodePoolName}
	id := strings.Join(parts, NKSNodePoolIDSeparator)
	return id
}

func NodePoolParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, NKSNodePoolIDSeparator)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}
	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-name%[2]snode-pool-name", id, NKSNodePoolIDSeparator)
}

func getNKSWorkerNodes(ctx context.Context, config *ProviderConfig, uuid string) ([]*vnks.WorkerNode, error) {
	resp, err := config.Client.vnks.V2Api.ClustersUuidNodesGet(ctx, ncloud.String(uuid))
	if err != nil {
		return nil, err
	}
	return resp.Nodes, nil
}

func getNKSNodePoolWorkerNodes(ctx context.Context, config *ProviderConfig, uuid string, nodePoolName string) ([]*vnks.WorkerNode, error) {
	var res []*vnks.WorkerNode
	wns, err := getNKSWorkerNodes(ctx, config, uuid)
	if err != nil {
		return nil, err
	}

	for _, wn := range wns {
		if ncloud.StringValue(wn.NodePoolName) == nodePoolName {
			res = append(res, wn)
		}
	}
	return res, nil
}
