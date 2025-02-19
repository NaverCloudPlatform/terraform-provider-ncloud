
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

type PrimitivePOSTProductsProductidApisApiidModelsRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
ModelSchema string `json:"modelSchema"`
ModelDescription string `json:"modelDescription"`
ModelName string `json:"modelName"`

}

type StringifiedPOSTProductsProductidApisApiidModelsRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
ModelSchema string `json:"modelSchema"`
ModelDescription string `json:"modelDescription"`
ModelName string `json:"modelName"`

}

func (n *NClient) POSTProductsProductidApisApiidModels(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidModelsRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["modelSchema"] = r.ModelSchema

			if r.ModelDescription != "" {
				initBody["modelDescription"] = r.ModelDescription
			}
initBody["modelName"] = r.ModelName


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"models"

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

func (n *NClient) POSTProductsProductidApisApiidModels_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidModelsRequest) (*POSTProductsProductidApisApiidModelsResponse, error) {
	t, err := n.POSTProductsProductidApisApiidModels(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidModels(context.TODO(), t)
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

type POSTProductsProductidApisApiidModelsResponse struct {
    Modelschema         types.String `tfsdk:"model_schema"`
Modelname         types.String `tfsdk:"model_name"`
Modelid         types.String `tfsdk:"model_id"`
Modeldescription         types.String `tfsdk:"model_description"`
Apiid         types.String `tfsdk:"api_id"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidModels(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidModelsResponse, error) {
	var dto POSTProductsProductidApisApiidModelsResponse

    
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

func convertToObject_POSTProductsProductidApisApiidModels(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidModels(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidModels(value interface{}) (attr.Value, error) {
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

