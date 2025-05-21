package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
)

func flattenAccessControlRules(accessControlRules []*server.AccessControlRule) []string {
	var s []string

	for _, accessControlRule := range accessControlRules {
		s = append(s, ncloud.StringValue(accessControlRule.AccessControlRuleConfigurationNo))
	}

	return s
}

func expandTagListParams(tl []interface{}) ([]*server.InstanceTagParameter, error) {
	tagList := make([]*server.InstanceTagParameter, 0, len(tl))

	for _, v := range tl {
		tag := new(server.InstanceTagParameter)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "tag_key":
				tag.TagKey = ncloud.String(value.(string))
			case "tag_value":
				tag.TagValue = ncloud.String(value.(string))
			}
		}
		tagList = append(tagList, tag)
	}

	return tagList, nil
}

func expandBlockDevicePartitionListParams(bl []interface{}) ([]*vserver.BlockDevicePartition, error) {
	blockDevicePartitionList := make([]*vserver.BlockDevicePartition, 0, len(bl))

	for _, v := range bl {
		blockDevicePartition := &vserver.BlockDevicePartition{}
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "mount_point":
				blockDevicePartition.MountPoint = ncloud.String(value.(string))
			case "partition_size":
				blockDevicePartition.PartitionSize = ncloud.String(value.(string))
			}
		}
		blockDevicePartitionList = append(blockDevicePartitionList, blockDevicePartition)
	}

	return blockDevicePartitionList, nil
}

func flattenMapByKey(i interface{}, key string) *string {
	m := ConvertToMap(i)
	if m[key] != nil {
		return ncloud.String(m[key].(string))
	} else {
		return nil
	}
}
