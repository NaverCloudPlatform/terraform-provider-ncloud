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

type PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesRequest struct {
	Productid           string `json:"product-id"`
	Apiid               string `json:"api-id"`
	Resourceid          string `json:"resource-id"`
	Methodname          string `json:"method-name"`
	ResponseDescription string `json:"responseDescription"`
	StatusCode          int32  `json:"statusCode"`
}

type StringifiedPOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesRequest struct {
	Productid           string `json:"product-id"`
	Apiid               string `json:"api-id"`
	Resourceid          string `json:"resource-id"`
	Methodname          string `json:"method-name"`
	ResponseDescription string `json:"responseDescription"`
	StatusCode          string `json:"statusCode"`
}

func (n *NClient) POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.ResponseDescription != "" {
		initBody["responseDescription"] = r.ResponseDescription
	}
	initBody["statusCode"] = r.StatusCode

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "apis" + "/" + ClearDoubleQuote(r.Apiid) + "/" + "resources" + "/" + ClearDoubleQuote(r.Resourceid) + "/" + "methods" + "/" + ClearDoubleQuote(r.Methodname) + "/" + "responses"

	response, err := n.MakeRequestWithContext(ctx, "POST", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesRequest) (*POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesResponse, error) {
	t, err := n.POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(context.TODO(), t)
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

type POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesResponse struct {
	MethodResponseDto types.Object `tfsdk:"method_response_dto"`
}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesResponse, error) {
	var dto POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesResponse

	if data["method_response_dto"] != nil {
		tempMethodResponseDto := data["method_response_dto"].(map[string]interface{})

		allFields := []string{
			"tenant_id",
			"status_code",
			"response_description",
			"resource_id",
			"modifier",
			"method_code",
			"api_id",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempMethodResponseDto[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempMethodResponseDto, err := convertToObject_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.MethodResponseDto = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"tenant_id":            types.StringType,
			"status_code":          types.Int32Type,
			"response_description": types.StringType,
			"resource_id":          types.StringType,
			"modifier":             types.StringType,
			"method_code":          types.StringType,
			"api_id":               types.StringType,
		}}.AttributeTypes(), convertedTempMethodResponseDto)
	}

	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"tenant_id":            types.StringType,
		"status_code":          types.Int32Type,
		"response_description": types.StringType,
		"resource_id":          types.StringType,
		"modifier":             types.StringType,
		"method_code":          types.StringType,
		"api_id":               types.StringType,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponses(value interface{}) (attr.Value, error) {
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
