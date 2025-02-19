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

type PrimitiveGETDashboardsApikeysApikeyidProductidsRequest struct {
	Apikeyid string     `json:"api-key-id"`
	From     string     `json:"from"`
	Limit    int64      `json:"limit"`
	Offset   int64      `json:"offset"`
	Regions  types.List `json:"regions"`
	TimeZone string     `json:"timeZone"`
	To       string     `json:"to"`
}

type StringifiedGETDashboardsApikeysApikeyidProductidsRequest struct {
	Apikeyid string `json:"api-key-id"`
	From     string `json:"from"`
	Limit    string `json:"limit"`
	Offset   string `json:"offset"`
	Regions  string `json:"regions"`
	TimeZone string `json:"timeZone"`
	To       string `json:"to"`
}

func (n *NClient) GETDashboardsApikeysApikeyidProductids(ctx context.Context, primitiveReq *PrimitiveGETDashboardsApikeysApikeyidProductidsRequest) (map[string]interface{}, error) {
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

	url := n.BaseURL + "/" + "dashboards" + "/" + "api-keys" + "/" + ClearDoubleQuote(r.Apikeyid) + "/" + "product-ids"

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

func (n *NClient) GETDashboardsApikeysApikeyidProductids_TF(ctx context.Context, r *PrimitiveGETDashboardsApikeysApikeyidProductidsRequest) (*GETDashboardsApikeysApikeyidProductidsResponse, error) {
	t, err := n.GETDashboardsApikeysApikeyidProductids(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETDashboardsApikeysApikeyidProductids(context.TODO(), t)
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

type GETDashboardsApikeysApikeyidProductidsResponse struct {
	Products types.List `tfsdk:"products"`
}

func ConvertToFrameworkTypes_GETDashboardsApikeysApikeyidProductids(ctx context.Context, data map[string]interface{}) (*GETDashboardsApikeysApikeyidProductidsResponse, error) {
	var dto GETDashboardsApikeysApikeyidProductidsResponse

	if data["products"] != nil {
		tempProducts := data["products"].([]interface{})
		dto.Products = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"product_name": types.StringType,
			"product_id":   types.StringType,
			"permission":   types.StringType,
			"is_deleted":   types.BoolType,
			"disabled":     types.BoolType,
			"action_name":  types.StringType,
		},
		}}.ElementType(), tempProducts)
	}

	return &dto, nil
}

func convertToObject_GETDashboardsApikeysApikeyidProductids(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETDashboardsApikeysApikeyidProductids(value)
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

func convertValueToAttr_GETDashboardsApikeysApikeyidProductids(value interface{}) (attr.Value, error) {
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
