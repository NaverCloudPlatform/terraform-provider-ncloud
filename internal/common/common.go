package common

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

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

	ApiErrorNetworkInterfaceAtLeastOneAcgMustRemain = "1002035"

	ApiErrorAcgCantChangeSameTime           = "1007009"
	ApiErrorNetworkAclCantAccessaApropriate = "1011002"
	ApiErrorNetworkAclRuleChangeIngRules    = "1012005"

	ApiErrorASGIsUsingPolicyOrLaunchConfiguration      = "50150" // This is returned when you cannot delete a launch configuration, scaling policy, or auto scaling group because it is being used.
	ApiErrorASGScalingIsActive                         = "50160" // You cannot request actions while there are scaling activities in progress for that group.
	ApiErrorASGIsUsingPolicyOrLaunchConfigurationOnVpc = "1250700"
)

const (
	InstanceStatusInit        = "INIT"
	InstanceStatusCreate      = "CREATING"
	InstanceStatusRunning     = "RUN"
	InstanceStatusSetting     = "SET"
	InstanceStatusTerminating = "TERMTING"
	InstanceStatusTerminated  = "TERMINATED"
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

// CommonError response error body
type CommonError struct {
	ReturnCode    string
	ReturnMessage string
}

func LogErrorResponse(tag string, err error, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[ERROR] %s error params=%s, err=%s", tag, param, err)
}

func LogCommonRequest(tag string, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[INFO] %s params=%s", tag, param)
}

func LogResponse(tag string, args interface{}) {
	resp, _ := json.Marshal(args)
	log.Printf("[INFO] %s response=%s", tag, resp)
}

func LogCommonResponse(tag string, commonResponse *CommonResponse, logs ...string) {
	result := fmt.Sprintf("RequestID: %s, ReturnCode: %s, ReturnMessage: %s", ncloud.StringValue(commonResponse.RequestId), ncloud.StringValue(commonResponse.ReturnCode), ncloud.StringValue(commonResponse.ReturnMessage))
	log.Printf("[INFO] %s success response=%s %s", tag, result, strings.Join(logs, " "))
}

func ContainsInStringList(str string, s []string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ExtractEngineVersion(input string) string {
	re := regexp.MustCompile(`\d+\.\d+(\.\d+)?`)
	version := re.FindString(input)

	return version
}
