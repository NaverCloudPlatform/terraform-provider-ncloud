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

func NewAPIClient(apiKey *ncloud.APIKey, site string) *APIClient {
	var endpoint string
	switch site {
	case "gov":
		endpoint = "https://subaccount.apigw.gov-ntruss.com"
	case "fin":
		endpoint = "https://subaccount.apigw.fin-ntruss.com"
	default:
		endpoint = "https://subaccount.apigw.ntruss.com"
	}

	return &APIClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewAPIClientWithEndpoint is used by tests to point the client at a mock server.
func NewAPIClientWithEndpoint(apiKey *ncloud.APIKey, endpoint string) *APIClient {
	return &APIClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Sub Account API returned status %d: %s", e.StatusCode, e.Body)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
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
		return &APIError{StatusCode: resp.StatusCode, Body: string(respBytes)}
	}

	if respBody != nil && len(respBytes) > 0 {
		return json.Unmarshal(respBytes, respBody)
	}

	return nil
}
