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
	"strconv"
	"time"
)

func init() {
	RegisterResource("ncloud_vcdss_cluster", resourceNcloudVCDSSCluster())
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

func resourceNcloudVCDSSCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudVCDSSClusterCreate,
		ReadContext:   resourceNcloudVCDSSClusterRead,
		DeleteContext: resourceNcloudVCDSSClusterDelete,
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
			"cluster_name": {
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
			"kafka_manager_user_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 15)),
			},
			"kafka_manager_user_password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"software_product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_node_product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_node_subnet_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"broker_node_product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"broker_node_subnet_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"broker_node_count": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"broker_node_storage_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_group_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudVCDSSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_cluster`"))
	}

	reqParams := vcdss.CreateCluster{
		ClusterName:              *StringPtrOrNil(d.GetOk("cluster_name")),
		KafkaVersionCode:         *StringPtrOrNil(d.GetOk("kafka_version_code")),
		KafkaManagerUserName:     *StringPtrOrNil(d.GetOk("kafka_manager_user_name")),
		KafkaManagerUserPassword: *StringPtrOrNil(d.GetOk("kafka_manager_user_password")),
		SoftwareProductCode:      *StringPtrOrNil(d.GetOk("software_product_code")),
		VpcNo:                    *getInt32FromString(d.GetOk("vpc_no")),
		ManagerNodeProductCode:   *StringPtrOrNil(d.GetOk("manager_node_product_code")),
		ManagerNodeSubnetNo:      *getInt32FromString(d.GetOk("manager_node_subnet_no")),
		BrokerNodeProductCode:    *StringPtrOrNil(d.GetOk("broker_node_product_code")),
		BrokerNodeCount:          *getInt32FromString(d.GetOk("broker_node_count")),
		BrokerNodeSubnetNo:       *getInt32FromString(d.GetOk("broker_node_subnet_no")),
		BrokerNodeStorageSize:    *getInt32FromString(d.GetOk("broker_node_storage_size")),
		ConfigGroupNo:            *getInt32FromString(d.GetOk("config_group_no")),
	}

	logCommonRequest("resourceNcloudVCDSSClusterCreate", reqParams)
	resp, _, err := config.Client.vcdss.V1Api.ClusterCreateCDSSClusterReturnServiceGroupInstanceNoPost(ctx, reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudVCDSSClusterCreate", err, reqParams)
		return diag.FromErr(err)
	}
	logResponse("resourceNcloudVCDSSClusterCreate", resp)

	uuid := strconv.Itoa(int(ncloud.Int32Value(&resp.Result.ServiceGroupInstanceNo)))
	if err := waitForVCDSSClusterActive(ctx, d, config, uuid); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)
	return resourceNcloudVCDSSClusterRead(ctx, d, meta)
}

func resourceNcloudVCDSSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_cluster`"))
	}

	cluster, err := getVCDSSCluster(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if cluster == nil {
		d.SetId("")
		return nil
	}

	d.Set("cluster_name", cluster.ClusterName)
	d.Set("kafka_version_code", cluster.KafkaVersionCode)
	d.Set("kafka_manager_user_name", cluster.KafkaManagerUserName)
	d.Set("software_product_code", cluster.SoftwareProductCode)
	d.Set("vpc_no", cluster.VpcNo)
	d.Set("manager_node_product_code", cluster.ManagerNodeProductCode)
	d.Set("manager_node_subnet_no", cluster.ManagerNodeSubnetNo)
	d.Set("broker_node_product_code", cluster.ManagerNodeProductCode)
	d.Set("broker_node_count", cluster.BrokerNodeCount)
	d.Set("broker_node_subnet_no", cluster.BrokerNodeSubnetNo)
	d.Set("broker_node_storage_size", cluster.BrokerNodeStorageSize)

	return nil
}

func resourceNcloudVCDSSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	if err := waitForVCDSSClusterActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudVCDSClusterDelete", d.Id())
	if _, _, err := config.Client.vcdss.V1Api.ClusterDeleteCDSSClusterServiceGroupInstanceNoDelete(ctx, d.Id()); err != nil {
		logErrorResponse("resourceNcloudVCDSSClusterDelete", err, d.Id())
		return diag.FromErr(err)
	}

	if err := waitForVCDSSClusterDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForVCDSSClusterDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{StatusDeleting},
		Target:  []string{StatusReturn},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getVCDSSCluster(ctx, config, d.Id())
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

func waitForVCDSSClusterActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, uuid string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{StatusCreating, StatusChanging},
		Target:  []string{StatusRunning},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := getVCDSSCluster(ctx, config, uuid)
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
		return fmt.Errorf("error waiting for VCDSS Cluster (%s) to become activating: %s", uuid, err)
	}
	return nil
}

func getVCDSSCluster(ctx context.Context, config *ProviderConfig, uuid string) (*vcdss.OpenApiGetClusterInfoResponseVo, error) {
	resp, _, err := config.Client.vcdss.V1Api.ClusterGetClusterInfoListServiceGroupInstanceNoPost(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
