package ncloud

import (
	"fmt"
	"testing"
)

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
