package nks

import (
	"context"
	"fmt"
	"log"
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
	NKSStatusCreatingCode = "CREATING"
	NKSStatusWorkingCode  = "WORKING"
	NKSStatusRunningCode  = "RUNNING"
	NKSStatusDeletingCode = "DELETING"
	NKSStatusNoNodeCode   = "NO_NODE"
	NKSStatusNullCode     = "NULL"
)

func ResourceNcloudNKSCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudNKSClusterCreate,
		ReadContext:   resourceNcloudNKSClusterRead,
		DeleteContext: resourceNcloudNKSClusterDelete,
		UpdateContext: resourceNcloudNKSClusterUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Update: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultCreateTimeout),
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("subnet_no_list", func(ctx context.Context, old, new, meta any) bool {
				_, removed, _ := getSubnetDiff(old, new)
				return len(removed) > 0
			}),
			customdiff.ValidateValue("ip_acl_default_action", func(ctx context.Context, value, meta interface{}) error {
				config := meta.(*conn.ProviderConfig)
				if value != "" && checkFinSite(config) {
					return fmt.Errorf("ip_acl_default_action is not supported on fin site")
				}
				return nil
			}),
			customdiff.ValidateValue("ip_acl", func(ctx context.Context, value, meta interface{}) error {
				set := value.(*schema.Set)
				config := meta.(*conn.ProviderConfig)
				if set.Len() > 0 && checkFinSite(config) {
					return fmt.Errorf("ip_acl is not supported on fin site")
				}
				return nil
			}),
		),
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(3, 20)),
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hypervisor_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				MaxItems: 5,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lb_private_subnet_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lb_public_subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
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
			"oidc": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"username_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"groups_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"groups_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"required_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ip_acl_default_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"allow", "deny"}, false)),
			},
			"ip_acl": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"return_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"kms_key_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudNKSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	reqParams := &vnks.ClusterInputBody{
		RegionCode: &config.RegionCode,
		//Required
		Name:              StringPtrOrNil(d.GetOk("name")),
		ClusterType:       StringPtrOrNil(d.GetOk("cluster_type")),
		HypervisorCode:    StringPtrOrNil(d.GetOk("hypervisor_code")),
		LoginKeyName:      StringPtrOrNil(d.GetOk("login_key_name")),
		K8sVersion:        StringPtrOrNil(d.GetOk("k8s_version")),
		ZoneCode:          StringPtrOrNil(d.GetOk("zone")),
		VpcNo:             GetInt32FromString(d.GetOk("vpc_no")),
		SubnetLbNo:        GetInt32FromString(d.GetOk("lb_private_subnet_no")),
		LbPublicSubnetNo:  GetInt32FromString(d.GetOk("lb_public_subnet_no")),
		KubeNetworkPlugin: StringPtrOrNil(d.GetOk("kube_network_plugin")),
		KmsKeyTag:         StringPtrOrNil(d.GetOk("kms_key_tag")),
	}

	if publicNetwork, ok := d.GetOk("public_network"); ok {
		reqParams.PublicNetwork = ncloud.Bool(publicNetwork.(bool))
	}

	if list, ok := d.GetOk("subnet_no_list"); ok {
		reqParams.SubnetNoList = ExpandStringInterfaceListToInt32List(list.([]interface{}))
	}

	if log, ok := d.GetOk("log"); ok {
		reqParams.Log = expandNKSClusterLogInput(log.([]interface{}), reqParams.Log)
	}

	var oidcReq *vnks.UpdateOidcDto
	if oidc, ok := d.GetOk("oidc"); ok {
		oidcReq = expandNKSClusterOIDCSpec(oidc.([]interface{}))
	}

	var ipAclReq *vnks.IpAclsDto
	ipAclDefaultAction, ipAclDefaultActionExist := d.GetOk("ip_acl_default_action")
	ipAcl, ipAclExist := d.GetOk("ip_acl")
	if ipAclDefaultActionExist || ipAclExist {
		ipAclReq = &vnks.IpAclsDto{
			DefaultAction: StringPtrOrNil(ipAclDefaultAction, ipAclDefaultActionExist),
			Entries:       expandNKSClusterIPAcl(ipAcl),
		}
	}

	var returnProtectionReq *vnks.ReturnProtectionDto
	if returnProtection, ok := d.GetOk("return_protection"); ok {
		returnProtectionReq = &vnks.ReturnProtectionDto{
			ReturnProtection: ncloud.Bool(returnProtection.(bool)),
		}
	}
	LogCommonRequest("resourceNcloudNKSClusterCreate", reqParams)
	resp, err := config.Client.Vnks.V2Api.ClustersPost(ctx, reqParams)
	if err != nil {
		LogErrorResponse("resourceNcloudNKSClusterCreate", err, reqParams)
		return diag.FromErr(err)
	}
	uuid := ncloud.StringValue(resp.Uuid)

	LogResponse("resourceNcloudNKSClusterCreate", resp)
	if err := waitForNKSClusterActive(ctx, d, config, uuid); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)

	if oidcReq != nil {
		_, err = config.Client.Vnks.V2Api.ClustersUuidOidcPatch(ctx, oidcReq, resp.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterCreate:oidc", err, oidcReq)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSClusterCreateoidc:oidc", oidcReq)
		if err := waitForNKSClusterActive(ctx, d, config, uuid); err != nil {
			return diag.FromErr(err)
		}
	}

	if ipAclReq != nil && !checkFinSite(config) {
		_, err = config.Client.Vnks.V2Api.ClustersUuidIpAclPatch(ctx, ipAclReq, resp.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterCreate:ipAcl", err, ipAclReq)
			return diag.FromErr(err)
		}
	}

	if returnProtectionReq != nil && *returnProtectionReq.ReturnProtection {
		_, err = config.Client.Vnks.V2Api.ClustersUuidReturnProtectionPatch(ctx, returnProtectionReq, resp.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterCreate:returnProtection", err, returnProtectionReq)
			return diag.FromErr(err)
		}
	}

	return resourceNcloudNKSClusterRead(ctx, d, meta)
}

func resourceNcloudNKSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	cluster, err := GetNKSCluster(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	oidcSpec, err := getOIDCSpec(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ipAcl, err := getIPAcl(ctx, config, d.Id())
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
	d.Set("hypervisor_code", cluster.HypervisorCode)
	d.Set("endpoint", cluster.Endpoint)
	d.Set("login_key_name", cluster.LoginKeyName)
	d.Set("k8s_version", cluster.K8sVersion)
	d.Set("zone", cluster.ZoneCode)
	d.Set("vpc_no", strconv.Itoa(int(ncloud.Int32Value(cluster.VpcNo))))
	d.Set("lb_private_subnet_no", strconv.Itoa(int(ncloud.Int32Value(cluster.SubnetLbNo))))
	d.Set("kube_network_plugin", cluster.KubeNetworkPlugin)
	d.Set("acg_no", strconv.Itoa(int(ncloud.Int32Value(cluster.AcgNo))))
	d.Set("return_protection", cluster.ReturnProtection)
	d.Set("kms_key_tag", cluster.KmsKeyTag)

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

	if oidcSpec != nil {
		oidc := flattenNKSClusterOIDCSpec(oidcSpec)
		if err := d.Set("oidc", oidc); err != nil {
			log.Printf("[WARN] Error setting OIDCSpec set for (%s): %s", d.Id(), err)
		}
	}

	if ipAcl != nil {
		d.Set("ip_acl_default_action", ipAcl.DefaultAction)

		if err := d.Set("ip_acl", flattenNKSClusterIPAclEntries(ipAcl).List()); err != nil {
			log.Printf("[WARN] Error setting ip_acl list set for (%s): %s", d.Id(), err)
		}

	}

	return nil
}

func resourceNcloudNKSClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	cluster, err := GetNKSCluster(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("k8s_version") {
		newVersion := StringPtrOrNil(d.GetOk("k8s_version"))
		_, err := config.Client.Vnks.V2Api.ClustersUuidUpgradePatch(ctx, cluster.Uuid, newVersion, map[string]interface{}{})
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterUpgrade", err, newVersion)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSClusterUpgrade", newVersion)
		if err := waitForNKSClusterActive(ctx, d, config, *cluster.Uuid); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("oidc") {
		var oidcSpec *vnks.UpdateOidcDto
		oidc, _ := d.GetOk("oidc")
		oidcSpec = expandNKSClusterOIDCSpec(oidc.([]interface{}))

		_, err = config.Client.Vnks.V2Api.ClustersUuidOidcPatch(ctx, oidcSpec, cluster.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterOIDCPatch", err, oidcSpec)
			return diag.FromErr(err)
		}

		LogResponse("resourceNcloudNKSClusterOIDCPatch", oidcSpec)
		if err := waitForNKSClusterActive(ctx, d, config, *cluster.Uuid); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("ip_acl", "ip_acl_default_action") && !checkFinSite(config) {

		ipAclReq := &vnks.IpAclsDto{
			DefaultAction: StringPtrOrNil(d.GetOk("ip_acl_default_action")),
			Entries:       []*vnks.IpAclsEntriesDto{},
		}
		if ipAcl, ok := d.GetOk("ip_acl"); ok {
			ipAclReq.Entries = expandNKSClusterIPAcl(ipAcl)
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidIpAclPatch(ctx, ipAclReq, cluster.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterIPAclPatch", err, ipAclReq)
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("log") {

		var logDto *vnks.AuditLogDto
		if log, ok := d.GetOk("log"); ok {
			logDto = expandNKSClusterLogInput(log.([]interface{}), logDto)
		} else {
			logDto.Audit = ncloud.Bool(false)
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidLogPatch(ctx, logDto, cluster.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterLogPatch", err, logDto)
			return diag.FromErr(err)
		}

	}

	if d.HasChanges("lb_private_subnet_no") {

		lbPrivateSubnetNo, _ := strconv.Atoi(d.Get("lb_private_subnet_no").(string))
		_, err = config.Client.Vnks.V2Api.ClustersUuidLbSubnetPatch(ctx, cluster.Uuid, ncloud.Int32(int32(lbPrivateSubnetNo)), map[string]interface{}{"igwYn": ncloud.String("N")})
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterLbPrivateSubnetPatch", err, lbPrivateSubnetNo)
			return diag.FromErr(err)
		}

	}

	if d.HasChanges("lb_public_subnet_no") {

		lbPrivateSubnetNo, _ := strconv.Atoi(d.Get("lb_public_subnet_no").(string))
		_, err = config.Client.Vnks.V2Api.ClustersUuidLbSubnetPatch(ctx, cluster.Uuid, ncloud.Int32(int32(lbPrivateSubnetNo)), map[string]interface{}{"igwYn": ncloud.String("Y")})
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterLbPublicSubnetPatch", err, lbPrivateSubnetNo)
			return diag.FromErr(err)
		}

	}

	if d.HasChanges("subnet_no_list") {

		oldList, newList := d.GetChange("subnet_no_list")
		added, _, _ := getSubnetDiff(oldList, newList)

		subnets := &vnks.AddSubnetDto{
			Subnets: []*vnks.SubnetDto{},
		}

		for _, subnetNo := range added {
			subnets.Subnets = append(subnets.Subnets, &vnks.SubnetDto{Number: subnetNo})
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidAddSubnetPatch(ctx, subnets, cluster.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterAddSubnetsPatch", err, subnets)
			return diag.FromErr(err)
		}

	}

	if d.HasChanges("return_protection") {

		returnProtectionReq := &vnks.ReturnProtectionDto{}
		if returnProtection, ok := d.GetOk("return_protection"); ok {
			returnProtectionReq.ReturnProtection = ncloud.Bool(returnProtection.(bool))
		} else {
			returnProtectionReq.ReturnProtection = ncloud.Bool(false)
		}

		_, err = config.Client.Vnks.V2Api.ClustersUuidReturnProtectionPatch(ctx, returnProtectionReq, cluster.Uuid)
		if err != nil {
			LogErrorResponse("resourceNcloudNKSClusterReturnProtectionPatch", err, returnProtectionReq)
			return diag.FromErr(err)
		}
	}

	return resourceNcloudNKSClusterRead(ctx, d, config)
}

func resourceNcloudNKSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_nks_cluster`"))
	}

	if err := waitForNKSClusterActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	LogCommonRequest("resourceNcloudNKSClusterDelete", d.Id())
	if err := config.Client.Vnks.V2Api.ClustersUuidDelete(ctx, ncloud.String(d.Id())); err != nil {
		LogErrorResponse("resourceNcloudNKSClusterDelete", err, d.Id())
		return diag.FromErr(err)
	}

	if err := waitForNKSClusterDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNKSClusterDeletion(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSStatusDeletingCode},
		Target:  []string{NKSStatusNullCode, NKSStatusRunningCode}, // ToDo: remove runnig status after external autoscaler callback removed.
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
		Delay:      5 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for NKS Cluster (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForNKSClusterActive(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig, uuid string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{NKSStatusCreatingCode, NKSStatusWorkingCode},
		Target:  []string{NKSStatusRunningCode, NKSStatusNoNodeCode},
		Refresh: func() (result interface{}, state string, err error) {
			cluster, err := GetNKSCluster(ctx, config, uuid)
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
		Delay:      5 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for NKS Cluster (%s) to become activating: %s", uuid, err)
	}
	return nil
}

func GetNKSCluster(ctx context.Context, config *conn.ProviderConfig, uuid string) (*vnks.Cluster, error) {

	resp, err := config.Client.Vnks.V2Api.ClustersUuidGet(ctx, &uuid)
	if err != nil {
		return nil, err
	}
	return resp.Cluster, nil
}

func getOIDCSpec(ctx context.Context, config *conn.ProviderConfig, uuid string) (*vnks.OidcRes, error) {

	resp, err := config.Client.Vnks.V2Api.ClustersUuidOidcGet(ctx, &uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getIPAcl(ctx context.Context, config *conn.ProviderConfig, uuid string) (*vnks.IpAclsRes, error) {

	if checkFinSite(config) {
		return &vnks.IpAclsRes{}, nil
	}

	resp, err := config.Client.Vnks.V2Api.ClustersUuidIpAclGet(ctx, &uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getNKSClusterFromList(ctx context.Context, config *conn.ProviderConfig, uuid string) (*vnks.Cluster, error) {
	clusters, err := GetNKSClusters(ctx, config)
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

func GetNKSClusters(ctx context.Context, config *conn.ProviderConfig) ([]*vnks.Cluster, error) {
	resp, err := config.Client.Vnks.V2Api.ClustersGet(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Clusters, nil
}

func getSubnetDiff(oldList interface{}, newList interface{}) (added []*int32, removed []*int32, autoSelect bool) {
	oldMap := make(map[string]int)
	newMap := make(map[string]int)
	autoSelect = true

	for _, v := range ExpandStringInterfaceList(oldList.(([]interface{}))) {
		oldMap[*v] += 1
		autoSelect = false
	}

	for _, v := range ExpandStringInterfaceList(newList.(([]interface{}))) {
		newMap[*v] += 1
	}

	for subnet := range oldMap {
		if _, exist := newMap[subnet]; !exist {
			intV, err := strconv.Atoi(subnet)
			if err == nil {
				removed = append(removed, ncloud.Int32(int32(intV)))
			}
		}
	}

	for subnet := range newMap {
		if _, exist := oldMap[subnet]; !exist {
			intV, err := strconv.Atoi(subnet)
			if err == nil {
				added = append(added, ncloud.Int32(int32(intV)))
			}
		}
	}
	return
}

func checkFinSite(config *conn.ProviderConfig) bool {
	return strings.HasPrefix(config.RegionCode, "F")
}
