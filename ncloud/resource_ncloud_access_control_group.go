package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterResource("ncloud_access_control_group", resourceNcloudAccessControlGroup())
}

func resourceNcloudAccessControlGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAccessControlGroupCreate,
		Read:   resourceNcloudAccessControlGroupRead,
		Update: resourceNcloudAccessControlGroupUpdate,
		Delete: resourceNcloudAccessControlGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: ncloudVpcCommonCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"inbound": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
						},
						"ip_block": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 32),
						},
						"source_access_control_group_no": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"port_range": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validatePortRange,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 1000),
						},
					},
				},
			},
			"outbound": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
						},
						"ip_block": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 32),
						},
						"source_access_control_group_no": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"port_range": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validatePortRange,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 1000),
						},
					},
				},
			},

			"access_control_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNcloudAccessControlGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := createAccessControlGroup(d, config)

	if err != nil {
		return err
	}

	d.SetId(*instance.AccessControlGroupNo)
	log.Printf("[INFO] ACG ID: %s", d.Id())

	return resourceNcloudAccessControlGroupUpdate(d, meta)
}

func resourceNcloudAccessControlGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getAccessControlGroup(config, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.AccessControlGroupNo)
	d.Set("access_control_group_no", instance.AccessControlGroupNo)
	d.Set("name", instance.AccessControlGroupName)
	d.Set("description", instance.AccessControlGroupDescription)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("is_default", instance.IsDefault)

	rules, err := getAccessControlGroupRuleList(d, config)
	if err != nil {
		return err
	}

	var inbound []map[string]interface{}
	var outbound []map[string]interface{}

	for _, r := range rules {
		m := map[string]interface{}{
			"protocol":                       *r.ProtocolType.Code,
			"port_range":                     *r.PortRange,
			"ip_block":                       *r.IpBlock,
			"source_access_control_group_no": *r.AccessControlGroupSequence,
			"description":                    *r.AccessControlGroupRuleDescription,
		}

		if *r.AccessControlGroupRuleType.Code == "INBND" {
			inbound = append(inbound, m)
		} else {
			outbound = append(outbound, m)
		}
	}

	if err := d.Set("inbound", inbound); err != nil {
		log.Printf("[WARN] Error setting inbound rule set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("outbound", outbound); err != nil {
		log.Printf("[WARN] Error setting outbound rule set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudAccessControlGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("inbound") {
		if err := updateAccessControlGroupRule(d, config, "inbound"); err != nil {
			return err
		}
	}

	if d.HasChange("outbound") {
		if err := updateAccessControlGroupRule(d, config, "outbound"); err != nil {
			return err
		}
	}

	return resourceNcloudAccessControlGroupRead(d, meta)
}

func resourceNcloudAccessControlGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if err := deleteAccessControlGroup(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func getAccessControlGroup(config *ProviderConfig, id string) (*vserver.AccessControlGroup, error) {
	if config.SupportVPC {
		return getVpcAccessControlGroup(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_access_control_group`")
}

func getVpcAccessControlGroup(config *ProviderConfig, id string) (*vserver.AccessControlGroup, error) {
	reqParams := &vserver.GetAccessControlGroupDetailRequest{
		RegionCode:           &config.RegionCode,
		AccessControlGroupNo: ncloud.String(id),
	}

	logCommonRequest("getVpcAccessControlGroup", reqParams)
	resp, err := config.Client.vserver.V2Api.GetAccessControlGroupDetail(reqParams)
	if err != nil {
		logErrorResponse("getVpcAccessControlGroup", err, reqParams)
		return nil, err
	}
	logResponse("getVpcAccessControlGroup", resp)

	if len(resp.AccessControlGroupList) > 0 {
		return resp.AccessControlGroupList[0], nil
	}

	return nil, nil
}

func createAccessControlGroup(d *schema.ResourceData, config *ProviderConfig) (*vserver.AccessControlGroup, error) {
	if config.SupportVPC {
		return createVpcAccessControlGroup(d, config)
	}

	return nil, NotSupportClassic("resource `ncloud_access_control_group`")
}

func createVpcAccessControlGroup(d *schema.ResourceData, config *ProviderConfig) (*vserver.AccessControlGroup, error) {
	reqParams := &vserver.CreateAccessControlGroupRequest{
		RegionCode:                    &config.RegionCode,
		VpcNo:                         ncloud.String(d.Get("vpc_no").(string)),
		AccessControlGroupName:        StringPtrOrNil(d.GetOk("name")),
		AccessControlGroupDescription: StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("createVpcAccessControlGroup", reqParams)
	resp, err := config.Client.vserver.V2Api.CreateAccessControlGroup(reqParams)
	if err != nil {
		logErrorResponse("createVpcAccessControlGroup", err, reqParams)
		return nil, err
	}
	logResponse("createVpcAccessControlGroup", resp)

	return resp.AccessControlGroupList[0], nil
}

func deleteAccessControlGroup(config *ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcAccessControlGroup(config, id)
	}

	return NotSupportClassic("resource `ncloud_access_control_group`")
}

func deleteVpcAccessControlGroup(config *ProviderConfig, id string) error {
	accessControlGroup, err := getAccessControlGroup(config, id)
	if err != nil {
		return err
	}

	if accessControlGroup == nil {
		return fmt.Errorf("no matching Access Control Group: %s", id)
	}

	reqParams := &vserver.DeleteAccessControlGroupRequest{
		RegionCode:           &config.RegionCode,
		VpcNo:                accessControlGroup.VpcNo,
		AccessControlGroupNo: ncloud.String(id),
	}

	logCommonRequest("deleteVpcAccessControlGroup", reqParams)
	resp, err := config.Client.vserver.V2Api.DeleteAccessControlGroup(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcAccessControlGroup", err, reqParams)
		return err
	}
	logResponse("deleteVpcAccessControlGroup", resp)

	if err := waitForVpcAccessControlGroupDeletion(config, id); err != nil {
		return err
	}

	return nil
}

func waitForVpcAccessControlGroupDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getAccessControlGroup(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "AccessControlGroupStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Access Control Group (%s) to become terminated: %s", id, err)
	}

	return nil
}

func waitForVpcAccessControlGroupRunning(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getAccessControlGroup(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "AccessControlGroupStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Access Control Group (%s) to become running: %s", id, err)
	}

	return nil
}

func updateAccessControlGroupRule(d *schema.ResourceData, meta interface{}, ruleType string) error {
	config := meta.(*ProviderConfig)

	o, n := d.GetChange(ruleType)

	if o == nil {
		o = new(schema.Set)
	}
	if n == nil {
		n = new(schema.Set)
	}

	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	add := ns.Difference(os).List()
	remove := os.Difference(ns).List()

	removeAccessControlGroupRuleList := expandRemoveAccessControlGroupRule(remove)
	addAccessControlGroupRuleList, err := expandAddAccessControlGroupRule(add)
	if err != nil {
		return err
	}

	if len(removeAccessControlGroupRuleList) > 0 {
		if err := removeAccessControlGroupRule(d, config, ruleType, removeAccessControlGroupRuleList); err != nil {
			return err
		}
	}

	if len(addAccessControlGroupRuleList) > 0 {
		if err := addAccessControlGroupRule(d, config, ruleType, addAccessControlGroupRuleList); err != nil {
			return err
		}
	}

	return nil
}

func addAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, ruleType string, accessControlGroupRule []*vserver.AddAccessControlGroupRuleParameter) error {
	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		var reqParams interface{}
		if ruleType == "inbound" {
			reqParams = &vserver.AddAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      ncloud.String(d.Get("vpc_no").(string)),
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("AddAccessControlGroupInboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.AddAccessControlGroupInboundRule(reqParams.(*vserver.AddAccessControlGroupInboundRuleRequest))
		} else {
			reqParams = &vserver.AddAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      ncloud.String(d.Get("vpc_no").(string)),
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("AddAccessControlGroupOutboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.AddAccessControlGroupOutboundRule(reqParams.(*vserver.AddAccessControlGroupOutboundRuleRequest))
		}

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == ApiErrorAcgCantChangeSameTime {
			logErrorResponse("retry AddAccessControlGroupRule", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("AddAccessControlGroupRule", err, reqParams)
		return err
	}

	logResponse("AddAccessControlGroupRule", resp)

	if err = waitForVpcAccessControlGroupRunning(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func removeAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, ruleType string, accessControlGroupRule []*vserver.RemoveAccessControlGroupRuleParameter) error {
	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		var reqParams interface{}
		if ruleType == "inbound" {
			reqParams = &vserver.RemoveAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      ncloud.String(d.Get("vpc_no").(string)),
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("RemoveAccessControlGroupInboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.RemoveAccessControlGroupInboundRule(reqParams.(*vserver.RemoveAccessControlGroupInboundRuleRequest))
		} else {
			reqParams = &vserver.RemoveAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      ncloud.String(d.Get("vpc_no").(string)),
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("RemoveAccessControlGroupOutboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.RemoveAccessControlGroupOutboundRule(reqParams.(*vserver.RemoveAccessControlGroupOutboundRuleRequest))
		}

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == ApiErrorAcgCantChangeSameTime {
			logErrorResponse("retry RemoveAccessControlGroupRule", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("RemoveAccessControlGroupRule", err, reqParams)
		return err
	}

	logResponse("RemoveAccessControlGroupRule", resp)

	if err = waitForVpcAccessControlGroupRunning(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandAddAccessControlGroupRule(rules []interface{}) ([]*vserver.AddAccessControlGroupRuleParameter, error) {
	var acgRuleList []*vserver.AddAccessControlGroupRuleParameter

	for _, vi := range rules {
		m := vi.(map[string]interface{})

		if len(m["ip_block"].(string)) == 0 && len(m["source_access_control_group_no"].(string)) == 0 {
			return nil, fmt.Errorf("one of either `ip_block` or `source_access_control_group_no` is required")
		}

		if len(m["ip_block"].(string)) > 0 && len(m["source_access_control_group_no"].(string)) > 0 {
			return nil, fmt.Errorf("cannot be specified with `ip_block` and `source_access_control_group_no`")
		}

		acgRule := &vserver.AddAccessControlGroupRuleParameter{
			AccessControlGroupRuleDescription: ncloud.String(m["description"].(string)),
			IpBlock:                           ncloud.String(m["ip_block"].(string)),
			AccessControlGroupSequence:        ncloud.String(m["source_access_control_group_no"].(string)),
			ProtocolTypeCode:                  ncloud.String(m["protocol"].(string)),
			PortRange:                         ncloud.String(m["port_range"].(string)),
		}

		acgRuleList = append(acgRuleList, acgRule)
	}

	return acgRuleList, nil
}

func expandRemoveAccessControlGroupRule(rules []interface{}) []*vserver.RemoveAccessControlGroupRuleParameter {
	var acgRuleList []*vserver.RemoveAccessControlGroupRuleParameter

	for _, vi := range rules {
		m := vi.(map[string]interface{})

		acgRule := &vserver.RemoveAccessControlGroupRuleParameter{
			IpBlock:                    ncloud.String(m["ip_block"].(string)),
			AccessControlGroupSequence: ncloud.String(m["source_access_control_group_no"].(string)),
			ProtocolTypeCode:           ncloud.String(m["protocol"].(string)),
			PortRange:                  ncloud.String(m["port_range"].(string)),
		}

		acgRuleList = append(acgRuleList, acgRule)
	}

	return acgRuleList
}

func getAccessControlGroupRuleList(d *schema.ResourceData, config *ProviderConfig) ([]*vserver.AccessControlGroupRule, error) {
	reqParams := &vserver.GetAccessControlGroupRuleListRequest{
		RegionCode:           &config.RegionCode,
		AccessControlGroupNo: ncloud.String(d.Id()),
	}

	logCommonRequest("getAccessControlGroupRuleList", reqParams)
	resp, err := config.Client.vserver.V2Api.GetAccessControlGroupRuleList(reqParams)
	if err != nil {
		logErrorResponse("getAccessControlGroupRuleList", err, reqParams)
		return nil, err
	}
	logResponse("getAccessControlGroupRuleList", resp)

	return resp.AccessControlGroupRuleList, nil
}
