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

type PrimitiveGETProductsRequest struct {
	HasDeployedStage                     bool   `json:"hasDeployedStage"`
	HasStageNotAssociatedWithUsagePlanId string `json:"hasStageNotAssociatedWithUsagePlanId"`
	IsPublished                          bool   `json:"isPublished"`
	Limit                                int64  `json:"limit"`
	Offset                               int64  `json:"offset"`
	ProductName                          string `json:"productName"`
	SubscriptionCode                     string `json:"subscriptionCode"`
}

type StringifiedGETProductsRequest struct {
	HasDeployedStage                     string `json:"hasDeployedStage"`
	HasStageNotAssociatedWithUsagePlanId string `json:"hasStageNotAssociatedWithUsagePlanId"`
	IsPublished                          string `json:"isPublished"`
	Limit                                string `json:"limit"`
	Offset                               string `json:"offset"`
	ProductName                          string `json:"productName"`
	SubscriptionCode                     string `json:"subscriptionCode"`
}

func (n *NClient) GETProducts(ctx context.Context, primitiveReq *PrimitiveGETProductsRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.HasDeployedStage != "" {
		query["hasDeployedStage"] = r.HasDeployedStage
	}

	if r.HasStageNotAssociatedWithUsagePlanId != "" {
		query["hasStageNotAssociatedWithUsagePlanId"] = r.HasStageNotAssociatedWithUsagePlanId
	}

	if r.IsPublished != "" {
		query["isPublished"] = r.IsPublished
	}

	if r.Limit != "" {
		query["limit"] = r.Limit
	}

	if r.Offset != "" {
		query["offset"] = r.Offset
	}

	if r.ProductName != "" {
		query["productName"] = r.ProductName
	}

	if r.SubscriptionCode != "" {
		query["subscriptionCode"] = r.SubscriptionCode
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products"

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

func (n *NClient) GETProducts_TF(ctx context.Context, r *PrimitiveGETProductsRequest) (*GETProductsResponse, error) {
	t, err := n.GETProducts(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProducts(context.TODO(), t)
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

type GETProductsResponse struct {
	ProductPage  types.Object `tfsdk:"product_page"`
	Initialcount types.Int64  `tfsdk:"initial_count"`
}

func ConvertToFrameworkTypes_GETProducts(ctx context.Context, data map[string]interface{}) (*GETProductsResponse, error) {
	var dto GETProductsResponse

	if data["product_page"] != nil {
		tempProductPage := data["product_page"].(map[string]interface{})

		allFields := []string{
			"total",
			"content",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempProductPage[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempProductPage, err := convertToObject_GETProducts(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.ProductPage = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"total": types.Int64Type,

			"content": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

				"tenant_id":           types.StringType,
				"subscription_code":   types.StringType,
				"product_name":        types.StringType,
				"product_id":          types.StringType,
				"product_description": types.StringType,
				"permission":          types.StringType,
				"modifier":            types.StringType,
				"mod_time":            types.StringType,
				"is_published":        types.BoolType,
				"is_deleted":          types.BoolType,
				"invoke_id":           types.StringType,
				"has_deployed_stage":  types.BoolType,
				"domain_code":         types.StringType,
				"disabled":            types.BoolType,
				"action_name":         types.StringType,
			},
			}},
		}}.AttributeTypes(), convertedTempProductPage)
	}

	if data["initial_count"] != nil {
		dto.Initialcount = types.Int64Value(data["initial_count"].(int64))
	}

	return &dto, nil
}

func convertToObject_GETProducts(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"total": types.Int64Type,

		"content": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"tenant_id":           types.StringType,
			"subscription_code":   types.StringType,
			"product_name":        types.StringType,
			"product_id":          types.StringType,
			"product_description": types.StringType,
			"permission":          types.StringType,
			"modifier":            types.StringType,
			"mod_time":            types.StringType,
			"is_published":        types.BoolType,
			"is_deleted":          types.BoolType,
			"invoke_id":           types.StringType,
			"has_deployed_stage":  types.BoolType,
			"domain_code":         types.StringType,
			"disabled":            types.BoolType,
			"action_name":         types.StringType,
		},
		}},
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			if field == "content" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"tenant_id":           types.StringType,
					"subscription_code":   types.StringType,
					"product_name":        types.StringType,
					"product_id":          types.StringType,
					"product_description": types.StringType,
					"permission":          types.StringType,
					"modifier":            types.StringType,
					"mod_time":            types.StringType,
					"is_published":        types.BoolType,
					"is_deleted":          types.BoolType,
					"invoke_id":           types.StringType,
					"has_deployed_stage":  types.BoolType,
					"domain_code":         types.StringType,
					"disabled":            types.BoolType,
					"action_name":         types.StringType,
				}).Type(ctx))
				attrValues[field] = listV
				continue
			}

			attrValue, err := convertValueToAttr_GETProducts(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {

			if field == "content" {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"tenant_id":           types.StringType,
					"subscription_code":   types.StringType,
					"product_name":        types.StringType,
					"product_id":          types.StringType,
					"product_description": types.StringType,
					"permission":          types.StringType,
					"modifier":            types.StringType,
					"mod_time":            types.StringType,
					"is_published":        types.BoolType,
					"is_deleted":          types.BoolType,
					"invoke_id":           types.StringType,
					"has_deployed_stage":  types.BoolType,
					"domain_code":         types.StringType,
					"disabled":            types.BoolType,
					"action_name":         types.StringType,
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

func convertValueToAttr_GETProducts(value interface{}) (attr.Value, error) {
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
