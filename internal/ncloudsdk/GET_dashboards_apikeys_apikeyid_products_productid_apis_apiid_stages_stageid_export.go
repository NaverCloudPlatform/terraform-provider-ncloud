
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

type PrimitiveGETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportRequest struct {
    Apikeyid string `json:"api-key-id"`
Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Stageid string `json:"stage-id"`
From string `json:"from"`
Limit int64 `json:"limit"`
Offset int64 `json:"offset"`
Regions types.List `json:"regions"`
TimeZone string `json:"timeZone"`
To string `json:"to"`

}

type StringifiedGETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportRequest struct {
	Apikeyid string `json:"api-key-id"`
Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Stageid string `json:"stage-id"`
From string `json:"from"`
Limit string `json:"limit"`
Offset string `json:"offset"`
Regions string `json:"regions"`
TimeZone string `json:"timeZone"`
To string `json:"to"`

}

func (n *NClient) GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(ctx context.Context, primitiveReq *PrimitiveGETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	
				query["from"] = r.From

				if r.Limit!= "" {
					query["limit"] = r.Limit
				}

				if r.Offset!= "" {
					query["offset"] = r.Offset
				}

				if r.Regions!= "" {
					query["regions"] = r.Regions
				}

				query["timeZone"] = r.TimeZone

				query["to"] = r.To


	

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"dashboards"+"/"+"api-keys"+"/"+ClearDoubleQuote(r.Apikeyid)+"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"stages"+"/"+ClearDoubleQuote(r.Stageid)+"/"+"export"

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

func (n *NClient) GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport_TF(ctx context.Context, r *PrimitiveGETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportRequest) (*GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportResponse, error) {
	t, err := n.GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(context.TODO(), t)
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

type GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportResponse struct {
    
}

func ConvertToFrameworkTypes_GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(ctx context.Context, data map[string]interface{}) (*GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportResponse, error) {
	var dto GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExportResponse

    

	return &dto, nil
}

func convertToObject_GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(value)
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

func convertValueToAttr_GETDashboardsApikeysApikeyidProductsProductidApisApiidStagesStageidExport(value interface{}) (attr.Value, error) {
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

