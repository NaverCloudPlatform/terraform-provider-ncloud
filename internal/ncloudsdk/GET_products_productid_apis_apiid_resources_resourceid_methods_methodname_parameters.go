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

type PrimitiveGETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersRequest struct {
	Productid  string `json:"product-id"`
	Apiid      string `json:"api-id"`
	Resourceid string `json:"resource-id"`
	Methodname string `json:"method-name"`
}

type StringifiedGETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersRequest struct {
	Productid  string `json:"product-id"`
	Apiid      string `json:"api-id"`
	Resourceid string `json:"resource-id"`
	Methodname string `json:"method-name"`
}

func (n *NClient) GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "apis" + "/" + ClearDoubleQuote(r.Apiid) + "/" + "resources" + "/" + ClearDoubleQuote(r.Resourceid) + "/" + "methods" + "/" + ClearDoubleQuote(r.Methodname) + "/" + "parameters"

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

func (n *NClient) GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters_TF(ctx context.Context, r *PrimitiveGETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersRequest) (*GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersResponse, error) {
	t, err := n.GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(context.TODO(), t)
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

type GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersResponse struct {
	Usebodywhenformdata types.Bool `tfsdk:"use_body_when_form_data"`
	QueryStrings        types.List `tfsdk:"query_strings"`
	Isncptenant         types.Bool `tfsdk:"is_ncp_tenant"`
	Headers             types.List `tfsdk:"headers"`
	FormDatas           types.List `tfsdk:"form_datas"`
}

func ConvertToFrameworkTypes_GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(ctx context.Context, data map[string]interface{}) (*GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersResponse, error) {
	var dto GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParametersResponse

	if data["use_body_when_form_data"] != nil {
		dto.Usebodywhenformdata = types.BoolValue(data["use_body_when_form_data"].(bool))
	}

	if data["query_strings"] != nil {
		tempQueryStrings := data["query_strings"].([]interface{})
		dto.QueryStrings = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

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
		},
		}}.ElementType(), tempQueryStrings)
	}
	if data["is_ncp_tenant"] != nil {
		dto.Isncptenant = types.BoolValue(data["is_ncp_tenant"].(bool))
	}

	if data["headers"] != nil {
		tempHeaders := data["headers"].([]interface{})
		dto.Headers = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

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
		},
		}}.ElementType(), tempHeaders)
	}
	if data["form_datas"] != nil {
		tempFormDatas := data["form_datas"].([]interface{})
		dto.FormDatas = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

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
		},
		}}.ElementType(), tempFormDatas)
	}

	return &dto, nil
}

func convertToObject_GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(value)
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

func convertValueToAttr_GETProductsProductidApisApiidResourcesResourceidMethodsMethodnameParameters(value interface{}) (attr.Value, error) {
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
