package ncloud

import (
	"reflect"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

type TestResponse struct {
	RequestId     *string `json:"requestId,omitempty"`
	ReturnCode    *string `json:"returnCode,omitempty"`
	ReturnMessage *string `json:"returnMessage,omitempty"`
	TotalRows     *int32  `json:"totalRows,omitempty"`
}

type TestCommonCode struct {
	Code     *string `json:"code,omitempty"`
	CodeName *string `json:"codeName,omitempty"`
}

type TestRegion struct {
	RegionNo   *string `json:"regionNo,omitempty"`
	RegionCode *string `json:"regionCode,omitempty"`
	RegionName *string `json:"regionName,omitempty"`
}

type TestZone struct {
	ZoneNo          *string `json:"zoneNo,omitempty"`
	ZoneName        *string `json:"zoneName,omitempty"`
	ZoneDescription *string `json:"zoneDescription,omitempty"`
}

func TestGetCommonResponse(t *testing.T) {
	requestId := ncloud.String("RequestId")
	returnCode := ncloud.String("ReturnCode")
	returnMessage := ncloud.String("ReturnMessage")
	resp := &TestResponse{
		RequestId:     requestId,
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}
	commonResponse := GetCommonResponse(resp)

	if !reflect.DeepEqual(requestId, commonResponse.RequestId) {
		t.Fatalf("Expected: %s, Actual: %s", ncloud.StringValue(requestId), ncloud.StringValue(commonResponse.RequestId))
	}
}

func TestGetCommonResponse_convertNil(t *testing.T) {
	returnCode := ncloud.String("ReturnCode")
	returnMessage := ncloud.String("ReturnMessage")
	resp := &TestResponse{
		RequestId:     nil,
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}
	commonResponse := GetCommonResponse(resp)

	if commonResponse.RequestId != nil {
		t.Fatalf("Expected: nil, Actual: %s", ncloud.StringValue(commonResponse.RequestId))
	}
}

// func TestGetCommonCode(t *testing.T) {
// 	code := ncloud.String("code")
// 	codeName := ncloud.String("codeName")
// 	resp := &TestCommonCode{
// 		Code:     code,
// 		CodeName: codeName,
// 	}
// 	commonCode := GetCommonCode(resp)

// 	if !reflect.DeepEqual(code, commonCode.Code) {
// 		t.Fatalf("Expected: %s, Actual: %s", ncloud.StringValue(code), ncloud.StringValue(commonCode.Code))
// 	}
// }

// func TestGetCommonCode_convertNil(t *testing.T) {
// 	codeName := ncloud.String("codeName")
// 	resp := &TestCommonCode{
// 		Code:     nil,
// 		CodeName: codeName,
// 	}
// 	commonCode := GetCommonCode(resp)

// 	if commonCode.Code != nil {
// 		t.Fatalf("Expected: nil, Actual: %s", ncloud.StringValue(commonCode.Code))
// 	}
// }

// func TestGetZone(t *testing.T) {
// 	code := ncloud.String("code")
// 	codeName := ncloud.String("codeName")
// 	resp := &TestCommonCode{
// 		Code:     code,
// 		CodeName: codeName,
// 	}
// 	commonCode := GetCommonCode(resp)

// 	if !reflect.DeepEqual(code, commonCode.Code) {
// 		t.Fatalf("Expected: %s, Actual: %s", ncloud.StringValue(code), ncloud.StringValue(commonCode.Code))
// 	}
// }

// func TestGetZone_convertNil(t *testing.T) {
// 	codeName := ncloud.String("codeName")
// 	resp := &TestCommonCode{
// 		Code:     nil,
// 		CodeName: codeName,
// 	}
// 	commonCode := GetCommonCode(resp)

// 	if commonCode.Code != nil {
// 		t.Fatalf("Expected: nil, Actual: %s", ncloud.StringValue(commonCode.Code))
// 	}
// }
