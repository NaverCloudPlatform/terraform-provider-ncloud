
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

type PrimitivePOSTPublishedproductsSearchRequest struct {
    ProductType types.List `json:"productType"`
Subscribed bool `json:"subscribed"`
PublisherId string `json:"publisherId"`
ProductName string `json:"productName"`
Offset int64 `json:"offset"`
Limit int64 `json:"limit"`

}

type StringifiedPOSTPublishedproductsSearchRequest struct {
	ProductType string `json:"productType"`
Subscribed string `json:"subscribed"`
PublisherId string `json:"publisherId"`
ProductName string `json:"productName"`
Offset string `json:"offset"`
Limit string `json:"limit"`

}

func (n *NClient) POSTPublishedproductsSearch(ctx context.Context, primitiveReq *PrimitivePOSTPublishedproductsSearchRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	
			if r.ProductType != "" {
				initBody["productType"] = r.ProductType
			}

			if r.Subscribed != "" {
				initBody["subscribed"] = r.Subscribed
			}

			if r.PublisherId != "" {
				initBody["publisherId"] = r.PublisherId
			}

			if r.ProductName != "" {
				initBody["productName"] = r.ProductName
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

	url := n.BaseURL +"/"+"published-products"+"/"+"search"

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

func (n *NClient) POSTPublishedproductsSearch_TF(ctx context.Context, r *PrimitivePOSTPublishedproductsSearchRequest) (*POSTPublishedproductsSearchResponse, error) {
	t, err := n.POSTPublishedproductsSearch(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTPublishedproductsSearch(context.TODO(), t)
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

type POSTPublishedproductsSearchResponse struct {
    ProductPage         types.Object `tfsdk:"product_page"`
Initialcount         types.Int64`tfsdk:"initial_count"`

}

func ConvertToFrameworkTypes_POSTPublishedproductsSearch(ctx context.Context, data map[string]interface{}) (*POSTPublishedproductsSearchResponse, error) {
	var dto POSTPublishedproductsSearchResponse

    
			if data["product_page"] != nil {
				tempProductPage := data["product_page"].(map[string]interface{})

				allFields := []string{
					"total",
"content",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempProductPage[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempProductPage, err := convertToObject_POSTPublishedproductsSearch(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.ProductPage = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"tenant_id": types.StringType,
"subscription_code": types.StringType,
"subscribed": types.BoolType,
"product_name": types.StringType,
"product_id": types.StringType,
"product_description": types.StringType,
"domain_code": types.StringType,

	},
			}},

				}}.AttributeTypes(), convertedTempProductPage)
			}

				if data["initial_count"] != nil {
					dto.Initialcount = types.Int64Value(data["initial_count"].(int64))
				}


	return &dto, nil
}

func convertToObject_POSTPublishedproductsSearch(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"tenant_id": types.StringType,
"subscription_code": types.StringType,
"subscribed": types.BoolType,
"product_name": types.StringType,
"product_id": types.StringType,
"product_description": types.StringType,
"domain_code": types.StringType,

	},
			}},


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			
			if field == "content" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"tenant_id": types.StringType,
"subscription_code": types.StringType,
"subscribed": types.BoolType,
"product_name": types.StringType,
"product_id": types.StringType,
"product_description": types.StringType,
"domain_code": types.StringType,

				}).Type(ctx))
				attrValues[field] = listV
				continue
			}


			attrValue, err := convertValueToAttr_POSTPublishedproductsSearch(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            
				if field == "content" {
					listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
						"tenant_id": types.StringType,
"subscription_code": types.StringType,
"subscribed": types.BoolType,
"product_name": types.StringType,
"product_id": types.StringType,
"product_description": types.StringType,
"domain_code": types.StringType,

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

func convertValueToAttr_POSTPublishedproductsSearch(value interface{}) (attr.Value, error) {
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

