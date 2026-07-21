package subaccount

import (
	"context"
	"net/http"
	"net/url"
)

type AccessKey struct {
	AccessKey  string `json:"accessKey"`
	KeySecret  string `json:"keySecret,omitempty"`
	Active     bool   `json:"active,omitempty"`
	CreateTime string `json:"createTime"`
}

// CreateAccessKey issues a new API access key for the sub account.
// KeySecret is returned only in this response and can never be retrieved again.
func (c *APIClient) CreateAccessKey(ctx context.Context, subAccountId string) (*AccessKey, error) {
	var resp AccessKey
	if err := c.do(ctx, http.MethodPost, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId)+"/access-keys", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *APIClient) ListAccessKeys(ctx context.Context, subAccountId string) ([]AccessKey, error) {
	var resp []AccessKey
	if err := c.do(ctx, http.MethodGet, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId)+"/access-keys", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *APIClient) DeleteAccessKey(ctx context.Context, subAccountId, accessKeyId string) error {
	return c.do(ctx, http.MethodDelete, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId)+"/access-keys/"+url.PathEscape(accessKeyId), nil, nil)
}
