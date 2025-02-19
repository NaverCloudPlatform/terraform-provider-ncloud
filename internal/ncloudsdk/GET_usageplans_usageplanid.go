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

type PrimitiveGETUsageplansUsageplanidRequest struct {
	Usageplanid string `json:"usage-plan-id"`
}

type StringifiedGETUsageplansUsageplanidRequest struct {
	Usageplanid string `json:"usage-plan-id"`
}

func (n *NClient) GETUsageplansUsageplanid(ctx context.Context, primitiveReq *PrimitiveGETUsageplansUsageplanidRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "usage-plans" + "/" + ClearDoubleQuote(r.Usageplanid)

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

func (n *NClient) GETUsageplansUsageplanid_TF(ctx context.Context, r *PrimitiveGETUsageplansUsageplanidRequest) (*GETUsageplansUsageplanidResponse, error) {
	t, err := n.GETUsageplansUsageplanid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETUsageplansUsageplanid(context.TODO(), t)
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

type GETUsageplansUsageplanidResponse struct {
	UsagePlan types.Object `tfsdk:"usage_plan"`
}

func ConvertToFrameworkTypes_GETUsageplansUsageplanid(ctx context.Context, data map[string]interface{}) (*GETUsageplansUsageplanidResponse, error) {
	var dto GETUsageplansUsageplanidResponse

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

		convertedTempUsagePlan, err := convertToObject_GETUsageplansUsageplanid(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.UsagePlan = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
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
		}}.AttributeTypes(), convertedTempUsagePlan)
	}

	return &dto, nil
}

func convertToObject_GETUsageplansUsageplanid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
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
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETUsageplansUsageplanid(value)
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

func convertValueToAttr_GETUsageplansUsageplanid(value interface{}) (attr.Value, error) {
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
