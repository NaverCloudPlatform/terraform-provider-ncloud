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

type PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersRequest struct {
	Productid            string `json:"product-id"`
	Apiid                string `json:"api-id"`
	Resourceid           string `json:"resource-id"`
	Methodname           string `json:"method-name"`
	Statuscode           string `json:"status-code"`
	ParameterType        string `json:"parameterType"`
	IsArray              bool   `json:"isArray"`
	ParameterDescription string `json:"parameterDescription"`
	ParameterName        string `json:"parameterName"`
}

type StringifiedPOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersRequest struct {
	Productid            string `json:"product-id"`
	Apiid                string `json:"api-id"`
	Resourceid           string `json:"resource-id"`
	Methodname           string `json:"method-name"`
	Statuscode           string `json:"status-code"`
	ParameterType        string `json:"parameterType"`
	IsArray              string `json:"isArray"`
	ParameterDescription string `json:"parameterDescription"`
	ParameterName        string `json:"parameterName"`
}

func (n *NClient) POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.ParameterType != "" {
		initBody["parameterType"] = r.ParameterType
	}

	if r.IsArray != "" {
		initBody["isArray"] = r.IsArray
	}

	if r.ParameterDescription != "" {
		initBody["parameterDescription"] = r.ParameterDescription
	}
	initBody["parameterName"] = r.ParameterName

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "apis" + "/" + ClearDoubleQuote(r.Apiid) + "/" + "resources" + "/" + ClearDoubleQuote(r.Resourceid) + "/" + "methods" + "/" + ClearDoubleQuote(r.Methodname) + "/" + "responses" + "/" + ClearDoubleQuote(r.Statuscode) + "/" + "parameters"

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

func (n *NClient) POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersRequest) (*POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersResponse, error) {
	t, err := n.POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(context.TODO(), t)
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

type POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersResponse struct {
	MethodParameterDto types.Object `tfsdk:"method_parameter_dto"`
}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersResponse, error) {
	var dto POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParametersResponse

	if data["method_parameter_dto"] != nil {
		tempMethodParameterDto := data["method_parameter_dto"].(map[string]interface{})

		allFields := []string{
			"tenant_id",
			"status_code",
			"resource_id",
			"parameter_type",
			"parameter_no",
			"parameter_name",
			"parameter_description",
			"parameter_code",
			"modifier",
			"method_code",
			"is_required",
			"is_logged",
			"is_array",
			"api_id",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempMethodParameterDto[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempMethodParameterDto, err := convertToObject_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.MethodParameterDto = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"tenant_id":             types.StringType,
			"status_code":           types.Int32Type,
			"resource_id":           types.StringType,
			"parameter_type":        types.StringType,
			"parameter_no":          types.Int64Type,
			"parameter_name":        types.StringType,
			"parameter_description": types.StringType,
			"parameter_code":        types.StringType,
			"modifier":              types.StringType,
			"method_code":           types.StringType,
			"is_required":           types.BoolType,
			"is_logged":             types.BoolType,
			"is_array":              types.BoolType,
			"api_id":                types.StringType,
		}}.AttributeTypes(), convertedTempMethodParameterDto)
	}

	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"tenant_id":             types.StringType,
		"status_code":           types.Int32Type,
		"resource_id":           types.StringType,
		"parameter_type":        types.StringType,
		"parameter_no":          types.Int64Type,
		"parameter_name":        types.StringType,
		"parameter_description": types.StringType,
		"parameter_code":        types.StringType,
		"modifier":              types.StringType,
		"method_code":           types.StringType,
		"is_required":           types.BoolType,
		"is_logged":             types.BoolType,
		"is_array":              types.BoolType,
		"api_id":                types.StringType,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidResourcesResourceidMethodsMethodnameResponsesStatuscodeParameters(value interface{}) (attr.Value, error) {
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
