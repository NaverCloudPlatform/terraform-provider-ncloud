package autoscaling

import "github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"

func flattenAccessControlGroupList(asgs []*autoscaling.AccessControlGroup) []*string {
	l := make([]*string, 0)
	for _, asg := range asgs {
		l = append(l, asg.AccessControlGroupConfigurationNo)
	}
	return l
}
