package ncloud

import (
	"fmt"
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

func TestGetCommonErrorBody(t *testing.T) {
	err := fmt.Errorf(`Status: 400 Bad Request, Body: {"responseError": {
  "returnCode": "1007009",
  "returnMessage": "If the Acg settings are being changed, you cannot change other settings at the same time."
}}`)

	e, err := GetCommonErrorBody(err)

	if err != nil {
		t.Fatalf("Got error: %s", err)
	}

	if e.ReturnCode != "1007009" {
		t.Fatalf("Return code expected '1007009' but %s", e.ReturnCode)
	}

	if e.ReturnMessage != "If the Acg settings are being changed, you cannot change other settings at the same time." {
		t.Fatalf("Return code expected 'If the Acg settings are being changed, you cannot change other settings at the same time.' but %s", e.ReturnMessage)
	}

}
