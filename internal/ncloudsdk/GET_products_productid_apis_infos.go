
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

type PrimitiveGETProductsProductidApisInfosRequest struct {
    Productid string `json:"product-id"`
ApiName string `json:"apiName"`
HasStage bool `json:"hasStage"`
HasStageNotAssociatedWithUsagePlanId string `json:"hasStageNotAssociatedWithUsagePlanId"`
Limit int64 `json:"limit"`
Offset int64 `json:"offset"`
WithStage bool `json:"withStage"`

}

type StringifiedGETProductsProductidApisInfosRequest struct {
	Productid string `json:"product-id"`
ApiName string `json:"apiName"`
HasStage string `json:"hasStage"`
HasStageNotAssociatedWithUsagePlanId string `json:"hasStageNotAssociatedWithUsagePlanId"`
Limit string `json:"limit"`
Offset string `json:"offset"`
WithStage string `json:"withStage"`

}

func (n *NClient) GETProductsProductidApisInfos(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidApisInfosRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	
				if r.ApiName!= "" {
					query["apiName"] = r.ApiName
				}

				if r.HasStage!= "" {
					query["hasStage"] = r.HasStage
				}

				if r.HasStageNotAssociatedWithUsagePlanId!= "" {
					query["hasStageNotAssociatedWithUsagePlanId"] = r.HasStageNotAssociatedWithUsagePlanId
				}

				if r.Limit!= "" {
					query["limit"] = r.Limit
				}

				if r.Offset!= "" {
					query["offset"] = r.Offset
				}

				if r.WithStage!= "" {
					query["withStage"] = r.WithStage
				}


	

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+"infos"

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

func (n *NClient) GETProductsProductidApisInfos_TF(ctx context.Context, r *PrimitiveGETProductsProductidApisInfosRequest) (*GETProductsProductidApisInfosResponse, error) {
	t, err := n.GETProductsProductidApisInfos(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidApisInfos(context.TODO(), t)
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

type GETProductsProductidApisInfosResponse struct {
    Total         types.Int64`tfsdk:"total"`
Apis         types.List `tfsdk:"apis"`

}

func ConvertToFrameworkTypes_GETProductsProductidApisInfos(ctx context.Context, data map[string]interface{}) (*GETProductsProductidApisInfosResponse, error) {
	var dto GETProductsProductidApisInfosResponse

    
				if data["total"] != nil {
					dto.Total = types.Int64Value(data["total"].(int64))
				}

				if data["apis"] != nil {
					tempApis := data["apis"].([]interface{})
					dto.Apis = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
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

	},

					}}.ElementType(), tempApis)
				}

	return &dto, nil
}

func convertToObject_GETProductsProductidApisInfos(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_GETProductsProductidApisInfos(value)
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

func convertValueToAttr_GETProductsProductidApisInfos(value interface{}) (attr.Value, error) {
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

