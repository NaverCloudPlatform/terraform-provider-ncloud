
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

type PrimitivePOSTProductsProductidApisApiidResourcesRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
ResourcePath string `json:"resourcePath"`
CorsMaxAge string `json:"corsMaxAge"`
CorsExposeHeaders string `json:"corsExposeHeaders"`
CorsAllowOrigin string `json:"corsAllowOrigin"`
CorsAllowMethods string `json:"corsAllowMethods"`
CorsAllowHeaders string `json:"corsAllowHeaders"`
CorsAllowCredentials string `json:"corsAllowCredentials"`

}

type StringifiedPOSTProductsProductidApisApiidResourcesRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
ResourcePath string `json:"resourcePath"`
CorsMaxAge string `json:"corsMaxAge"`
CorsExposeHeaders string `json:"corsExposeHeaders"`
CorsAllowOrigin string `json:"corsAllowOrigin"`
CorsAllowMethods string `json:"corsAllowMethods"`
CorsAllowHeaders string `json:"corsAllowHeaders"`
CorsAllowCredentials string `json:"corsAllowCredentials"`

}

func (n *NClient) POSTProductsProductidApisApiidResources(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidResourcesRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["resourcePath"] = r.ResourcePath

			if r.CorsMaxAge != "" {
				initBody["corsMaxAge"] = r.CorsMaxAge
			}

			if r.CorsExposeHeaders != "" {
				initBody["corsExposeHeaders"] = r.CorsExposeHeaders
			}

			if r.CorsAllowOrigin != "" {
				initBody["corsAllowOrigin"] = r.CorsAllowOrigin
			}

			if r.CorsAllowMethods != "" {
				initBody["corsAllowMethods"] = r.CorsAllowMethods
			}

			if r.CorsAllowHeaders != "" {
				initBody["corsAllowHeaders"] = r.CorsAllowHeaders
			}

			if r.CorsAllowCredentials != "" {
				initBody["corsAllowCredentials"] = r.CorsAllowCredentials
			}


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"resources"

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

func (n *NClient) POSTProductsProductidApisApiidResources_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidResourcesRequest) (*POSTProductsProductidApisApiidResourcesResponse, error) {
	t, err := n.POSTProductsProductidApisApiidResources(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidResources(context.TODO(), t)
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

type POSTProductsProductidApisApiidResourcesResponse struct {
    Resource         types.Object `tfsdk:"resource"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidResources(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidResourcesResponse, error) {
	var dto POSTProductsProductidApisApiidResourcesResponse

    
			if data["resource"] != nil {
				tempResource := data["resource"].(map[string]interface{})

				allFields := []string{
					"resource_path",
"resource_id",
"methods",
"cors_max_age",
"cors_expose_headers",
"cors_allow_origin",
"cors_allow_methods",
"cors_allow_headers",
"cors_allow_credentials",
"api_id",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempResource[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempResource, err := convertToObject_POSTProductsProductidApisApiidResources(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.Resource = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"resource_path": types.StringType,
"resource_id": types.StringType,

			"methods": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"method_name": types.StringType,
"method_code": types.StringType,

	},
			}},
"cors_max_age": types.StringType,
"cors_expose_headers": types.StringType,
"cors_allow_origin": types.StringType,
"cors_allow_methods": types.StringType,
"cors_allow_headers": types.StringType,
"cors_allow_credentials": types.StringType,
"api_id": types.StringType,

				}}.AttributeTypes(), convertedTempResource)
			}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidResources(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "resource_path": types.StringType,
"resource_id": types.StringType,

			"methods": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"method_name": types.StringType,
"method_code": types.StringType,

	},
			}},
"cors_max_age": types.StringType,
"cors_expose_headers": types.StringType,
"cors_allow_origin": types.StringType,
"cors_allow_methods": types.StringType,
"cors_allow_headers": types.StringType,
"cors_allow_credentials": types.StringType,
"api_id": types.StringType,


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			
			if field == "methods" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"method_name": types.StringType,
"method_code": types.StringType,

				}).Type(ctx))
				attrValues[field] = listV
				continue
			}


			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidResources(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            
				if field == "methods" {
					listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
						"method_name": types.StringType,
"method_code": types.StringType,

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

func convertValueToAttr_POSTProductsProductidApisApiidResources(value interface{}) (attr.Value, error) {
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

