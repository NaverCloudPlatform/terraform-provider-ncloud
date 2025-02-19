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

type PrimitiveGETProductsProductidRequest struct {
	Productid string `json:"product-id"`
}

type StringifiedGETProductsProductidRequest struct {
	Productid string `json:"product-id"`
}

func (n *NClient) GETProductsProductid(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid)

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

func (n *NClient) GETProductsProductid_TF(ctx context.Context, r *PrimitiveGETProductsProductidRequest) (*GETProductsProductidResponse, error) {
	t, err := n.GETProductsProductid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductid(context.TODO(), t)
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

type GETProductsProductidResponse struct {
	Product types.Object `tfsdk:"product"`
}

func ConvertToFrameworkTypes_GETProductsProductid(ctx context.Context, data map[string]interface{}) (*GETProductsProductidResponse, error) {
	var dto GETProductsProductidResponse

	if data["product"] != nil {
		tempProduct := data["product"].(map[string]interface{})

		allFields := []string{
			"tenant_id",
			"subscription_code",
			"product_name",
			"product_id",
			"product_description",
			"permission",
			"modifier",
			"mod_time",
			"is_published",
			"is_deleted",
			"invoke_id",
			"domain_code",
			"disabled",
			"action_name",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempProduct[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempProduct, err := convertToObject_GETProductsProductid(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.Product = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
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
			"domain_code":         types.StringType,
			"disabled":            types.BoolType,
			"action_name":         types.StringType,
		}}.AttributeTypes(), convertedTempProduct)
	}

	return &dto, nil
}

func convertToObject_GETProductsProductid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
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
		"domain_code":         types.StringType,
		"disabled":            types.BoolType,
		"action_name":         types.StringType,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETProductsProductid(value)
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

func convertValueToAttr_GETProductsProductid(value interface{}) (attr.Value, error) {
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
