
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

type PrimitivePOSTProductsProductidApisRequest struct {
    Productid string `json:"product-id"`
ApiName string `json:"apiName"`
ApiDescription string `json:"apiDescription"`

}

type StringifiedPOSTProductsProductidApisRequest struct {
	Productid string `json:"product-id"`
ApiName string `json:"apiName"`
ApiDescription string `json:"apiDescription"`

}

func (n *NClient) POSTProductsProductidApis(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["apiName"] = r.ApiName

			if r.ApiDescription != "" {
				initBody["apiDescription"] = r.ApiDescription
			}


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"

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

func (n *NClient) POSTProductsProductidApis_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisRequest) (*POSTProductsProductidApisResponse, error) {
	t, err := n.POSTProductsProductidApis(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApis(context.TODO(), t)
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

type POSTProductsProductidApisResponse struct {
    Api         types.Object `tfsdk:"api"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApis(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisResponse, error) {
	var dto POSTProductsProductidApisResponse

    
			if data["api"] != nil {
				tempApi := data["api"].(map[string]interface{})

				allFields := []string{
					"tenant_id",
"stages",
"product_id",
"permission",
"modifier",
"mod_time",
"is_deleted",
"domain_code",
"disabled",
"api_name",
"api_id",
"api_description",
"action_name",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempApi[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempApi, err := convertToObject_POSTProductsProductidApis(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.Api = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"tenant_id": types.StringType,

			"stages": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"stage_name": types.StringType,
"stage_id": types.StringType,
"is_published": types.BoolType,
"api_id": types.StringType,

	},
			}},
"product_id": types.StringType,
"permission": types.StringType,
"modifier": types.StringType,
"mod_time": types.StringType,
"is_deleted": types.BoolType,
"domain_code": types.StringType,
"disabled": types.BoolType,
"api_name": types.StringType,
"api_id": types.StringType,
"api_description": types.StringType,
"action_name": types.StringType,

				}}.AttributeTypes(), convertedTempApi)
			}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApis(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "tenant_id": types.StringType,

			"stages": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"stage_name": types.StringType,
"stage_id": types.StringType,
"is_published": types.BoolType,
"api_id": types.StringType,

	},
			}},
"product_id": types.StringType,
"permission": types.StringType,
"modifier": types.StringType,
"mod_time": types.StringType,
"is_deleted": types.BoolType,
"domain_code": types.StringType,
"disabled": types.BoolType,
"api_name": types.StringType,
"api_id": types.StringType,
"api_description": types.StringType,
"action_name": types.StringType,


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			
			if field == "stages" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"stage_name": types.StringType,
"stage_id": types.StringType,
"is_published": types.BoolType,
"api_id": types.StringType,

				}).Type(ctx))
				attrValues[field] = listV
				continue
			}


			attrValue, err := convertValueToAttr_POSTProductsProductidApis(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            
				if field == "stages" {
					listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
						"stage_name": types.StringType,
"stage_id": types.StringType,
"is_published": types.BoolType,
"api_id": types.StringType,

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

func convertValueToAttr_POSTProductsProductidApis(value interface{}) (attr.Value, error) {
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

