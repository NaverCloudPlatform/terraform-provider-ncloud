package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vses2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
	"strconv"
	"time"
)

func init() {
	RegisterResource("ncloud_ses_cluster", resourceNcloudSESCluster())
}

const (
	SESStatusCreatingCode = "creating"
	SESStatusChangingCode = "changing"
	SESStatusWorkingCode  = "working"
	SESStatusRunningCode  = "running"
	SESStatusDeletingCode = "deleting"
	SESStatusReturnCode   = "return"
	SESStatusNullCode     = "null"
)

func resourceNcloudSESCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudSESClusterCreate,
		ReadContext:   resourceNcloudSESClusterRead,
		UpdateContext: resourceNcloudSESClusterUpdate,
		DeleteContext: resourceNcloudSESClusterDelete,
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
			"service_group_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 15),
					validation.StringMatch(regexp.MustCompile(`^[a-z]+[a-z0-9]*(-[a-z0-9]+)*[a-z0-9]+$`), "Composed of alphabets(lower-case), numbers, non-consecutive hyphen (-)."),
				)),
			},
			"search_engine": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_code": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dashboard_port": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"user_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 15)),
						},
						"user_password": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"os_image_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_no": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"manager_node": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_dual_manager": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"subnet_no": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"acg_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"data_node": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_no": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"count": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IntAtLeast(3)),
						},
						"acg_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
							ValidateDiagFunc: ToDiagFunc(validation.All(
								validation.IntBetween(10, 2000),
								validation.IntDivisibleBy(10)),
							),
						},
					},
				},
			},
			"master_node": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_no": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"count": {
							Type:             schema.TypeInt,
							Optional:         true,
							ForceNew:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IntInSlice([]int{3, 5})),
						},
						"acg_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acg_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"manager_node_instance_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cluster_node_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_instance_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"compute_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudSESClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_ses_cluster`"))
	}

	searchEngineParamsMap := d.Get("search_engine").([]interface{})[0].(map[string]interface{})
	dataNodeParamsMap := d.Get("data_node").([]interface{})[0].(map[string]interface{})
	managerNodeParamsMap := d.Get("manager_node").([]interface{})[0].(map[string]interface{})

	isMasterOnlyNodeActivated := false
	var masterNodeProductCode *string
	var masterNodeSubnetNo *int32
	var masterNodeCount *int32

	masterNodeParams := d.Get("master_node")
	if masterNodeParams != nil && len(masterNodeParams.([]interface{})) > 0 {
		masterNodeParamsMap := masterNodeParams.([]interface{})[0].(map[string]interface{})
		isMasterOnlyNodeActivated = true
		masterNodeProductCode = StringPtrOrNil(masterNodeParamsMap["product_code"], true)
		masterNodeSubnetNo = Int32PtrOrNil(masterNodeParamsMap["subnet_no"], true)
		masterNodeCount = Int32PtrOrNil(masterNodeParamsMap["count"], true)
	}

	var reqParams = &vses2.CreateClusterRequestVo{
		ClusterName:               StringPtrOrNil(d.GetOk("cluster_name")),
		SearchEngineVersionCode:   StringPtrOrNil(searchEngineParamsMap["version_code"], true),
		SearchEngineUserName:      StringPtrOrNil(searchEngineParamsMap["user_name"], true),
		SearchEngineUserPassword:  StringPtrOrNil(searchEngineParamsMap["user_password"], true),
		SearchEngineDashboardPort: StringPtrOrNil(searchEngineParamsMap["dashboard_port"], true),
		SoftwareProductCode:       StringPtrOrNil(d.GetOk("os_image_code")),
		VpcNo:                     Int32PtrOrNil(d.GetOk("vpc_no")),
		IsDualManager:             BoolPtrOrNil(managerNodeParamsMap["is_dual_manager"], true),
		ManagerNodeProductCode:    StringPtrOrNil(managerNodeParamsMap["product_code"], true),
		ManagerNodeSubnetNo:       Int32PtrOrNil(managerNodeParamsMap["subnet_no"], true),
		DataNodeProductCode:       StringPtrOrNil(dataNodeParamsMap["product_code"], true),
		DataNodeSubnetNo:          Int32PtrOrNil(dataNodeParamsMap["subnet_no"], true),
		DataNodeCount:             Int32PtrOrNil(dataNodeParamsMap["count"], true),
		DataNodeStorageSize:       Int32PtrOrNil(dataNodeParamsMap["storage_size"], true),
		IsMasterOnlyNodeActivated: BoolPtrOrNil(isMasterOnlyNodeActivated, true),
		MasterNodeProductCode:     masterNodeProductCode,
		MasterNodeSubnetNo:        masterNodeSubnetNo,
		MasterNodeCount:           masterNodeCount,
		LoginKeyName:              StringPtrOrNil(d.GetOk("login_key_name")),
	}
	logCommonRequest("resourceNcloudSESClusterCreate", reqParams)
	resp, _, err := config.Client.vses.V2Api.CreateClusterUsingPOST(ctx, *reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudSESClusterCreate", err, reqParams)
		return diag.FromErr(err)
	}
	id := strconv.Itoa(int(ncloud.Int32Value(resp.Result.ServiceGroupInstanceNo)))

	logResponse("resourceNcloudSESClusterCreate", resp)
	if err := waitForSESClusterActive(ctx, d, config, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return resourceNcloudSESClusterRead(ctx, d, meta)
}

func resourceNcloudSESClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_ses_cluster`"))
	}

	cluster, err := getSESCluster(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(cluster.ServiceGroupInstanceNo))
	d.Set("id", cluster.ServiceGroupInstanceNo)
	d.Set("service_group_instance_no", cluster.ServiceGroupInstanceNo)
	d.Set("cluster_name", cluster.ClusterName)
	d.Set("os_image_code", cluster.SoftwareProductCode)
	d.Set("vpc_no", cluster.VpcNo)
	d.Set("login_key_name", cluster.LoginKeyName)
	d.Set("manager_node_instance_no_list", cluster.ManagerNodeInstanceNoList)

	searchEngineMap := d.Get("search_engine").([]interface{})[0].(map[string]interface{})
	searchEngineSet := schema.NewSet(schema.HashResource(resourceNcloudSESCluster().Schema["search_engine"].Elem.(*schema.Resource)), []interface{}{})
	searchEngineSet.Add(map[string]interface{}{
		"version_code":   *cluster.SearchEngineVersionCode,
		"user_name":      *cluster.SearchEngineUserName,
		"user_password":  searchEngineMap["user_password"],
		"port":           *cluster.SearchEnginePort,
		"dashboard_port": *cluster.SearchEngineDashboardPort,
	})
	if err := d.Set("search_engine", searchEngineSet.List()); err != nil {
		log.Printf("[WARN] Error setting search_engine set for (%s): %s", d.Id(), err)
	}

	managerNodeSet := schema.NewSet(schema.HashResource(resourceNcloudSESCluster().Schema["manager_node"].Elem.(*schema.Resource)), []interface{}{})
	managerNodeSet.Add(map[string]interface{}{
		"is_dual_manager": *cluster.IsDualManager,
		"count":           *cluster.ManagerNodeCount,
		"subnet_no":       *cluster.ManagerNodeSubnetNo,
		"product_code":    *cluster.ManagerNodeProductCode,
		"acg_id":          *cluster.ManagerNodeAcgId,
		"acg_name":        *cluster.ManagerNodeAcgName,
	})
	if err := d.Set("manager_node", managerNodeSet.List()); err != nil {
		log.Printf("[WARN] Error setting manager_node set for (%s): %s", d.Id(), err)
	}

	dataNodeSet := schema.NewSet(schema.HashResource(resourceNcloudSESCluster().Schema["data_node"].Elem.(*schema.Resource)), []interface{}{})
	storageSize, _ := strconv.Atoi(*cluster.DataNodeStorageSize)
	dataNodeSet.Add(map[string]interface{}{
		"count":        *cluster.DataNodeCount,
		"subnet_no":    *cluster.DataNodeSubnetNo,
		"product_code": *cluster.DataNodeProductCode,
		"acg_id":       *cluster.DataNodeAcgId,
		"acg_name":     *cluster.DataNodeAcgName,
		"storage_size": storageSize,
	})
	if err := d.Set("data_node", dataNodeSet.List()); err != nil {
		log.Printf("[WARN] Error setting data_node set for (%s): %s", d.Id(), err)
	}

	if cluster.IsMasterOnlyNodeActivated != nil && *cluster.IsMasterOnlyNodeActivated {
		masterNodeSet := schema.NewSet(schema.HashResource(resourceNcloudSESCluster().Schema["master_node"].Elem.(*schema.Resource)), []interface{}{})
		masterNodeSet.Add(map[string]interface{}{
			"count":        *cluster.MasterNodeCount,
			"subnet_no":    *cluster.MasterNodeSubnetNo,
			"product_code": *cluster.MasterNodeProductCode,
			"acg_id":       *cluster.MasterNodeAcgId,
			"acg_name":     *cluster.MasterNodeAcgName,
		})

		if err := d.Set("master_node", masterNodeSet.List()); err != nil {
			log.Printf("[WARN] Error setting master_node set for (%s): %s", d.Id(), err)
		}
	}

	clusterNodeList := schema.NewSet(schema.HashResource(resourceNcloudSESCluster().Schema["cluster_node_list"].Elem.(*schema.Resource)), []interface{}{})
	if cluster.ClusterNodeList != nil {
		for _, clusterNode := range cluster.ClusterNodeList {
			clusterNodeList.Add(map[string]interface{}{
				"compute_instance_no":   clusterNode.ComputeInstanceNo,
				"compute_instance_name": clusterNode.ComputeInstanceName,
				"private_ip":            clusterNode.PrivateIp,
				"server_status":         clusterNode.ServerStatus,
				"node_type":             clusterNode.NodeType,
				"subnet":                clusterNode.Subnet,
			})
		}
	}
	if err := d.Set("cluster_node_list", clusterNodeList.List()); err != nil {
		log.Printf("[WARN] Error setting cluster node list for (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceNcloudSESClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_ses_cluster`"))
	}

	if err := checkSearchEngineChanged(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}
	if err := checkDataNodeChanged(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}
	if err := checkNodeProductCodeChanged(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func checkSearchEngineChanged(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	if d.HasChanges("search_engine") {
		o, n := d.GetChange("search_engine")

		oldSearchEngineMap := o.([]interface{})[0].(map[string]interface{})
		newSearchEngineMap := n.([]interface{})[0].(map[string]interface{})
		if oldSearchEngineMap["user_password"] != newSearchEngineMap["user_password"] {
			logCommonRequest("resourceNcloudSESClusterUpdate", d.Id())
			if err := waitForSESClusterActive(ctx, d, config, d.Id()); err != nil {
				return fmt.Errorf("error waiting for SES Cluster (%s) to become activating: %s", d.Id(), err)
			}

			reqParams := &vses2.ResetSearchEngineUserPasswordRequestVo{
				SearchEngineUserPassword: StringPtrOrNil(newSearchEngineMap["user_password"], true),
			}

			if _, _, err := config.Client.vses.V2Api.ResetSearchEngineUserPasswordUsingPOST(ctx, d.Id(), reqParams); err != nil {
				logErrorResponse("resourceNcloudSESClusterResetSearchEngineUserPassword", err, d.Id())
				return fmt.Errorf("error Reset Search Engine User Password with Cluster (%s) : %s", d.Id(), err)
			}

			if err := waitForSESClusterActive(ctx, d, config, d.Id()); err != nil {
				return fmt.Errorf("error waiting for SES Cluster (%s) to become activating: %s", d.Id(), err)
			}
		}
	}
	return nil
}

func checkDataNodeChanged(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	if d.HasChanges("data_node") {
		o, n := d.GetChange("data_node")

		oldDataNodeMap := o.([]interface{})[0].(map[string]interface{})
		newDataNodeMap := n.([]interface{})[0].(map[string]interface{})

		oldDataNodeCount := *Int32PtrOrNil(oldDataNodeMap["count"], true)
		newDataNodeCount := *Int32PtrOrNil(newDataNodeMap["count"], true)

		if oldDataNodeCount < newDataNodeCount {
			logCommonRequest("resourceNcloudSESClusterUpdate", d.Id())
			if err := waitForSESClusterActive(ctx, d, config, d.Id()); err != nil {
				return fmt.Errorf("error waiting for SES Cluster (%s) to become activating: %s", d.Id(), err)
			}

			reqParams := &vses2.AddNodesInClusterRequestVo{
				NewDataNodeCount: StringPtrOrNil(strconv.Itoa(int(newDataNodeCount-oldDataNodeCount)), true),
			}

			if _, _, err := config.Client.vses.V2Api.AddNodesInClusterUsingPOST(ctx, d.Id(), reqParams); err != nil {
				logErrorResponse("resourceNcloudSESClusterAddNodes", err, d.Id())
				return fmt.Errorf("error Add Nodes to SES Cluster (%s) : %s", d.Id(), err)
			}

			if err := waitForSESClusterActive(ctx, d, config, d.Id()); err != nil {
				return fmt.Errorf("error waiting for SES Cluster (%s) to become activating: %s", d.Id(), err)
			}
		} else if oldDataNodeCount > newDataNodeCount {
			logErrorResponse("resourceNcloudSESClusterAddNodes", nil, d.Id())
			return fmt.Errorf("data node count cannot be decreased")
		}
	}
	return nil
}

func checkNodeProductCodeChanged(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	managerNodeProductCode := getChangedNodeProductCode("manager_node", d)
	dataNodeProductCode := getChangedNodeProductCode("data_node", d)
	masterNodeProductCode := getChangedNodeProductCode("master_node", d)

	if managerNodeProductCode != nil || dataNodeProductCode != nil || masterNodeProductCode != nil {
		reqParams := &vses2.ChangeSpecNodeRequestVo{
			ManagerNodeProductCode: managerNodeProductCode,
			DataNodeProductCode:    dataNodeProductCode,
			MasterNodeProductCode:  masterNodeProductCode,
		}

		if _, _, err := config.Client.vses.V2Api.ChangeSpecNodeUsingPOST1(ctx, d.Id(), reqParams); err != nil {
			logErrorResponse("resourceNcloudSESClusterChangeSpec", nil, d.Id())
			return fmt.Errorf("error Change Node Product Code (%s) : %s", d.Id(), err)
		}

		if err := waitForSESClusterActive(ctx, d, config, d.Id()); err != nil {
			return fmt.Errorf("error waiting for SES Cluster (%s) to become activating: %s", d.Id(), err)
		}
	}
	return nil
}

func getChangedNodeProductCode(nodeType string, d *schema.ResourceData) *string {
	nodeParams := d.Get(nodeType)
	if nodeParams != nil && len(nodeParams.([]interface{})) > 0 {
		if d.HasChanges(nodeType) {
			o, n := d.GetChange(nodeType)
			oldNodeMap := o.([]interface{})[0].(map[string]interface{})
			newNodeMap := n.([]interface{})[0].(map[string]interface{})

			if oldNodeMap["product_code"] != newNodeMap["product_code"] {
				return StringPtrOrNil(newNodeMap["product_code"], true)
			}
		}
	}
	return nil
}

func resourceNcloudSESClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_ses_cluster`"))
	}

	if err := waitForSESClusterActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudSESClusterDelete", d.Id())
	if _, _, err := config.Client.vses.V2Api.DeleteClusterUsingDELETE(ctx, d.Id()); err != nil {
		logErrorResponse("resourceNcloudSESClusterDelete", err, d.Id())
		return diag.FromErr(err)
	}

	if err := waitForSESClusterDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSESClusterDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{SESStatusRunningCode, SESStatusDeletingCode},
		Target:  []string{SESStatusReturnCode},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getSESCluster(ctx, config, d.Id())
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return d.Id(), SESStatusNullCode, nil
			}
			return cluster, ncloud.StringValue(cluster.ClusterStatus), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for SES Cluster (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForSESClusterActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{SESStatusCreatingCode, SESStatusChangingCode},
		Target:  []string{SESStatusRunningCode},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getSESCluster(ctx, config, id)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return id, SESStatusNullCode, nil
			}
			return cluster, ncloud.StringValue(cluster.ClusterStatus), nil

		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for SES Cluster (%s) to become activating: %s", id, err)
	}
	return nil
}

func getSESCluster(ctx context.Context, config *ProviderConfig, id string) (*vses2.OpenApiGetClusterInfoResponseVo, error) {

	resp, _, err := config.Client.vses.V2Api.GetClusterInfoUsingGET(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func getSESClusters(ctx context.Context, config *ProviderConfig) (*vses2.GetSearchEngineClusterInfoListResponse, error) {

	resp, _, err := config.Client.vses.V2Api.GetClusterInfoListUsingGET(ctx, nil)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
