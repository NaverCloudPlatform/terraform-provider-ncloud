
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

type PrimitiveGETAuthorizersRequest struct {
    Limit int64 `json:"limit"`
Name string `json:"name"`
Offset int64 `json:"offset"`

}

type StringifiedGETAuthorizersRequest struct {
	Limit string `json:"limit"`
Name string `json:"name"`
Offset string `json:"offset"`

}

func (n *NClient) GETAuthorizers(ctx context.Context, primitiveReq *PrimitiveGETAuthorizersRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	
				if r.Limit!= "" {
					query["limit"] = r.Limit
				}

				if r.Name!= "" {
					query["name"] = r.Name
				}

				if r.Offset!= "" {
					query["offset"] = r.Offset
				}


	

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"authorizers"

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

func (n *NClient) GETAuthorizers_TF(ctx context.Context, r *PrimitiveGETAuthorizersRequest) (*GETAuthorizersResponse, error) {
	t, err := n.GETAuthorizers(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETAuthorizers(context.TODO(), t)
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

type GETAuthorizersResponse struct {
    Total         types.Int64`tfsdk:"total"`
Initialcount         types.Int64`tfsdk:"initial_count"`
Authorizers         types.List `tfsdk:"authorizers"`

}

func ConvertToFrameworkTypes_GETAuthorizers(ctx context.Context, data map[string]interface{}) (*GETAuthorizersResponse, error) {
	var dto GETAuthorizersResponse

    
				if data["total"] != nil {
					dto.Total = types.Int64Value(data["total"].(int64))
				}

				if data["initial_count"] != nil {
					dto.Initialcount = types.Int64Value(data["initial_count"].(int64))
				}

				if data["authorizers"] != nil {
					tempAuthorizers := data["authorizers"].([]interface{})
					dto.Authorizers = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"permission": types.StringType,
"disabled": types.BoolType,
"authorizer_name": types.StringType,
"authorizer_id": types.StringType,
"authorizer_description": types.StringType,
"action_name": types.StringType,

	},

					}}.ElementType(), tempAuthorizers)
				}

	return &dto, nil
}

func convertToObject_GETAuthorizers(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_GETAuthorizers(value)
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

func convertValueToAttr_GETAuthorizers(value interface{}) (attr.Value, error) {
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

