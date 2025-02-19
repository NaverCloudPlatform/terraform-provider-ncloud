
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

type PrimitivePOSTApikeysSearchRequest struct {
    ApiKeyName string `json:"apiKeyName"`
StatusSet types.List `json:"statusSet"`
ProductId string `json:"productId"`
Offset int64 `json:"offset"`
Limit int64 `json:"limit"`

}

type StringifiedPOSTApikeysSearchRequest struct {
	ApiKeyName string `json:"apiKeyName"`
StatusSet string `json:"statusSet"`
ProductId string `json:"productId"`
Offset string `json:"offset"`
Limit string `json:"limit"`

}

func (n *NClient) POSTApikeysSearch(ctx context.Context, primitiveReq *PrimitivePOSTApikeysSearchRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	
			if r.ApiKeyName != "" {
				initBody["apiKeyName"] = r.ApiKeyName
			}

			if r.StatusSet != "" {
				initBody["statusSet"] = r.StatusSet
			}
initBody["productId"] = r.ProductId

			if r.Offset != "" {
				initBody["offset"] = r.Offset
			}

			if r.Limit != "" {
				initBody["limit"] = r.Limit
			}


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"api-keys"+"/"+"search"

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

func (n *NClient) POSTApikeysSearch_TF(ctx context.Context, r *PrimitivePOSTApikeysSearchRequest) (*POSTApikeysSearchResponse, error) {
	t, err := n.POSTApikeysSearch(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTApikeysSearch(context.TODO(), t)
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

type POSTApikeysSearchResponse struct {
    Total         types.Int64`tfsdk:"total"`
ApiKeys         types.List `tfsdk:"api_keys"`

}

func ConvertToFrameworkTypes_POSTApikeysSearch(ctx context.Context, data map[string]interface{}) (*POSTApikeysSearchResponse, error) {
	var dto POSTApikeysSearchResponse

    
				if data["total"] != nil {
					dto.Total = types.Int64Value(data["total"].(int64))
				}

				if data["api_keys"] != nil {
					tempApiKeys := data["api_keys"].([]interface{})
					dto.ApiKeys = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"status": types.StringType,
"permission": types.StringType,
"disabled": types.BoolType,
"api_key_name": types.StringType,
"api_key_id": types.StringType,
"api_key_description": types.StringType,
"action_name": types.StringType,

	},

					}}.ElementType(), tempApiKeys)
				}

	return &dto, nil
}

func convertToObject_POSTApikeysSearch(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTApikeysSearch(value)
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

func convertValueToAttr_POSTApikeysSearch(value interface{}) (attr.Value, error) {
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

