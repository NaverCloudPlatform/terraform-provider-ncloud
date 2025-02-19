
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

type PrimitivePOSTProductsProductidApisApiidResponsesResponsetypeRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Responsetype string `json:"response-type"`
StatusCode int32 `json:"statusCode"`

}

type StringifiedPOSTProductsProductidApisApiidResponsesResponsetypeRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Responsetype string `json:"response-type"`
StatusCode string `json:"statusCode"`

}

func (n *NClient) POSTProductsProductidApisApiidResponsesResponsetype(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidResponsesResponsetypeRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	
			if r.StatusCode != "" {
				initBody["statusCode"] = r.StatusCode
			}


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"responses"+"/"+ClearDoubleQuote(r.Responsetype)

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

func (n *NClient) POSTProductsProductidApisApiidResponsesResponsetype_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidResponsesResponsetypeRequest) (*POSTProductsProductidApisApiidResponsesResponsetypeResponse, error) {
	t, err := n.POSTProductsProductidApisApiidResponsesResponsetype(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidResponsesResponsetype(context.TODO(), t)
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

type POSTProductsProductidApisApiidResponsesResponsetypeResponse struct {
    Templates         types.List `tfsdk:"templates"`
Statuscode         types.Int32`tfsdk:"status_code"`
Responsetype         types.String `tfsdk:"response_type"`
Responsename         types.String `tfsdk:"response_name"`
Isdefault         types.Bool `tfsdk:"is_default"`
Headers         types.List `tfsdk:"headers"`
Apiid         types.String `tfsdk:"api_id"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidResponsesResponsetype(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidResponsesResponsetypeResponse, error) {
	var dto POSTProductsProductidApisApiidResponsesResponsetypeResponse

    
				if data["templates"] != nil {
					tempTemplates := data["templates"].([]interface{})
					dto.Templates = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"response_type": types.StringType,
"mapping_template": types.StringType,
"content_type": types.StringType,
"api_id": types.StringType,

	},

					}}.ElementType(), tempTemplates)
				}
				if data["status_code"] != nil {
					dto.Statuscode = types.Int32Value(data["status_code"].(int32))
				}

			if data["response_type"] != nil {
				dto.Responsetype = types.StringValue(data["response_type"].(string))
			}

			if data["response_name"] != nil {
				dto.Responsename = types.StringValue(data["response_name"].(string))
			}

			if data["is_default"] != nil {
				dto.Isdefault = types.BoolValue(data["is_default"].(bool))
			}

				if data["headers"] != nil {
					tempHeaders := data["headers"].([]interface{})
					dto.Headers = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"response_type": types.StringType,
"header_value": types.StringType,
"header_name": types.StringType,
"api_id": types.StringType,

	},

					}}.ElementType(), tempHeaders)
				}
			if data["api_id"] != nil {
				dto.Apiid = types.StringValue(data["api_id"].(string))
			}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidResponsesResponsetype(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidResponsesResponsetype(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidResponsesResponsetype(value interface{}) (attr.Value, error) {
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

