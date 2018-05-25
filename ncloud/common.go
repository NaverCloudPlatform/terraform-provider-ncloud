package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"log"
)

func convertToStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

func logCommonResponse(tag string, resp common.CommonResponse) {
	log.Printf("[DEBUG] %s Response: %s", tag, fmt.Sprintf("RequestID: %s, ReturnCode: %d, ReturnMessage: %s", resp.RequestID, resp.ReturnCode, resp.ReturnMessage))
}
