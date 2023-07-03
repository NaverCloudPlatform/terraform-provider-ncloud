package autoscaling

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudAutoScalingSchedule() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"auto_scaling_group_no": {
			Type:     schema.TypeString,
			Required: true,
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudAutoScalingSchedule(), fieldMap, dataSourceNcloudAutoScalingScheduleRead)
}

func dataSourceNcloudAutoScalingScheduleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	scheduleList, err := getAutoScalingScheduleList(d, config)
	if err != nil {
		return err
	}

	scheduleListMap := ConvertToArrayMap(scheduleList)
	if f, ok := d.GetOk("filter"); ok {
		scheduleListMap = ApplyFilters(f.(*schema.Set), scheduleListMap, DataSourceNcloudAutoScalingSchedule().Schema)
	}

	if err := ValidateOneResult(len(scheduleListMap)); err != nil {
		return err
	}

	d.SetId(scheduleListMap[0]["name"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudAutoScalingSchedule(), d, scheduleListMap[0])
	return nil
}

func getAutoScalingScheduleList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*AutoScalingSchedule, error) {
	if config.SupportVPC {
		return getVpcAutoScalingScheduleList(d, config)
	} else {
		return getClassicAutoScalingScheduleList(d, config)
	}
}

func getVpcAutoScalingScheduleList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*AutoScalingSchedule, error) {
	reqParams := &vautoscaling.GetScheduledActionListRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: ncloud.String(d.Get("auto_scaling_group_no").(string)),
	}
	if d.Id() != "" {
		reqParams.ScheduledActionNameList = []*string{ncloud.String(d.Id())}
	}

	resp, err := config.Client.Vautoscaling.V2Api.GetScheduledActionList(reqParams)
	if err != nil {
		return nil, err
	}

	list := make([]*AutoScalingSchedule, 0)
	for _, s := range resp.ScheduledUpdateGroupActionList {
		schedule := &AutoScalingSchedule{
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
		}
		list = append(list, schedule)
	}
	if len(list) < 1 {
		return nil, nil
	}

	return list, nil
}

func getClassicAutoScalingScheduleList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*AutoScalingSchedule, error) {
	reqParams := &autoscaling.GetScheduledActionListRequest{}

	if d.Id() != "" {
		reqParams.ScheduledActionNameList = []*string{ncloud.String(d.Id())}
	}

	resp, err := config.Client.Autoscaling.V2Api.GetScheduledActionList(reqParams)
	if err != nil {
		return nil, err
	}

	list := make([]*AutoScalingSchedule, 0)
	for _, s := range resp.ScheduledUpdateGroupActionList {
		asg, err := getClassicAutoScalingGroupByName(config, *s.AutoScalingGroupName)
		if err != nil {
			return nil, err
		}
		schedule := &AutoScalingSchedule{
			ScheduledActionName: s.ScheduledActionName,
			AutoScalingGroupNo:  asg.AutoScalingGroupNo,
			DesiredCapacity:     s.DesiredCapacity,
			MinSize:             s.MinSize,
			MaxSize:             s.MaxSize,
			StartTime:           s.StartTime,
			EndTime:             s.EndTime,
			RecurrenceInKST:     s.RecurrenceInKST,
		}
		list = append(list, schedule)
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list, nil
}
