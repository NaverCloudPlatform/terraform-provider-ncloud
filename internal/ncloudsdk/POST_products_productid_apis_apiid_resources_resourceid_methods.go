
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

type PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Resourceid string `json:"resource-id"`
MethodDescription string `json:"methodDescription"`
RequiredApiKey types.Object `json:"requiredApiKey"`
Authentication types.Object `json:"authentication"`
Validation types.Object `json:"validation"`
NcpEndPoint types.Object `json:"ncpEndPoint"`
MockEndPoint types.Object `json:"mockEndPoint"`
HttpEndPoint types.Object `json:"httpEndPoint"`
MethodName string `json:"methodName"`

}

type StringifiedPOSTProductsProductidApisApiidResourcesResourceidMethodsRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Resourceid string `json:"resource-id"`
MethodDescription string `json:"methodDescription"`
RequiredApiKey string `json:"requiredApiKey"`
Authentication string `json:"authentication"`
Validation string `json:"validation"`
NcpEndPoint string `json:"ncpEndPoint"`
MockEndPoint string `json:"mockEndPoint"`
HttpEndPoint string `json:"httpEndPoint"`
MethodName string `json:"methodName"`

}

func (n *NClient) POSTProductsProductidApisApiidResourcesResourceidMethods(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	
			if r.MethodDescription != "" {
				initBody["methodDescription"] = r.MethodDescription
			}
initBody["requiredApiKey"] = r.RequiredApiKey
initBody["authentication"] = r.Authentication
initBody["validation"] = r.Validation

			if r.NcpEndPoint != "" {
				initBody["ncpEndPoint"] = r.NcpEndPoint
			}

			if r.MockEndPoint != "" {
				initBody["mockEndPoint"] = r.MockEndPoint
			}

			if r.HttpEndPoint != "" {
				initBody["httpEndPoint"] = r.HttpEndPoint
			}
initBody["methodName"] = r.MethodName


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"resources"+"/"+ClearDoubleQuote(r.Resourceid)+"/"+"methods"

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

func (n *NClient) POSTProductsProductidApisApiidResourcesResourceidMethods_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidResourcesResourceidMethodsRequest) (*POSTProductsProductidApisApiidResourcesResourceidMethodsResponse, error) {
	t, err := n.POSTProductsProductidApisApiidResourcesResourceidMethods(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidResourcesResourceidMethods(context.TODO(), t)
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

type POSTProductsProductidApisApiidResourcesResourceidMethodsResponse struct {
    Methods         types.List `tfsdk:"methods"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidResourcesResourceidMethods(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidResourcesResourceidMethodsResponse, error) {
	var dto POSTProductsProductidApisApiidResourcesResourceidMethodsResponse

    
				if data["methods"] != nil {
					tempMethods := data["methods"].([]interface{})
					dto.Methods = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
			"validation": types.ObjectType{AttrTypes: map[string]attr.Type{
				"headers": types.ListType{ElemType: types.StringType},
"query_strings": types.ListType{ElemType: types.StringType},
"type": types.StringType,

			}},

			"required_api_key": types.ObjectType{AttrTypes: map[string]attr.Type{
				"required": types.BoolType,

			}},

			"ncp_end_point": types.ObjectType{AttrTypes: map[string]attr.Type{
				"url": types.StringType,
"method": types.StringType,
"stream": types.BoolType,
"action_name": types.StringType,
"action_id": types.StringType,
"region": types.StringType,
"service": types.StringType,

			}},

			"mock_end_point": types.ObjectType{AttrTypes: map[string]attr.Type{
				
				"headers": types.ObjectType{AttrTypes: map[string]attr.Type{
				}},
"response": types.StringType,
"http_status": types.Int32Type,

			}},

			"http_end_point": types.ObjectType{AttrTypes: map[string]attr.Type{
				"url": types.StringType,
"method": types.StringType,
"stream": types.BoolType,

			}},

			"authentication": types.ObjectType{AttrTypes: map[string]attr.Type{
				"authorizer_id": types.StringType,
"platform": types.StringType,

			}},

		"use_body_when_form_data": types.BoolType,
"tenant_id": types.StringType,
"resource_path": types.StringType,
"resource_id": types.StringType,
"produces": types.StringType,
"modifier": types.StringType,
"method_name": types.StringType,
"method_description": types.StringType,
"consumers": types.StringType,
"api_id": types.StringType,

	},

					}}.ElementType(), tempMethods)
				}

	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidResourcesResourceidMethods(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidResourcesResourceidMethods(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidResourcesResourceidMethods(value interface{}) (attr.Value, error) {
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

