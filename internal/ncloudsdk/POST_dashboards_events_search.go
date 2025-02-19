
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

type PrimitivePOSTDashboardsEventsSearchRequest struct {
    TimeZone string `json:"timeZone"`
To string `json:"to"`
From string `json:"from"`
Level types.List `json:"level"`
Type types.List `json:"type"`
Offset int64 `json:"offset"`
Limit int64 `json:"limit"`

}

type StringifiedPOSTDashboardsEventsSearchRequest struct {
	TimeZone string `json:"timeZone"`
To string `json:"to"`
From string `json:"from"`
Level string `json:"level"`
Type string `json:"type"`
Offset string `json:"offset"`
Limit string `json:"limit"`

}

func (n *NClient) POSTDashboardsEventsSearch(ctx context.Context, primitiveReq *PrimitivePOSTDashboardsEventsSearchRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["timeZone"] = r.TimeZone
initBody["to"] = r.To
initBody["from"] = r.From

			if r.Level != "" {
				initBody["level"] = r.Level
			}

			if r.Type != "" {
				initBody["type"] = r.Type
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

	url := n.BaseURL +"/"+"dashboards"+"/"+"events"+"/"+"search"

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

func (n *NClient) POSTDashboardsEventsSearch_TF(ctx context.Context, r *PrimitivePOSTDashboardsEventsSearchRequest) (*POSTDashboardsEventsSearchResponse, error) {
	t, err := n.POSTDashboardsEventsSearch(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTDashboardsEventsSearch(context.TODO(), t)
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

type POSTDashboardsEventsSearchResponse struct {
    Events         types.Object `tfsdk:"events"`

}

func ConvertToFrameworkTypes_POSTDashboardsEventsSearch(ctx context.Context, data map[string]interface{}) (*POSTDashboardsEventsSearchResponse, error) {
	var dto POSTDashboardsEventsSearchResponse

    
			if data["events"] != nil {
				tempEvents := data["events"].(map[string]interface{})

				allFields := []string{
					"total",
"content",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempEvents[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempEvents, err := convertToObject_POSTDashboardsEventsSearch(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.Events = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"type": types.StringType,
"time": types.StringType,
"message": types.StringType,
"level": types.StringType,

	},
			}},

				}}.AttributeTypes(), convertedTempEvents)
			}


	return &dto, nil
}

func convertToObject_POSTDashboardsEventsSearch(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "total": types.Int64Type,

			"content": types.ListType{ElemType:
				
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"type": types.StringType,
"time": types.StringType,
"message": types.StringType,
"level": types.StringType,

	},
			}},


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			
			if field == "content" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"type": types.StringType,
"time": types.StringType,
"message": types.StringType,
"level": types.StringType,

				}).Type(ctx))
				attrValues[field] = listV
				continue
			}


			attrValue, err := convertValueToAttr_POSTDashboardsEventsSearch(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            
				if field == "content" {
					listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
						"type": types.StringType,
"time": types.StringType,
"message": types.StringType,
"level": types.StringType,

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

func convertValueToAttr_POSTDashboardsEventsSearch(value interface{}) (attr.Value, error) {
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

