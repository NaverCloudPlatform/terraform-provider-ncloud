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
	UseConsolePermitIp   bool             `json:"useConsolePermitIp,omitempty"`
	ConsolePermitIps     []string         `json:"consolePermitIps,omitempty"`
	UseApiAllowSource    bool             `json:"useApiAllowSource,omitempty"`
	ApiAllowSources      []ApiAllowSource `json:"apiAllowSources,omitempty"`
}

type CreateSubAccountResponse struct {
	Id                string `json:"id"`
	Success           bool   `json:"success"`
	GeneratedPassword string `json:"generatedPassword"`
}

// UpdateSubAccountRequest intentionally sends every mutable field on each PUT:
// the Edit Sub Account API treats the request as a full replacement, so
// omitting a field the user removed from configuration would leave the old
// value behind and show up as permanent drift.
type UpdateSubAccountRequest struct {
	Name                string           `json:"name"`
	Email               string           `json:"email"`
	Memo                string           `json:"memo"`
	IsMfaMandatory      *bool            `json:"isMfaMandatory,omitempty"`
	CanConsoleAccess    bool             `json:"canConsoleAccess"`
	CanAPIGatewayAccess bool             `json:"canAPIGatewayAccess"`
	UseConsolePermitIp  bool             `json:"useConsolePermitIp"`
	ConsolePermitIps    []string         `json:"consolePermitIps"`
	UseApiAllowSource   bool             `json:"useApiAllowSource"`
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

func (c *APIClient) CreateSubAccount(ctx context.Context, req *CreateSubAccountRequest) (*CreateSubAccountResponse, error) {
	var resp CreateSubAccountResponse
	if err := c.do(ctx, http.MethodPost, "/api/v1/sub-accounts", req, &resp); err != nil {
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
	return c.do(ctx, http.MethodPut, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId), req, nil)
}

func (c *APIClient) DeleteSubAccount(ctx context.Context, subAccountId string) error {
	return c.do(ctx, http.MethodDelete, "/api/v1/sub-accounts/"+url.PathEscape(subAccountId), nil, nil)
}
