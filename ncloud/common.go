package ncloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

type CommonResponse struct {
	RequestId     *string `json:"requestId,omitempty"`
	ReturnCode    *string `json:"returnCode,omitempty"`
	ReturnMessage *string `json:"returnMessage,omitempty"`
}

type CommonCode struct {
	Code     *string `json:"code,omitempty"`
	CodeName *string `json:"codeName,omitempty"`
}

func logErrorResponse(tag string, err error, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[ERROR] %s error params=%s, err=%s", tag, param, err)
}

func logCommonResponse(tag string, args interface{}, commonResponse *CommonResponse) {
	param, _ := json.Marshal(args)
	result := fmt.Sprintf("RequestID: %s, ReturnCode: %s, ReturnMessage: %s", ncloud.StringValue(commonResponse.RequestId), ncloud.StringValue(commonResponse.ReturnCode), ncloud.StringValue(commonResponse.ReturnMessage))
	log.Printf("[DEBUG] %s success params=%s, response=%s", tag, param, result)
}

func isRetryableErr(commResp *CommonResponse, code []string) bool {
	for _, c := range code {
		if commResp != nil && commResp.ReturnCode != nil && *commResp.ReturnCode == c {
			return true
		}
	}

	return false
}
