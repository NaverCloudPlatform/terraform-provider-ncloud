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
