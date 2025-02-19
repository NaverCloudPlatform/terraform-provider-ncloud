
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

type PrimitiveGETProductsProductidApisApiidModelsModelidRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Modelid string `json:"model-id"`

}

type StringifiedGETProductsProductidApisApiidModelsModelidRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Modelid string `json:"model-id"`

}

func (n *NClient) GETProductsProductidApisApiidModelsModelid(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidApisApiidModelsModelidRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"models"+"/"+ClearDoubleQuote(r.Modelid)

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

func (n *NClient) GETProductsProductidApisApiidModelsModelid_TF(ctx context.Context, r *PrimitiveGETProductsProductidApisApiidModelsModelidRequest) (*GETProductsProductidApisApiidModelsModelidResponse, error) {
	t, err := n.GETProductsProductidApisApiidModelsModelid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidApisApiidModelsModelid(context.TODO(), t)
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

type GETProductsProductidApisApiidModelsModelidResponse struct {
    Modelschema         types.String `tfsdk:"model_schema"`
Modelname         types.String `tfsdk:"model_name"`
Modelid         types.String `tfsdk:"model_id"`
Modeldescription         types.String `tfsdk:"model_description"`
Apiid         types.String `tfsdk:"api_id"`

}

func ConvertToFrameworkTypes_GETProductsProductidApisApiidModelsModelid(ctx context.Context, data map[string]interface{}) (*GETProductsProductidApisApiidModelsModelidResponse, error) {
	var dto GETProductsProductidApisApiidModelsModelidResponse

    
			if data["model_schema"] != nil {
				dto.Modelschema = types.StringValue(data["model_schema"].(string))
			}

			if data["model_name"] != nil {
				dto.Modelname = types.StringValue(data["model_name"].(string))
			}

			if data["model_id"] != nil {
				dto.Modelid = types.StringValue(data["model_id"].(string))
			}

			if data["model_description"] != nil {
				dto.Modeldescription = types.StringValue(data["model_description"].(string))
			}

			if data["api_id"] != nil {
				dto.Apiid = types.StringValue(data["api_id"].(string))
			}


	return &dto, nil
}

func convertToObject_GETProductsProductidApisApiidModelsModelid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_GETProductsProductidApisApiidModelsModelid(value)
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

func convertValueToAttr_GETProductsProductidApisApiidModelsModelid(value interface{}) (attr.Value, error) {
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

