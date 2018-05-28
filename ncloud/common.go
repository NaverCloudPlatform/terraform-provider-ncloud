package ncloud

import (
	"encoding/json"
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

func logCommonResponse(tag string, err error, args interface{}, resp common.CommonResponse) {
	param, _ := json.Marshal(args)

	if err != nil {
		log.Printf("[DEBUG] %s error params=%s, err=%s", tag, param, err)
	} else {
		result := fmt.Sprintf("RequestID: %s, ReturnCode: %d, ReturnMessage: %s", resp.RequestID, resp.ReturnCode, resp.ReturnMessage)
		log.Printf("[DEBUG] %s success params=%s, response=%s", tag, param, result)
	}
}
