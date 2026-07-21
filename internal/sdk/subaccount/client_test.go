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

	return NewAPIClientWithEndpoint(&ncloud.APIKey{AccessKey: "testAccessKey", SecretKey: "testSecretKey"}, server.URL)
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

		var req CreateSubAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decoding request body: %s", err)
		}
		if req.LoginId != "tf-test" {
			t.Errorf("expected loginId tf-test, got %s", req.LoginId)
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

func TestGetSubAccountNotFound(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
	})

	_, err := client.GetSubAccount(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound to be true, got %s", err)
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
