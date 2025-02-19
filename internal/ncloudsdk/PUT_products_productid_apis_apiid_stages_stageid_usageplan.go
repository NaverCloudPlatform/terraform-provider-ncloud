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

type PrimitivePUTProductsProductidApisApiidStagesStageidUsageplanRequest struct {
	Productid         string `json:"product-id"`
	Apiid             string `json:"api-id"`
	Stageid           string `json:"stage-id"`
	QuotaCondition    string `json:"quotaCondition"`
	MonthQuotaRequest int64  `json:"monthQuotaRequest"`
	DayQuotaRequest   int64  `json:"dayQuotaRequest"`
	RateRps           int32  `json:"rateRps"`
}

type StringifiedPUTProductsProductidApisApiidStagesStageidUsageplanRequest struct {
	Productid         string `json:"product-id"`
	Apiid             string `json:"api-id"`
	Stageid           string `json:"stage-id"`
	QuotaCondition    string `json:"quotaCondition"`
	MonthQuotaRequest string `json:"monthQuotaRequest"`
	DayQuotaRequest   string `json:"dayQuotaRequest"`
	RateRps           string `json:"rateRps"`
}

func (n *NClient) PUTProductsProductidApisApiidStagesStageidUsageplan(ctx context.Context, primitiveReq *PrimitivePUTProductsProductidApisApiidStagesStageidUsageplanRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.QuotaCondition != "" {
		initBody["quotaCondition"] = r.QuotaCondition
	}

	if r.MonthQuotaRequest != "" {
		initBody["monthQuotaRequest"] = r.MonthQuotaRequest
	}

	if r.DayQuotaRequest != "" {
		initBody["dayQuotaRequest"] = r.DayQuotaRequest
	}

	if r.RateRps != "" {
		initBody["rateRps"] = r.RateRps
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "apis" + "/" + ClearDoubleQuote(r.Apiid) + "/" + "stages" + "/" + ClearDoubleQuote(r.Stageid) + "/" + "usage-plan"

	response, err := n.MakeRequestWithContext(ctx, "PUT", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) PUTProductsProductidApisApiidStagesStageidUsageplan_TF(ctx context.Context, r *PrimitivePUTProductsProductidApisApiidStagesStageidUsageplanRequest) (*PUTProductsProductidApisApiidStagesStageidUsageplanResponse, error) {
	t, err := n.PUTProductsProductidApisApiidStagesStageidUsageplan(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_PUTProductsProductidApisApiidStagesStageidUsageplan(context.TODO(), t)
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

type PUTProductsProductidApisApiidStagesStageidUsageplanResponse struct {
	Stage types.Object `tfsdk:"stage"`
}

func ConvertToFrameworkTypes_PUTProductsProductidApisApiidStagesStageidUsageplan(ctx context.Context, data map[string]interface{}) (*PUTProductsProductidApisApiidStagesStageidUsageplanResponse, error) {
	var dto PUTProductsProductidApisApiidStagesStageidUsageplanResponse

	if data["stage"] != nil {
		tempStage := data["stage"].(map[string]interface{})

		allFields := []string{
			"minimum_compression_size",
			"enabled_content_encoding",
			"use_distribution_rate",
			"throttle_rps",
			"tenant_id",
			"status_code",
			"stage_name",
			"stage_id",
			"response",
			"reg_time",
			"rate_rps",
			"month_quota_request",
			"modifier",
			"is_maintenance",
			"ip_acl_list",
			"ip_acl_code",
			"endpoint_domain",
			"deployed_stage_deployment_no",
			"day_quota_request",
			"canary_throttle_rps",
			"canary_endpoint_domain",
			"canary_distribution_rate",
			"canary_deployment_no",
			"canary_deployment_description",
			"canary_deployed_time",
			"canary_conditions",
			"canary_cache_ttl_sec",
			"cache_ttl_sec",
			"api_id",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempStage[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempStage, err := convertToObject_PUTProductsProductidApisApiidStagesStageidUsageplan(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.Stage = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"minimum_compression_size":      types.Int32Type,
			"enabled_content_encoding":      types.BoolType,
			"use_distribution_rate":         types.BoolType,
			"throttle_rps":                  types.Int32Type,
			"tenant_id":                     types.StringType,
			"status_code":                   types.Int32Type,
			"stage_name":                    types.StringType,
			"stage_id":                      types.StringType,
			"response":                      types.StringType,
			"reg_time":                      types.StringType,
			"rate_rps":                      types.Int32Type,
			"month_quota_request":           types.Int64Type,
			"modifier":                      types.StringType,
			"is_maintenance":                types.BoolType,
			"ip_acl_list":                   types.StringType,
			"ip_acl_code":                   types.StringType,
			"endpoint_domain":               types.StringType,
			"deployed_stage_deployment_no":  types.Int64Type,
			"day_quota_request":             types.Int64Type,
			"canary_throttle_rps":           types.Int32Type,
			"canary_endpoint_domain":        types.StringType,
			"canary_distribution_rate":      types.Float64Type,
			"canary_deployment_no":          types.Int64Type,
			"canary_deployment_description": types.StringType,
			"canary_deployed_time":          types.StringType,

			"canary_conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

				"parameter_value": types.StringType,
				"parameter_name":  types.StringType,
				"parameter_code":  types.StringType,
			},
			}},
			"canary_cache_ttl_sec": types.Int32Type,
			"cache_ttl_sec":        types.Int32Type,
			"api_id":               types.StringType,
		}}.AttributeTypes(), convertedTempStage)
	}

	return &dto, nil
}

func convertToObject_PUTProductsProductidApisApiidStagesStageidUsageplan(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"minimum_compression_size":      types.Int32Type,
		"enabled_content_encoding":      types.BoolType,
		"use_distribution_rate":         types.BoolType,
		"throttle_rps":                  types.Int32Type,
		"tenant_id":                     types.StringType,
		"status_code":                   types.Int32Type,
		"stage_name":                    types.StringType,
		"stage_id":                      types.StringType,
		"response":                      types.StringType,
		"reg_time":                      types.StringType,
		"rate_rps":                      types.Int32Type,
		"month_quota_request":           types.Int64Type,
		"modifier":                      types.StringType,
		"is_maintenance":                types.BoolType,
		"ip_acl_list":                   types.StringType,
		"ip_acl_code":                   types.StringType,
		"endpoint_domain":               types.StringType,
		"deployed_stage_deployment_no":  types.Int64Type,
		"day_quota_request":             types.Int64Type,
		"canary_throttle_rps":           types.Int32Type,
		"canary_endpoint_domain":        types.StringType,
		"canary_distribution_rate":      types.Float64Type,
		"canary_deployment_no":          types.Int64Type,
		"canary_deployment_description": types.StringType,
		"canary_deployed_time":          types.StringType,

		"canary_conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"parameter_value": types.StringType,
			"parameter_name":  types.StringType,
			"parameter_code":  types.StringType,
		},
		}},
		"canary_cache_ttl_sec": types.Int32Type,
		"cache_ttl_sec":        types.Int32Type,
		"api_id":               types.StringType,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			if field == "canary_conditions" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"parameter_value": types.StringType,
					"parameter_name":  types.StringType,
					"parameter_code":  types.StringType,
				}).Type(ctx))
				attrValues[field] = listV
				continue
			}

			attrValue, err := convertValueToAttr_PUTProductsProductidApisApiidStagesStageidUsageplan(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {

			if field == "canary_conditions" {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"parameter_value": types.StringType,
					"parameter_name":  types.StringType,
					"parameter_code":  types.StringType,
				}).Type(ctx))
				attrValues[field] = listV
				continue
			}

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

func convertValueToAttr_PUTProductsProductidApisApiidStagesStageidUsageplan(value interface{}) (attr.Value, error) {
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
