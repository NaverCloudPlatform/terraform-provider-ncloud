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

type PrimitiveGETDashboardsApikeysRequest struct {
	From     string     `json:"from"`
	Limit    int64      `json:"limit"`
	Offset   int64      `json:"offset"`
	Regions  types.List `json:"regions"`
	TimeZone string     `json:"timeZone"`
	To       string     `json:"to"`
}

type StringifiedGETDashboardsApikeysRequest struct {
	From     string `json:"from"`
	Limit    string `json:"limit"`
	Offset   string `json:"offset"`
	Regions  string `json:"regions"`
	TimeZone string `json:"timeZone"`
	To       string `json:"to"`
}

func (n *NClient) GETDashboardsApikeys(ctx context.Context, primitiveReq *PrimitiveGETDashboardsApikeysRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	query["from"] = r.From

	if r.Limit != "" {
		query["limit"] = r.Limit
	}

	if r.Offset != "" {
		query["offset"] = r.Offset
	}

	if r.Regions != "" {
		query["regions"] = r.Regions
	}

	query["timeZone"] = r.TimeZone

	query["to"] = r.To

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "dashboards" + "/" + "api-keys"

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

func (n *NClient) GETDashboardsApikeys_TF(ctx context.Context, r *PrimitiveGETDashboardsApikeysRequest) (*GETDashboardsApikeysResponse, error) {
	t, err := n.GETDashboardsApikeys(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETDashboardsApikeys(context.TODO(), t)
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

type GETDashboardsApikeysResponse struct {
	ApiKeyIds types.List `tfsdk:"api_key_ids"`
}

func ConvertToFrameworkTypes_GETDashboardsApikeys(ctx context.Context, data map[string]interface{}) (*GETDashboardsApikeysResponse, error) {
	var dto GETDashboardsApikeysResponse

	if data["api_key_ids"] != nil {
		tempApiKeyIds := data["api_key_ids"].([]interface{})
		dto.ApiKeyIds = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"permission":   types.StringType,
			"is_enabled":   types.BoolType,
			"is_deleted":   types.BoolType,
			"disabled":     types.BoolType,
			"api_key_name": types.StringType,
			"api_key_id":   types.StringType,
			"action_name":  types.StringType,
		},
		}}.ElementType(), tempApiKeyIds)
	}

	return &dto, nil
}

func convertToObject_GETDashboardsApikeys(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETDashboardsApikeys(value)
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

func convertValueToAttr_GETDashboardsApikeys(value interface{}) (attr.Value, error) {
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
