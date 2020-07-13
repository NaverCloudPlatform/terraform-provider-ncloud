package ncloud

import (
	"encoding/json"
	"fmt"
	"log"

	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

const (
	ApiErrorAuthorityParameter = "800"
	ApiErrorUnknown            = "1300"

	ApiErrorObjectInOperation                            = "25013"
	ApiErrorPortForwardingObjectInOperation              = "25033"
	ApiErrorServerObjectInOperation                      = "23006" // Unable to request server termination and creation simultaneously
	ApiErrorServerObjectInOperation2                     = "25017"
	ApiErrorPreviousServersHaveNotBeenEntirelyTerminated = "23003"

	ApiErrorDetachingMountedStorage = "24002"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
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

func logCommonRequest(tag string, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[INFO] %s params=%s", tag, param)
}

func logResponse(tag string, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[INFO] %s response=%s", tag, param)
}

func logCommonResponse(tag string, commonResponse *CommonResponse, logs ...string) {
	result := fmt.Sprintf("RequestID: %s, ReturnCode: %s, ReturnMessage: %s", ncloud.StringValue(commonResponse.RequestId), ncloud.StringValue(commonResponse.ReturnCode), ncloud.StringValue(commonResponse.ReturnMessage))
	log.Printf("[INFO] %s success response=%s %s", tag, result, strings.Join(logs, " "))
}

func isRetryableErr(commResp *CommonResponse, code []string) bool {
	for _, c := range code {
		if commResp != nil && commResp.ReturnCode != nil && ncloud.StringValue(commResp.ReturnCode) == c {
			return true
		}
	}

	return false
}
