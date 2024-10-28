package nks

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

const (
	NKSNodePoolStatusRunCode             = "RUN"
	NKSNodePoolStatusNodeScaleDown       = "NODE_SCALE_DOWN"
	NKSNodePoolStatusNodeScaleOut        = "NODE_SCALE_OUT"
	NKSNodePoolStatusRotateNodeScaleOut  = "ROTATE_NODE_SCALE_OUT"
	NKSNodePoolStatusRotateNodeScaleDown = "ROTATE_NODE_SCALE_DOWN"
	NKSNodePoolStatusUpgrade             = "UPGRADE"
	NKSNodePoolStatusUpdate              = "UPDATING"
	NKSNodePoolIDSeparator               = ":"
)

func ResourceNcloudNKSNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudNKSNodePoolCreate,
		ReadContext:   resourceNcloudNKSNodePoolRead,
		UpdateContext: resourceNcloudNKSNodePoolUpdate,
		DeleteContext: resourceNcloudNKSNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Update: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultCreateTimeout),
		},
		CustomizeDiff: customdiff.Sequence(
			// add subnet nubmer to subnet_no_list when using deprecated subnet_no parameter.
			customdiff.IfValue(
				"subnet_no",
				func(ctx context.Context, subnetNo, meta interface{}) bool {
					return subnetNo.(string) != ""
				},
				func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
					if _, ok := d.GetOk("subnet_no_list"); !ok {
						subnetNo := d.Get("subnet_no").(string)
						return d.SetNew("subnet_no_list", []*string{ncloud.String(subnetNo)})
					}
					return nil
				}),
			customdiff.ForceNewIfChange("subnet_no_list", func(ctx context.Context, old, new, meta any) bool {
				// force new if removed subnet or subnet auto select(emtpy sunbnet_no_list)
				_, removed, autoSelect := getSubnetDiff(old, new)
				return len(removed) > 0 || autoSelect
			}),
		),

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
				Optional: true,
			},
			"node_pool_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`), "Allows only lowercase letters(a-z), numbers, hyphen (-). Must start with an alphabetic character, must end with an English letter or number"))),
			},
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"subnet_no": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "use 'subnet_no_list' instead",
				ConflictsWith: []string{"subnet_no_list"},
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 5,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"software_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_role_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"label": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"taint": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
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
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterUuid := d.Get("cluster_uuid").(string)
	nodePoolName := d.Get("node_pool_name").(string)
	id := NodePoolCreateResourceID(clusterUuid, nodePoolName)

	reqParams := &vnks.NodePoolCreationBody{
		Name:           ncloud.String(nodePoolName),
		NodeCount:      Int32PtrOrNil(d.GetOk("node_count")),
		ProductCode:    StringPtrOrNil(d.GetOk("product_code")),
		SoftwareCode:   StringPtrOrNil(d.GetOk("software_code")),
		ServerSpecCode: StringPtrOrNil(d.GetOk("server_spec_code")),
		StorageSize:    Int32PtrOrNil(d.GetOk("storage_size")),
		ServerRoleId:   StringPtrOrNil(d.GetOk("server_role_id")),
	}

	if list, ok := d.GetOk("subnet_no_list"); ok {
		reqParams.SubnetNoList = ExpandStringInterfaceListToInt32List(list.([]interface{}))
	}

	if _, ok := d.GetOk("autoscale"); ok {
		reqParams.Autoscale = expandNKSNodePoolAutoScale(d.Get("autoscale").([]interface{}))
	}

	LogCommonRequest("resourceNcloudNKSNodePoolCreate", reqParams)
	_, err := config.Client.Vnks.V2Api.ClustersUuidNodePoolPost(ctx, reqParams, ncloud.String(clusterUuid))
	if err != nil {
		LogErrorResponse("resourceNcloudNKSNodePoolCreate", err, reqParams)
		return diag.FromErr(err)
	}

	LogResponse("resourceNcloudNKSNodePoolCreate", reqParams)
	if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, ncloud.StringValue(reqParams.Name)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	if taints, ok := d.GetOk("taint"); ok {
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}

		nodePool, err := GetNKSNodePool(ctx, config, clusterUuid, nodePoolName)
		if err != nil {
			return diag.FromErr(err)
		}

		instanceNo := strconv.Itoa(int(ncloud.Int32Value(nodePool.InstanceNo)))

		nodePoolTaintReq := &vnks.UpdateNodepoolTaintDto{
			Taints: expandNKSNodePoolTaints(taints),
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoTaintsPut(ctx, nodePoolTaintReq, &clusterUuid, &instanceNo)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodePoolCreate - put taints", err, nodePoolTaintReq)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSNodePoolCreate - put taints", reqParams)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, ncloud.StringValue(reqParams.Name)); err != nil {
			return diag.FromErr(err)
		}
	}

	if labels, ok := d.GetOk("label"); ok {
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}

		nodePool, err := GetNKSNodePool(ctx, config, clusterUuid, nodePoolName)
		if err != nil {
			return diag.FromErr(err)
		}

		instanceNo := strconv.Itoa(int(ncloud.Int32Value(nodePool.InstanceNo)))

		labelsReq := &vnks.UpdateNodepoolLabelDto{
			Labels: expandNKSNodePoolLabels(labels),
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoLabelsPut(ctx, labelsReq, &clusterUuid, &instanceNo)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodePoolCreate - put labels", err, labelsReq)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSNodePoolCreate - put labels", reqParams)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, ncloud.StringValue(reqParams.Name)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNcloudNKSNodePoolRead(ctx, d, meta)
}

func resourceNcloudNKSNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nodePool, err := GetNKSNodePool(ctx, config, clusterUuid, nodePoolName)
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
	d.Set("software_code", nodePool.SoftwareCode)
	d.Set("node_count", nodePool.NodeCount)
	d.Set("k8s_version", nodePool.K8sVersion)
	d.Set("server_spec_code", nodePool.ServerSpecCode)
	d.Set("storage_size", nodePool.StorageSize)
	d.Set("server_role_id", nodePool.ServerRoleId)

	if err := d.Set("autoscale", flattenNKSNodePoolAutoScale(nodePool.Autoscale)); err != nil {
		log.Printf("[WARN] Error setting Autoscale set for (%s): %s", d.Id(), err)
	}

	if len(nodePool.SubnetNoList) > 0 {
		if err := d.Set("subnet_no_list", flattenInt32ListToStringList(nodePool.SubnetNoList)); err != nil {
			log.Printf("[WARN] Error setting subnet no list set for (%s): %s", d.Id(), err)
		}
	} else {
		d.Set("subnet_no_list", nil)
	}

	if err := d.Set("taint", flattenNKSNodePoolTaints(nodePool.Taints)); err != nil {
		log.Printf("[WARN] Error setting taints set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("label", flattenNKSNodePoolLabels(nodePool.Labels)); err != nil {
		log.Printf("[WARN] Error setting labels set for (%s): %s", d.Id(), err)
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
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceNo := StringPtrOrNil(d.GetOk("instance_no"))
	k8sVersion := StringPtrOrNil(d.GetOk("k8s_version"))

	if d.HasChanges("k8s_version") {
		_, err = config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoUpgradePatch(ctx, ncloud.String(clusterUuid), instanceNo, k8sVersion, map[string]interface{}{})
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodepoolUpgrade", err, k8sVersion)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSNodepoolUpgrade", k8sVersion)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("node_count", "autoscale") {
		reqParams := &vnks.NodePoolUpdateBody{
			NodeCount: Int32PtrOrNil(d.GetOk("node_count")),
		}

		if _, ok := d.GetOk("autoscale"); ok {
			reqParams.Autoscale = expandNKSNodePoolAutoScale(d.Get("autoscale").([]interface{}))
		}

		err := config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoPatch(ctx, reqParams, ncloud.String(clusterUuid), instanceNo)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodePoolUpdate", err, reqParams)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSNodePoolUpdate", reqParams)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("taint") {
		nodePoolTaintReq := &vnks.UpdateNodepoolTaintDto{
			Taints: expandNKSNodePoolTaints(d.Get("taint")),
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoTaintsPut(ctx, nodePoolTaintReq, &clusterUuid, instanceNo)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodePoolUpdate - put taints", err, nodePoolTaintReq)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSNodePoolUpdate - put taints", nodePoolTaintReq)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("label") {
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}

		labelsReq := &vnks.UpdateNodepoolLabelDto{
			Labels: expandNKSNodePoolLabels(d.Get("label")),
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoLabelsPut(ctx, labelsReq, &clusterUuid, instanceNo)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodePoolUpdate - put labels", err, labelsReq)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSNodePoolUpdate - put labels", labelsReq)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterUuid, nodePoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("subnet_no_list") {

		oldList, newList := d.GetChange("subnet_no_list")
		added, _, _ := getSubnetDiff(oldList, newList)

		subnetReq := &vnks.UpdateNodepoolSubnetDto{
			Subnets: added,
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoSubnetsPatch(ctx, subnetReq, &clusterUuid, instanceNo)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSNodePoolUpdate - addSubnets", err, subnetReq)
			return diag.FromErr(err)
		}

	}

	return resourceNcloudNKSNodePoolRead(ctx, d, config)
}

func resourceNcloudNKSNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
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

	LogCommonRequest("resourceNcloudNKSNodePoolDelete", d.Id())
	if err := config.Client.Vnks.V2Api.ClustersUuidNodePoolInstanceNoDelete(ctx, ncloud.String(clusterUuid), instanceNo); err != nil {
		LogErrorResponse("resourceNcloudNKSNodePoolDelete", err, instanceNo)
		return diag.FromErr(err)
	}

	if err := waitForNKSNodePoolDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNKSNodePoolDeletion(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSNodePoolStatusNodeScaleDown, NKSStatusDeletingCode},
		Target:  []string{NKSStatusNullCode},
		Refresh: func() (result interface{}, state string, err error) {

			clusterUuid, nodePoolName, err := NodePoolParseResourceID(d.Id())
			if err != nil {
				return nil, "", err
			}

			np, err := GetNKSNodePool(ctx, config, clusterUuid, nodePoolName)
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
		Delay:      5 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for NKS NodePool (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForNKSNodePoolActive(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig, clusterUuid string, nodePoolName string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSStatusCreatingCode, NKSNodePoolStatusNodeScaleOut, NKSNodePoolStatusNodeScaleDown, NKSNodePoolStatusUpgrade, NKSNodePoolStatusRotateNodeScaleOut, NKSNodePoolStatusRotateNodeScaleDown, NKSNodePoolStatusUpdate},
		Target:  []string{NKSNodePoolStatusRunCode},
		Refresh: func() (result interface{}, state string, err error) {
			np, err := GetNKSNodePool(ctx, config, clusterUuid, nodePoolName)
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
		Delay:      5 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for NKS NodePool (%s) to become activating: %s", nodePoolName, err)
	}
	return nil
}

func GetNKSNodePool(ctx context.Context, config *conn.ProviderConfig, uuid string, nodePoolName string) (*vnks.NodePool, error) {
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

func getNKSNodePools(ctx context.Context, config *conn.ProviderConfig, uuid string) ([]*vnks.NodePool, error) {
	resp, err := config.Client.Vnks.V2Api.ClustersUuidNodePoolGet(ctx, ncloud.String(uuid))
	if err != nil {
		return nil, err
	}
	LogResponse("getNKSNodePools", resp)

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
	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-uuid%[2]snode-pool-name", id, NKSNodePoolIDSeparator)
}

func getNKSWorkerNodes(ctx context.Context, config *conn.ProviderConfig, uuid string) ([]*vnks.WorkerNode, error) {
	resp, err := config.Client.Vnks.V2Api.ClustersUuidNodesGet(ctx, ncloud.String(uuid))
	if err != nil {
		return nil, err
	}
	return resp.Nodes, nil
}

func getNKSNodePoolWorkerNodes(ctx context.Context, config *conn.ProviderConfig, uuid string, nodePoolName string) ([]*vnks.WorkerNode, error) {
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
