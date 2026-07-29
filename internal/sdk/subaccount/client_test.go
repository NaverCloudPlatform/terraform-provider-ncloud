package subaccount

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

func testClient(t *testing.T, handler http.HandlerFunc) *APIClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return NewAPIClient(&ncloud.APIKey{AccessKey: "testAccessKey", SecretKey: "testSecretKey"}, server.URL)
}

func TestCreateSubAccountSignsRequest(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/sub-accounts" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		for _, header := range []string{"x-ncp-apigw-timestamp", "x-ncp-iam-access-key", "x-ncp-apigw-signature-v2"} {
			if r.Header.Get(header) == "" {
				t.Errorf("missing header %s", header)
			}
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decoding request body: %s", err)
		}
		if req["loginId"] != "tf-test" {
			t.Errorf("expected loginId tf-test, got %v", req["loginId"])
		}
		if use, ok := req["useConsolePermitIp"]; ok && use != false {
			t.Errorf("expected useConsolePermitIp to be absent or false without permit ips, got %v", use)
		}

		_ = json.NewEncoder(w).Encode(CreateSubAccountResponse{Id: "sub-account-id", Success: true})
	})

	resp, err := client.CreateSubAccount(context.Background(), &CreateSubAccountRequest{LoginId: "tf-test", Name: "tf-test"})
	if err != nil {
		t.Fatalf("CreateSubAccount: %s", err)
	}
	if resp.Id != "sub-account-id" {
		t.Errorf("expected id sub-account-id, got %s", resp.Id)
	}
}

func TestUpdateSubAccountDerivesWireInvariants(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decoding request body: %s", err)
		}
		// nil slices must be sent as [] (full-replacement PUT), with the
		// paired use-flags derived as false.
		if ips, ok := req["consolePermitIps"].([]any); !ok || len(ips) != 0 {
			t.Errorf("expected consolePermitIps to be [], got %v", req["consolePermitIps"])
		}
		if req["useConsolePermitIp"] != false {
			t.Errorf("expected useConsolePermitIp false, got %v", req["useConsolePermitIp"])
		}
		if req["useApiAllowSource"] != false {
			t.Errorf("expected useApiAllowSource false, got %v", req["useApiAllowSource"])
		}
		if _, ok := req["isMfaMandatory"]; !ok {
			t.Error("expected isMfaMandatory to always be sent")
		}
		if _, ok := req["active"]; !ok {
			t.Error("expected active to always be sent")
		}

		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	})

	err := client.UpdateSubAccount(context.Background(), "abc", &UpdateSubAccountRequest{Name: "tf-test", Active: true})
	if err != nil {
		t.Fatalf("UpdateSubAccount: %s", err)
	}
}

func TestIsNotFoundMatchesErrorCode30(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"errorCode": 30, "message": "잘못된 subAccountId입니다."}`, http.StatusBadRequest)
	})

	_, err := client.GetSubAccount(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for code 30 response")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound to be true for 400 + errorCode 30, got %s", err)
	}
}

func TestIsNotFoundIgnoresGateway404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"errorCode":"200","message":"Not Found Product"}}`, http.StatusNotFound)
	})

	_, err := client.GetSubAccount(context.Background(), "abc")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if IsNotFound(err) {
		t.Error("a gateway 404 without errorCode 30 must not be treated as a deleted resource")
	}
}

func TestSuccessFalseIsAnError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false})
	})

	err := client.UpdateSubAccount(context.Background(), "abc", &UpdateSubAccountRequest{Name: "tf-test"})
	if err == nil {
		t.Fatal("expected HTTP 200 with success=false to be an error")
	}
}

func TestListAccessKeys(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/sub-accounts/abc/access-keys" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]AccessKey{{AccessKey: "AKIA", Active: true, CreateTime: "2026-01-01T00:00:00Z"}})
	})

	keys, err := client.ListAccessKeys(context.Background(), "abc")
	if err != nil {
		t.Fatalf("ListAccessKeys: %s", err)
	}
	if len(keys) != 1 || keys[0].AccessKey != "AKIA" {
		t.Errorf("unexpected keys: %+v", keys)
	}
}

func TestDeleteAccessKeySendsKeyInBody(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/sub-accounts/abc/access-keys" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decoding request body: %s", err)
		}
		if req["accessKey"] != "AKIA" {
			t.Errorf("expected accessKey AKIA in body, got %v", req)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	})

	if err := client.DeleteAccessKey(context.Background(), "abc", "AKIA"); err != nil {
		t.Fatalf("DeleteAccessKey: %s", err)
	}
}
