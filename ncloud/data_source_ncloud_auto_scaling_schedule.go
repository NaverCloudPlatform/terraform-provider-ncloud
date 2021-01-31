package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_auto_scaling_schedule", dataSourceNcloudAutoScalingSchedule())
}

func dataSourceNcloudAutoScalingSchedule() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudAutoScalingSchedule(), fieldMap, dataSourceNcloudAutoScalingScheduleRead)
}

func dataSourceNcloudAutoScalingScheduleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if v, ok := d.GetOk("name"); ok {
		d.SetId(v.(string))
	}

	scheduleList, err := getAutoScalingScheduleList(config, d.Id())
	if err != nil {
		return err
	}

	scheduleListMap := ConvertToArrayMap(scheduleList)
	if f, ok := d.GetOk("filter"); ok {
		scheduleListMap = ApplyFilters(f.(*schema.Set), scheduleListMap, dataSourceNcloudAutoScalingSchedule().Schema)
	}

	if err := validateOneResult(len(scheduleListMap)); err != nil {
		return err
	}

	d.SetId(scheduleListMap[0]["name"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudAutoScalingSchedule(), d, scheduleListMap[0])
	return nil
}

func getAutoScalingScheduleList(config *ProviderConfig, id string) ([]*AutoScalingSchedule, error) {
	if config.SupportVPC {
		return getVpcAutoScalingScheduleList(config, id)
	} else {
		return getClassicAutoScalingScheduleList(config, id)
	}
}

func getVpcAutoScalingScheduleList(config *ProviderConfig, id string) ([]*AutoScalingSchedule, error) {
	reqParams := &vautoscaling.GetScheduledActionListRequest{
		RegionCode: &config.RegionCode,
	}
	if id != "" {
		reqParams.ScheduledActionNameList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.vautoscaling.V2Api.GetScheduledActionList(reqParams)
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

func getClassicAutoScalingScheduleList(config *ProviderConfig, id string) ([]*AutoScalingSchedule, error) {
	reqParams := &autoscaling.GetScheduledActionListRequest{}

	if id != "" {
		reqParams.ScheduledActionNameList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.autoscaling.V2Api.GetScheduledActionList(reqParams)
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
