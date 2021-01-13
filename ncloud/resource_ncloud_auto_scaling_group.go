package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"
)

func init() {
	RegisterResource("ncloud_auto_scaling_group", resourceNcloudAutoScalingGroup())
}

func resourceNcloudAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAutoScalingGroupCreate,
		Read:   resourceNcloudAutoScalingGroupRead,
		Update: resourceNcloudAutoScalingGroupUpdate,
		Delete: resourceNcloudAutoScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 255)),
			},
			"launch_configuration_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desired_capacity": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 30)),
			},
			"min_size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 30)),
			},
			"max_size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 30)),
			},
			"default_cooldown": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 2147483647)),
			},
			// TODO: healthCheckTypeCode 값이 LOADB(로드밸런서) 경우 필수
			"health_check_grace_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 2147483647)),
			},
			"health_check_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"SVR", "LOADB"}, false)),
			},
			// required Only Classic
			"zone_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// required Only VPC
			"subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// required Only VPC
			"access_control_group_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// TODO: healthCheckTypeCode 가 LOADB 인 경우에만 유효
			"target_group_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"server_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_instance_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"wait_for_capacity_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "10m",
				ValidateDiagFunc: ToDiagFunc(validateParseDuration),
			},
		},
	}
}

func resourceNcloudAutoScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	id, err := createAutoScalingGroup(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	if err := waitForAutoScalingGroupCapacity(d, config); err != nil {
		return err
	}

	return resourceNcloudAutoScalingGroupRead(d, meta)
}

func createAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcAutoScalingGroup(d, config)
	} else {
		return createClassicAutoScalingGroup(d, config)
	}
}

func createVpcAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) (*string, error) {

	subnetNo := d.Get("subnet_no").(string)
	subnet, err := getSubnetInstance(config, subnetNo)
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		return nil, fmt.Errorf("Not fount subnet(%s)", subnetNo)
	}

	reqParams := &vautoscaling.CreateAutoScalingGroupRequest{
		RegionCode:               &config.RegionCode,
		LaunchConfigurationNo:    ncloud.String(d.Get("launch_configuration_no").(string)),
		AutoScalingGroupName:     StringPtrOrNil(d.GetOk("name")),
		VpcNo:                    subnet.VpcNo,
		SubnetNo:                 subnet.SubnetNo,
		AccessControlGroupNoList: expandStringInterfaceList(d.Get("access_control_group_no_list").([]interface{})),
		ServerNamePrefix:         StringPtrOrNil(d.GetOk("server_name_prefix")),
		MinSize:                  ncloud.Int32(int32(d.Get("min_size").(int))),
		MaxSize:                  ncloud.Int32(int32(d.Get("max_size").(int))),
		DesiredCapacity:          Int32PtrOrNil(d.GetOk("desired_capacity")),
		DefaultCoolDown:          Int32PtrOrNil(d.GetOk("default_cooldown")),
		HealthCheckGracePeriod:   Int32PtrOrNil(d.GetOk("health_check_grace_period")),
		HealthCheckTypeCode:      StringPtrOrNil(d.GetOk("health_check_type_code")),
		TargetGroupNoList:        StringListPtrOrNil(d.GetOk("target_group_list")),
	}

	resp, err := config.Client.vautoscaling.V2Api.CreateAutoScalingGroup(reqParams)
	if err != nil {
		return nil, err
	}

	return resp.AutoScalingGroupList[0].AutoScalingGroupNo, nil
}

func createClassicAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	// TODO : Zero value 핸들링
	l, err := getClassicLaunchConfigurationNameByNo(StringPtrOrNil(d.GetOk("launch_configuration_no")), config)
	if err != nil {
		return nil, err
	}
	reqParams := &autoscaling.CreateAutoScalingGroupRequest{
		AutoScalingGroupName:    StringPtrOrNil(d.GetOk("name")),
		LaunchConfigurationName: l.LaunchConfigurationName,
		DesiredCapacity:         Int32PtrOrNil(d.GetOk("desired_capacity")),
		MinSize:                 ncloud.Int32(int32(d.Get("min_size").(int))),
		MaxSize:                 ncloud.Int32(int32(d.Get("max_size").(int))),
		DefaultCooldown:         Int32PtrOrNil(d.GetOk("default_cooldown")),
		//LoadBalancerNameList:
		HealthCheckGracePeriod: Int32PtrOrNil(d.GetOk("health_check_grace_period")),
		HealthCheckTypeCode:    StringPtrOrNil(d.GetOk("health_check_type_code")),
		ZoneNoList:             expandStringInterfaceList(d.Get("zone_no_list").([]interface{})),
	}
	logCommonRequest("createClassicAutoScalingGroup", reqParams)
	resp, err := config.Client.autoscaling.V2Api.CreateAutoScalingGroup(reqParams)
	if err != nil {
		logErrorResponse("createClassicAutoScalingGroup", err, reqParams)
		return nil, err
	}
	logResponse("createClassicAutoScalingGroup", resp)
	return resp.AutoScalingGroupList[0].AutoScalingGroupNo, nil
}

func resourceNcloudAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	autoScalingGroup, err := getAutoScalingGroup(config, d.Id())
	if err != nil {
		return err
	}

	if autoScalingGroup == nil {
		d.SetId("")
		return nil
	}

	autoScalingGroupMap := ConvertToMap(autoScalingGroup)
	SetSingularResourceDataFromMapSchema(resourceNcloudAutoScalingGroup(), d, autoScalingGroupMap)
	return nil
}

func getAutoScalingGroup(config *ProviderConfig, id string) (*AutoScalingGroup, error) {
	if config.SupportVPC {
		return getVpcAutoScalingGroup(config, id)
	} else {
		return getClassicAutoScalingGroup(config, id)
	}
}

func getVpcAutoScalingGroup(config *ProviderConfig, id string) (*AutoScalingGroup, error) {
	reqParams := &vautoscaling.GetAutoScalingGroupListRequest{
		RegionCode:             &config.RegionCode,
		AutoScalingGroupNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("getVpcAutoScalingGroup", reqParams)
	resp, err := config.Client.vautoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		logErrorResponse("getVpcAutoScalingGroup", err, reqParams)
		return nil, err
	}
	logResponse("getVpcAutoScalingGroup", resp)

	if len(resp.AutoScalingGroupList) < 1 {
		return nil, nil
	}

	asg := resp.AutoScalingGroupList[0]

	return &AutoScalingGroup{
		AutoScalingGroupNo:                   asg.AutoScalingGroupNo,
		AutoScalingGroupName:                 asg.AutoScalingGroupName,
		LaunchConfigurationNo:                asg.LaunchConfigurationNo,
		DesiredCapacity:                      asg.DesiredCapacity,
		MinSize:                              asg.MinSize,
		MaxSize:                              asg.MaxSize,
		DefaultCooldown:                      asg.DefaultCoolDown,
		HealthCheckGracePeriod:               asg.HealthCheckGracePeriod,
		HealthCheckTypeCode:                  asg.HealthCheckType.Code,
		InAutoScalingGroupServerInstanceList: flattenVpcAutoScalingGroupServerInstanceList(asg.InAutoScalingGroupServerInstanceList),
		SuspendedProcessList:                 flattenVpcSuspendedProcessList(asg.SuspendedProcessList),
		VpcNo:                                asg.VpcNo,
		SubnetNo:                             asg.SubnetNo,
		ServerNamePrefix:                     asg.ServerNamePrefix,
		TargetGroupNoList:                    asg.TargetGroupNoList,
		AccessControlGroupNoList:             asg.AccessControlGroupNoList,
	}, nil
}

func getClassicAutoScalingGroup(config *ProviderConfig, id string) (*AutoScalingGroup, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetAutoScalingGroupListRequest{
		RegionNo: &config.RegionNo,
	}

	logCommonRequest("getClassicAutoScalingGroup", reqParams)
	resp, err := config.Client.autoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		logErrorResponse("getClassicAutoScalingGroup", err, reqParams)
		return nil, err
	}
	logResponse("getClassicAutoScalingGroup", resp)

	for _, a := range resp.AutoScalingGroupList {
		if *a.AutoScalingGroupNo == *no {
			return &AutoScalingGroup{
				AutoScalingGroupNo:                   a.AutoScalingGroupNo,
				AutoScalingGroupName:                 a.AutoScalingGroupName,
				LaunchConfigurationNo:                a.LaunchConfigurationNo,
				DesiredCapacity:                      a.DesiredCapacity,
				MinSize:                              a.MinSize,
				MaxSize:                              a.MaxSize,
				DefaultCooldown:                      a.DefaultCooldown,
				LoadBalancerInstanceSummaryList:      flattenLoadBalancerInstanceSummaryList(a.LoadBalancerInstanceSummaryList),
				HealthCheckGracePeriod:               a.HealthCheckGracePeriod,
				HealthCheckTypeCode:                  a.HealthCheckType.Code,
				InAutoScalingGroupServerInstanceList: flattenClassicAutoScalingGroupServerInstanceList(a.InAutoScalingGroupServerInstanceList),
				SuspendedProcessList:                 flattenClassicSuspendedProcessList(a.SuspendedProcessList),
				ZoneList:                             flattenZoneList(a.ZoneList),
			}, nil
		}
	}

	return nil, nil
}

func resourceNcloudAutoScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if err := updateAutoScalingGroup(d, config); err != nil {
		return err
	}

	return resourceNcloudAutoScalingGroupRead(d, config)
}

func updateAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) error {
	if config.SupportVPC {
		return changeVpcAutoScalingGroup(d, config)
	} else {
		return changeClassicAutoScalingGroup(d, config)
	}
}

func changeVpcAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) error {
	asg, err := getAutoScalingGroup(config, d.Id())
	if err != nil {
		return nil
	}

	// TODO: LaunchConfiguration
	reqParams := &vautoscaling.UpdateAutoScalingGroupRequest{
		AutoScalingGroupNo: asg.AutoScalingGroupNo,
	}

	if d.HasChange("desired_capacity") {
		reqParams.DesiredCapacity = Int32PtrOrNil(d.GetOk("desired_capacity"))
	}

	if d.HasChange("min_size") || d.HasChange("max_size") {
		min := ncloud.Int32(int32(d.Get("min_size").(int)))
		max := ncloud.Int32(int32(d.Get("max_size").(int)))
		if *min > *max {
			return fmt.Errorf("min_size is must be at least 0 and less than or equal to max_size")
		}
		reqParams.MinSize = min
		reqParams.MaxSize = max
	}

	if d.HasChange("default_cooldown") {
		reqParams.DefaultCoolDown = Int32PtrOrNil(d.GetOk("default_cooldown"))
	}

	if d.HasChange("health_check_grace_period") {
		reqParams.HealthCheckGracePeriod = Int32PtrOrNil(d.GetOk("health_check_grace_period"))
	}

	if d.HasChange("health_check_type_code") {
		reqParams.HealthCheckTypeCode = StringPtrOrNil(d.GetOk("health_check_type_code"))
	}

	if d.HasChange("server_name_prefix") {
		reqParams.ServerNamePrefix = StringPtrOrNil(d.GetOk("server_name_prefix"))
	}

	logCommonRequest("changeVpcAutoScalingGroup", reqParams)
	resp, err := config.Client.vautoscaling.V2Api.UpdateAutoScalingGroup(reqParams)
	logResponse("changeVpcAutoScalingGroup", resp)
	if err != nil {
		return err
	}

	return nil
}

func changeClassicAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) error {
	asg, err := getAutoScalingGroup(config, d.Id())
	if err != nil {
		return err
	}
	reqParams := &autoscaling.UpdateAutoScalingGroupRequest{
		AutoScalingGroupName: asg.AutoScalingGroupName,
	}
	if d.HasChange("desired_capacity") {
		reqParams.DesiredCapacity = Int32PtrOrNil(d.GetOk("desired_capacity"))
	}

	if d.HasChange("min_size") || d.HasChange("max_size") {
		min := ncloud.Int32(int32(d.Get("min_size").(int)))
		max := ncloud.Int32(int32(d.Get("max_size").(int)))
		if *min > *max {
			return fmt.Errorf("min_size is must be at least 0 and less than or equal to max_size")
		}
		reqParams.MinSize = min
		reqParams.MaxSize = max
	}

	if d.HasChange("default_cooldown") {
		reqParams.DefaultCooldown = Int32PtrOrNil(d.GetOk("default_cooldown"))
	}

	if d.HasChange("health_check_grace_period") {
		reqParams.HealthCheckGracePeriod = Int32PtrOrNil(d.GetOk("health_check_grace_period"))
	}

	if d.HasChange("health_check_type_code") {
		reqParams.HealthCheckTypeCode = StringPtrOrNil(d.GetOk("health_check_type_code"))
	}

	if d.HasChange("zone_no_list") {
		reqParams.ZoneNoList = expandStringInterfaceList(d.Get("zone_no_list").([]interface{}))
	}

	logCommonRequest("changeClassicAutoScalingGroup", reqParams)
	resp, err := config.Client.autoscaling.V2Api.UpdateAutoScalingGroup(reqParams)
	logResponse("changeClassicAutoScalingGroup", resp)
	if err != nil {
		return err
	}
	return nil
}

func resourceNcloudAutoScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if err := deleteAutoScalingGroup(d, config); err != nil {
		return err
	}
	return nil
}

func deleteAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) error {
	d.Timeout(schema.TimeoutDelete)
	if config.SupportVPC {
		return deleteVpcAutoScalingGroup(config, d.Id())
	} else {
		return deleteClassicAutoScalingGroup(config, d.Id())
	}
}

func deleteVpcAutoScalingGroup(config *ProviderConfig, id string) error {
	asg, err := getAutoScalingGroup(config, id)
	if err != nil {
		return err
	}

	// 1. Set max_size, min_size, desired_capacity to 0
	cReqParams := &vautoscaling.UpdateAutoScalingGroupRequest{
		AutoScalingGroupNo: asg.AutoScalingGroupNo,
		DesiredCapacity:    ncloud.Int32(0),
		MinSize:            ncloud.Int32(0),
		MaxSize:            ncloud.Int32(0),
	}

	if _, err := config.Client.vautoscaling.V2Api.UpdateAutoScalingGroup(cReqParams); err != nil {
		return err
	}

	if err := waitForVpcInAutoScalingGroupServerInstanceListDeletion(config, id); err != nil {
		return err
	}

	if err := waitForVpcAutoScalingGroupDeletion(config, id); err != nil {
		return err
	}

	return nil
}

func deleteClassicAutoScalingGroup(config *ProviderConfig, id string) error {
	asg, err := getAutoScalingGroup(config, id)
	if err != nil {
		return err
	}

	// 1. Set max_size, min_size, desired_capacity to 0
	cReqParams := &autoscaling.UpdateAutoScalingGroupRequest{
		AutoScalingGroupName: asg.AutoScalingGroupName,
		DesiredCapacity:      ncloud.Int32(0),
		MinSize:              ncloud.Int32(0),
		MaxSize:              ncloud.Int32(0),
	}

	if _, err := config.Client.autoscaling.V2Api.UpdateAutoScalingGroup(cReqParams); err != nil {
		return err
	}

	// 2. Delete Server Instance List in AutoScalingGroup
	if err := waitForClassicInAutoScalingGroupServerInstanceListDeletion(config, id); err != nil {
		return err
	}

	// 3. Delete Auto Scaling Group
	if err := waitForClassicAutoScalingGroupDeletion(config, ncloud.StringValue(asg.AutoScalingGroupName)); err != nil {
		return err
	}

	return nil
}

func getVpcInAutoScalingGroupServerInstanceList(config *ProviderConfig, id string) ([]*InAutoScalingGroupServerInstance, error) {
	reqParams := &vautoscaling.GetAutoScalingGroupListRequest{
		RegionCode:             &config.RegionCode,
		AutoScalingGroupNoList: []*string{ncloud.String(id)},
	}

	resp, err := config.Client.vautoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		return nil, err
	}

	asg := resp.AutoScalingGroupList[0]
	list := make([]*InAutoScalingGroupServerInstance, 0)
	for _, i := range asg.InAutoScalingGroupServerInstanceList {
		list = append(list, &InAutoScalingGroupServerInstance{
			HealthStatus:     i.HealthStatus.Code,
			LifecycleState:   i.LifecycleState.Code,
			ServerInstanceNo: i.ServerInstanceNo,
		})
	}
	return list, nil
}

func getClassicInAutoScalingGroupServerInstanceList(config *ProviderConfig, id string) ([]*InAutoScalingGroupServerInstance, error) {
	tmpAsg, err := getClassicAutoScalingGroup(config, id)
	if err != nil {
		return nil, err
	}

	reqParams := &autoscaling.GetAutoScalingGroupListRequest{
		AutoScalingGroupNameList: []*string{tmpAsg.AutoScalingGroupName},
		RegionNo:                 &config.RegionNo,
	}

	resp, err := config.Client.autoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		return nil, err
	}

	asg := resp.AutoScalingGroupList[0]
	list := make([]*InAutoScalingGroupServerInstance, 0)
	for _, i := range asg.InAutoScalingGroupServerInstanceList {
		list = append(list, &InAutoScalingGroupServerInstance{
			HealthStatus:     i.HealthStatus.Code,
			LifecycleState:   i.LifecycleState.Code,
			ServerInstanceNo: i.ServerInstanceNo,
		})
	}
	return list, nil
}

func waitForClassicInAutoScalingGroupServerInstanceListDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSVC"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			asg, err := getAutoScalingGroup(config, id)
			if err != nil {
				return 0, "", err
			}
			if len(asg.InAutoScalingGroupServerInstanceList) > 0 {
				return asg, "INSVC", nil
			} else {
				return asg, "TERMT", nil
			}
		},
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
		Timeout:    DefaultStopTimeout * 3,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for InAutoScalingGroupServerInstanceList (%s) to become deleting: %s", id, err)
	}
	return nil
}

func waitForVpcInAutoScalingGroupServerInstanceListDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSVC"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			asg, err := getAutoScalingGroup(config, id)
			if err != nil {
				return 0, "", err
			}
			if len(asg.InAutoScalingGroupServerInstanceList) > 0 {
				return asg, "INSVC", nil
			} else {
				return asg, "TERMT", nil
			}
		},
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
		Timeout:    DefaultStopTimeout * 3,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for InAutoScalingGroupServerInstanceList (%s) to become deleting: %s", id, err)
	}
	return nil
}

func waitForClassicAutoScalingGroupDeletion(config *ProviderConfig, name string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"DELETE"},
		Refresh: func() (interface{}, string, error) {
			client := config.Client
			reqParams := &autoscaling.DeleteAutoScalingGroupRequest{
				AutoScalingGroupName: ncloud.String(name),
			}
			resp, err := client.autoscaling.V2Api.DeleteAutoScalingGroup(reqParams)
			if err != nil {
				errBody, _ := GetCommonErrorBody(err)
				if errBody.ReturnCode == ApiErrorASGScalingIsActive || errBody.ReturnCode == ApiErrorASGIsUsingPolicyOrLaunchConfiguration {
					return resp, "RUN", nil
				} else {
					return 0, "", err
				}
			} else {
				return resp, "DELETE", nil
			}
		},
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
		Timeout:    DefaultTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for AutoScalingGroup (%s) to become deleting: %s", name, err)
	}
	return nil
}

func waitForVpcAutoScalingGroupDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"DELETE"},
		Refresh: func() (interface{}, string, error) {
			client := config.Client
			reqParams := &vautoscaling.DeleteAutoScalingGroupRequest{
				AutoScalingGroupNo: ncloud.String(id),
			}
			resp, err := client.vautoscaling.V2Api.DeleteAutoScalingGroup(reqParams)
			if err != nil {
				errBody, _ := GetCommonErrorBody(err)
				if errBody.ReturnCode == ApiErrorASGIsUsingPolicyOrLaunchConfigurationOnVpc {
					return resp, "RUN", nil
				} else {
					return 0, "", err
				}
			} else {
				return resp, "DELETE", nil
			}
		},
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
		Timeout:    DefaultTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for AutoScalingGroup (%s) to become deleting: %s", id, err)
	}
	return nil
}

func waitForAutoScalingGroupCapacity(d *schema.ResourceData, config *ProviderConfig) error {
	wait, err := time.ParseDuration(d.Get("wait_for_capacity_timeout").(string))
	if err != nil {
		return err
	}

	if wait == 0 {
		return nil
	}

	if config.SupportVPC {
		return waitForVpcAutoScalingGroupCapacity(d, config, wait)
	} else {
		return waitForClassicAutoScalingGroupCapacity(d, config, wait)
	}
}

func waitForVpcAutoScalingGroupCapacity(d *schema.ResourceData, config *ProviderConfig, wait time.Duration) error {
	return resource.Retry(wait, func() *resource.RetryError {
		asg, err := getVpcAutoScalingGroup(config, d.Id())
		asgServerInstanceList, err := getVpcInAutoScalingGroupServerInstanceList(config, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}

		var currentServerInstanceCnt int32
		for _, i := range asgServerInstanceList {
			if !strings.EqualFold(*i.HealthStatus, "HLTHY") {
				continue
			}

			if !strings.EqualFold(*i.LifecycleState, "INSVC") {
				continue
			}

			currentServerInstanceCnt++
		}

		minASG := asg.MinSize
		if asg.DesiredCapacity != nil {
			minASG = asg.DesiredCapacity
		}

		if currentServerInstanceCnt < *minASG {
			return resource.RetryableError(fmt.Errorf("Wait for the server instances in the AutoScalingGroup(%s) to be serviced. : Need at least %d healthy instances in ASG, have %d", d.Id(), *minASG, currentServerInstanceCnt))
		}
		return nil
	})
}

func waitForClassicAutoScalingGroupCapacity(d *schema.ResourceData, config *ProviderConfig, wait time.Duration) error {
	return resource.Retry(wait, func() *resource.RetryError {
		asg, err := getClassicAutoScalingGroup(config, d.Id())
		asgServerInstanceList, err := getClassicInAutoScalingGroupServerInstanceList(config, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}

		var currentServerInstanceCnt int32
		for _, i := range asgServerInstanceList {
			if !strings.EqualFold(*i.HealthStatus, "HLTHY") {
				continue
			}

			if !strings.EqualFold(*i.LifecycleState, "INSVC") {
				continue
			}

			currentServerInstanceCnt++
		}

		minASG := asg.MinSize
		if asg.DesiredCapacity != nil {
			minASG = asg.DesiredCapacity
		}

		if currentServerInstanceCnt < *minASG {
			return resource.RetryableError(fmt.Errorf("%q: Waiting up to %s: Need at least %d healthy instances in ASG, have %d", d.Id(), DefaultCreateTimeout, *minASG, currentServerInstanceCnt))
		}
		return nil
	})
}
