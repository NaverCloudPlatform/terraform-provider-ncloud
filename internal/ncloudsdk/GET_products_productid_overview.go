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

type PrimitiveGETProductsProductidOverviewRequest struct {
	Productid string `json:"product-id"`
}

type StringifiedGETProductsProductidOverviewRequest struct {
	Productid string `json:"product-id"`
}

func (n *NClient) GETProductsProductidOverview(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidOverviewRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "overview"

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

func (n *NClient) GETProductsProductidOverview_TF(ctx context.Context, r *PrimitiveGETProductsProductidOverviewRequest) (*GETProductsProductidOverviewResponse, error) {
	t, err := n.GETProductsProductidOverview(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidOverview(context.TODO(), t)
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

type GETProductsProductidOverviewResponse struct {
	Apis            types.List   `tfsdk:"apis"`
	ApiKeyCountInfo types.Object `tfsdk:"api_key_count_info"`
}

func ConvertToFrameworkTypes_GETProductsProductidOverview(ctx context.Context, data map[string]interface{}) (*GETProductsProductidOverviewResponse, error) {
	var dto GETProductsProductidOverviewResponse

	if data["apis"] != nil {
		tempApis := data["apis"].([]interface{})
		dto.Apis = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"stages": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"stage_name":                   types.StringType,
				"stage_id":                     types.StringType,
				"host":                         types.StringType,
				"deployed_stage_deployment_no": types.Int64Type,
			},
			}},
			"methods_count": types.Int64Type,
			"domain_code":   types.StringType,
			"api_name":      types.StringType,
			"api_id":        types.StringType,
		},
		}}.ElementType(), tempApis)
	}
	if data["api_key_count_info"] != nil {
		tempApiKeyCountInfo := data["api_key_count_info"].(map[string]interface{})

		allFields := []string{
			"total",
			"request",
			"rejected",
			"denied",
			"accepted",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempApiKeyCountInfo[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempApiKeyCountInfo, err := convertToObject_GETProductsProductidOverview(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.ApiKeyCountInfo = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"total":    types.Int64Type,
			"request":  types.Int64Type,
			"rejected": types.Int64Type,
			"denied":   types.Int64Type,
			"accepted": types.Int64Type,
		}}.AttributeTypes(), convertedTempApiKeyCountInfo)
	}

	return &dto, nil
}

func convertToObject_GETProductsProductidOverview(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"total":    types.Int64Type,
		"request":  types.Int64Type,
		"rejected": types.Int64Type,
		"denied":   types.Int64Type,
		"accepted": types.Int64Type,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETProductsProductidOverview(value)
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

func convertValueToAttr_GETProductsProductidOverview(value interface{}) (attr.Value, error) {
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
