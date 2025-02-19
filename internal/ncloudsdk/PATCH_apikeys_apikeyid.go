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

type PrimitivePATCHApikeysApikeyidRequest struct {
	Apikeyid string `json:"api-key-id"`
	KeyType  string `json:"keyType"`
}

type StringifiedPATCHApikeysApikeyidRequest struct {
	Apikeyid string `json:"api-key-id"`
	KeyType  string `json:"keyType"`
}

func (n *NClient) PATCHApikeysApikeyid(ctx context.Context, primitiveReq *PrimitivePATCHApikeysApikeyidRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	initBody["keyType"] = r.KeyType

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "api-keys" + "/" + ClearDoubleQuote(r.Apikeyid)

	response, err := n.MakeRequestWithContext(ctx, "PATCH", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) PATCHApikeysApikeyid_TF(ctx context.Context, r *PrimitivePATCHApikeysApikeyidRequest) (*PATCHApikeysApikeyidResponse, error) {
	t, err := n.PATCHApikeysApikeyid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_PATCHApikeysApikeyid(context.TODO(), t)
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

type PATCHApikeysApikeyidResponse struct {
	ApiKey types.Object `tfsdk:"api_key"`
}

func ConvertToFrameworkTypes_PATCHApikeysApikeyid(ctx context.Context, data map[string]interface{}) (*PATCHApikeysApikeyidResponse, error) {
	var dto PATCHApikeysApikeyidResponse

	if data["api_key"] != nil {
		tempApiKey := data["api_key"].(map[string]interface{})

		allFields := []string{
			"tenant_id",
			"secondary_key",
			"primary_key",
			"modifier",
			"mod_time",
			"is_enabled",
			"domain_code",
			"api_key_name",
			"api_key_id",
			"api_key_description",
		}

		convertedMap := make(map[string]interface{})
		for _, field := range allFields {
			if val, ok := tempApiKey[field]; ok {
				convertedMap[field] = val
			}
		}

		convertedTempApiKey, err := convertToObject_PATCHApikeysApikeyid(ctx, convertedMap)
		if err != nil {
			return nil, err
		}

		dto.ApiKey = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"tenant_id":           types.StringType,
			"secondary_key":       types.StringType,
			"primary_key":         types.StringType,
			"modifier":            types.StringType,
			"mod_time":            types.StringType,
			"is_enabled":          types.BoolType,
			"domain_code":         types.StringType,
			"api_key_name":        types.StringType,
			"api_key_id":          types.StringType,
			"api_key_description": types.StringType,
		}}.AttributeTypes(), convertedTempApiKey)
	}

	return &dto, nil
}

func convertToObject_PATCHApikeysApikeyid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{
		"tenant_id":           types.StringType,
		"secondary_key":       types.StringType,
		"primary_key":         types.StringType,
		"modifier":            types.StringType,
		"mod_time":            types.StringType,
		"is_enabled":          types.BoolType,
		"domain_code":         types.StringType,
		"api_key_name":        types.StringType,
		"api_key_id":          types.StringType,
		"api_key_description": types.StringType,
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_PATCHApikeysApikeyid(value)
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

func convertValueToAttr_PATCHApikeysApikeyid(value interface{}) (attr.Value, error) {
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
