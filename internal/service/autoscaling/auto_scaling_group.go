package autoscaling

import (
	"fmt"

	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func ResourceNcloudAutoScalingGroup() *schema.Resource {
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
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 255)),
				ForceNew:         true,
			},
			"launch_configuration_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desired_capacity": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 30)),
			},
			"min_size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 30)),
			},
			"max_size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 30)),
			},
			"default_cooldown": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          300,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 2147483647)),
			},
			"health_check_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SVR", "LOADB"}, false)),
			},
			"health_check_grace_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          300,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 2147483647)),
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
				ValidateDiagFunc: validation.ToDiagFunc(ValidateParseDuration),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_control_group_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				},
				ForceNew: true,
			},
			"target_group_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"server_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ignore_capacity_changes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNcloudAutoScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

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

func createAutoScalingGroup(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if _, ok := d.GetOk("subnet_no"); !ok {
		return nil, ErrorRequiredArgOnVpc("subnet_no")
	}

	if _, ok := d.GetOk("access_control_group_no_list"); !ok {
		return nil, ErrorRequiredArgOnVpc("access_control_group_no_list")
	}

	subnetNo := d.Get("subnet_no").(string)
	subnet, err := vpc.GetSubnetInstance(config, subnetNo)
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
		AccessControlGroupNoList: ExpandStringInterfaceList(d.Get("access_control_group_no_list").([]interface{})),
		ServerNamePrefix:         StringPtrOrNil(d.GetOk("server_name_prefix")),
		MinSize:                  ncloud.Int32(int32(d.Get("min_size").(int))),
		MaxSize:                  ncloud.Int32(int32(d.Get("max_size").(int))),
		DesiredCapacity:          Int32PtrOrNil(d.GetOk("desired_capacity")),
		DefaultCoolDown:          ncloud.Int32(int32(d.Get("default_cooldown").(int))),
		HealthCheckGracePeriod:   ncloud.Int32(int32(d.Get("health_check_grace_period").(int))),
		HealthCheckTypeCode:      StringPtrOrNil(d.GetOk("health_check_type_code")),
		TargetGroupNoList:        StringListPtrOrNil(d.GetOk("target_group_list")),
	}

	resp, err := config.Client.Vautoscaling.V2Api.CreateAutoScalingGroup(reqParams)
	if err != nil {
		return nil, err
	}

	return resp.AutoScalingGroupList[0].AutoScalingGroupNo, nil
}

func resourceNcloudAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	autoScalingGroup, err := GetAutoScalingGroup(config, d.Id())
	if err != nil {
		return err
	}

	if autoScalingGroup == nil {
		d.SetId("")
		return nil
	}

	max_size := d.Get("max_size")
	min_size := d.Get("min_size")
	desired_capacity := d.Get("desired_capacity")

	autoScalingGroupMap := ConvertToMap(autoScalingGroup)
	SetSingularResourceDataFromMapSchema(ResourceNcloudAutoScalingGroup(), d, autoScalingGroupMap)

	if d.Get("ignore_capacity_changes").(bool) {
		if err := d.Set("max_size", max_size); err != nil {
			return err
		}
		if err := d.Set("min_size", min_size); err != nil {
			return err
		}
		if err := d.Set("desired_capacity", desired_capacity); err != nil {
			return err
		}
	}

	if err := d.Set("server_instance_no_list", autoScalingGroup.InAutoScalingGroupServerInstanceList); err != nil {
		return err
	}

	return nil
}

func GetAutoScalingGroup(config *conn.ProviderConfig, id string) (*AutoScalingGroup, error) {
	reqParams := &vautoscaling.GetAutoScalingGroupListRequest{
		RegionCode:             &config.RegionCode,
		AutoScalingGroupNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("getVpcAutoScalingGroup", reqParams)
	resp, err := config.Client.Vautoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcAutoScalingGroup", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcAutoScalingGroup", resp)

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

func resourceNcloudAutoScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if err := updateAutoScalingGroup(d, config); err != nil {
		return err
	}

	if err := waitForAutoScalingGroupCapacity(d, config); err != nil {
		return err
	}

	return resourceNcloudAutoScalingGroupRead(d, config)
}

func updateAutoScalingGroup(d *schema.ResourceData, config *conn.ProviderConfig) error {
	asg, err := GetAutoScalingGroup(config, d.Id())
	if err != nil {
		return nil
	}

	reqParams := &vautoscaling.UpdateAutoScalingGroupRequest{
		AutoScalingGroupNo: asg.AutoScalingGroupNo,
	}

	if d.HasChange("launch_configuration_no") {
		reqParams.LaunchConfigurationNo = ncloud.String(d.Get("launch_configuration_no").(string))
	}

	if !d.Get("ignore_capacity_changes").(bool) {
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

	LogCommonRequest("changeVpcAutoScalingGroup", reqParams)
	resp, err := config.Client.Vautoscaling.V2Api.UpdateAutoScalingGroup(reqParams)
	LogResponse("changeVpcAutoScalingGroup", resp)
	if err != nil {
		return err
	}

	return nil
}

func resourceNcloudAutoScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if err := deleteAutoScalingGroup(d, config); err != nil {
		return err
	}
	return nil
}

func deleteAutoScalingGroup(d *schema.ResourceData, config *conn.ProviderConfig) error {
	d.Timeout(schema.TimeoutDelete)

	asg, err := GetAutoScalingGroup(config, d.Id())
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

	if _, err := config.Client.Vautoscaling.V2Api.UpdateAutoScalingGroup(cReqParams); err != nil {
		return err
	}

	if err := waitForVpcInAutoScalingGroupServerInstanceListDeletion(config, d.Id()); err != nil {
		return err
	}

	if err := waitForVpcAutoScalingGroupDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func getVpcInAutoScalingGroupServerInstanceList(config *conn.ProviderConfig, id string) ([]*InAutoScalingGroupServerInstance, error) {
	reqParams := &vautoscaling.GetAutoScalingGroupListRequest{
		RegionCode:             &config.RegionCode,
		AutoScalingGroupNoList: []*string{ncloud.String(id)},
	}

	resp, err := config.Client.Vautoscaling.V2Api.GetAutoScalingGroupList(reqParams)
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

func waitForVpcInAutoScalingGroupServerInstanceListDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSVC"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			asg, err := GetAutoScalingGroup(config, id)
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
		Timeout:    conn.DefaultStopTimeout * 3,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for InAutoScalingGroupServerInstanceList (%s) to become deleting: %s", id, err)
	}
	return nil
}

func waitForVpcAutoScalingGroupDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"DELETE"},
		Refresh: func() (interface{}, string, error) {
			client := config.Client
			reqParams := &vautoscaling.DeleteAutoScalingGroupRequest{
				AutoScalingGroupNo: ncloud.String(id),
			}
			resp, err := client.Vautoscaling.V2Api.DeleteAutoScalingGroup(reqParams)
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
		Timeout:    conn.DefaultTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for AutoScalingGroup (%s) to become deleting: %s", id, err)
	}
	return nil
}

func waitForAutoScalingGroupCapacity(d *schema.ResourceData, config *conn.ProviderConfig) error {
	wait, err := time.ParseDuration(d.Get("wait_for_capacity_timeout").(string))
	if err != nil {
		return err
	}

	if wait == 0 {
		return nil
	}

	return waitForVpcAutoScalingGroupCapacity(d, config, wait)
}

func waitForVpcAutoScalingGroupCapacity(d *schema.ResourceData, config *conn.ProviderConfig, wait time.Duration) error {
	return resource.Retry(wait, func() *resource.RetryError {
		asg, err := GetAutoScalingGroup(config, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}

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
