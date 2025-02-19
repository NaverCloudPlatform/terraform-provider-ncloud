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

type PrimitivePOSTAuthorizersAuthorizeridTestRequest struct {
	Authorizerid string       `json:"authorizer-id"`
	Params       types.Object `json:"params"`
}

type StringifiedPOSTAuthorizersAuthorizeridTestRequest struct {
	Authorizerid string `json:"authorizer-id"`
	Params       string `json:"params"`
}

func (n *NClient) POSTAuthorizersAuthorizeridTest(ctx context.Context, primitiveReq *PrimitivePOSTAuthorizersAuthorizeridTestRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.Params != "" {
		initBody["params"] = r.Params
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "authorizers" + "/" + ClearDoubleQuote(r.Authorizerid) + "/" + "test"

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

func (n *NClient) POSTAuthorizersAuthorizeridTest_TF(ctx context.Context, r *PrimitivePOSTAuthorizersAuthorizeridTestRequest) (*POSTAuthorizersAuthorizeridTestResponse, error) {
	t, err := n.POSTAuthorizersAuthorizeridTest(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTAuthorizersAuthorizeridTest(context.TODO(), t)
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

type POSTAuthorizersAuthorizeridTestResponse struct {
	Response types.Object `tfsdk:"response"`
}

func ConvertToFrameworkTypes_POSTAuthorizersAuthorizeridTest(ctx context.Context, data map[string]interface{}) (*POSTAuthorizersAuthorizeridTestResponse, error) {
	var dto POSTAuthorizersAuthorizeridTestResponse

	if data["response"] != nil {
		tempResponse := data["response"].(map[string]interface{})

		allFields := []string{
			"status",
			"latency",
			"body",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempResponse[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempResponse, err := convertToObject_POSTAuthorizersAuthorizeridTest(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.Response = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"status":  types.Int32Type,
			"latency": types.StringType,

			"body": types.ObjectType{AttrTypes: map[string]attr.Type{}},
		}}.AttributeTypes(), convertedTempResponse)
	}

	return &dto, nil
}

func convertToObject_POSTAuthorizersAuthorizeridTest(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"status":  types.Int32Type,
		"latency": types.StringType,

		"body": types.ObjectType{AttrTypes: map[string]attr.Type{}},
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_POSTAuthorizersAuthorizeridTest(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {

			if field == "body" {
				listV := types.ObjectNull(map[string]attr.Type{})
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

func convertValueToAttr_POSTAuthorizersAuthorizeridTest(value interface{}) (attr.Value, error) {
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
