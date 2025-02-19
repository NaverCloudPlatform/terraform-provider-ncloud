
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

type PrimitivePOSTUsageplansSearchRequest struct {
    UsagePlanName string `json:"usagePlanName"`
Offset int64 `json:"offset"`
Limit int64 `json:"limit"`

}

type StringifiedPOSTUsageplansSearchRequest struct {
	UsagePlanName string `json:"usagePlanName"`
Offset string `json:"offset"`
Limit string `json:"limit"`

}

func (n *NClient) POSTUsageplansSearch(ctx context.Context, primitiveReq *PrimitivePOSTUsageplansSearchRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	
			if r.UsagePlanName != "" {
				initBody["usagePlanName"] = r.UsagePlanName
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

	url := n.BaseURL +"/"+"usage-plans"+"/"+"search"

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

func (n *NClient) POSTUsageplansSearch_TF(ctx context.Context, r *PrimitivePOSTUsageplansSearchRequest) (*POSTUsageplansSearchResponse, error) {
	t, err := n.POSTUsageplansSearch(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTUsageplansSearch(context.TODO(), t)
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

type POSTUsageplansSearchResponse struct {
    UsagePlans         types.List `tfsdk:"usage_plans"`
Total         types.Int64`tfsdk:"total"`
Initialcount         types.Int64`tfsdk:"initial_count"`

}

func ConvertToFrameworkTypes_POSTUsageplansSearch(ctx context.Context, data map[string]interface{}) (*POSTUsageplansSearchResponse, error) {
	var dto POSTUsageplansSearchResponse

    
				if data["usage_plans"] != nil {
					tempUsagePlans := data["usage_plans"].([]interface{})
					dto.UsagePlans = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"usage_plan_name": types.StringType,
"usage_plan_id": types.StringType,
"usage_plan_description": types.StringType,
"tenant_id": types.StringType,
"rate_rps": types.Int32Type,
"quota_condition": types.StringType,
"permission": types.StringType,
"month_quota_request": types.Int64Type,
"modifier": types.StringType,
"domain_code": types.StringType,
"disabled": types.BoolType,
"day_quota_request": types.Int64Type,
"associated_stages_count": types.Int64Type,
"action_name": types.StringType,

	},

					}}.ElementType(), tempUsagePlans)
				}
				if data["total"] != nil {
					dto.Total = types.Int64Value(data["total"].(int64))
				}

				if data["initial_count"] != nil {
					dto.Initialcount = types.Int64Value(data["initial_count"].(int64))
				}


	return &dto, nil
}

func convertToObject_POSTUsageplansSearch(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTUsageplansSearch(value)
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

func convertValueToAttr_POSTUsageplansSearch(value interface{}) (attr.Value, error) {
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

