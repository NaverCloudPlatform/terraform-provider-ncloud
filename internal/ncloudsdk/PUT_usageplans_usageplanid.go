
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

type PrimitivePUTUsageplansUsageplanidRequest struct {
    Usageplanid string `json:"usage-plan-id"`
UsagePlanName string `json:"usagePlanName"`
UsagePlanDescription string `json:"usagePlanDescription"`
RateRps int32 `json:"rateRps"`
QuotaCondition string `json:"quotaCondition"`
MonthQuotaRequest int64 `json:"monthQuotaRequest"`
DayQuotaRequest int64 `json:"dayQuotaRequest"`

}

type StringifiedPUTUsageplansUsageplanidRequest struct {
	Usageplanid string `json:"usage-plan-id"`
UsagePlanName string `json:"usagePlanName"`
UsagePlanDescription string `json:"usagePlanDescription"`
RateRps string `json:"rateRps"`
QuotaCondition string `json:"quotaCondition"`
MonthQuotaRequest string `json:"monthQuotaRequest"`
DayQuotaRequest string `json:"dayQuotaRequest"`

}

func (n *NClient) PUTUsageplansUsageplanid(ctx context.Context, primitiveReq *PrimitivePUTUsageplansUsageplanidRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["usagePlanName"] = r.UsagePlanName

			if r.UsagePlanDescription != "" {
				initBody["usagePlanDescription"] = r.UsagePlanDescription
			}

			if r.RateRps != "" {
				initBody["rateRps"] = r.RateRps
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


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"usage-plans"+"/"+ClearDoubleQuote(r.Usageplanid)

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

func (n *NClient) PUTUsageplansUsageplanid_TF(ctx context.Context, r *PrimitivePUTUsageplansUsageplanidRequest) (*PUTUsageplansUsageplanidResponse, error) {
	t, err := n.PUTUsageplansUsageplanid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_PUTUsageplansUsageplanid(context.TODO(), t)
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

type PUTUsageplansUsageplanidResponse struct {
    UsagePlan         types.Object `tfsdk:"usage_plan"`

}

func ConvertToFrameworkTypes_PUTUsageplansUsageplanid(ctx context.Context, data map[string]interface{}) (*PUTUsageplansUsageplanidResponse, error) {
	var dto PUTUsageplansUsageplanidResponse

    
			if data["usage_plan"] != nil {
				tempUsagePlan := data["usage_plan"].(map[string]interface{})

				allFields := []string{
					"usage_plan_name",
"usage_plan_id",
"usage_plan_description",
"tenant_id",
"rate_rps",
"quota_condition",
"permission",
"month_quota_request",
"modifier",
"domain_code",
"disabled",
"day_quota_request",
"associated_stages_count",
"action_name",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempUsagePlan[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempUsagePlan, err := convertToObject_PUTUsageplansUsageplanid(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.UsagePlan = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"usage_plan_name": types.StringType,
"usage_plan_id": types.StringType,
"usage_plan_description": types.StringType,
"tenant_id": types.StringType,
"rate_rps": types.Int32Type,
"quota_condition": types.StringType,
"permission": types.StringType,
"month_quota_request": types.Int64Type,
"modifier": types.StringType,
"domain_code": types.StringType,
"disabled": types.BoolType,
"day_quota_request": types.Int64Type,
"associated_stages_count": types.Int64Type,
"action_name": types.StringType,

				}}.AttributeTypes(), convertedTempUsagePlan)
			}


	return &dto, nil
}

func convertToObject_PUTUsageplansUsageplanid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "usage_plan_name": types.StringType,
"usage_plan_id": types.StringType,
"usage_plan_description": types.StringType,
"tenant_id": types.StringType,
"rate_rps": types.Int32Type,
"quota_condition": types.StringType,
"permission": types.StringType,
"month_quota_request": types.Int64Type,
"modifier": types.StringType,
"domain_code": types.StringType,
"disabled": types.BoolType,
"day_quota_request": types.Int64Type,
"associated_stages_count": types.Int64Type,
"action_name": types.StringType,


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_PUTUsageplansUsageplanid(value)
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

func convertValueToAttr_PUTUsageplansUsageplanid(value interface{}) (attr.Value, error) {
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

