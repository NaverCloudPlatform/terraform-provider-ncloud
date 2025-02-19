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

type PrimitivePUTProductsProductidApikeysApikeyidStagesStageidRequest struct {
	Productid   string `json:"product-id"`
	Apikeyid    string `json:"api-key-id"`
	Stageid     string `json:"stage-id"`
	UsagePlanId string `json:"usagePlanId"`
}

type StringifiedPUTProductsProductidApikeysApikeyidStagesStageidRequest struct {
	Productid   string `json:"product-id"`
	Apikeyid    string `json:"api-key-id"`
	Stageid     string `json:"stage-id"`
	UsagePlanId string `json:"usagePlanId"`
}

func (n *NClient) PUTProductsProductidApikeysApikeyidStagesStageid(ctx context.Context, primitiveReq *PrimitivePUTProductsProductidApikeysApikeyidStagesStageidRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.UsagePlanId != "" {
		initBody["usagePlanId"] = r.UsagePlanId
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "api-keys" + "/" + ClearDoubleQuote(r.Apikeyid) + "/" + "stages" + "/" + ClearDoubleQuote(r.Stageid)

	response, err := n.MakeRequestWithContext(ctx, "PUT", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) PUTProductsProductidApikeysApikeyidStagesStageid_TF(ctx context.Context, r *PrimitivePUTProductsProductidApikeysApikeyidStagesStageidRequest) (*PUTProductsProductidApikeysApikeyidStagesStageidResponse, error) {
	t, err := n.PUTProductsProductidApikeysApikeyidStagesStageid(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_PUTProductsProductidApikeysApikeyidStagesStageid(context.TODO(), t)
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

type PUTProductsProductidApikeysApikeyidStagesStageidResponse struct {
}

func ConvertToFrameworkTypes_PUTProductsProductidApikeysApikeyidStagesStageid(ctx context.Context, data map[string]interface{}) (*PUTProductsProductidApikeysApikeyidStagesStageidResponse, error) {
	var dto PUTProductsProductidApikeysApikeyidStagesStageidResponse

	return &dto, nil
}

func convertToObject_PUTProductsProductidApikeysApikeyidStagesStageid(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_PUTProductsProductidApikeysApikeyidStagesStageid(value)
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

func convertValueToAttr_PUTProductsProductidApikeysApikeyidStagesStageid(value interface{}) (attr.Value, error) {
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
