/* =================================================================================
 * NCLOUD SDK LAYER FOR TERRAFORM CODEGEN - DO NOT EDIT
 * =================================================================================
 * Refresh Template
 * Required data are as follows
 *
 *		MethodName         string
 *		PrimitiveRequest   string
 *		StringifiedRequest string
 *		Query              string
 *		Body               string
 *		Path               string
 *		Method             string
 * ================================================================================= */

package ncloudsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrimitiveGETProductsProductidApisApiidStagesStageidUsageplansRequest struct {
	Productid     string `json:"product-id"`
	Apiid         string `json:"api-id"`
	Stageid       string `json:"stage-id"`
	Limit         int64  `json:"limit"`
	Offset        int64  `json:"offset"`
	UsagePlanName string `json:"usagePlanName"`
}

type StringifiedGETProductsProductidApisApiidStagesStageidUsageplansRequest struct {
	Productid     string `json:"product-id"`
	Apiid         string `json:"api-id"`
	Stageid       string `json:"stage-id"`
	Limit         string `json:"limit"`
	Offset        string `json:"offset"`
	UsagePlanName string `json:"usagePlanName"`
}

func (n *NClient) GETProductsProductidApisApiidStagesStageidUsageplans(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidApisApiidStagesStageidUsageplansRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.Limit != "" {
		query["limit"] = r.Limit
	}

	if r.Offset != "" {
		query["offset"] = r.Offset
	}

	if r.UsagePlanName != "" {
		query["usagePlanName"] = r.UsagePlanName
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "apis" + "/" + ClearDoubleQuote(r.Apiid) + "/" + "stages" + "/" + ClearDoubleQuote(r.Stageid) + "/" + "usage-plans"

	response, err := n.MakeRequestWithContext(ctx, "GET", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) GETProductsProductidApisApiidStagesStageidUsageplans_TF(ctx context.Context, r *PrimitiveGETProductsProductidApisApiidStagesStageidUsageplansRequest) (*GETProductsProductidApisApiidStagesStageidUsageplansResponse, error) {
	t, err := n.GETProductsProductidApisApiidStagesStageidUsageplans(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidApisApiidStagesStageidUsageplans(context.TODO(), t)
	if err != nil {
		return nil, err
	}

	return res, nil
}

/* =================================================================================
 * NCLOUD SDK LAYER FOR TERRAFORM CODEGEN - DO NOT EDIT
 * =================================================================================
 * Refresh Template
 * Required data are as follows
 *
 *		Model             string
 *		MethodName        string
 *		RefreshLogic      string
 *		PossibleTypes     string
 *		ConditionalObjectFieldsWithNull string
 * ================================================================================= */

type GETProductsProductidApisApiidStagesStageidUsageplansResponse struct {
	Total             types.Int64  `tfsdk:"total"`
	Raterps           types.Int32  `tfsdk:"rate_rps"`
	Quotacondition    types.String `tfsdk:"quota_condition"`
	Monthquotarequest types.Int64  `tfsdk:"month_quota_request"`
	Dayquotarequest   types.Int64  `tfsdk:"day_quota_request"`
	Content           types.List   `tfsdk:"content"`
}

func ConvertToFrameworkTypes_GETProductsProductidApisApiidStagesStageidUsageplans(ctx context.Context, data map[string]interface{}) (*GETProductsProductidApisApiidStagesStageidUsageplansResponse, error) {
	var dto GETProductsProductidApisApiidStagesStageidUsageplansResponse

	if data["total"] != nil {
		dto.Total = types.Int64Value(data["total"].(int64))
	}

	if data["rate_rps"] != nil {
		dto.Raterps = types.Int32Value(data["rate_rps"].(int32))
	}

	if data["quota_condition"] != nil {
		dto.Quotacondition = types.StringValue(data["quota_condition"].(string))
	}

	if data["month_quota_request"] != nil {
		dto.Monthquotarequest = types.Int64Value(data["month_quota_request"].(int64))
	}

	if data["day_quota_request"] != nil {
		dto.Dayquotarequest = types.Int64Value(data["day_quota_request"].(int64))
	}

	if data["content"] != nil {
		tempContent := data["content"].([]interface{})
		dto.Content = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"usage_plan_name":         types.StringType,
			"usage_plan_id":           types.StringType,
			"usage_plan_description":  types.StringType,
			"tenant_id":               types.StringType,
			"rate_rps":                types.Int32Type,
			"quota_condition":         types.StringType,
			"permission":              types.StringType,
			"month_quota_request":     types.Int64Type,
			"modifier":                types.StringType,
			"domain_code":             types.StringType,
			"disabled":                types.BoolType,
			"day_quota_request":       types.Int64Type,
			"associated_stages_count": types.Int64Type,
			"action_name":             types.StringType,
		},
		}}.ElementType(), tempContent)
	}

	return &dto, nil
}

func convertToObject_GETProductsProductidApisApiidStagesStageidUsageplans(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETProductsProductidApisApiidStagesStageidUsageplans(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {

			switch fieldType {
			case types.StringType:
				attrValues[field] = types.StringNull()
			case types.Int64Type:
				attrValues[field] = types.Int64Null()
			case types.BoolType:
				attrValues[field] = types.BoolNull()
			}
		}
	}

	r, diag := types.ObjectValue(attrTypes, attrValues)
	if diag.HasError() {
		return types.Object{}, fmt.Errorf("error from converting object: %v", diag)
	}

	// OK
	return r, nil
}

func convertValueToAttr_GETProductsProductidApisApiidStagesStageidUsageplans(value interface{}) (attr.Value, error) {
	switch v := value.(type) {
	case string:
		return types.StringValue(v), nil
	case int32:
		return types.Int32Value(v), nil
	case int64:
		return types.Int64Value(v), nil
	case float64:
		return types.Float64Value(v), nil
	case bool:
		return types.BoolValue(v), nil
	case nil:
		return types.StringNull(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
