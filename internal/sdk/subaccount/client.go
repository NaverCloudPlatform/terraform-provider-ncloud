// Package subaccount is a minimal client for the Sub Account REST API
// (https://api.ncloud-docs.com/docs/en/management-subaccount), which is not
// yet covered by ncloud-sdk-go-v2. It follows the same API Gateway HMAC
// signing scheme as the SDK so that it can be replaced with a generated SDK
// client once one is published.
package subaccount

import (
	"bytes"
	"context"
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/hmac"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

type APIClient struct {
	endpoint   string
	apiKey     *ncloud.APIKey
	httpClient *http.Client
}

func NewAPIClient(apiKey *ncloud.APIKey, endpoint string) *APIClient {
	return &APIClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type APIError struct {
	StatusCode int
	// ErrorCode is the service error code parsed from the response body
	// (e.g. "30" for a nonexistent subAccountId), empty when absent.
	ErrorCode string
	Body      string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Sub Account API returned status %d: %s", e.StatusCode, e.Body)
}

// notFoundErrorCode is how the Sub Account API signals a nonexistent
// subAccountId; it arrives with HTTP 400 or 401, never 404. A bare 404 is
// an API Gateway routing failure and must NOT be treated as deletion, or a
// misrouted request would silently wipe resources from state.
const notFoundErrorCode = "30"

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.ErrorCode == notFoundErrorCode
}

func newAPIError(statusCode int, body []byte) *APIError {
	var parsed struct {
		ErrorCode json.Number `json:"errorCode"`
		Code      json.Number `json:"code"`
		Error     *struct {
			ErrorCode json.Number `json:"errorCode"`
		} `json:"error"`
	}
	_ = json.Unmarshal(body, &parsed)

	errorCode := parsed.ErrorCode.String()
	if errorCode == "" {
		errorCode = parsed.Code.String()
	}
	if errorCode == "" && parsed.Error != nil {
		errorCode = parsed.Error.ErrorCode.String()
	}

	return &APIError{StatusCode: statusCode, ErrorCode: errorCode, Body: string(body)}
}

func (c *APIClient) do(ctx context.Context, method, path string, reqBody, respBody any) error {
	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, bodyReader)
	if err != nil {
		return err
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	signer := hmac.NewSigner(c.apiKey.SecretKey, crypto.SHA256)
	signature, err := signer.Sign(method, path, c.apiKey.AccessKey, timestamp)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", c.apiKey.AccessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp.StatusCode, respBytes)
	}

	// The API reports some failures (e.g. exceeding a quota) with HTTP 200
	// and {"success": false} in the body.
	if len(respBytes) > 0 {
		var envelope struct {
			Success *bool `json:"success"`
		}
		if err := json.Unmarshal(respBytes, &envelope); err == nil && envelope.Success != nil && !*envelope.Success {
			return newAPIError(resp.StatusCode, respBytes)
		}
	}

	if respBody != nil && len(respBytes) > 0 {
		return json.Unmarshal(respBytes, respBody)
	}

	return nil
}
