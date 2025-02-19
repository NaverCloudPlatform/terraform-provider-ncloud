
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

type PrimitiveGETPublishedproductsProductidDocumentsRequest struct {
    Productid string `json:"product-id"`
ApiName string `json:"apiName"`
Limit int64 `json:"limit"`
Offset int64 `json:"offset"`

}

type StringifiedGETPublishedproductsProductidDocumentsRequest struct {
	Productid string `json:"product-id"`
ApiName string `json:"apiName"`
Limit string `json:"limit"`
Offset string `json:"offset"`

}

func (n *NClient) GETPublishedproductsProductidDocuments(ctx context.Context, primitiveReq *PrimitiveGETPublishedproductsProductidDocumentsRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	
				if r.ApiName!= "" {
					query["apiName"] = r.ApiName
				}

				if r.Limit!= "" {
					query["limit"] = r.Limit
				}

				if r.Offset!= "" {
					query["offset"] = r.Offset
				}


	

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"published-products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"documents"

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

func (n *NClient) GETPublishedproductsProductidDocuments_TF(ctx context.Context, r *PrimitiveGETPublishedproductsProductidDocumentsRequest) (*GETPublishedproductsProductidDocumentsResponse, error) {
	t, err := n.GETPublishedproductsProductidDocuments(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETPublishedproductsProductidDocuments(context.TODO(), t)
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

type GETPublishedproductsProductidDocumentsResponse struct {
    ApiPage         types.Object `tfsdk:"api_page"`

}

func ConvertToFrameworkTypes_GETPublishedproductsProductidDocuments(ctx context.Context, data map[string]interface{}) (*GETPublishedproductsProductidDocumentsResponse, error) {
	var dto GETPublishedproductsProductidDocumentsResponse

    
			if data["api_page"] != nil {
				tempApiPage := data["api_page"].(map[string]interface{})

				allFields := []string{
					"total",
"content",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempApiPage[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempApiPage, err := convertToObject_GETPublishedproductsProductidDocuments(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.ApiPage = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		
				"stages": types.ListType{ElemType:
					types.ObjectType{AttrTypes: map[string]attr.Type{
						"stage_name": types.StringType,
"stage_id": types.StringType,
"host": types.StringType,
"deployed_stage_deployment_no": types.Int64Type,

					},
				}},
"domain_code": types.StringType,
"api_name": types.StringType,
"api_id": types.StringType,
"api_description": types.StringType,

	},
			}},

				}}.AttributeTypes(), convertedTempApiPage)
			}


	return &dto, nil
}

func convertToObject_GETPublishedproductsProductidDocuments(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		
				"stages": types.ListType{ElemType:
					types.ObjectType{AttrTypes: map[string]attr.Type{
						"stage_name": types.StringType,
"stage_id": types.StringType,
"host": types.StringType,
"deployed_stage_deployment_no": types.Int64Type,

					},
				}},
"domain_code": types.StringType,
"api_name": types.StringType,
"api_id": types.StringType,
"api_description": types.StringType,

	},
			}},


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			
			if field == "content" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					
			"stages": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"stage_name": types.StringType,
"stage_id": types.StringType,
"host": types.StringType,
"deployed_stage_deployment_no": types.Int64Type,

	},
			}},
"domain_code": types.StringType,
"api_name": types.StringType,
"api_id": types.StringType,
"api_description": types.StringType,

				}).Type(ctx))
				attrValues[field] = listV
				continue
			}


			attrValue, err := convertValueToAttr_GETPublishedproductsProductidDocuments(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            
				if field == "content" {
					listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
						
			"stages": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"stage_name": types.StringType,
"stage_id": types.StringType,
"host": types.StringType,
"deployed_stage_deployment_no": types.Int64Type,

	},
			}},
"domain_code": types.StringType,
"api_name": types.StringType,
"api_id": types.StringType,
"api_description": types.StringType,

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

func convertValueToAttr_GETPublishedproductsProductidDocuments(value interface{}) (attr.Value, error) {
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

