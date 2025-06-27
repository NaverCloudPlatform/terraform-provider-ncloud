package autoscaling

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudAutoScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAutoScalingPolicyCreate,
		Read:   resourceNcloudAutoScalingPolicyRead,
		Update: resourceNcloudAutoScalingPolicyUpdate,
		Delete: resourceNcloudAutoScalingPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected auto_scaling_group_no:id", d.Id())
				}
				AutoScalingGroupNo := idParts[0]
				id := idParts[1]
				d.SetId(id)
				d.Set("auto_scaling_group_no", AutoScalingGroupNo)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 255),
					validation.StringMatch(regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`), "Allows only lowercase letters(a-z), numbers, hyphen (-). Must start with an alphabetic character, must end with an English letter or number"))),
			},
			"adjustment_type_code": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CHANG", "EXACT", "PRCNT"}, false)),
			},
			"scaling_adjustment": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(-2147483648, 2147483647)),
			},
			"cooldown": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          300,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 2147483647)),
			},
			"min_adjustment_step": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 2147483647)),
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
	config := meta.(*conn.ProviderConfig)

	autoscaling_group_no, id, err := createAutoScalingPolicy(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	d.Set("auto_scaling_group_no", autoscaling_group_no)
	return resourceNcloudAutoScalingPolicyRead(d, meta)
}

func createAutoScalingPolicy(d *schema.ResourceData, config *conn.ProviderConfig) (*string, *string, error) {
	reqParams := &vautoscaling.PutScalingPolicyRequest{
		RegionCode: &config.RegionCode,
		// Required
		AdjustmentTypeCode: ncloud.String(d.Get("adjustment_type_code").(string)),
		ScalingAdjustment:  ncloud.Int32(int32(d.Get("scaling_adjustment").(int))),
		AutoScalingGroupNo: ncloud.String(d.Get("auto_scaling_group_no").(string)),
		PolicyName:         ncloud.String(d.Get("name").(string)),
		// Optional
		MinAdjustmentStep: Int32PtrOrNil(d.GetOk("min_adjustment_step")),
		CoolDown:          ncloud.Int32(int32(d.Get("cooldown").(int))),
	}
	LogCommonRequest("createVpcAutoScalingPolicy", reqParams)

	resp, err := config.Client.Vautoscaling.V2Api.PutScalingPolicy(reqParams)
	if err != nil {
		return nil, nil, err
	}
	LogResponse("createVpcAutoScalingPolicy", resp)

	policy := resp.ScalingPolicyList[0]
	return policy.AutoScalingGroupNo, policy.PolicyNo, nil
}

func resourceNcloudAutoScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	policy, err := GetAutoScalingPolicy(config, d.Id(), d.Get("auto_scaling_group_no").(string))
	if err != nil {
		return err
	}

	if policy == nil {
		d.SetId("")
		return nil
	}

	policyMap := ConvertToMap(policy)
	SetSingularResourceDataFromMapSchema(ResourceNcloudAutoScalingPolicy(), d, policyMap)
	return nil
}

func GetAutoScalingPolicy(config *conn.ProviderConfig, id string, autoScalingGroupNo string) (*AutoScalingPolicy, error) {
	reqParams := &vautoscaling.GetAutoScalingPolicyListRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: ncloud.String(autoScalingGroupNo),
		PolicyNoList:       []*string{ncloud.String(id)},
	}
	LogCommonRequest("getVpcAutoScalingPolicy", reqParams)

	resp, err := config.Client.Vautoscaling.V2Api.GetAutoScalingPolicyList(reqParams)
	if err != nil {
		return nil, err
	}
	LogResponse("getVpcAutoScalingPolicy", resp)

	if len(resp.ScalingPolicyList) == 0 {
		return nil, nil
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

func resourceNcloudAutoScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	_, _, err := createAutoScalingPolicy(d, config)
	if err != nil {
		return err
	}
	return resourceNcloudAutoScalingPolicyRead(d, meta)
}

func resourceNcloudAutoScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if err := deleteAutoScalingPolicy(config, d.Id(), d.Get("auto_scaling_group_no").(string)); err != nil {
		return err
	}
	return nil
}

func deleteAutoScalingPolicy(config *conn.ProviderConfig, id string, autoScalingGroupNo string) error {
	p, err := GetAutoScalingPolicy(config, id, autoScalingGroupNo)
	if err != nil {
		return err
	}
	reqParams := &vautoscaling.DeleteScalingPolicyRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: p.AutoScalingGroupNo,
		PolicyNo:           p.AutoScalingPolicyNo,
	}
	LogCommonRequest("deleteVpcAutoScalingPolicy", reqParams)

	resp, err := config.Client.Vautoscaling.V2Api.DeleteScalingPolicy(reqParams)
	if err != nil {
		return err
	}
	LogResponse("deleteVpcAutoScalingPolicy", resp)

	return nil
}
