
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

type PrimitivePOSTProductsProductidApikeysSearchRequest struct {
    Productid string `json:"product-id"`
ApiKeyId string `json:"apiKeyId"`
ApiKeyValue string `json:"apiKeyValue"`
StatusSet types.List `json:"statusSet"`
Offset int64 `json:"offset"`
Limit int64 `json:"limit"`

}

type StringifiedPOSTProductsProductidApikeysSearchRequest struct {
	Productid string `json:"product-id"`
ApiKeyId string `json:"apiKeyId"`
ApiKeyValue string `json:"apiKeyValue"`
StatusSet string `json:"statusSet"`
Offset string `json:"offset"`
Limit string `json:"limit"`

}

func (n *NClient) POSTProductsProductidApikeysSearch(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApikeysSearchRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	
			if r.ApiKeyId != "" {
				initBody["apiKeyId"] = r.ApiKeyId
			}

			if r.ApiKeyValue != "" {
				initBody["apiKeyValue"] = r.ApiKeyValue
			}

			if r.StatusSet != "" {
				initBody["statusSet"] = r.StatusSet
			}

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

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"api-keys"+"/"+"search"

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

func (n *NClient) POSTProductsProductidApikeysSearch_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApikeysSearchRequest) (*POSTProductsProductidApikeysSearchResponse, error) {
	t, err := n.POSTProductsProductidApikeysSearch(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApikeysSearch(context.TODO(), t)
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

type POSTProductsProductidApikeysSearchResponse struct {
    ApiKeyPage         types.Object `tfsdk:"api_key_page"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApikeysSearch(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApikeysSearchResponse, error) {
	var dto POSTProductsProductidApikeysSearchResponse

    
			if data["api_key_page"] != nil {
				tempApiKeyPage := data["api_key_page"].(map[string]interface{})

				allFields := []string{
					"total",
"content",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempApiKeyPage[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempApiKeyPage, err := convertToObject_POSTProductsProductidApikeysSearch(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.ApiKeyPage = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"status": types.StringType,
"secondary_key": types.StringType,
"reg_time": types.StringType,
"primary_key": types.StringType,
"mod_time": types.StringType,
"domain_code": types.StringType,
"api_key_name": types.StringType,
"api_key_id": types.StringType,

	},
			}},

				}}.AttributeTypes(), convertedTempApiKeyPage)
			}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApikeysSearch(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"status": types.StringType,
"secondary_key": types.StringType,
"reg_time": types.StringType,
"primary_key": types.StringType,
"mod_time": types.StringType,
"domain_code": types.StringType,
"api_key_name": types.StringType,
"api_key_id": types.StringType,

	},
			}},


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			
			if field == "content" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"status": types.StringType,
"secondary_key": types.StringType,
"reg_time": types.StringType,
"primary_key": types.StringType,
"mod_time": types.StringType,
"domain_code": types.StringType,
"api_key_name": types.StringType,
"api_key_id": types.StringType,

				}).Type(ctx))
				attrValues[field] = listV
				continue
			}


			attrValue, err := convertValueToAttr_POSTProductsProductidApikeysSearch(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            
				if field == "content" {
					listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
						"status": types.StringType,
"secondary_key": types.StringType,
"reg_time": types.StringType,
"primary_key": types.StringType,
"mod_time": types.StringType,
"domain_code": types.StringType,
"api_key_name": types.StringType,
"api_key_id": types.StringType,

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

func convertValueToAttr_POSTProductsProductidApikeysSearch(value interface{}) (attr.Value, error) {
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

