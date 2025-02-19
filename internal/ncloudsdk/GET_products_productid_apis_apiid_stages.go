
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

type PrimitiveGETProductsProductidApisApiidStagesRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`

}

type StringifiedGETProductsProductidApisApiidStagesRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`

}

func (n *NClient) GETProductsProductidApisApiidStages(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidApisApiidStagesRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"stages"

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

func (n *NClient) GETProductsProductidApisApiidStages_TF(ctx context.Context, r *PrimitiveGETProductsProductidApisApiidStagesRequest) (*GETProductsProductidApisApiidStagesResponse, error) {
	t, err := n.GETProductsProductidApisApiidStages(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidApisApiidStages(context.TODO(), t)
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

type GETProductsProductidApisApiidStagesResponse struct {
    Stages         types.List `tfsdk:"stages"`
Host         types.String `tfsdk:"host"`

}

func ConvertToFrameworkTypes_GETProductsProductidApisApiidStages(ctx context.Context, data map[string]interface{}) (*GETProductsProductidApisApiidStagesResponse, error) {
	var dto GETProductsProductidApisApiidStagesResponse

    
				if data["stages"] != nil {
					tempStages := data["stages"].([]interface{})
					dto.Stages = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"stage_name": types.StringType,
"stage_id": types.StringType,

				"resource_dtos": types.ListType{ElemType:
					types.ObjectType{AttrTypes: map[string]attr.Type{
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

					},
				}},
"is_canary": types.BoolType,
"default_deployment_no": types.Int64Type,
"canary_deployment_no": types.Int64Type,

	},

					}}.ElementType(), tempStages)
				}
			if data["host"] != nil {
				dto.Host = types.StringValue(data["host"].(string))
			}


	return &dto, nil
}

func convertToObject_GETProductsProductidApisApiidStages(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_GETProductsProductidApisApiidStages(value)
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

func convertValueToAttr_GETProductsProductidApisApiidStages(value interface{}) (attr.Value, error) {
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

