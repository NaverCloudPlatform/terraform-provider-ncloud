package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_auto_scaling_policy", resourceNcloudAutoScalingPolicy())
}

func resourceNcloudAutoScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAutoScalingPolicyCreate,
		Read:   resourceNcloudAutoScalingPolicyRead,
		Update: resourceNcloudAutoScalingPolicyUpdate,
		Delete: resourceNcloudAutoScalingPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"adjustment_type_code": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"CHANG", "EXACT", "PRCNT"}, false)),
			},
			"scaling_adjustment": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cooldown": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_adjustment_step": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_scaling_group_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudAutoScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	autoscaling_group_no, id, err := createAutoScalingPolicy(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	d.Set("auto_scaling_group_no", autoscaling_group_no)
	return resourceNcloudAutoScalingPolicyRead(d, meta)
}

func createAutoScalingPolicy(d *schema.ResourceData, config *ProviderConfig) (*string, *string, error) {
	if config.SupportVPC {
		return createVpcAutoScalingPolicy(d, config)
	} else {
		return createClassicAutoScalingPolicy(d, config)
	}
}

func createVpcAutoScalingPolicy(d *schema.ResourceData, config *ProviderConfig) (*string, *string, error) {
	reqParams := &vautoscaling.PutScalingPolicyRequest{
		RegionCode: &config.RegionCode,
		// Required
		AdjustmentTypeCode: ncloud.String(d.Get("adjustment_type_code").(string)),
		ScalingAdjustment:  ncloud.Int32(int32(d.Get("scaling_adjustment").(int))),
		AutoScalingGroupNo: ncloud.String(d.Get("auto_scaling_group_no").(string)),
		PolicyName:         ncloud.String(d.Get("name").(string)),
		// Optional
		MinAdjustmentStep: Int32PtrOrNil(d.GetOk("min_adjustment_step")),
		CoolDown:          Int32PtrOrNil(d.GetOk("cooldown")),
	}
	resp, err := config.Client.vautoscaling.V2Api.PutScalingPolicy(reqParams)
	if err != nil {
		return nil, nil, err
	}

	policy := resp.ScalingPolicyList[0]
	return policy.AutoScalingGroupNo, policy.PolicyNo, nil
}

func createClassicAutoScalingPolicy(d *schema.ResourceData, config *ProviderConfig) (*string, *string, error) {
	no := d.Get("auto_scaling_group_no").(string)
	name := ncloud.String(d.Get("name").(string))
	asg, err := getClassicAutoScalingGroup(config, no)
	if err != nil {
		return nil, nil, err
	}
	reqParams := &autoscaling.PutScalingPolicyRequest{
		// Required
		AdjustmentTypeCode:   ncloud.String(d.Get("adjustment_type_code").(string)),
		ScalingAdjustment:    ncloud.Int32(int32(d.Get("scaling_adjustment").(int))),
		AutoScalingGroupName: asg.AutoScalingGroupName,
		PolicyName:           name,
		// Optional
		MinAdjustmentStep: Int32PtrOrNil(d.GetOk("min_adjustment_step")),
		Cooldown:          Int32PtrOrNil(d.GetOk("cooldown")),
	}

	if _, err := config.Client.autoscaling.V2Api.PutScalingPolicy(reqParams); err != nil {
		return nil, nil, err
	}

	return ncloud.String(no), name, nil
}

func resourceNcloudAutoScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	policy, err := getAutoScalingPolicy(config, d.Id(), d.Get("auto_scaling_group_no").(string))
	if err != nil {
		return err
	}

	if policy == nil {
		d.SetId("")
		return nil
	}

	policyMap := ConvertToMap(policy)
	SetSingularResourceDataFromMapSchema(resourceNcloudAutoScalingPolicy(), d, policyMap)
	return nil
}

func getAutoScalingPolicy(config *ProviderConfig, id string, autoScalingGroupNo string) (*AutoScalingPolicy, error) {
	if config.SupportVPC {
		return getVpcAutoScalingPolicy(config, id, autoScalingGroupNo)
	} else {
		return getClassicAutoScalingPolicy(config, id, autoScalingGroupNo)
	}
}

func getVpcAutoScalingPolicy(config *ProviderConfig, id string, autoScalingGroupNo string) (*AutoScalingPolicy, error) {
	reqParams := &vautoscaling.GetAutoScalingPolicyListRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: ncloud.String(autoScalingGroupNo),
		PolicyNoList:       []*string{ncloud.String(id)},
	}
	resp, err := config.Client.vautoscaling.V2Api.GetAutoScalingPolicyList(reqParams)
	if err != nil {
		return nil, err
	}

	p := resp.ScalingPolicyList[0]
	return &AutoScalingPolicy{
		AutoScalingPolicyNo:   p.PolicyNo,
		AutoScalingPolicyName: p.PolicyName,
		AutoScalingGroupNo:    p.AutoScalingGroupNo,
		AdjustmentTypeCode:    p.AdjustmentType.Code,
		ScalingAdjustment:     p.ScalingAdjustment,
		Cooldown:              p.CoolDown,
		MinAdjustmentStep:     p.MinAdjustmentStep,
	}, nil

}

func getClassicAutoScalingPolicy(config *ProviderConfig, id string, autoScalingGroupNo string) (*AutoScalingPolicy, error) {
	asg, err := getClassicAutoScalingGroup(config, autoScalingGroupNo)
	if err != nil {
		return nil, err
	}

	reqParams := &autoscaling.GetAutoScalingPolicyListRequest{
		PolicyNameList:       []*string{ncloud.String(id)},
		AutoScalingGroupName: asg.AutoScalingGroupName,
	}
	resp, err := config.Client.autoscaling.V2Api.GetAutoScalingPolicyList(reqParams)
	if err != nil {
		return nil, err
	}
	if len(resp.ScalingPolicyList) < 1 {
		return nil, nil
	}

	p := resp.ScalingPolicyList[0]
	return &AutoScalingPolicy{
		AutoScalingPolicyName: p.PolicyName,
		AdjustmentTypeCode:    p.AdjustmentType.Code,
		ScalingAdjustment:     p.ScalingAdjustment,
		Cooldown:              p.Cooldown,
		MinAdjustmentStep:     p.MinAdjustmentStep,
		AutoScalingGroupNo:    asg.AutoScalingGroupNo,
	}, nil
}

func resourceNcloudAutoScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	_, _, err := createAutoScalingPolicy(d, config)
	if err != nil {
		return err
	}
	return resourceNcloudAutoScalingPolicyRead(d, meta)
}

func resourceNcloudAutoScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if err := deleteAutoScalingPolicy(config, d.Id(), d.Get("auto_scaling_group_no").(string)); err != nil {
		return err
	}
	return nil
}

func deleteAutoScalingPolicy(config *ProviderConfig, id string, autoScalingGroupNo string) error {
	if config.SupportVPC {
		return deleteVpcAutoScalingPolicy(config, id, autoScalingGroupNo)
	} else {
		return deleteClassicAutoScalingPolicy(config, id, autoScalingGroupNo)
	}
}

func deleteVpcAutoScalingPolicy(config *ProviderConfig, id string, autoScalingGroupNo string) error {
	p, err := getVpcAutoScalingPolicy(config, id, autoScalingGroupNo)
	if err != nil {
		return err
	}
	reqParams := &vautoscaling.DeleteScalingPolicyRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: p.AutoScalingGroupNo,
		PolicyNo:           p.AutoScalingPolicyNo,
	}

	if _, err := config.Client.vautoscaling.V2Api.DeleteScalingPolicy(reqParams); err != nil {
		return err
	}
	return nil
}

func deleteClassicAutoScalingPolicy(config *ProviderConfig, id string, autoScalingGroupNo string) error {
	asg, err := getClassicAutoScalingGroup(config, autoScalingGroupNo)
	if err != nil {
		return err
	}
	reqParams := &autoscaling.DeletePolicyRequest{
		AutoScalingGroupName: asg.AutoScalingGroupName,
		PolicyName:           ncloud.String(id),
	}
	if _, err := config.Client.autoscaling.V2Api.DeletePolicy(reqParams); err != nil {
		return err
	}
	return nil
}
