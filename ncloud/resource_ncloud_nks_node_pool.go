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
			"cluster_id": {
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
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_name_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Computed: true,
				Optional: true,
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

	clusterId := StringPtrOrNil(d.GetOk("cluster_id"))
	reqParams := &vnks.NodePoolCreationBody{
		Name:        StringPtrOrNil(d.GetOk("name")),
		NodeCount:   Int32PtrOrNil(d.GetOk("node_count")),
		SubnetNo:    nil,
		ProductCode: StringPtrOrNil(d.GetOk("product_code")),
	}

	logCommonRequest("resourceNcloudNKSNodePoolCreate", reqParams)
	err := config.Client.vnks.V2Api.ClustersUuidNodePoolPost(ctx, reqParams, clusterId)

	if err != nil {
		logErrorResponse("resourceNcloudNKSNodePoolCreate", err, reqParams)
		return diag.FromErr(err)
	}

	logResponse("resourceNcloudNKSNodePoolCreate", reqParams)
	if err := waitForNKSNodePoolActive(ctx, d, config, clusterId, reqParams.Name); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ncloud.StringValue(reqParams.Name))
	return resourceNcloudNKSNodePoolRead(ctx, d, meta)
}

func resourceNcloudNKSNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_node_pool`"))
	}

	clusterId := StringPtrOrNil(d.GetOk("cluster_id"))
	name := ncloud.String(d.Id())
	nodePool, err := getNKSNodePool(ctx, config, clusterId, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if nodePool == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(nodePool.Name))

	d.Set("instance_no", nodePool.InstanceNo)
	d.Set("name", nodePool.Name)
	d.Set("status", nodePool.Status)
	d.Set("product_code", nodePool.ProductCode)

	if err := d.Set("subnet_no_list", flattenSubnetNoList(nodePool.SubnetNoList)); err != nil {
		log.Printf("[WARN] Error setting subet no list set for (%s): %s", d.Id(), err)
	}
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

	clusterId := StringPtrOrNil(d.GetOk("cluster_id"))
	name := ncloud.String(d.Id())
	instanceNo := Int32PtrOrNil(d.GetOk("instance_no"))
	strInstanceNo := ncloud.String(fmt.Sprintf("%d", *instanceNo))

	if d.HasChanges("node_count", "autoscale") {
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterId, name); err != nil {
			return diag.FromErr(err)
		}
		reqParams := &vnks.NodePoolUpdateBody{
			NodeCount: Int32PtrOrNil(d.GetOk("node_count")),
		}

		if _, ok := d.GetOk("autoscale"); ok {
			reqParams.Autoscale = expandAutoScaleUpdate(d.Get("autoscale").([]interface{}))
		}

		err := config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoPatch(ctx, reqParams, clusterId, strInstanceNo)
		if err != nil {
			logErrorResponse("resourceNcloudNKSNodePoolUpdate", err, reqParams)
			return diag.FromErr(err)
		}

		logResponse("resourceNcloudNKSNodePoolUpdate", reqParams)
		if err := waitForNKSNodePoolActive(ctx, d, config, clusterId, name); err != nil {
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

	clusterId := StringPtrOrNil(d.GetOk("cluster_id"))
	name := ncloud.String(d.Id())
	instanceNo := Int32PtrOrNil(d.GetOk("instance_no"))
	strInstanceNo := ncloud.String(fmt.Sprintf("%d", *instanceNo))

	if err := waitForNKSNodePoolActive(ctx, d, config, clusterId, name); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudNKSNodePoolDelete", d.Id())
	if err := config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoDelete(ctx, clusterId, strInstanceNo); err != nil {
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
		Pending: []string{"NODE_SCALE_DOWN"},
		Target:  []string{"NULL"},
		Refresh: func() (result interface{}, state string, err error) {
			clusterId := StringPtrOrNil(d.GetOk("cluster_id"))
			name := ncloud.String(d.Id())
			np, err := getNKSNodePool(ctx, config, clusterId, name)
			if err != nil {
				return nil, "", err
			}
			if np == nil {
				return name, "NULL", nil
			}
			return np, ncloud.StringValue(np.Status), nil

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

func flattenSubnetNoList(list []*int32) (res []*string) {
	for _, v := range list {
		res = append(res, ncloud.IntString(int(ncloud.Int32Value(v))))
	}
	return
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

func expandAutoScaleUpdate(as []interface{}) *vnks.AutoscalerUpdate {
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
