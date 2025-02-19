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

type PrimitiveGETProductsProductidSubscriptiondetailRequest struct {
	Productid string `json:"product-id"`
	ApiId     string `json:"apiId"`
	ApiKeyId  string `json:"apiKeyId"`
	Limit     int64  `json:"limit"`
	Offset    int64  `json:"offset"`
	StageId   string `json:"stageId"`
}

type StringifiedGETProductsProductidSubscriptiondetailRequest struct {
	Productid string `json:"product-id"`
	ApiId     string `json:"apiId"`
	ApiKeyId  string `json:"apiKeyId"`
	Limit     string `json:"limit"`
	Offset    string `json:"offset"`
	StageId   string `json:"stageId"`
}

func (n *NClient) GETProductsProductidSubscriptiondetail(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidSubscriptiondetailRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.ApiId != "" {
		query["apiId"] = r.ApiId
	}

	query["apiKeyId"] = r.ApiKeyId

	if r.Limit != "" {
		query["limit"] = r.Limit
	}

	if r.Offset != "" {
		query["offset"] = r.Offset
	}

	if r.StageId != "" {
		query["stageId"] = r.StageId
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "subscription-detail"

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

func (n *NClient) GETProductsProductidSubscriptiondetail_TF(ctx context.Context, r *PrimitiveGETProductsProductidSubscriptiondetailRequest) (*GETProductsProductidSubscriptiondetailResponse, error) {
	t, err := n.GETProductsProductidSubscriptiondetail(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidSubscriptiondetail(context.TODO(), t)
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

type GETProductsProductidSubscriptiondetailResponse struct {
	ApiKeyRelations types.Object `tfsdk:"api_key_relations"`
}

func ConvertToFrameworkTypes_GETProductsProductidSubscriptiondetail(ctx context.Context, data map[string]interface{}) (*GETProductsProductidSubscriptiondetailResponse, error) {
	var dto GETProductsProductidSubscriptiondetailResponse

	if data["api_key_relations"] != nil {
		tempApiKeyRelations := data["api_key_relations"].(map[string]interface{})

		allFields := []string{
			"total",
			"content",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempApiKeyRelations[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempApiKeyRelations, err := convertToObject_GETProductsProductidSubscriptiondetail(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.ApiKeyRelations = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"total": types.Int64Type,

			"content": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

				"usage_plan_name":  types.StringType,
				"usage_plan_id":    types.StringType,
				"stage_name":       types.StringType,
				"stage_id":         types.StringType,
				"reg_time":         types.StringType,
				"monthly_usage":    types.Int64Type,
				"month_call_count": types.Int64Type,
				"mod_time":         types.StringType,
				"domain_code":      types.StringType,
				"day_call_count":   types.Int64Type,
				"daily_usage":      types.Int64Type,
				"api_name":         types.StringType,
				"api_key_name":     types.StringType,
				"api_key_id":       types.StringType,
				"api_id":           types.StringType,
			},
			}},
		}}.AttributeTypes(), convertedTempApiKeyRelations)
	}

	return &dto, nil
}

func convertToObject_GETProductsProductidSubscriptiondetail(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"total": types.Int64Type,

		"content": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"usage_plan_name":  types.StringType,
			"usage_plan_id":    types.StringType,
			"stage_name":       types.StringType,
			"stage_id":         types.StringType,
			"reg_time":         types.StringType,
			"monthly_usage":    types.Int64Type,
			"month_call_count": types.Int64Type,
			"mod_time":         types.StringType,
			"domain_code":      types.StringType,
			"day_call_count":   types.Int64Type,
			"daily_usage":      types.Int64Type,
			"api_name":         types.StringType,
			"api_key_name":     types.StringType,
			"api_key_id":       types.StringType,
			"api_id":           types.StringType,
		},
		}},
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			if field == "content" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"usage_plan_name":  types.StringType,
					"usage_plan_id":    types.StringType,
					"stage_name":       types.StringType,
					"stage_id":         types.StringType,
					"reg_time":         types.StringType,
					"monthly_usage":    types.Int64Type,
					"month_call_count": types.Int64Type,
					"mod_time":         types.StringType,
					"domain_code":      types.StringType,
					"day_call_count":   types.Int64Type,
					"daily_usage":      types.Int64Type,
					"api_name":         types.StringType,
					"api_key_name":     types.StringType,
					"api_key_id":       types.StringType,
					"api_id":           types.StringType,
				}).Type(ctx))
				attrValues[field] = listV
				continue
			}

			attrValue, err := convertValueToAttr_GETProductsProductidSubscriptiondetail(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {

			if field == "content" {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"usage_plan_name":  types.StringType,
					"usage_plan_id":    types.StringType,
					"stage_name":       types.StringType,
					"stage_id":         types.StringType,
					"reg_time":         types.StringType,
					"monthly_usage":    types.Int64Type,
					"month_call_count": types.Int64Type,
					"mod_time":         types.StringType,
					"domain_code":      types.StringType,
					"day_call_count":   types.Int64Type,
					"daily_usage":      types.Int64Type,
					"api_name":         types.StringType,
					"api_key_name":     types.StringType,
					"api_key_id":       types.StringType,
					"api_id":           types.StringType,
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

func convertValueToAttr_GETProductsProductidSubscriptiondetail(value interface{}) (attr.Value, error) {
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
