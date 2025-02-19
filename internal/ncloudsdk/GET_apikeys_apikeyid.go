
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

type PrimitiveGETApikeysApikeyidRequest struct {
    Apikeyid string `json:"api-key-id"`

}

type StringifiedGETApikeysApikeyidRequest struct {
	Apikeyid string `json:"api-key-id"`

}

func (n *NClient) GETApikeysApikeyid(ctx context.Context, primitiveReq *PrimitiveGETApikeysApikeyidRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL +"/"+"api-keys"+"/"+ClearDoubleQuote(r.Apikeyid)

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

func (n *NClient) GETApikeysApikeyid_TF(ctx context.Context, r *PrimitiveGETApikeysApikeyidRequest) (*GETApikeysApikeyidResponse, error) {
	t, err := n.GETApikeysApikeyid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETApikeysApikeyid(context.TODO(), t)
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

type GETApikeysApikeyidResponse struct {
    Tenantid         types.String `tfsdk:"tenant_id"`
Secondarykey         types.String `tfsdk:"secondary_key"`
Primarykey         types.String `tfsdk:"primary_key"`
Modifier         types.String `tfsdk:"modifier"`
Modtime         types.String `tfsdk:"mod_time"`
Isenabled         types.Bool `tfsdk:"is_enabled"`
Domaincode         types.String `tfsdk:"domain_code"`
Apikeyname         types.String `tfsdk:"api_key_name"`
Apikeyid         types.String `tfsdk:"api_key_id"`
Apikeydescription         types.String `tfsdk:"api_key_description"`

}

func ConvertToFrameworkTypes_GETApikeysApikeyid(ctx context.Context, data map[string]interface{}) (*GETApikeysApikeyidResponse, error) {
	var dto GETApikeysApikeyidResponse

    
			if data["tenant_id"] != nil {
				dto.Tenantid = types.StringValue(data["tenant_id"].(string))
			}

			if data["secondary_key"] != nil {
				dto.Secondarykey = types.StringValue(data["secondary_key"].(string))
			}

			if data["primary_key"] != nil {
				dto.Primarykey = types.StringValue(data["primary_key"].(string))
			}

			if data["modifier"] != nil {
				dto.Modifier = types.StringValue(data["modifier"].(string))
			}

			if data["mod_time"] != nil {
				dto.Modtime = types.StringValue(data["mod_time"].(string))
			}

			if data["is_enabled"] != nil {
				dto.Isenabled = types.BoolValue(data["is_enabled"].(bool))
			}

			if data["domain_code"] != nil {
				dto.Domaincode = types.StringValue(data["domain_code"].(string))
			}

			if data["api_key_name"] != nil {
				dto.Apikeyname = types.StringValue(data["api_key_name"].(string))
			}

			if data["api_key_id"] != nil {
				dto.Apikeyid = types.StringValue(data["api_key_id"].(string))
			}

			if data["api_key_description"] != nil {
				dto.Apikeydescription = types.StringValue(data["api_key_description"].(string))
			}


	return &dto, nil
}

func convertToObject_GETApikeysApikeyid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_GETApikeysApikeyid(value)
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

func convertValueToAttr_GETApikeysApikeyid(value interface{}) (attr.Value, error) {
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

