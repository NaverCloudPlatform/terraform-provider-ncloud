package subaccount

import (
	"context"
	"net/http"
	"net/url"
)

type ApiAllowSource struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

type CreateSubAccountRequest struct {
	LoginId              string           `json:"loginId"`
	Name                 string           `json:"name"`
	CanConsoleAccess     bool             `json:"canConsoleAccess"`
	CanAPIGatewayAccess  bool             `json:"canAPIGatewayAccess"`
	NeedPasswordReset    bool             `json:"needPasswordReset"`
	NeedPasswordGenerate bool             `json:"needPasswordGenerate,omitempty"`
	Email                string           `json:"email,omitempty"`
	IsMfaMandatory       bool             `json:"isMfaMandatory,omitempty"`
	Memo                 string           `json:"memo,omitempty"`
	ConsolePermitIps     []string         `json:"consolePermitIps,omitempty"`
	ApiAllowSources      []ApiAllowSource `json:"apiAllowSources,omitempty"`
}

type CreateSubAccountResponse struct {
	Id                string `json:"id"`
	Success           bool   `json:"success"`
	GeneratedPassword string `json:"generatedPassword"`
}

// UpdateSubAccountRequest fields are all sent on each PUT: the Edit Sub
// Account API treats the request as a full replacement, so omitting a field
// would leave the old server-side value behind as invisible drift.
// The useConsolePermitIp/useApiAllowSource wire flags are derived from the
// slices by the client and are not part of this struct.
type UpdateSubAccountRequest struct {
	Name                string           `json:"name"`
	Email               string           `json:"email"`
	Memo                string           `json:"memo"`
	Active              bool             `json:"active"`
	IsMfaMandatory      bool             `json:"isMfaMandatory"`
	CanConsoleAccess    bool             `json:"canConsoleAccess"`
	CanAPIGatewayAccess bool             `json:"canAPIGatewayAccess"`
	ConsolePermitIps    []string         `json:"consolePermitIps"`
	ApiAllowSources     []ApiAllowSource `json:"apiAllowSources"`
}

type SubAccountDetail struct {
	SubAccountId        string           `json:"subAccountId"`
	SubAccountNo        int64            `json:"subAccountNo"`
	LoginId             string           `json:"loginId"`
	Name                string           `json:"name"`
	Email               string           `json:"email"`
	CanAPIGatewayAccess bool             `json:"canAPIGatewayAccess"`
	CanConsoleAccess    bool             `json:"canConsoleAccess"`
	UseConsolePermitIp  bool             `json:"useConsolePermitIp"`
	ConsolePermitIps    []string         `json:"consolePermitIps"`
	UseApiAllowSource   bool             `json:"useApiAllowSource"`
	ApiAllowSources     []ApiAllowSource `json:"apiAllowSources"`
	Memo                string           `json:"memo"`
	Active              bool             `json:"active"`
	CreateTime          string           `json:"createTime"`
	Nrn                 string           `json:"nrn"`
}

// createSubAccountPayload adds the wire flags the API pairs with the
// allowlist slices; they are derived here so the flag/list invariant
// (flag true iff list non-empty) holds for every caller.
type createSubAccountPayload struct {
	*CreateSubAccountRequest
	UseConsolePermitIp bool `json:"useConsolePermitIp,omitempty"`
	UseApiAllowSource  bool `json:"useApiAllowSource,omitempty"`
}

type updateSubAccountPayload struct {
	*UpdateSubAccountRequest
	UseConsolePermitIp bool `json:"useConsolePermitIp"`
	UseApiAllowSource  bool `json:"useApiAllowSource"`
}

func (c *APIClient) CreateSubAccount(ctx context.Context, req *CreateSubAccountRequest) (*CreateSubAccountResponse, error) {
	payload := &createSubAccountPayload{
		CreateSubAccountRequest: req,
		UseConsolePermitIp:      len(req.ConsolePermitIps) > 0,
		UseApiAllowSource:       len(req.ApiAllowSources) > 0,
	}

	var resp CreateSubAccountResponse
	if err := c.do(ctx, http.MethodPost, "/api/v1/sub-accounts", payload, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *APIClient) GetSubAccount(ctx context.Context, subAccountId string) (*SubAccountDetail, error) {
	var resp SubAccountDetail
	if err := c.do(ctx, http.MethodGet, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *APIClient) UpdateSubAccount(ctx context.Context, subAccountId string, req *UpdateSubAccountRequest) error {
	// The PUT is a full replacement: nil slices must become [] so removed
	// allowlists are cleared server-side instead of marshaling to null.
	normalized := *req
	if normalized.ConsolePermitIps == nil {
		normalized.ConsolePermitIps = []string{}
	}
	if normalized.ApiAllowSources == nil {
		normalized.ApiAllowSources = []ApiAllowSource{}
	}

	payload := &updateSubAccountPayload{
		UpdateSubAccountRequest: &normalized,
		UseConsolePermitIp:      len(normalized.ConsolePermitIps) > 0,
		UseApiAllowSource:       len(normalized.ApiAllowSources) > 0,
	}

	return c.do(ctx, http.MethodPut, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId), payload, nil)
}

func (c *APIClient) DeleteSubAccount(ctx context.Context, subAccountId string) error {
	return c.do(ctx, http.MethodDelete, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId), nil, nil)
}
