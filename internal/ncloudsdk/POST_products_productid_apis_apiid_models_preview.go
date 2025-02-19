
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

type PrimitivePOSTProductsProductidApisApiidModelsPreviewRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
ModelSchema string `json:"modelSchema"`
ModelId string `json:"modelId"`

}

type StringifiedPOSTProductsProductidApisApiidModelsPreviewRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
ModelSchema string `json:"modelSchema"`
ModelId string `json:"modelId"`

}

func (n *NClient) POSTProductsProductidApisApiidModelsPreview(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidModelsPreviewRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["modelSchema"] = r.ModelSchema

			if r.ModelId != "" {
				initBody["modelId"] = r.ModelId
			}


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"models"+"/"+"preview"

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

func (n *NClient) POSTProductsProductidApisApiidModelsPreview_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidModelsPreviewRequest) (*POSTProductsProductidApisApiidModelsPreviewResponse, error) {
	t, err := n.POSTProductsProductidApisApiidModelsPreview(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidModelsPreview(context.TODO(), t)
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

type POSTProductsProductidApisApiidModelsPreviewResponse struct {
    ModelSchema         types.Object `tfsdk:"model_schema"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidModelsPreview(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidModelsPreviewResponse, error) {
	var dto POSTProductsProductidApisApiidModelsPreviewResponse

    
			if data["model_schema"] != nil {
				tempModelSchema := data["model_schema"].(map[string]interface{})

				allFields := []string{
					
				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempModelSchema[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempModelSchema, err := convertToObject_POSTProductsProductidApisApiidModelsPreview(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.ModelSchema = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					
				}}.AttributeTypes(), convertedTempModelSchema)
			}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidModelsPreview(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        

	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidModelsPreview(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidModelsPreview(value interface{}) (attr.Value, error) {
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

