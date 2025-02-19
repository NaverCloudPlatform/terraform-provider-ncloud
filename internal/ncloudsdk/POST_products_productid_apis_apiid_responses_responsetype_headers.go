
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

type PrimitivePOSTProductsProductidApisApiidResponsesResponsetypeHeadersRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Responsetype string `json:"response-type"`
HeaderValue string `json:"headerValue"`
HeaderName string `json:"headerName"`

}

type StringifiedPOSTProductsProductidApisApiidResponsesResponsetypeHeadersRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Responsetype string `json:"response-type"`
HeaderValue string `json:"headerValue"`
HeaderName string `json:"headerName"`

}

func (n *NClient) POSTProductsProductidApisApiidResponsesResponsetypeHeaders(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidResponsesResponsetypeHeadersRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["headerValue"] = r.HeaderValue
initBody["headerName"] = r.HeaderName


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"responses"+"/"+ClearDoubleQuote(r.Responsetype)+"/"+"headers"

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

func (n *NClient) POSTProductsProductidApisApiidResponsesResponsetypeHeaders_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidResponsesResponsetypeHeadersRequest) (*POSTProductsProductidApisApiidResponsesResponsetypeHeadersResponse, error) {
	t, err := n.POSTProductsProductidApisApiidResponsesResponsetypeHeaders(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidResponsesResponsetypeHeaders(context.TODO(), t)
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

type POSTProductsProductidApisApiidResponsesResponsetypeHeadersResponse struct {
    Responsetype         types.String `tfsdk:"response_type"`
Headervalue         types.String `tfsdk:"header_value"`
Headername         types.String `tfsdk:"header_name"`
Apiid         types.String `tfsdk:"api_id"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidResponsesResponsetypeHeaders(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidResponsesResponsetypeHeadersResponse, error) {
	var dto POSTProductsProductidApisApiidResponsesResponsetypeHeadersResponse

    
			if data["response_type"] != nil {
				dto.Responsetype = types.StringValue(data["response_type"].(string))
			}

			if data["header_value"] != nil {
				dto.Headervalue = types.StringValue(data["header_value"].(string))
			}

			if data["header_name"] != nil {
				dto.Headername = types.StringValue(data["header_name"].(string))
			}

			if data["api_id"] != nil {
				dto.Apiid = types.StringValue(data["api_id"].(string))
			}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidResponsesResponsetypeHeaders(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidResponsesResponsetypeHeaders(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidResponsesResponsetypeHeaders(value interface{}) (attr.Value, error) {
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

