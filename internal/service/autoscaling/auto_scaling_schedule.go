package autoscaling

import (
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

const SCHEDULE_TIME_FORMAT = "2006-01-02T15:04:05Z0700"

func ResourceNcloudAutoScalingSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAutoScalingScheduleCreate,
		Read:   resourceNcloudAutoScalingScheduleRead,
		Update: resourceNcloudAutoScalingScheduleUpdate,
		Delete: resourceNcloudAutoScalingScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 255),
					validation.StringMatch(regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`), "Allows only lowercase letters(a-z), numbers, hyphen (-). Must start with an alphabetic character, must end with an English letter or number"))),
			},
			"desired_capacity": {
				Type:             schema.TypeInt,
				Required:         true,
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
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: ToDiagFunc(validation.Any(
					validation.IsRFC3339Time,
					ValidateDateISO8601,
				)),
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: ToDiagFunc(validation.Any(
					validation.IsRFC3339Time,
					ValidateDateISO8601,
				)),
			},
			"recurrence": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_scaling_group_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNcloudAutoScalingScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	id, err := createAutoScalingSchedule(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	return resourceNcloudAutoScalingScheduleRead(d, meta)
}

func createAutoScalingSchedule(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcAutoScalingSchedule(d, config)
	} else {
		return createClassicAutoScalingSchedule(d, config)
	}
}

func createVpcAutoScalingSchedule(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &vautoscaling.PutScheduledUpdateGroupActionRequest{
		RegionCode: &config.RegionCode,
		// Required
		AutoScalingGroupNo:  ncloud.String(d.Get("auto_scaling_group_no").(string)),
		MaxSize:             ncloud.Int32(int32(d.Get("max_size").(int))),
		MinSize:             ncloud.Int32(int32(d.Get("min_size").(int))),
		ScheduledActionName: ncloud.String(d.Get("name").(string)),
		DesiredCapacity:     ncloud.Int32(int32(d.Get("desired_capacity").(int))),
		// Optional
		StartTime:  StringPtrOrNil(d.GetOk("start_time")),
		EndTime:    StringPtrOrNil(d.GetOk("end_time")),
		Recurrence: StringPtrOrNil(d.GetOk("recurrence")),
		TimeZone:   StringPtrOrNil(d.GetOk("time_zone")),
	}

	resp, err := config.Client.Vautoscaling.V2Api.PutScheduledUpdateGroupAction(reqParams)
	if err != nil {
		return nil, err
	}

	return resp.ScheduledUpdateGroupActionList[0].ScheduledActionNo, nil
}

func createClassicAutoScalingSchedule(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	asgNo := d.Get("auto_scaling_group_no").(string)
	asg, err := getClassicAutoScalingGroup(config, asgNo)
	if err != nil {
		return nil, err
	}
	reqParams := &autoscaling.PutScheduledUpdateGroupActionRequest{
		// Required
		AutoScalingGroupName: asg.AutoScalingGroupName,
		MaxSize:              ncloud.Int32(int32(d.Get("max_size").(int))),
		MinSize:              ncloud.Int32(int32(d.Get("min_size").(int))),
		ScheduledActionName:  ncloud.String(d.Get("name").(string)),
		//Optional
		DesiredCapacity: Int32PtrOrNil(d.GetOk("desired_capacity")),
		StartTime:       StringPtrOrNil(d.GetOk("start_time")),
		EndTime:         StringPtrOrNil(d.GetOk("end_time")),
		RecurrenceInKST: StringPtrOrNil(d.GetOk("recurrence")),
	}

	resp, err := config.Client.Autoscaling.V2Api.PutScheduledUpdateGroupAction(reqParams)
	if err != nil {
		return nil, err
	}

	return resp.ScheduledUpdateGroupActionList[0].ScheduledActionName, nil
}

func resourceNcloudAutoScalingScheduleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	schedule, err := GetAutoScalingSchedule(config, d.Id(), d.Get("auto_scaling_group_no").(string))
	if err != nil {
		return err
	}

	if schedule == nil {
		d.SetId("")
		return nil
	}

	scheduleMap := ConvertToMap(schedule)
	SetSingularResourceDataFromMapSchema(ResourceNcloudAutoScalingSchedule(), d, scheduleMap)
	return nil
}

func GetAutoScalingSchedule(config *conn.ProviderConfig, id string, asgNo string) (*AutoScalingSchedule, error) {
	if config.SupportVPC {
		return getVpcAutoScalingSchedule(config, id, asgNo)
	} else {
		return getClassicAutoScalingSchedule(config, id, asgNo)
	}
}

func getVpcAutoScalingSchedule(config *conn.ProviderConfig, id string, asgNo string) (*AutoScalingSchedule, error) {
	reqParams := &vautoscaling.GetScheduledActionListRequest{
		RegionCode:            &config.RegionCode,
		AutoScalingGroupNo:    ncloud.String(asgNo),
		ScheduledActionNoList: []*string{ncloud.String(id)},
	}
	resp, err := config.Client.Vautoscaling.V2Api.GetScheduledActionList(reqParams)
	if err != nil {
		return nil, err
	}

	if len(resp.ScheduledUpdateGroupActionList) < 1 {
		return nil, nil
	}

	s := resp.ScheduledUpdateGroupActionList[0]
	return &AutoScalingSchedule{
		ScheduledActionNo:   s.ScheduledActionNo,
		ScheduledActionName: s.ScheduledActionName,
		AutoScalingGroupNo:  s.AutoScalingGroupNo,
		DesiredCapacity:     s.DesiredCapacity,
		MinSize:             s.MinSize,
		MaxSize:             s.MaxSize,
		StartTime:           s.StartTime,
		EndTime:             s.EndTime,
		RecurrenceInKST:     s.Recurrence,
		TimeZone:            s.TimeZone,
	}, nil
}

func getClassicAutoScalingSchedule(config *conn.ProviderConfig, id string, asgNo string) (*AutoScalingSchedule, error) {
	asg, err := getClassicAutoScalingGroup(config, asgNo)
	if err != nil {
		return nil, err
	}
	reqParams := &autoscaling.GetScheduledActionListRequest{
		AutoScalingGroupName:    asg.AutoScalingGroupName,
		ScheduledActionNameList: []*string{ncloud.String(id)},
	}
	resp, err := config.Client.Autoscaling.V2Api.GetScheduledActionList(reqParams)
	if err != nil {
		return nil, err
	}

	if len(resp.ScheduledUpdateGroupActionList) < 1 {
		return nil, nil
	}

	s := resp.ScheduledUpdateGroupActionList[0]
	return &AutoScalingSchedule{
		AutoScalingGroupNo:  asg.AutoScalingGroupNo,
		ScheduledActionName: s.ScheduledActionName,
		DesiredCapacity:     s.DesiredCapacity,
		MinSize:             s.MinSize,
		MaxSize:             s.MaxSize,
		StartTime:           s.StartTime,
		EndTime:             s.EndTime,
		RecurrenceInKST:     s.RecurrenceInKST,
	}, nil
}

func resourceNcloudAutoScalingScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if _, err := createAutoScalingSchedule(d, config); err != nil {
		return err
	}
	return resourceNcloudAutoScalingScheduleRead(d, meta)
}

func resourceNcloudAutoScalingScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if err := deleteAutoScalingSchedule(config, d.Id(), d.Get("auto_scaling_group_no").(string)); err != nil {
		return err
	}
	return nil
}

func deleteAutoScalingSchedule(config *conn.ProviderConfig, id string, asgNo string) error {
	if config.SupportVPC {
		return deleteVpcAutoScalingSchedule(config, id, asgNo)
	} else {
		return deleteClassicAutoScalingSchedule(config, id, asgNo)
	}
}

func deleteVpcAutoScalingSchedule(config *conn.ProviderConfig, id string, asgNo string) error {
	schedule, err := getVpcAutoScalingSchedule(config, id, asgNo)
	if err != nil {
		return err
	}
	reqParams := &vautoscaling.DeleteScheduledActionRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: ncloud.String(asgNo),
		ScheduledActionNo:  schedule.ScheduledActionNo,
	}
	_, err = config.Client.Vautoscaling.V2Api.DeleteScheduledAction(reqParams)
	if err != nil {
		return err
	}
	return nil
}

func deleteClassicAutoScalingSchedule(config *conn.ProviderConfig, id string, asgNo string) error {
	asg, err := getClassicAutoScalingGroup(config, asgNo)
	if err != nil {
		return err
	}
	reqParams := &autoscaling.DeleteScheduledActionRequest{
		AutoScalingGroupName: asg.AutoScalingGroupName,
		ScheduledActionName:  ncloud.String(id),
	}
	if _, err := config.Client.Autoscaling.V2Api.DeleteScheduledAction(reqParams); err != nil {
		return err
	}
	return nil
}
