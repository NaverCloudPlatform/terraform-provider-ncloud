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

type PrimitiveGETAuthorizersAuthorizeridRequest struct {
	Authorizerid string `json:"authorizer-id"`
}

type StringifiedGETAuthorizersAuthorizeridRequest struct {
	Authorizerid string `json:"authorizer-id"`
}

func (n *NClient) GETAuthorizersAuthorizerid(ctx context.Context, primitiveReq *PrimitiveGETAuthorizersAuthorizeridRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL + "/" + "authorizers" + "/" + ClearDoubleQuote(r.Authorizerid)

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

func (n *NClient) GETAuthorizersAuthorizerid_TF(ctx context.Context, r *PrimitiveGETAuthorizersAuthorizeridRequest) (*GETAuthorizersAuthorizeridResponse, error) {
	t, err := n.GETAuthorizersAuthorizerid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETAuthorizersAuthorizerid(context.TODO(), t)
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

type GETAuthorizersAuthorizeridResponse struct {
	Tenantid              types.String `tfsdk:"tenant_id"`
	Modifier              types.String `tfsdk:"modifier"`
	Modtime               types.String `tfsdk:"mod_time"`
	Domaincode            types.String `tfsdk:"domain_code"`
	Cachettlsec           types.Int32  `tfsdk:"cache_ttl_sec"`
	Authorizertype        types.String `tfsdk:"authorizer_type"`
	Authorizername        types.String `tfsdk:"authorizer_name"`
	Authorizerid          types.String `tfsdk:"authorizer_id"`
	Authorizerdescription types.String `tfsdk:"authorizer_description"`
	AuthorizerConfig      types.Object `tfsdk:"authorizer_config"`
}

func ConvertToFrameworkTypes_GETAuthorizersAuthorizerid(ctx context.Context, data map[string]interface{}) (*GETAuthorizersAuthorizeridResponse, error) {
	var dto GETAuthorizersAuthorizeridResponse

	if data["tenant_id"] != nil {
		dto.Tenantid = types.StringValue(data["tenant_id"].(string))
	}

	if data["modifier"] != nil {
		dto.Modifier = types.StringValue(data["modifier"].(string))
	}

	if data["mod_time"] != nil {
		dto.Modtime = types.StringValue(data["mod_time"].(string))
	}

	if data["domain_code"] != nil {
		dto.Domaincode = types.StringValue(data["domain_code"].(string))
	}

	if data["cache_ttl_sec"] != nil {
		dto.Cachettlsec = types.Int32Value(data["cache_ttl_sec"].(int32))
	}

	if data["authorizer_type"] != nil {
		dto.Authorizertype = types.StringValue(data["authorizer_type"].(string))
	}

	if data["authorizer_name"] != nil {
		dto.Authorizername = types.StringValue(data["authorizer_name"].(string))
	}

	if data["authorizer_id"] != nil {
		dto.Authorizerid = types.StringValue(data["authorizer_id"].(string))
	}

	if data["authorizer_description"] != nil {
		dto.Authorizerdescription = types.StringValue(data["authorizer_description"].(string))
	}

	if data["authorizer_config"] != nil {
		tempAuthorizerConfig := data["authorizer_config"].(map[string]interface{})

		allFields := []string{
			"payload",
			"function_id",
			"region",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempAuthorizerConfig[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempAuthorizerConfig, err := convertToObject_GETAuthorizersAuthorizerid(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.AuthorizerConfig = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{

			"payload": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

				"name": types.StringType,
				"in":   types.StringType,
			},
			}},
			"function_id": types.StringType,
			"region":      types.StringType,
		}}.AttributeTypes(), convertedTempAuthorizerConfig)
	}

	return &dto, nil
}

func convertToObject_GETAuthorizersAuthorizerid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{

		"payload": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"name": types.StringType,
			"in":   types.StringType,
		},
		}},
		"function_id": types.StringType,
		"region":      types.StringType,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			if field == "payload" && len(value.([]interface{})) == 0 {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"name": types.StringType,
					"in":   types.StringType,
				}).Type(ctx))
				attrValues[field] = listV
				continue
			}

			attrValue, err := convertValueToAttr_GETAuthorizersAuthorizerid(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {

			if field == "payload" {
				listV := types.ListNull(types.ObjectNull(map[string]attr.Type{
					"name": types.StringType,
					"in":   types.StringType,
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

func convertValueToAttr_GETAuthorizersAuthorizerid(value interface{}) (attr.Value, error) {
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
