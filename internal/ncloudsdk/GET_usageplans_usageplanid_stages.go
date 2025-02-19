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

type PrimitiveGETUsageplansUsageplanidStagesRequest struct {
	Usageplanid string `json:"usage-plan-id"`
	Limit       int64  `json:"limit"`
	Name        string `json:"name"`
	Offset      int64  `json:"offset"`
}

type StringifiedGETUsageplansUsageplanidStagesRequest struct {
	Usageplanid string `json:"usage-plan-id"`
	Limit       string `json:"limit"`
	Name        string `json:"name"`
	Offset      string `json:"offset"`
}

func (n *NClient) GETUsageplansUsageplanidStages(ctx context.Context, primitiveReq *PrimitiveGETUsageplansUsageplanidStagesRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.Limit != "" {
		query["limit"] = r.Limit
	}

	if r.Name != "" {
		query["name"] = r.Name
	}

	if r.Offset != "" {
		query["offset"] = r.Offset
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "usage-plans" + "/" + ClearDoubleQuote(r.Usageplanid) + "/" + "stages"

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

func (n *NClient) GETUsageplansUsageplanidStages_TF(ctx context.Context, r *PrimitiveGETUsageplansUsageplanidStagesRequest) (*GETUsageplansUsageplanidStagesResponse, error) {
	t, err := n.GETUsageplansUsageplanidStages(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETUsageplansUsageplanidStages(context.TODO(), t)
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

type GETUsageplansUsageplanidStagesResponse struct {
	UsagePlan types.Object `tfsdk:"usage_plan"`
	Total     types.Int64  `tfsdk:"total"`
	Stages    types.List   `tfsdk:"stages"`
}

func ConvertToFrameworkTypes_GETUsageplansUsageplanidStages(ctx context.Context, data map[string]interface{}) (*GETUsageplansUsageplanidStagesResponse, error) {
	var dto GETUsageplansUsageplanidStagesResponse

	if data["usage_plan"] != nil {
		tempUsagePlan := data["usage_plan"].(map[string]interface{})

		allFields := []string{
			"rate_rps",
			"quota_condition",
			"month_quota_request",
			"domain_code",
			"day_quota_request",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempUsagePlan[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempUsagePlan, err := convertToObject_GETUsageplansUsageplanidStages(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.UsagePlan = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"rate_rps":            types.Int32Type,
			"quota_condition":     types.StringType,
			"month_quota_request": types.Int64Type,
			"domain_code":         types.StringType,
			"day_quota_request":   types.Int64Type,
		}}.AttributeTypes(), convertedTempUsagePlan)
	}

	if data["total"] != nil {
		dto.Total = types.Int64Value(data["total"].(int64))
	}

	if data["stages"] != nil {
		tempStages := data["stages"].([]interface{})
		dto.Stages = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"stage_name":   types.StringType,
			"stage_id":     types.StringType,
			"product_name": types.StringType,
			"product_id":   types.StringType,
			"api_name":     types.StringType,
			"api_id":       types.StringType,
		},
		}}.ElementType(), tempStages)
	}

	return &dto, nil
}

func convertToObject_GETUsageplansUsageplanidStages(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"rate_rps":            types.Int32Type,
		"quota_condition":     types.StringType,
		"month_quota_request": types.Int64Type,
		"domain_code":         types.StringType,
		"day_quota_request":   types.Int64Type,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETUsageplansUsageplanidStages(value)
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

func convertValueToAttr_GETUsageplansUsageplanidStages(value interface{}) (attr.Value, error) {
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
