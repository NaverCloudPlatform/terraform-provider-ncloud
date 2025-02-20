/* =================================================================================
 * NCLOUD SDK LAYER FOR TERRAFORM CODEGEN - DO NOT EDIT
 * =================================================================================
 * Refresh Template
 * Required data are as follows
 *
 *		MethodName             string
 *		RequestQueryParameters string
 *		RequestBodyParameters  string
 *		FunctionName           string
 *		Query                  string
 *		Body                   string
 *		Path                   string
 *		Method                 string
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

type PATCHProductsProductidRequestQuery struct {
	Productid *string `json:"product-id,omitempty"`
}

type PATCHProductsProductidRequestBody struct {
	SubscriptionCode *string `json:"subscriptionCode,,omitempty"`
	Description      *string `json:"description,,omitempty"`
	ProductName      *string `json:"productName,,omitempty"`
}

func (n *NClient) PATCHProductsProductid(ctx context.Context, q *PATCHProductsProductidRequestQuery, b *PATCHProductsProductidRequestBody) (map[string]interface{}, error) {

	query := map[string]string{}
	initBody := map[string]string{}

	initBody["subscriptionCode"] = *b.SubscriptionCode

	if b.Description != nil {
		initBody["description"] = *b.Description
	}
	initBody["productName"] = *b.ProductName

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(*q.Productid)

	response, err := n.MakeRequestWithContext(ctx, "PATCH", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
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

type PATCHProductsProductidResponse struct {
}

func ConvertToFrameworkTypes_PATCHProductsProductid(ctx context.Context, data map[string]interface{}) (*PATCHProductsProductidResponse, error) {
	var dto PATCHProductsProductidResponse

	return &dto, nil
}

func convertToObject_PATCHProductsProductid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_PATCHProductsProductid(value)
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

func convertValueToAttr_PATCHProductsProductid(value interface{}) (attr.Value, error) {
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
